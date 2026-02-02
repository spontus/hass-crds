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

// MQTTEventSpec defines the desired state of MQTTEvent.
type MQTTEventSpec struct {
	CommonSpec `json:",inline"`

	// StateTopic is the topic to subscribe to for events
	StateTopic string `json:"stateTopic"`

	// EventTypes is the list of supported event types
	EventTypes []string `json:"eventTypes"`

	// ValueTemplate is the template to extract event type from payload
	// +optional
	ValueTemplate string `json:"valueTemplate,omitempty"`

	// DeviceClass is the event device class
	// +kubebuilder:validation:Enum=button;doorbell;motion
	// +optional
	DeviceClass string `json:"deviceClass,omitempty"`
}

// MQTTEventStatus defines the observed state of MQTTEvent.
type MQTTEventStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTEvent is the Schema for the mqttevents API.
// It is an event entity for stateless events such as button presses or doorbell rings.
type MQTTEvent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTEventSpec   `json:"spec,omitempty"`
	Status MQTTEventStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTEventList contains a list of MQTTEvent.
type MQTTEventList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTEvent `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTEvent{}, &MQTTEventList{})
}
