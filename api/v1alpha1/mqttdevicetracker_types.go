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

// MQTTDeviceTrackerSpec defines the desired state of MQTTDeviceTracker.
type MQTTDeviceTrackerSpec struct {
	CommonSpec `json:",inline"`

	// StateTopic is the topic to read tracker state (home/not_home or zone name)
	StateTopic string `json:"stateTopic"`

	// ValueTemplate is the template to extract state from payload
	// +optional
	ValueTemplate string `json:"valueTemplate,omitempty"`

	// PayloadHome is the payload representing home (default: home)
	// +optional
	PayloadHome string `json:"payloadHome,omitempty"`

	// PayloadNotHome is the payload representing not home (default: not_home)
	// +optional
	PayloadNotHome string `json:"payloadNotHome,omitempty"`

	// PayloadReset is the payload that resets the tracker to unknown
	// +optional
	PayloadReset string `json:"payloadReset,omitempty"`

	// SourceType is the source type (e.g. gps, router, bluetooth, bluetooth_le)
	// +optional
	SourceType string `json:"sourceType,omitempty"`
}

// MQTTDeviceTrackerStatus defines the observed state of MQTTDeviceTracker.
type MQTTDeviceTrackerStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTDeviceTracker is the Schema for the mqttdevicetrackers API.
// It is a device tracker entity for presence detection and location tracking via MQTT.
type MQTTDeviceTracker struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTDeviceTrackerSpec   `json:"spec,omitempty"`
	Status MQTTDeviceTrackerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTDeviceTrackerList contains a list of MQTTDeviceTracker.
type MQTTDeviceTrackerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTDeviceTracker `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTDeviceTracker{}, &MQTTDeviceTrackerList{})
}
