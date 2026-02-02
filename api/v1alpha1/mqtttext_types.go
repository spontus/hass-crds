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

// MQTTTextSpec defines the desired state of MQTTText.
type MQTTTextSpec struct {
	CommonSpec `json:",inline"`

	// CommandTopic is the topic to publish text value
	CommandTopic string `json:"commandTopic"`

	// CommandTemplate is the template for the command payload
	// +optional
	CommandTemplate string `json:"commandTemplate,omitempty"`

	// StateTopic is the topic to read current value
	// +optional
	StateTopic string `json:"stateTopic,omitempty"`

	// ValueTemplate is the template to extract value from payload
	// +optional
	ValueTemplate string `json:"valueTemplate,omitempty"`

	// Min is the minimum text length (default: 0)
	// +optional
	Min *int `json:"min,omitempty"`

	// Max is the maximum text length (default: 255)
	// +optional
	Max *int `json:"max,omitempty"`

	// Pattern is the regex pattern for validation
	// +optional
	Pattern string `json:"pattern,omitempty"`

	// Mode is the input mode
	// +kubebuilder:validation:Enum=text;password
	// +optional
	Mode string `json:"mode,omitempty"`
}

// MQTTTextStatus defines the observed state of MQTTText.
type MQTTTextStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTText is the Schema for the mqtttexts API.
// It is a free-text input entity.
type MQTTText struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTTextSpec   `json:"spec,omitempty"`
	Status MQTTTextStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTTextList contains a list of MQTTText.
type MQTTTextList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTText `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTText{}, &MQTTTextList{})
}
