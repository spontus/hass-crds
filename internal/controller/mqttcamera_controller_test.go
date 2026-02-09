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

package controller

import (
	"encoding/json"
	"testing"

	mqttv1alpha1 "github.com/spontus/hass-crds/api/v1alpha1"
	"github.com/spontus/hass-crds/internal/payload"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMQTTCameraPayload_MatchesLegacyFormat(t *testing.T) {
	// This test verifies that MQTTCamera can produce a payload matching the legacy format:
	// {
	//   "name": "Skaftared Weather Forecast Forecast",
	//   "unique_id": "2-2678647_2-2678647_forecast_chart",
	//   "topic": "auto/2-2678647/2-2678647_forecast_chart/state",
	//   "state_class": "measurement",
	//   "availability_topic": "mqttAuto/availability",
	//   "expire_after": 86400,
	//   "image_encoding": "b64",
	//   "device": {
	//     "manufacturer": "yr.no",
	//     "identifiers": ["2-2678647"],
	//     "name": "Skaftared Weather Forecast"
	//   }
	// }

	expireAfter := 86400
	camera := &mqttv1alpha1.MQTTCamera{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "weather-forecast-chart",
			Namespace: "hass-resources",
		},
		Spec: mqttv1alpha1.MQTTCameraSpec{
			CommonSpec: mqttv1alpha1.CommonSpec{
				EntityMetadata: mqttv1alpha1.EntityMetadata{
					Name:     "Skaftared Weather Forecast Forecast",
					UniqueId: "2-2678647_2-2678647_forecast_chart",
				},
				Device: &mqttv1alpha1.DeviceBlock{
					Name:         "Skaftared Weather Forecast",
					Manufacturer: "yr.no",
					Identifiers:  []string{"2-2678647"},
				},
				AvailabilityTopic: "mqttAuto/availability",
			},
			Topic:         "auto/2-2678647/2-2678647_forecast_chart/state",
			ImageEncoding: "b64",
			StateClass:    "measurement",
			ExpireAfter:   &expireAfter,
		},
	}

	wrapper := &mqttCameraWrapper{camera}
	spec := camera.Spec

	r := &MQTTCameraReconciler{}
	pb, err := r.buildPayload(wrapper, "2-2678647_2-2678647_forecast_chart")
	if err != nil {
		t.Fatalf("buildPayload failed: %v", err)
	}

	// Add unique_id (normally done by base reconciler)
	pb.Set("uniqueId", "2-2678647_2-2678647_forecast_chart")

	// Add device block (normally done by base reconciler)
	if spec.Device != nil {
		device := payload.DeviceBlockToMap(
			spec.Device.Name,
			spec.Device.Identifiers,
			spec.Device.Connections,
			spec.Device.Manufacturer,
			spec.Device.Model,
			spec.Device.ModelId,
			spec.Device.SerialNumber,
			spec.Device.HwVersion,
			spec.Device.SwVersion,
			spec.Device.SuggestedArea,
			spec.Device.ConfigurationUrl,
			spec.Device.ViaDevice,
		)
		pb.SetDevice(device)
	}

	jsonBytes, err := pb.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Verify required fields
	tests := []struct {
		key      string
		expected interface{}
	}{
		{"name", "Skaftared Weather Forecast Forecast"},
		{"unique_id", "2-2678647_2-2678647_forecast_chart"},
		{"topic", "auto/2-2678647/2-2678647_forecast_chart/state"},
		{"state_class", "measurement"},
		{"availability_topic", "mqttAuto/availability"},
		{"expire_after", float64(86400)}, // JSON numbers are float64
		{"image_encoding", "b64"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if result[tt.key] != tt.expected {
				t.Errorf("%s = %v, want %v", tt.key, result[tt.key], tt.expected)
			}
		})
	}

	// Verify device block
	device, ok := result["device"].(map[string]interface{})
	if !ok {
		t.Fatal("device block not found or not an object")
	}

	if device["name"] != "Skaftared Weather Forecast" {
		t.Errorf("device.name = %v, want 'Skaftared Weather Forecast'", device["name"])
	}
	if device["manufacturer"] != "yr.no" {
		t.Errorf("device.manufacturer = %v, want 'yr.no'", device["manufacturer"])
	}

	identifiers, ok := device["identifiers"].([]interface{})
	if !ok || len(identifiers) != 1 || identifiers[0] != "2-2678647" {
		t.Errorf("device.identifiers = %v, want ['2-2678647']", device["identifiers"])
	}
}
