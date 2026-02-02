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

// MQTTSwitchSpec defines the desired state of MQTTSwitch.
type MQTTSwitchSpec struct {
	CommonSpec `json:",inline"`

	// CommandTopic is the topic to publish on/off commands
	CommandTopic string `json:"commandTopic"`

	// StateTopic is the topic to read current state
	// +optional
	StateTopic string `json:"stateTopic,omitempty"`

	// CommandTemplate is the template for the command payload
	// +optional
	CommandTemplate string `json:"commandTemplate,omitempty"`

	// ValueTemplate is the template to extract state from payload
	// +optional
	ValueTemplate string `json:"valueTemplate,omitempty"`

	// PayloadOn is the payload representing on (default: ON)
	// +optional
	PayloadOn string `json:"payloadOn,omitempty"`

	// PayloadOff is the payload representing off (default: OFF)
	// +optional
	PayloadOff string `json:"payloadOff,omitempty"`

	// StateOn is the state value that means on (if different from payloadOn)
	// +optional
	StateOn string `json:"stateOn,omitempty"`

	// StateOff is the state value that means off (if different from payloadOff)
	// +optional
	StateOff string `json:"stateOff,omitempty"`

	// DeviceClass is the switch device class
	// +kubebuilder:validation:Enum=outlet;switch
	// +optional
	DeviceClass string `json:"deviceClass,omitempty"`

	// Optimistic indicates whether to assume state changes immediately
	// +optional
	Optimistic *bool `json:"optimistic,omitempty"`
}

// MQTTSwitchStatus defines the observed state of MQTTSwitch.
type MQTTSwitchStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTSwitch is the Schema for the mqttswitches API.
// It is an on/off toggle entity with state feedback.
type MQTTSwitch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTSwitchSpec   `json:"spec,omitempty"`
	Status MQTTSwitchStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTSwitchList contains a list of MQTTSwitch.
type MQTTSwitchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTSwitch `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTSwitch{}, &MQTTSwitchList{})
}
