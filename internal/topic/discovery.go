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

package topic

import (
	"fmt"
	"strings"
)

// ComponentToKind is a reverse mapping from Home Assistant component types to Kubernetes kinds.
// Generated automatically from ComponentMapping.
var ComponentToKind map[string]string

func init() {
	ComponentToKind = make(map[string]string, len(ComponentMapping))
	for kind, component := range ComponentMapping {
		ComponentToKind[component] = kind
	}
}

// DiscoveryTopicInfo holds parsed information from a discovery topic.
type DiscoveryTopicInfo struct {
	Prefix    string
	Component string
	Namespace string
	Name      string
}

// ParseDiscoveryTopic parses a discovery topic string into its components.
// Expected format: <prefix>/<component>/<namespace>/<name>/config
func ParseDiscoveryTopic(topic string) (*DiscoveryTopicInfo, error) {
	parts := strings.Split(topic, "/")
	if len(parts) != 5 || parts[4] != "config" {
		return nil, fmt.Errorf("invalid discovery topic format: %s", topic)
	}

	return &DiscoveryTopicInfo{
		Prefix:    parts[0],
		Component: parts[1],
		Namespace: parts[2],
		Name:      parts[3],
	}, nil
}

// DefaultDiscoveryPrefix is the default Home Assistant MQTT discovery prefix.
const DefaultDiscoveryPrefix = "homeassistant"

// ComponentMapping maps Kubernetes kinds to Home Assistant component types.
var ComponentMapping = map[string]string{
	"MQTTButton":            "button",
	"MQTTSwitch":            "switch",
	"MQTTSensor":            "sensor",
	"MQTTBinarySensor":      "binary_sensor",
	"MQTTNumber":            "number",
	"MQTTSelect":            "select",
	"MQTTText":              "text",
	"MQTTScene":             "scene",
	"MQTTTag":               "tag",
	"MQTTLight":             "light",
	"MQTTCover":             "cover",
	"MQTTLock":              "lock",
	"MQTTValve":             "valve",
	"MQTTFan":               "fan",
	"MQTTSiren":             "siren",
	"MQTTCamera":            "camera",
	"MQTTImage":             "image",
	"MQTTNotify":            "notify",
	"MQTTUpdate":            "update",
	"MQTTClimate":           "climate",
	"MQTTHumidifier":        "humidifier",
	"MQTTWaterHeater":       "water_heater",
	"MQTTVacuum":            "vacuum",
	"MQTTLawnMower":         "lawn_mower",
	"MQTTAlarmControlPanel": "alarm_control_panel",
	"MQTTDeviceTracker":     "device_tracker",
	"MQTTDeviceTrigger":     "device_automation",
	"MQTTEvent":             "event",
}

// DiscoveryTopic generates the MQTT discovery topic for an entity.
// Format: <prefix>/<component>/<node_id>/<object_id>/config
func DiscoveryTopic(kind, namespace, name string) string {
	return DiscoveryTopicWithPrefix(DefaultDiscoveryPrefix, kind, namespace, name)
}

// DiscoveryTopicWithPrefix generates the discovery topic with a custom prefix.
func DiscoveryTopicWithPrefix(prefix, kind, namespace, name string) string {
	component, ok := ComponentMapping[kind]
	if !ok {
		// Default to lowercase kind with "mqtt" prefix removed
		component = strings.ToLower(strings.TrimPrefix(kind, "MQTT"))
	}

	// Node ID uses namespace to ensure uniqueness across namespaces
	nodeID := namespace

	// Object ID uses the resource name
	objectID := name

	return fmt.Sprintf("%s/%s/%s/%s/config", prefix, component, nodeID, objectID)
}

// UniqueID generates a unique identifier for Home Assistant entity registry.
// Format: <namespace>-<name>
func UniqueID(namespace, name string) string {
	return fmt.Sprintf("%s-%s", namespace, name)
}

// UniqueIDWithOverride returns the provided uniqueID if non-empty,
// otherwise generates one from namespace and name.
func UniqueIDWithOverride(uniqueID, namespace, name string) string {
	if uniqueID != "" {
		return uniqueID
	}
	return UniqueID(namespace, name)
}
