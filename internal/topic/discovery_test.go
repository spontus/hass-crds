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
	"testing"
)

func TestParseDiscoveryTopic(t *testing.T) {
	tests := []struct {
		name      string
		topic     string
		want      *DiscoveryTopicInfo
		expectErr bool
	}{
		{
			name:  "valid button topic",
			topic: "homeassistant/button/default/my-button/config",
			want: &DiscoveryTopicInfo{
				Prefix:    "homeassistant",
				Component: "button",
				Namespace: "default",
				Name:      "my-button",
			},
		},
		{
			name:  "valid sensor topic",
			topic: "homeassistant/sensor/monitoring/temp-sensor/config",
			want: &DiscoveryTopicInfo{
				Prefix:    "homeassistant",
				Component: "sensor",
				Namespace: "monitoring",
				Name:      "temp-sensor",
			},
		},
		{
			name:  "valid binary_sensor topic",
			topic: "homeassistant/binary_sensor/ns/door/config",
			want: &DiscoveryTopicInfo{
				Prefix:    "homeassistant",
				Component: "binary_sensor",
				Namespace: "ns",
				Name:      "door",
			},
		},
		{
			name:      "too few parts",
			topic:     "homeassistant/button/config",
			expectErr: true,
		},
		{
			name:      "too many parts",
			topic:     "homeassistant/button/ns/name/extra/config",
			expectErr: true,
		},
		{
			name:      "missing config suffix",
			topic:     "homeassistant/button/ns/name/state",
			expectErr: true,
		},
		{
			name:      "empty string",
			topic:     "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDiscoveryTopic(tt.topic)
			if tt.expectErr {
				if err == nil {
					t.Errorf("ParseDiscoveryTopic(%q) expected error, got nil", tt.topic)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseDiscoveryTopic(%q) unexpected error: %v", tt.topic, err)
			}
			if got.Prefix != tt.want.Prefix {
				t.Errorf("Prefix = %q, want %q", got.Prefix, tt.want.Prefix)
			}
			if got.Component != tt.want.Component {
				t.Errorf("Component = %q, want %q", got.Component, tt.want.Component)
			}
			if got.Namespace != tt.want.Namespace {
				t.Errorf("Namespace = %q, want %q", got.Namespace, tt.want.Namespace)
			}
			if got.Name != tt.want.Name {
				t.Errorf("Name = %q, want %q", got.Name, tt.want.Name)
			}
		})
	}
}

func TestComponentToKind(t *testing.T) {
	// Verify the reverse map is correctly generated
	for kind, component := range ComponentMapping {
		gotKind, ok := ComponentToKind[component]
		if !ok {
			t.Errorf("ComponentToKind missing component %q (kind %q)", component, kind)
			continue
		}
		if gotKind != kind {
			t.Errorf("ComponentToKind[%q] = %q, want %q", component, gotKind, kind)
		}
	}

	// Verify lengths match
	if len(ComponentToKind) != len(ComponentMapping) {
		t.Errorf("ComponentToKind has %d entries, ComponentMapping has %d", len(ComponentToKind), len(ComponentMapping))
	}
}

func TestDiscoveryTopicRoundTrip(t *testing.T) {
	// Verify that generating a topic and parsing it back gives consistent results
	tests := []struct {
		kind      string
		namespace string
		name      string
	}{
		{"MQTTButton", "default", "my-button"},
		{"MQTTSensor", "monitoring", "temp"},
		{"MQTTBinarySensor", "home", "door-sensor"},
	}

	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			topic := DiscoveryTopic(tt.kind, tt.namespace, tt.name)
			info, err := ParseDiscoveryTopic(topic)
			if err != nil {
				t.Fatalf("ParseDiscoveryTopic(%q) failed: %v", topic, err)
			}
			if info.Namespace != tt.namespace {
				t.Errorf("Namespace = %q, want %q", info.Namespace, tt.namespace)
			}
			if info.Name != tt.name {
				t.Errorf("Name = %q, want %q", info.Name, tt.name)
			}
			expectedComponent := ComponentMapping[tt.kind]
			if info.Component != expectedComponent {
				t.Errorf("Component = %q, want %q", info.Component, expectedComponent)
			}
		})
	}
}
