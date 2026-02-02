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

// MQTTLightSpec defines the desired state of MQTTLight.
type MQTTLightSpec struct {
	CommonSpec `json:",inline"`

	// CommandTopic is the topic to publish on/off commands
	CommandTopic string `json:"commandTopic"`

	// Schema is the light schema mode
	// +kubebuilder:validation:Enum=default;json;template
	// +optional
	Schema string `json:"schema,omitempty"`

	// StateTopic is the topic to read current state
	// +optional
	StateTopic string `json:"stateTopic,omitempty"`

	// PayloadOn is the payload for on (default: ON)
	// +optional
	PayloadOn string `json:"payloadOn,omitempty"`

	// PayloadOff is the payload for off (default: OFF)
	// +optional
	PayloadOff string `json:"payloadOff,omitempty"`

	// BrightnessCommandTopic is the topic for brightness commands
	// +optional
	BrightnessCommandTopic string `json:"brightnessCommandTopic,omitempty"`

	// BrightnessStateTopic is the topic for brightness state
	// +optional
	BrightnessStateTopic string `json:"brightnessStateTopic,omitempty"`

	// BrightnessScale is the max brightness value (default: 255)
	// +optional
	BrightnessScale *int `json:"brightnessScale,omitempty"`

	// BrightnessValueTemplate is the template to extract brightness
	// +optional
	BrightnessValueTemplate string `json:"brightnessValueTemplate,omitempty"`

	// ColorTempCommandTopic is the topic for color temperature commands
	// +optional
	ColorTempCommandTopic string `json:"colorTempCommandTopic,omitempty"`

	// ColorTempStateTopic is the topic for color temperature state
	// +optional
	ColorTempStateTopic string `json:"colorTempStateTopic,omitempty"`

	// ColorTempValueTemplate is the template to extract color temp
	// +optional
	ColorTempValueTemplate string `json:"colorTempValueTemplate,omitempty"`

	// RgbCommandTopic is the topic for RGB color commands
	// +optional
	RgbCommandTopic string `json:"rgbCommandTopic,omitempty"`

	// RgbStateTopic is the topic for RGB color state
	// +optional
	RgbStateTopic string `json:"rgbStateTopic,omitempty"`

	// RgbCommandTemplate is the template for RGB command payload
	// +optional
	RgbCommandTemplate string `json:"rgbCommandTemplate,omitempty"`

	// RgbValueTemplate is the template to extract RGB state
	// +optional
	RgbValueTemplate string `json:"rgbValueTemplate,omitempty"`

	// EffectCommandTopic is the topic for effect commands
	// +optional
	EffectCommandTopic string `json:"effectCommandTopic,omitempty"`

	// EffectStateTopic is the topic for effect state
	// +optional
	EffectStateTopic string `json:"effectStateTopic,omitempty"`

	// EffectList is the list of supported effects
	// +optional
	EffectList []string `json:"effectList,omitempty"`

	// EffectValueTemplate is the template to extract effect
	// +optional
	EffectValueTemplate string `json:"effectValueTemplate,omitempty"`

	// MinMireds is the minimum color temp in mireds
	// +optional
	MinMireds *int `json:"minMireds,omitempty"`

	// MaxMireds is the maximum color temp in mireds
	// +optional
	MaxMireds *int `json:"maxMireds,omitempty"`

	// Optimistic indicates whether to assume state changes immediately
	// +optional
	Optimistic *bool `json:"optimistic,omitempty"`

	// OnCommandType is the on command type
	// +kubebuilder:validation:Enum=last;first;brightness
	// +optional
	OnCommandType string `json:"onCommandType,omitempty"`

	// Brightness enables brightness support (JSON schema)
	// +optional
	Brightness *bool `json:"brightness,omitempty"`

	// ColorTemp enables color temperature (JSON schema)
	// +optional
	ColorTemp *bool `json:"colorTemp,omitempty"`

	// Effect enables effects (JSON schema)
	// +optional
	Effect *bool `json:"effect,omitempty"`

	// SupportedColorModes lists supported color modes (e.g. rgb, xy, hs, color_temp)
	// +optional
	SupportedColorModes []string `json:"supportedColorModes,omitempty"`

	// CommandOnTemplate is the template for on command (template schema)
	// +optional
	CommandOnTemplate string `json:"commandOnTemplate,omitempty"`

	// CommandOffTemplate is the template for off command (template schema)
	// +optional
	CommandOffTemplate string `json:"commandOffTemplate,omitempty"`

	// StateTemplate is the template to extract state (template schema)
	// +optional
	StateTemplate string `json:"stateTemplate,omitempty"`

	// BrightnessTemplate is the template to extract brightness (template schema)
	// +optional
	BrightnessTemplate string `json:"brightnessTemplate,omitempty"`

	// ColorTempTemplate is the template to extract color temp (template schema)
	// +optional
	ColorTempTemplate string `json:"colorTempTemplate,omitempty"`

	// RedTemplate is the template to extract red value (template schema)
	// +optional
	RedTemplate string `json:"redTemplate,omitempty"`

	// GreenTemplate is the template to extract green value (template schema)
	// +optional
	GreenTemplate string `json:"greenTemplate,omitempty"`

	// BlueTemplate is the template to extract blue value (template schema)
	// +optional
	BlueTemplate string `json:"blueTemplate,omitempty"`
}

// MQTTLightStatus defines the observed state of MQTTLight.
type MQTTLightStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTLight is the Schema for the mqttlights API.
// It is a light entity with optional brightness, color temperature, and RGB color support.
type MQTTLight struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTLightSpec   `json:"spec,omitempty"`
	Status MQTTLightStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTLightList contains a list of MQTTLight.
type MQTTLightList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTLight `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTLight{}, &MQTTLightList{})
}
