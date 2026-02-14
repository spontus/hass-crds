/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gc

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/spontus/hass-crds/internal/mqtt"
	"github.com/spontus/hass-crds/internal/payload"
	"github.com/spontus/hass-crds/internal/topic"
)

// Config holds configuration for the OrphanCollector.
type Config struct {
	Enabled        bool
	Interval       time.Duration
	RunOnStartup   bool
	SilenceTimeout time.Duration
}

// NewConfigFromEnv creates a Config from environment variables.
func NewConfigFromEnv() Config {
	cfg := Config{
		Enabled:        true,
		Interval:       5 * time.Minute,
		RunOnStartup:   true,
		SilenceTimeout: 5 * time.Second,
	}

	if v := os.Getenv("GC_ENABLED"); v == "false" {
		cfg.Enabled = false
	}

	if v := os.Getenv("GC_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Interval = d
		}
	}

	if v := os.Getenv("GC_RUN_ON_STARTUP"); v == "false" {
		cfg.RunOnStartup = false
	}

	if v := os.Getenv("GC_SILENCE_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.SilenceTimeout = d
		}
	}

	return cfg
}

// OrphanCollector detects and removes orphaned MQTT discovery entities
// that have no corresponding Kubernetes Custom Resource.
type OrphanCollector struct {
	k8sClient  client.Client
	mqttClient mqtt.Client
	log        logr.Logger
	config     Config
}

// NewOrphanCollector creates a new OrphanCollector.
func NewOrphanCollector(k8sClient client.Client, mqttClient mqtt.Client, log logr.Logger, config Config) *OrphanCollector {
	return &OrphanCollector{
		k8sClient:  k8sClient,
		mqttClient: mqttClient,
		log:        log.WithName("gc"),
		config:     config,
	}
}

// Start implements manager.Runnable. It runs the garbage collection loop.
func (c *OrphanCollector) Start(ctx context.Context) error {
	if !c.config.Enabled {
		c.log.Info("Orphan garbage collector is disabled")
		return nil
	}

	c.log.Info("Starting orphan garbage collector",
		"interval", c.config.Interval,
		"runOnStartup", c.config.RunOnStartup,
		"silenceTimeout", c.config.SilenceTimeout,
	)

	// Wait for caches to sync before first run
	select {
	case <-ctx.Done():
		return nil
	case <-time.After(10 * time.Second):
	}

	if c.config.RunOnStartup {
		if err := c.Collect(ctx); err != nil {
			c.log.Error(err, "Initial garbage collection failed")
		}
	}

	ticker := time.NewTicker(c.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			c.log.Info("Stopping orphan garbage collector")
			return nil
		case <-ticker.C:
			if err := c.Collect(ctx); err != nil {
				c.log.Error(err, "Garbage collection cycle failed")
			}
		}
	}
}

// discoveredEntity represents an entity found via MQTT discovery.
type discoveredEntity struct {
	Topic   string
	Payload []byte
}

// Collect runs a single garbage collection cycle.
func (c *OrphanCollector) Collect(ctx context.Context) error {
	c.log.V(1).Info("Starting garbage collection cycle")

	// Step 1: Subscribe and collect retained discovery messages
	entities, err := c.collectDiscoveryMessages(ctx)
	if err != nil {
		return fmt.Errorf("collecting discovery messages: %w", err)
	}

	// Step 2: Filter to only entities we created (origin.name == "hass-crds")
	ours := filterOurEntities(entities)
	if len(ours) == 0 {
		c.log.V(1).Info("No entities with our origin found")
		return nil
	}

	c.log.V(1).Info("Found entities with our origin", "count", len(ours))

	// Step 3: Build set of expected discovery topics from existing CRs
	expected := c.buildExpectedTopics(ctx)

	// Step 4: Find orphans (topics with our origin that are not in expected set)
	orphans := findOrphans(ours, expected)
	if len(orphans) == 0 {
		c.log.V(1).Info("No orphaned entities found")
		return nil
	}

	c.log.Info("Found orphaned entities", "count", len(orphans))

	// Step 5: Publish empty payloads to remove orphans
	for _, orphanTopic := range orphans {
		c.log.Info("Removing orphaned entity", "topic", orphanTopic)
		if err := c.mqttClient.Publish(ctx, orphanTopic, []byte{}, 1, true); err != nil {
			c.log.Error(err, "Failed to remove orphaned entity", "topic", orphanTopic)
		}
	}

	return nil
}

