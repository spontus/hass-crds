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

// MQTTSceneSpec defines the desired state of MQTTScene.
type MQTTSceneSpec struct {
	CommonSpec `json:",inline"`

	// CommandTopic is the topic to publish when scene is activated
	CommandTopic string `json:"commandTopic"`

	// PayloadOn is the payload sent when scene is activated (default: ON)
	// +optional
	PayloadOn string `json:"payloadOn,omitempty"`
}

// MQTTSceneStatus defines the observed state of MQTTScene.
type MQTTSceneStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTScene is the Schema for the mqttscenes API.
// It is a scene entity that can be activated via MQTT.
type MQTTScene struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTSceneSpec   `json:"spec,omitempty"`
	Status MQTTSceneStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTSceneList contains a list of MQTTScene.
type MQTTSceneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTScene `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTScene{}, &MQTTSceneList{})
}
