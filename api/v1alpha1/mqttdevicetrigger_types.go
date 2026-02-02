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

// MQTTDeviceTriggerSpec defines the desired state of MQTTDeviceTrigger.
type MQTTDeviceTriggerSpec struct {
	CommonSpec `json:",inline"`

	// Topic is the MQTT topic to subscribe to for trigger events
	Topic string `json:"topic"`

	// Type is the trigger type (e.g. button_short_press, button_long_press)
	Type string `json:"type"`

	// Subtype is the trigger subtype (e.g. button_1, turn_on)
	Subtype string `json:"subtype"`

	// Payload is the specific payload that triggers the automation
	// +optional
	Payload string `json:"payload,omitempty"`

	// ValueTemplate is the template to extract value from payload
	// +optional
	ValueTemplate string `json:"valueTemplate,omitempty"`

	// AutomationType is the automation type (always trigger)
	// +optional
	AutomationType string `json:"automationType,omitempty"`
}

// MQTTDeviceTriggerStatus defines the observed state of MQTTDeviceTrigger.
type MQTTDeviceTriggerStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTDeviceTrigger is the Schema for the mqttdevicetriggers API.
// It is a device automation trigger that fires when a specific MQTT message is received.
type MQTTDeviceTrigger struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTDeviceTriggerSpec   `json:"spec,omitempty"`
	Status MQTTDeviceTriggerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTDeviceTriggerList contains a list of MQTTDeviceTrigger.
type MQTTDeviceTriggerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTDeviceTrigger `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTDeviceTrigger{}, &MQTTDeviceTriggerList{})
}
