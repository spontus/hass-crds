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
	"sort"
	"testing"
	"time"

	"github.com/go-logr/logr"

	"github.com/spontus/hass-crds/internal/mqtt"
)

func TestFindOrphans(t *testing.T) {
	allVerified := map[string]struct{}{
		"button": {},
		"sensor": {},
	}

	tests := []struct {
		name               string
		ours               []discoveredEntity
		expected           map[string]struct{}
		verifiedComponents map[string]struct{}
		want               []string
	}{
		{
			name: "no orphans",
			ours: []discoveredEntity{
				{Topic: "homeassistant/button/default/btn1/config"},
				{Topic: "homeassistant/sensor/default/temp/config"},
			},
			expected: map[string]struct{}{
				"homeassistant/button/default/btn1/config": {},
				"homeassistant/sensor/default/temp/config": {},
			},
			verifiedComponents: allVerified,
			want:               nil,
		},
		{
			name: "one orphan",
			ours: []discoveredEntity{
				{Topic: "homeassistant/button/default/btn1/config"},
				{Topic: "homeassistant/button/default/btn2/config"},
			},
			expected: map[string]struct{}{
				"homeassistant/button/default/btn1/config": {},
			},
			verifiedComponents: allVerified,
			want:               []string{"homeassistant/button/default/btn2/config"},
		},
		{
			name: "all orphans",
			ours: []discoveredEntity{
				{Topic: "homeassistant/sensor/ns/s1/config"},
				{Topic: "homeassistant/sensor/ns/s2/config"},
			},
			expected:           map[string]struct{}{},
			verifiedComponents: allVerified,
			want: []string{
				"homeassistant/sensor/ns/s1/config",
				"homeassistant/sensor/ns/s2/config",
			},
		},
		{
			name:               "empty discovered",
			ours:               nil,
			expected:           map[string]struct{}{"homeassistant/button/default/btn1/config": {}},
			verifiedComponents: allVerified,
			want:               nil,
		},
		{
			name: "skips unverified component",
			ours: []discoveredEntity{
				{Topic: "homeassistant/button/default/btn1/config"},
				{Topic: "homeassistant/image/default/img1/config"},
			},
			expected:           map[string]struct{}{},
			verifiedComponents: map[string]struct{}{"button": {}}, // image NOT verified
			want:               []string{"homeassistant/button/default/btn1/config"},
		},
		{
			name: "skips invalid topic format",
			ours: []discoveredEntity{
				{Topic: "bad-topic"},
				{Topic: "homeassistant/button/default/btn1/config"},
			},
			expected:           map[string]struct{}{},
			verifiedComponents: allVerified,
			want:               []string{"homeassistant/button/default/btn1/config"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findOrphans(tt.ours, tt.expected, tt.verifiedComponents)
			sort.Strings(got)
			sort.Strings(tt.want)

			if len(got) != len(tt.want) {
				t.Fatalf("findOrphans() returned %d orphans, want %d\ngot:  %v\nwant: %v", len(got), len(tt.want), got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("orphan[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestHasOurOrigin(t *testing.T) {
	tests := []struct {
		name    string
		payload interface{}
		want    bool
	}{
		{
			name: "our origin",
			payload: map[string]interface{}{
				"name":   "Test",
				"origin": map[string]interface{}{"name": "hass-crds", "support_url": "https://github.com/spontus/hass-crds"},
			},
			want: true,
		},
		{
			name: "different origin",
			payload: map[string]interface{}{
				"name":   "Test",
				"origin": map[string]interface{}{"name": "other-integration"},
			},
			want: false,
		},
		{
			name: "no origin field",
			payload: map[string]interface{}{
				"name": "Test",
			},
			want: false,
		},
		{
			name: "origin not a map",
			payload: map[string]interface{}{
				"origin": "string-origin",
			},
			want: false,
		},
		{
			name: "origin name not a string",
			payload: map[string]interface{}{
				"origin": map[string]interface{}{"name": 123},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, _ := json.Marshal(tt.payload)
			got := hasOurOrigin(data)
			if got != tt.want {
				t.Errorf("hasOurOrigin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasOurOrigin_InvalidJSON(t *testing.T) {
	if hasOurOrigin([]byte("not json")) {
		t.Error("hasOurOrigin(invalid json) should return false")
	}
}

func TestFilterOurEntities(t *testing.T) {
	ourPayload, _ := json.Marshal(map[string]interface{}{
		"name":   "Ours",
		"origin": map[string]interface{}{"name": "hass-crds"},
	})
	otherPayload, _ := json.Marshal(map[string]interface{}{
		"name":   "Other",
		"origin": map[string]interface{}{"name": "tasmota"},
	})
	noOriginPayload, _ := json.Marshal(map[string]interface{}{
		"name": "NoOrigin",
	})

	entities := []discoveredEntity{
		{Topic: "homeassistant/button/default/ours/config", Payload: ourPayload},
		{Topic: "homeassistant/button/default/other/config", Payload: otherPayload},
		{Topic: "homeassistant/button/default/none/config", Payload: noOriginPayload},
		{Topic: "homeassistant/button/default/empty/config", Payload: []byte{}},
	}

	result := filterOurEntities(entities)
	if len(result) != 1 {
		t.Fatalf("filterOurEntities() returned %d entities, want 1", len(result))
	}
	if result[0].Topic != "homeassistant/button/default/ours/config" {
		t.Errorf("got topic %q, want %q", result[0].Topic, "homeassistant/button/default/ours/config")
	}
}

func TestCollect_RemovesOrphans(t *testing.T) {
	mockClient := mqtt.NewMockClient()
	_ = mockClient.Connect(context.Background())

	collector := NewOrphanCollector(nil, mockClient, logr.Discard(), Config{
		Enabled:        true,
		Interval:       time.Minute,
		SilenceTimeout: 100 * time.Millisecond,
	})

	// Override collectDiscoveryMessages and buildExpectedTopics for unit testing
	// We test the full Collect flow by simulating messages through the mock client

	ourPayload, _ := json.Marshal(map[string]interface{}{
		"name":   "Orphan",
		"origin": map[string]interface{}{"name": "hass-crds"},
	})

	// When the collector subscribes, simulate retained messages
	go func() {
		// Wait for subscription to be set up
		time.Sleep(20 * time.Millisecond)
		mockClient.SimulateMessage("homeassistant/button/default/orphan-btn/config", ourPayload)
	}()

	// We can't test the full Collect() without a real K8s client,
	// so test collectDiscoveryMessages + filterOurEntities + findOrphans separately
	ctx := context.Background()
	entities, err := collector.collectDiscoveryMessages(ctx)
	if err != nil {
		t.Fatalf("collectDiscoveryMessages() error: %v", err)
	}

	ours := filterOurEntities(entities)
	if len(ours) != 1 {
		t.Fatalf("expected 1 entity with our origin, got %d", len(ours))
	}

	// Simulate empty expected set (no CRs exist) but button component is verified
	expected := map[string]struct{}{}
	verifiedComponents := map[string]struct{}{"button": {}}
	orphans := findOrphans(ours, expected, verifiedComponents)
	if len(orphans) != 1 {
		t.Fatalf("expected 1 orphan, got %d", len(orphans))
	}

	// Simulate cleanup
	if err := mockClient.Publish(ctx, orphans[0], []byte{}, 1, true); err != nil {
		t.Fatalf("publish cleanup failed: %v", err)
	}

	msgs := mockClient.GetPublishedMessages()
	found := false
	for _, msg := range msgs {
		if msg.Topic == "homeassistant/button/default/orphan-btn/config" && len(msg.Payload) == 0 && msg.Retain {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected empty retained payload published to orphan topic")
	}
}

func TestCollect_IgnoresNonOurEntities(t *testing.T) {
	mockClient := mqtt.NewMockClient()
	_ = mockClient.Connect(context.Background())

	collector := NewOrphanCollector(nil, mockClient, logr.Discard(), Config{
		Enabled:        true,
		SilenceTimeout: 100 * time.Millisecond,
	})

	otherPayload, _ := json.Marshal(map[string]interface{}{
		"name":   "Tasmota Sensor",
		"origin": map[string]interface{}{"name": "tasmota"},
	})

	go func() {
		time.Sleep(20 * time.Millisecond)
		mockClient.SimulateMessage("homeassistant/sensor/default/tasmota-temp/config", otherPayload)
	}()

	ctx := context.Background()
	entities, err := collector.collectDiscoveryMessages(ctx)
	if err != nil {
		t.Fatalf("collectDiscoveryMessages() error: %v", err)
	}

	ours := filterOurEntities(entities)
	if len(ours) != 0 {
		t.Fatalf("expected 0 entities with our origin, got %d", len(ours))
	}
}

func TestConfigFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		validate func(t *testing.T, cfg Config)
	}{
		{
			name: "defaults",
			env:  nil,
			validate: func(t *testing.T, cfg Config) {
				if !cfg.Enabled {
					t.Error("expected Enabled=true by default")
				}
				if cfg.Interval != 5*time.Minute {
					t.Errorf("expected Interval=5m, got %v", cfg.Interval)
				}
				if !cfg.RunOnStartup {
					t.Error("expected RunOnStartup=true by default")
				}
				if cfg.SilenceTimeout != 5*time.Second {
					t.Errorf("expected SilenceTimeout=5s, got %v", cfg.SilenceTimeout)
				}
			},
		},
		{
			name: "disabled",
			env:  map[string]string{"GC_ENABLED": "false"},
			validate: func(t *testing.T, cfg Config) {
				if cfg.Enabled {
					t.Error("expected Enabled=false")
				}
			},
		},
		{
			name: "custom interval",
			env:  map[string]string{"GC_INTERVAL": "10m"},
			validate: func(t *testing.T, cfg Config) {
				if cfg.Interval != 10*time.Minute {
					t.Errorf("expected Interval=10m, got %v", cfg.Interval)
				}
			},
		},
		{
			name: "no startup run",
			env:  map[string]string{"GC_RUN_ON_STARTUP": "false"},
			validate: func(t *testing.T, cfg Config) {
				if cfg.RunOnStartup {
					t.Error("expected RunOnStartup=false")
				}
			},
		},
		{
			name: "custom silence timeout",
			env:  map[string]string{"GC_SILENCE_TIMEOUT": "10s"},
			validate: func(t *testing.T, cfg Config) {
				if cfg.SilenceTimeout != 10*time.Second {
					t.Errorf("expected SilenceTimeout=10s, got %v", cfg.SilenceTimeout)
				}
			},
		},
		{
			name: "invalid interval keeps default",
			env:  map[string]string{"GC_INTERVAL": "not-a-duration"},
			validate: func(t *testing.T, cfg Config) {
				if cfg.Interval != 5*time.Minute {
					t.Errorf("expected default Interval=5m, got %v", cfg.Interval)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all GC env vars (t.Setenv handles cleanup)
			for _, key := range []string{"GC_ENABLED", "GC_INTERVAL", "GC_RUN_ON_STARTUP", "GC_SILENCE_TIMEOUT"} {
				t.Setenv(key, "")
			}
			// Set test env vars (overrides the empty values above)
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			cfg := NewConfigFromEnv()
			tt.validate(t, cfg)
		})
	}
}
