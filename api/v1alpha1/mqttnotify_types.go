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

// MQTTNotifySpec defines the desired state of MQTTNotify.
type MQTTNotifySpec struct {
	CommonSpec `json:",inline"`

	// CommandTopic is the topic to publish notification messages
	CommandTopic string `json:"commandTopic"`

	// CommandTemplate is the template for the notification payload
	// +optional
	CommandTemplate string `json:"commandTemplate,omitempty"`
}

// MQTTNotifyStatus defines the observed state of MQTTNotify.
type MQTTNotifyStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTNotify is the Schema for the mqttnotifys API.
// It is a notification service entity that sends messages to a device via MQTT.
type MQTTNotify struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTNotifySpec   `json:"spec,omitempty"`
	Status MQTTNotifyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTNotifyList contains a list of MQTTNotify.
type MQTTNotifyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTNotify `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTNotify{}, &MQTTNotifyList{})
}
