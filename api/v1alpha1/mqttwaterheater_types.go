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

// MQTTWaterHeaterSpec defines the desired state of MQTTWaterHeater.
type MQTTWaterHeaterSpec struct {
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

	// ModeCommandTopic is the topic to set operation mode
	// +optional
	ModeCommandTopic string `json:"modeCommandTopic,omitempty"`

	// ModeStateTopic is the topic to read operation mode
	// +optional
	ModeStateTopic string `json:"modeStateTopic,omitempty"`

	// ModeCommandTemplate is the template for mode command
	// +optional
	ModeCommandTemplate string `json:"modeCommandTemplate,omitempty"`

	// ModeStateTemplate is the template to extract mode
	// +optional
	ModeStateTemplate string `json:"modeStateTemplate,omitempty"`

	// Modes is the supported modes (e.g. off, eco, electric, gas, heat_pump, high_demand, performance)
	// +optional
	Modes []string `json:"modes,omitempty"`

	// PowerCommandTopic is the topic to publish on/off commands
	// +optional
	PowerCommandTopic string `json:"powerCommandTopic,omitempty"`

	// PayloadOn is the payload for on (default: ON)
	// +optional
	PayloadOn string `json:"payloadOn,omitempty"`

	// PayloadOff is the payload for off (default: OFF)
	// +optional
	PayloadOff string `json:"payloadOff,omitempty"`

	// MinTemp is the minimum target temperature (default: 110)
	// +optional
	MinTemp *float64 `json:"minTemp,omitempty"`

	// MaxTemp is the maximum target temperature (default: 140)
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

// MQTTWaterHeaterStatus defines the observed state of MQTTWaterHeater.
type MQTTWaterHeaterStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTWaterHeater is the Schema for the mqttwaterheaters API.
// It is a water heater entity with temperature control and operation modes.
type MQTTWaterHeater struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTWaterHeaterSpec   `json:"spec,omitempty"`
	Status MQTTWaterHeaterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTWaterHeaterList contains a list of MQTTWaterHeater.
type MQTTWaterHeaterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTWaterHeater `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTWaterHeater{}, &MQTTWaterHeaterList{})
}
