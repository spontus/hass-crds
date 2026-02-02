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

// MQTTLawnMowerSpec defines the desired state of MQTTLawnMower.
type MQTTLawnMowerSpec struct {
	CommonSpec `json:",inline"`

	// ActivityStateTopic is the topic to read mower activity state
	// +optional
	ActivityStateTopic string `json:"activityStateTopic,omitempty"`

	// ActivityValueTemplate is the template to extract activity from payload
	// +optional
	ActivityValueTemplate string `json:"activityValueTemplate,omitempty"`

	// DockCommandTopic is the topic to publish dock command
	// +optional
	DockCommandTopic string `json:"dockCommandTopic,omitempty"`

	// DockCommandTemplate is the template for dock command payload
	// +optional
	DockCommandTemplate string `json:"dockCommandTemplate,omitempty"`

	// PauseCommandTopic is the topic to publish pause command
	// +optional
	PauseCommandTopic string `json:"pauseCommandTopic,omitempty"`

	// PauseCommandTemplate is the template for pause command payload
	// +optional
	PauseCommandTemplate string `json:"pauseCommandTemplate,omitempty"`

	// StartMowingCommandTopic is the topic to publish start mowing command
	// +optional
	StartMowingCommandTopic string `json:"startMowingCommandTopic,omitempty"`

	// StartMowingCommandTemplate is the template for start mowing command payload
	// +optional
	StartMowingCommandTemplate string `json:"startMowingCommandTemplate,omitempty"`

	// Optimistic indicates whether to assume state changes immediately
	// +optional
	Optimistic *bool `json:"optimistic,omitempty"`
}

// MQTTLawnMowerStatus defines the observed state of MQTTLawnMower.
type MQTTLawnMowerStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTLawnMower is the Schema for the mqttlawnmowers API.
// It is a robot lawn mower entity with start mowing, pause, and dock commands.
type MQTTLawnMower struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTLawnMowerSpec   `json:"spec,omitempty"`
	Status MQTTLawnMowerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTLawnMowerList contains a list of MQTTLawnMower.
type MQTTLawnMowerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTLawnMower `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTLawnMower{}, &MQTTLawnMowerList{})
}
