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

// MQTTHumidifierSpec defines the desired state of MQTTHumidifier.
type MQTTHumidifierSpec struct {
	CommonSpec `json:",inline"`

	// CommandTopic is the topic to publish on/off commands
	CommandTopic string `json:"commandTopic"`

	// TargetHumidityCommandTopic is the topic to set target humidity
	TargetHumidityCommandTopic string `json:"targetHumidityCommandTopic"`

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

	// TargetHumidityStateTopic is the topic to read target humidity
	// +optional
	TargetHumidityStateTopic string `json:"targetHumidityStateTopic,omitempty"`

	// TargetHumidityCommandTemplate is the template for target humidity command
	// +optional
	TargetHumidityCommandTemplate string `json:"targetHumidityCommandTemplate,omitempty"`

	// TargetHumidityStateTemplate is the template to extract target humidity
	// +optional
	TargetHumidityStateTemplate string `json:"targetHumidityStateTemplate,omitempty"`

	// CurrentHumidityTopic is the topic to read current humidity
	// +optional
	CurrentHumidityTopic string `json:"currentHumidityTopic,omitempty"`

	// CurrentHumidityTemplate is the template to extract current humidity
	// +optional
	CurrentHumidityTemplate string `json:"currentHumidityTemplate,omitempty"`

	// ModeCommandTopic is the topic to set mode
	// +optional
	ModeCommandTopic string `json:"modeCommandTopic,omitempty"`

	// ModeStateTopic is the topic to read current mode
	// +optional
	ModeStateTopic string `json:"modeStateTopic,omitempty"`

	// ModeCommandTemplate is the template for mode command
	// +optional
	ModeCommandTemplate string `json:"modeCommandTemplate,omitempty"`

	// ModeStateTemplate is the template to extract mode
	// +optional
	ModeStateTemplate string `json:"modeStateTemplate,omitempty"`

	// Modes is the supported modes (e.g. normal, eco, boost, sleep)
	// +optional
	Modes []string `json:"modes,omitempty"`

	// ActionTopic is the topic to read current action
	// +optional
	ActionTopic string `json:"actionTopic,omitempty"`

	// ActionTemplate is the template to extract action
	// +optional
	ActionTemplate string `json:"actionTemplate,omitempty"`

	// MinHumidity is the minimum target humidity (default: 0)
	// +optional
	MinHumidity *float64 `json:"minHumidity,omitempty"`

	// MaxHumidity is the maximum target humidity (default: 100)
	// +optional
	MaxHumidity *float64 `json:"maxHumidity,omitempty"`

	// DeviceClass is the humidifier device class
	// +kubebuilder:validation:Enum=humidifier;dehumidifier
	// +optional
	DeviceClass string `json:"deviceClass,omitempty"`

	// Optimistic indicates whether to assume state changes immediately
	// +optional
	Optimistic *bool `json:"optimistic,omitempty"`
}

// MQTTHumidifierStatus defines the observed state of MQTTHumidifier.
type MQTTHumidifierStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTHumidifier is the Schema for the mqtthumidifiers API.
// It is a humidifier entity with target humidity and mode support.
type MQTTHumidifier struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTHumidifierSpec   `json:"spec,omitempty"`
	Status MQTTHumidifierStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTHumidifierList contains a list of MQTTHumidifier.
type MQTTHumidifierList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTHumidifier `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTHumidifier{}, &MQTTHumidifierList{})
}
