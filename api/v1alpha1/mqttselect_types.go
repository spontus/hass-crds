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

// MQTTSelectSpec defines the desired state of MQTTSelect.
type MQTTSelectSpec struct {
	CommonSpec `json:",inline"`

	// CommandTopic is the topic to publish selected option
	CommandTopic string `json:"commandTopic"`

	// Options is the list of selectable options
	Options []string `json:"options"`

	// CommandTemplate is the template for the command payload
	// +optional
	CommandTemplate string `json:"commandTemplate,omitempty"`

	// StateTopic is the topic to read current selection
	// +optional
	StateTopic string `json:"stateTopic,omitempty"`

	// ValueTemplate is the template to extract value from payload
	// +optional
	ValueTemplate string `json:"valueTemplate,omitempty"`

	// Optimistic indicates whether to assume state changes immediately
	// +optional
	Optimistic *bool `json:"optimistic,omitempty"`
}

// MQTTSelectStatus defines the observed state of MQTTSelect.
type MQTTSelectStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTSelect is the Schema for the mqttselects API.
// It is a dropdown selection entity with a fixed list of options.
type MQTTSelect struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTSelectSpec   `json:"spec,omitempty"`
	Status MQTTSelectStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTSelectList contains a list of MQTTSelect.
type MQTTSelectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTSelect `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTSelect{}, &MQTTSelectList{})
}
