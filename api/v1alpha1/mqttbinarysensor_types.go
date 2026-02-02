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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MQTTBinarySensorSpec defines the desired state of MQTTBinarySensor.
type MQTTBinarySensorSpec struct {
	CommonSpec `json:",inline"`

	// StateTopic is the topic to read sensor state
	StateTopic string `json:"stateTopic"`

	// ValueTemplate is the template to extract state from payload
	// +optional
	ValueTemplate string `json:"valueTemplate,omitempty"`

	// PayloadOn is the payload representing on/detected (default: ON)
	// +optional
	PayloadOn string `json:"payloadOn,omitempty"`

	// PayloadOff is the payload representing off/clear (default: OFF)
	// +optional
	PayloadOff string `json:"payloadOff,omitempty"`

	// DeviceClass is the HA device class (e.g. motion, door, window, moisture, smoke, occupancy)
	// +optional
	DeviceClass string `json:"deviceClass,omitempty"`

	// ExpireAfter is the seconds after which the state expires
	// +optional
	ExpireAfter *int `json:"expireAfter,omitempty"`

	// ForceUpdate indicates whether to update state even if unchanged
	// +optional
	ForceUpdate *bool `json:"forceUpdate,omitempty"`

	// OffDelay is the seconds after which the sensor auto-resets to off
	// +optional
	OffDelay *int `json:"offDelay,omitempty"`
}

// MQTTBinarySensorStatus defines the observed state of MQTTBinarySensor.
type MQTTBinarySensorStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTBinarySensor is the Schema for the mqttbinarysensors API.
// It is a read-only on/off sensor (e.g. motion detector, door contact).
type MQTTBinarySensor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTBinarySensorSpec   `json:"spec,omitempty"`
	Status MQTTBinarySensorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTBinarySensorList contains a list of MQTTBinarySensor.
type MQTTBinarySensorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTBinarySensor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTBinarySensor{}, &MQTTBinarySensorList{})
}
