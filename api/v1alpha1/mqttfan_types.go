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

// MQTTFanSpec defines the desired state of MQTTFan.
type MQTTFanSpec struct {
	CommonSpec `json:",inline"`

	// CommandTopic is the topic to publish on/off commands
	CommandTopic string `json:"commandTopic"`

	// StateTopic is the topic to read current on/off state
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

	// PercentageCommandTopic is the topic for speed percentage commands
	// +optional
	PercentageCommandTopic string `json:"percentageCommandTopic,omitempty"`

	// PercentageStateTopic is the topic to read speed percentage
	// +optional
	PercentageStateTopic string `json:"percentageStateTopic,omitempty"`

	// PercentageCommandTemplate is the template for percentage command
	// +optional
	PercentageCommandTemplate string `json:"percentageCommandTemplate,omitempty"`

	// PercentageValueTemplate is the template to extract percentage
	// +optional
	PercentageValueTemplate string `json:"percentageValueTemplate,omitempty"`

	// SpeedRangeMin is the minimum speed value (default: 1)
	// +optional
	SpeedRangeMin *int `json:"speedRangeMin,omitempty"`

	// SpeedRangeMax is the maximum speed value (default: 100)
	// +optional
	SpeedRangeMax *int `json:"speedRangeMax,omitempty"`

	// PresetModeCommandTopic is the topic for preset mode commands
	// +optional
	PresetModeCommandTopic string `json:"presetModeCommandTopic,omitempty"`

	// PresetModeStateTopic is the topic to read preset mode
	// +optional
	PresetModeStateTopic string `json:"presetModeStateTopic,omitempty"`

	// PresetModeCommandTemplate is the template for preset mode command
	// +optional
	PresetModeCommandTemplate string `json:"presetModeCommandTemplate,omitempty"`

	// PresetModeValueTemplate is the template to extract preset mode
	// +optional
	PresetModeValueTemplate string `json:"presetModeValueTemplate,omitempty"`

	// PresetModes is the list of supported preset modes
	// +optional
	PresetModes []string `json:"presetModes,omitempty"`

	// OscillationCommandTopic is the topic for oscillation commands
	// +optional
	OscillationCommandTopic string `json:"oscillationCommandTopic,omitempty"`

	// OscillationStateTopic is the topic to read oscillation state
	// +optional
	OscillationStateTopic string `json:"oscillationStateTopic,omitempty"`

	// OscillationCommandTemplate is the template for oscillation command
	// +optional
	OscillationCommandTemplate string `json:"oscillationCommandTemplate,omitempty"`

	// OscillationValueTemplate is the template to extract oscillation state
	// +optional
	OscillationValueTemplate string `json:"oscillationValueTemplate,omitempty"`

	// PayloadOscillationOn is the payload for oscillation on (default: oscillate_on)
	// +optional
	PayloadOscillationOn string `json:"payloadOscillationOn,omitempty"`

	// PayloadOscillationOff is the payload for oscillation off (default: oscillate_off)
	// +optional
	PayloadOscillationOff string `json:"payloadOscillationOff,omitempty"`

	// DirectionCommandTopic is the topic for direction commands
	// +optional
	DirectionCommandTopic string `json:"directionCommandTopic,omitempty"`

	// DirectionStateTopic is the topic to read direction state
	// +optional
	DirectionStateTopic string `json:"directionStateTopic,omitempty"`

	// DirectionValueTemplate is the template to extract direction
	// +optional
	DirectionValueTemplate string `json:"directionValueTemplate,omitempty"`

	// Optimistic indicates whether to assume state changes immediately
	// +optional
	Optimistic *bool `json:"optimistic,omitempty"`
}

// MQTTFanStatus defines the observed state of MQTTFan.
type MQTTFanStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTFan is the Schema for the mqttfans API.
// It is a fan entity with speed, direction, oscillation, and preset mode support.
type MQTTFan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTFanSpec   `json:"spec,omitempty"`
	Status MQTTFanStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTFanList contains a list of MQTTFan.
type MQTTFanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTFan `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTFan{}, &MQTTFanList{})
}
