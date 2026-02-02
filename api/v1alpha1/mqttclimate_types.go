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

// MQTTClimateSpec defines the desired state of MQTTClimate.
type MQTTClimateSpec struct {
	CommonSpec `json:",inline"`

	// TemperatureCommandTopic is the topic to set target temperature
	// +optional
	TemperatureCommandTopic string `json:"temperatureCommandTopic,omitempty"`

	// TemperatureStateTopic is the topic to read target temperature
	// +optional
	TemperatureStateTopic string `json:"temperatureStateTopic,omitempty"`

	// TemperatureCommandTemplate is the template for temperature command
	// +optional
	TemperatureCommandTemplate string `json:"temperatureCommandTemplate,omitempty"`

	// TemperatureStateTemplate is the template to extract target temp
	// +optional
	TemperatureStateTemplate string `json:"temperatureStateTemplate,omitempty"`

	// CurrentTemperatureTopic is the topic to read current temperature
	// +optional
	CurrentTemperatureTopic string `json:"currentTemperatureTopic,omitempty"`

	// CurrentTemperatureTemplate is the template to extract current temp
	// +optional
	CurrentTemperatureTemplate string `json:"currentTemperatureTemplate,omitempty"`

	// ModeCommandTopic is the topic to set HVAC mode
	// +optional
	ModeCommandTopic string `json:"modeCommandTopic,omitempty"`

	// ModeStateTopic is the topic to read HVAC mode
	// +optional
	ModeStateTopic string `json:"modeStateTopic,omitempty"`

	// ModeCommandTemplate is the template for mode command
	// +optional
	ModeCommandTemplate string `json:"modeCommandTemplate,omitempty"`

	// ModeStateTemplate is the template to extract mode
	// +optional
	ModeStateTemplate string `json:"modeStateTemplate,omitempty"`

	// Modes is the supported HVAC modes
	// +optional
	Modes []string `json:"modes,omitempty"`

	// FanModeCommandTopic is the topic to set fan mode
	// +optional
	FanModeCommandTopic string `json:"fanModeCommandTopic,omitempty"`

	// FanModeStateTopic is the topic to read fan mode
	// +optional
	FanModeStateTopic string `json:"fanModeStateTopic,omitempty"`

	// FanModeCommandTemplate is the template for fan mode command
	// +optional
	FanModeCommandTemplate string `json:"fanModeCommandTemplate,omitempty"`

	// FanModeStateTemplate is the template to extract fan mode
	// +optional
	FanModeStateTemplate string `json:"fanModeStateTemplate,omitempty"`

	// FanModes is the supported fan modes
	// +optional
	FanModes []string `json:"fanModes,omitempty"`

	// SwingModeCommandTopic is the topic to set swing mode
	// +optional
	SwingModeCommandTopic string `json:"swingModeCommandTopic,omitempty"`

	// SwingModeStateTopic is the topic to read swing mode
	// +optional
	SwingModeStateTopic string `json:"swingModeStateTopic,omitempty"`

	// SwingModes is the supported swing modes
	// +optional
	SwingModes []string `json:"swingModes,omitempty"`

	// PresetModeCommandTopic is the topic to set preset mode
	// +optional
	PresetModeCommandTopic string `json:"presetModeCommandTopic,omitempty"`

	// PresetModeStateTopic is the topic to read preset mode
	// +optional
	PresetModeStateTopic string `json:"presetModeStateTopic,omitempty"`

	// PresetModes is the supported preset modes (e.g. away, eco, boost)
	// +optional
	PresetModes []string `json:"presetModes,omitempty"`

	// ActionTopic is the topic to read current HVAC action
	// +optional
	ActionTopic string `json:"actionTopic,omitempty"`

	// ActionTemplate is the template to extract action
	// +optional
	ActionTemplate string `json:"actionTemplate,omitempty"`

	// TempStep is the step size for temperature adjustments (default: 1)
	// +optional
	TempStep *float64 `json:"tempStep,omitempty"`

	// MinTemp is the minimum setpoint temperature
	// +optional
	MinTemp *float64 `json:"minTemp,omitempty"`

	// MaxTemp is the maximum setpoint temperature
	// +optional
	MaxTemp *float64 `json:"maxTemp,omitempty"`

	// TemperatureUnit is the temperature unit
	// +kubebuilder:validation:Enum=C;F
	// +optional
	TemperatureUnit string `json:"temperatureUnit,omitempty"`

	// Precision is the temperature precision (default: 0.1)
	// +optional
	Precision *float64 `json:"precision,omitempty"`

	// Optimistic indicates whether to assume state changes immediately
	// +optional
	Optimistic *bool `json:"optimistic,omitempty"`
}

// MQTTClimateStatus defines the observed state of MQTTClimate.
type MQTTClimateStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTClimate is the Schema for the mqttclimates API.
// It is a thermostat/HVAC entity with temperature control, modes, and fan speed.
type MQTTClimate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTClimateSpec   `json:"spec,omitempty"`
	Status MQTTClimateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTClimateList contains a list of MQTTClimate.
type MQTTClimateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTClimate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTClimate{}, &MQTTClimateList{})
}
