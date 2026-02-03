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

package payload

import (
	"encoding/json"
	"testing"
)

func TestCamelToSnake(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"name", "name"},
		{"commandTopic", "command_topic"},
		{"uniqueId", "unique_id"},
		{"jsonAttributesTopic", "json_attributes_topic"},
		{"deviceClass", "device_class"},
		{"stateClass", "state_class"},
		{"RGB", "r_g_b"}, // edge case
		{"", ""},
		{"A", "a"},
		{"AB", "a_b"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := camelToSnake(tt.input)
			if result != tt.expected {
				t.Errorf("camelToSnake(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBuilder_Set(t *testing.T) {
	b := New()
	b.Set("commandTopic", "home/button/cmd")
	b.Set("name", "Test Button")
	b.Set("emptyString", "")        // Should be skipped
	b.Set("emptySlice", []string{}) // Should be skipped

	data := b.BuildMap()

	if data["command_topic"] != "home/button/cmd" {
		t.Errorf("expected command_topic = 'home/button/cmd', got %v", data["command_topic"])
	}
	if data["name"] != "Test Button" {
		t.Errorf("expected name = 'Test Button', got %v", data["name"])
	}
	if _, exists := data["empty_string"]; exists {
		t.Error("empty string should not be included")
	}
	if _, exists := data["empty_slice"]; exists {
		t.Error("empty slice should not be included")
	}
}

func TestBuilder_SetWithPointers(t *testing.T) {
	b := New()

	trueVal := true
	falseVal := false
	intVal := 42
	floatVal := 3.14

	b.Set("enabledByDefault", &trueVal)
	b.Set("retain", &falseVal)
	b.Set("qos", &intVal)
	b.Set("minTemp", &floatVal)
	b.Set("nilBool", (*bool)(nil))

	data := b.BuildMap()

	if data["enabled_by_default"] != true {
		t.Errorf("expected enabled_by_default = true, got %v", data["enabled_by_default"])
	}
	if data["retain"] != false {
		t.Errorf("expected retain = false, got %v", data["retain"])
	}
	if data["qos"] != 42 {
		t.Errorf("expected qos = 42, got %v", data["qos"])
	}
	if data["min_temp"] != 3.14 {
		t.Errorf("expected min_temp = 3.14, got %v", data["min_temp"])
	}
	if _, exists := data["nil_bool"]; exists {
		t.Error("nil pointer should not be included")
	}
}

func TestBuilder_Build(t *testing.T) {
	b := New()
	b.Set("name", "Test")
	b.Set("commandTopic", "test/cmd")

	jsonBytes, err := b.Build()
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if result["name"] != "Test" {
		t.Errorf("expected name = 'Test', got %v", result["name"])
	}
	if result["command_topic"] != "test/cmd" {
		t.Errorf("expected command_topic = 'test/cmd', got %v", result["command_topic"])
	}
}

func TestDeviceBlockToMap(t *testing.T) {
	device := DeviceBlockToMap(
		"My Device",
		[]string{"device-001"},
		[][]string{{"mac", "aa:bb:cc:dd:ee:ff"}},
		"ACME",
		"Widget",
		"widget-v1",
		"SN12345",
		"1.0",
		"2.0",
		"Living Room",
		"http://192.168.1.100",
		"hub-001",
	)

	if device["name"] != "My Device" {
		t.Errorf("expected name = 'My Device', got %v", device["name"])
	}
	if device["manufacturer"] != "ACME" {
		t.Errorf("expected manufacturer = 'ACME', got %v", device["manufacturer"])
	}
	if device["model_id"] != "widget-v1" {
		t.Errorf("expected model_id = 'widget-v1', got %v", device["model_id"])
	}
}

func TestAvailabilityToMap(t *testing.T) {
	avail := AvailabilityToMap(
		"home/device/avail",
		"online",
		"offline",
		"{{ value_json.status }}",
	)

	if avail["topic"] != "home/device/avail" {
		t.Errorf("expected topic = 'home/device/avail', got %v", avail["topic"])
	}
	if avail["payload_available"] != "online" {
		t.Errorf("expected payload_available = 'online', got %v", avail["payload_available"])
	}
	if avail["value_template"] != "{{ value_json.status }}" {
		t.Errorf("expected value_template, got %v", avail["value_template"])
	}
}
