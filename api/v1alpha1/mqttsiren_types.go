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

// MQTTSirenSpec defines the desired state of MQTTSiren.
type MQTTSirenSpec struct {
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

	// PayloadOn is the payload for on (default: ON)
	// +optional
	PayloadOn string `json:"payloadOn,omitempty"`

	// PayloadOff is the payload for off (default: OFF)
	// +optional
	PayloadOff string `json:"payloadOff,omitempty"`

	// StateOn is the state value meaning on (default: ON)
	// +optional
	StateOn string `json:"stateOn,omitempty"`

	// StateOff is the state value meaning off (default: OFF)
	// +optional
	StateOff string `json:"stateOff,omitempty"`

	// AvailableTones is the list of supported tones
	// +optional
	AvailableTones []string `json:"availableTones,omitempty"`

	// SupportTurnOn indicates whether the siren supports turn on (default: true)
	// +optional
	SupportTurnOn *bool `json:"supportTurnOn,omitempty"`

	// SupportTurnOff indicates whether the siren supports turn off (default: true)
	// +optional
	SupportTurnOff *bool `json:"supportTurnOff,omitempty"`

	// SupportDuration indicates whether duration is supported (default: true)
	// +optional
	SupportDuration *bool `json:"supportDuration,omitempty"`

	// SupportVolumeSet indicates whether volume level is supported (default: true)
	// +optional
	SupportVolumeSet *bool `json:"supportVolumeSet,omitempty"`

	// Optimistic indicates whether to assume state changes immediately
	// +optional
	Optimistic *bool `json:"optimistic,omitempty"`
}

// MQTTSirenStatus defines the observed state of MQTTSiren.
type MQTTSirenStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTSiren is the Schema for the mqttsirens API.
// It is a siren entity with optional tone, volume, and duration support.
type MQTTSiren struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTSirenSpec   `json:"spec,omitempty"`
	Status MQTTSirenStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTSirenList contains a list of MQTTSiren.
type MQTTSirenList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTSiren `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTSiren{}, &MQTTSirenList{})
}
