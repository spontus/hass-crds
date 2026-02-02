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

// MQTTNumberSpec defines the desired state of MQTTNumber.
type MQTTNumberSpec struct {
	CommonSpec `json:",inline"`

	// CommandTopic is the topic to publish number value
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

	// Min is the minimum value (default: 1)
	// +optional
	Min *float64 `json:"min,omitempty"`

	// Max is the maximum value (default: 100)
	// +optional
	Max *float64 `json:"max,omitempty"`

	// Step is the step size (default: 1)
	// +optional
	Step *float64 `json:"step,omitempty"`

	// Mode is the UI mode
	// +kubebuilder:validation:Enum=auto;box;slider
	// +optional
	Mode string `json:"mode,omitempty"`

	// UnitOfMeasurement is the unit displayed in HA
	// +optional
	UnitOfMeasurement string `json:"unitOfMeasurement,omitempty"`

	// DeviceClass is the HA device class (e.g. temperature, humidity, power_factor)
	// +optional
	DeviceClass string `json:"deviceClass,omitempty"`

	// Optimistic indicates whether to assume state changes immediately
	// +optional
	Optimistic *bool `json:"optimistic,omitempty"`
}

// MQTTNumberStatus defines the observed state of MQTTNumber.
type MQTTNumberStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTNumber is the Schema for the mqttnumbers API.
// It is a numeric input entity with min/max bounds and step size.
type MQTTNumber struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTNumberSpec   `json:"spec,omitempty"`
	Status MQTTNumberStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTNumberList contains a list of MQTTNumber.
type MQTTNumberList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTNumber `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTNumber{}, &MQTTNumberList{})
}
