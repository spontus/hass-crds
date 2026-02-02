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

// MQTTVacuumSpec defines the desired state of MQTTVacuum.
type MQTTVacuumSpec struct {
	CommonSpec `json:",inline"`

	// CommandTopic is the topic for basic commands (start, stop, return_to_base, etc.)
	// +optional
	CommandTopic string `json:"commandTopic,omitempty"`

	// StateTopic is the topic to read vacuum state
	// +optional
	StateTopic string `json:"stateTopic,omitempty"`

	// SendCommandTopic is the topic for custom commands
	// +optional
	SendCommandTopic string `json:"sendCommandTopic,omitempty"`

	// SetFanSpeedTopic is the topic for fan speed commands
	// +optional
	SetFanSpeedTopic string `json:"setFanSpeedTopic,omitempty"`

	// FanSpeedList is the list of supported fan speeds
	// +optional
	FanSpeedList []string `json:"fanSpeedList,omitempty"`

	// PayloadStart is the payload for start command (default: start)
	// +optional
	PayloadStart string `json:"payloadStart,omitempty"`

	// PayloadStop is the payload for stop command (default: stop)
	// +optional
	PayloadStop string `json:"payloadStop,omitempty"`

	// PayloadPause is the payload for pause command (default: pause)
	// +optional
	PayloadPause string `json:"payloadPause,omitempty"`

	// PayloadReturnToBase is the payload for return to base command (default: return_to_base)
	// +optional
	PayloadReturnToBase string `json:"payloadReturnToBase,omitempty"`

	// PayloadCleanSpot is the payload for clean spot command (default: clean_spot)
	// +optional
	PayloadCleanSpot string `json:"payloadCleanSpot,omitempty"`

	// PayloadLocate is the payload for locate command (default: locate)
	// +optional
	PayloadLocate string `json:"payloadLocate,omitempty"`

	// SupportedFeatures is the list of supported features (e.g. start, stop, pause, return_home, fan_speed, send_command, locate, clean_spot)
	// +optional
	SupportedFeatures []string `json:"supportedFeatures,omitempty"`

	// Schema is the vacuum schema
	// +kubebuilder:validation:Enum=legacy;state
	// +optional
	Schema string `json:"schema,omitempty"`
}

// MQTTVacuumStatus defines the observed state of MQTTVacuum.
type MQTTVacuumStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTVacuum is the Schema for the mqttvacuums API.
// It is a robot vacuum entity with start, stop, pause, return to base, and cleaning features.
type MQTTVacuum struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTVacuumSpec   `json:"spec,omitempty"`
	Status MQTTVacuumStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTVacuumList contains a list of MQTTVacuum.
type MQTTVacuumList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTVacuum `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTVacuum{}, &MQTTVacuumList{})
}
