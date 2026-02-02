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

// MQTTValveSpec defines the desired state of MQTTValve.
type MQTTValveSpec struct {
	CommonSpec `json:",inline"`

	// CommandTopic is the topic to publish open/close commands
	// +optional
	CommandTopic string `json:"commandTopic,omitempty"`

	// StateTopic is the topic to read current valve state
	// +optional
	StateTopic string `json:"stateTopic,omitempty"`

	// CommandTemplate is the template for the command payload
	// +optional
	CommandTemplate string `json:"commandTemplate,omitempty"`

	// ValueTemplate is the template to extract state from payload
	// +optional
	ValueTemplate string `json:"valueTemplate,omitempty"`

	// PositionTopic is the topic to read current position
	// +optional
	PositionTopic string `json:"positionTopic,omitempty"`

	// SetPositionTopic is the topic to publish position commands
	// +optional
	SetPositionTopic string `json:"setPositionTopic,omitempty"`

	// SetPositionTemplate is the template for position command payload
	// +optional
	SetPositionTemplate string `json:"setPositionTemplate,omitempty"`

	// PositionTemplate is the template to extract position from payload
	// +optional
	PositionTemplate string `json:"positionTemplate,omitempty"`

	// PayloadOpen is the payload for open command (default: OPEN)
	// +optional
	PayloadOpen string `json:"payloadOpen,omitempty"`

	// PayloadClose is the payload for close command (default: CLOSE)
	// +optional
	PayloadClose string `json:"payloadClose,omitempty"`

	// PayloadStop is the payload for stop command (default: STOP)
	// +optional
	PayloadStop string `json:"payloadStop,omitempty"`

	// StateOpen is the state value meaning open (default: open)
	// +optional
	StateOpen string `json:"stateOpen,omitempty"`

	// StateClosed is the state value meaning closed (default: closed)
	// +optional
	StateClosed string `json:"stateClosed,omitempty"`

	// StateOpening is the state value meaning opening (default: opening)
	// +optional
	StateOpening string `json:"stateOpening,omitempty"`

	// StateClosing is the state value meaning closing (default: closing)
	// +optional
	StateClosing string `json:"stateClosing,omitempty"`

	// DeviceClass is the valve device class
	// +kubebuilder:validation:Enum=water;gas
	// +optional
	DeviceClass string `json:"deviceClass,omitempty"`

	// ReportsPosition indicates whether the valve reports position
	// +optional
	ReportsPosition *bool `json:"reportsPosition,omitempty"`

	// Optimistic indicates whether to assume state changes immediately
	// +optional
	Optimistic *bool `json:"optimistic,omitempty"`
}

// MQTTValveStatus defines the observed state of MQTTValve.
type MQTTValveStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTValve is the Schema for the mqttvalves API.
// It is a valve entity for controlling water, gas, or irrigation valves.
type MQTTValve struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTValveSpec   `json:"spec,omitempty"`
	Status MQTTValveStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTValveList contains a list of MQTTValve.
type MQTTValveList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTValve `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTValve{}, &MQTTValveList{})
}