// collectDiscoveryMessages subscribes to discovery topics and collects retained messages.
func (c *OrphanCollector) collectDiscoveryMessages(ctx context.Context) ([]discoveredEntity, error) {
	var mu sync.Mutex
	var entities []discoveredEntity

	subscriptionTopic := topic.DefaultDiscoveryPrefix + "/+/+/+/config"

	err := c.mqttClient.Subscribe(ctx, subscriptionTopic, 0, func(t string, p []byte) {
		mu.Lock()
		entities = append(entities, discoveredEntity{Topic: t, Payload: p})
		mu.Unlock()
	})
	if err != nil {
		return nil, fmt.Errorf("subscribing to %s: %w", subscriptionTopic, err)
	}

	// Wait for retained messages to arrive. Reset timer on each new message.
	lastCount := -1
	timer := time.NewTimer(c.config.SilenceTimeout)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			_ = c.mqttClient.Unsubscribe(ctx, subscriptionTopic)
			return nil, ctx.Err()
		case <-timer.C:
			mu.Lock()
			currentCount := len(entities)
			mu.Unlock()

			if currentCount == lastCount || lastCount == -1 && currentCount == 0 {
				// No new messages arrived during the silence window
				_ = c.mqttClient.Unsubscribe(context.Background(), subscriptionTopic)
				return entities, nil
			}

			// New messages arrived, reset timer
			lastCount = currentCount
			timer.Reset(c.config.SilenceTimeout)
		}
	}
}

// filterOurEntities returns only entities whose payload has origin.name == "hass-crds".
func filterOurEntities(entities []discoveredEntity) []discoveredEntity {
	var result []discoveredEntity
	for _, e := range entities {
		if len(e.Payload) == 0 {
			continue
		}
		if hasOurOrigin(e.Payload) {
			result = append(result, e)
		}
	}
	return result
}

// hasOurOrigin checks if a JSON payload has origin.name matching our origin name.
func hasOurOrigin(data []byte) bool {
	var p map[string]interface{}
	if err := json.Unmarshal(data, &p); err != nil {
		return false
	}

	origin, ok := p["origin"]
	if !ok {
		return false
	}

	originMap, ok := origin.(map[string]interface{})
	if !ok {
		return false
	}

	name, ok := originMap["name"].(string)
	return ok && name == payload.OriginName
}

// buildExpectedTopics lists all CRs and returns the set of discovery topics that should exist.
func (c *OrphanCollector) buildExpectedTopics(ctx context.Context) map[string]struct{} {
	expected := make(map[string]struct{})

	for kind, component := range topic.ComponentMapping {
		resource := kindToResource(kind)

		gvr := schema.GroupVersionResource{
			Group:    "mqtt.home-assistant.io",
			Version:  "v1alpha1",
			Resource: resource,
		}

		list := &unstructured.UnstructuredList{}
		list.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   gvr.Group,
			Version: gvr.Version,
			Kind:    kind + "List",
		})

		if err := c.k8sClient.List(ctx, list); err != nil {
			c.log.V(1).Info("Failed to list CRs, skipping", "kind", kind, "error", err)
			continue
		}

		for _, item := range list.Items {
			t := fmt.Sprintf("%s/%s/%s/%s/config",
				topic.DefaultDiscoveryPrefix,
				component,
				item.GetNamespace(),
				item.GetName(),
			)
			expected[t] = struct{}{}
		}
	}

	return expected
}

// findOrphans returns topics from discovered entities that are not in the expected set.
func findOrphans(ours []discoveredEntity, expected map[string]struct{}) []string {
	var orphans []string
	for _, e := range ours {
		if _, ok := expected[e.Topic]; !ok {
			orphans = append(orphans, e.Topic)
		}
	}
	return orphans
}

// kindToResource converts a CRD kind to its plural resource name.
// e.g., "MQTTButton" -> "mqttbuttons", "MQTTBinarySensor" -> "mqttbinarysensors"
func kindToResource(kind string) string {
	return strings.ToLower(kind) + "s"
}
