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

// MQTTButtonSpec defines the desired state of MQTTButton.
type MQTTButtonSpec struct {
	CommonSpec `json:",inline"`

	// CommandTopic is the topic to publish when button is pressed
	CommandTopic string `json:"commandTopic"`

	// CommandTemplate is the template for the command payload
	// +optional
	CommandTemplate string `json:"commandTemplate,omitempty"`

	// PayloadPress is the payload sent when button is pressed (default: PRESS)
	// +optional
	PayloadPress string `json:"payloadPress,omitempty"`

	// DeviceClass is the button device class
	// +kubebuilder:validation:Enum=identify;restart;update
	// +optional
	DeviceClass string `json:"deviceClass,omitempty"`
}

// MQTTButtonStatus defines the observed state of MQTTButton.
type MQTTButtonStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTButton is the Schema for the mqttbuttons API.
// It is a stateless button entity - publishes to command topic when pressed.
type MQTTButton struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTButtonSpec   `json:"spec,omitempty"`
	Status MQTTButtonStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTButtonList contains a list of MQTTButton.
type MQTTButtonList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTButton `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTButton{}, &MQTTButtonList{})
}
