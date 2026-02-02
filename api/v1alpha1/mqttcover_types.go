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

// MQTTCoverSpec defines the desired state of MQTTCover.
type MQTTCoverSpec struct {
	CommonSpec `json:",inline"`

	// CommandTopic is the topic for open/close/stop commands
	// +optional
	CommandTopic string `json:"commandTopic,omitempty"`

	// StateTopic is the topic to read cover state
	// +optional
	StateTopic string `json:"stateTopic,omitempty"`

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

	// TiltCommandTopic is the topic for tilt commands
	// +optional
	TiltCommandTopic string `json:"tiltCommandTopic,omitempty"`

	// TiltStatusTopic is the topic to read tilt position
	// +optional
	TiltStatusTopic string `json:"tiltStatusTopic,omitempty"`

	// TiltStatusTemplate is the template to extract tilt from payload
	// +optional
	TiltStatusTemplate string `json:"tiltStatusTemplate,omitempty"`

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

	// StateStopped is the state value meaning stopped (default: stopped)
	// +optional
	StateStopped string `json:"stateStopped,omitempty"`

	// PositionOpen is the position value for fully open (default: 100)
	// +optional
	PositionOpen *int `json:"positionOpen,omitempty"`

	// PositionClosed is the position value for fully closed (default: 0)
	// +optional
	PositionClosed *int `json:"positionClosed,omitempty"`

	// TiltMin is the minimum tilt value (default: 0)
	// +optional
	TiltMin *int `json:"tiltMin,omitempty"`

	// TiltMax is the maximum tilt value (default: 100)
	// +optional
	TiltMax *int `json:"tiltMax,omitempty"`

	// DeviceClass is the cover device class
	// +kubebuilder:validation:Enum=awning;blind;curtain;damper;door;garage;gate;shade;shutter;window
	// +optional
	DeviceClass string `json:"deviceClass,omitempty"`

	// Optimistic indicates whether to assume state changes immediately
	// +optional
	Optimistic *bool `json:"optimistic,omitempty"`
}

// MQTTCoverStatus defines the observed state of MQTTCover.
type MQTTCoverStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTCover is the Schema for the mqttcovers API.
// It is a cover entity for garage doors, blinds, shutters, and similar devices.
type MQTTCover struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTCoverSpec   `json:"spec,omitempty"`
	Status MQTTCoverStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTCoverList contains a list of MQTTCover.
type MQTTCoverList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTCover `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTCover{}, &MQTTCoverList{})
}
