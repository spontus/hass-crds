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

// MQTTAlarmControlPanelSpec defines the desired state of MQTTAlarmControlPanel.
type MQTTAlarmControlPanelSpec struct {
	CommonSpec `json:",inline"`

	// CommandTopic is the topic to publish arm/disarm commands
	CommandTopic string `json:"commandTopic"`

	// StateTopic is the topic to read alarm state
	StateTopic string `json:"stateTopic"`

	// CommandTemplate is the template for the command payload
	// +optional
	CommandTemplate string `json:"commandTemplate,omitempty"`

	// ValueTemplate is the template to extract state from payload
	// +optional
	ValueTemplate string `json:"valueTemplate,omitempty"`

	// PayloadArmHome is the payload for arm home (default: ARM_HOME)
	// +optional
	PayloadArmHome string `json:"payloadArmHome,omitempty"`

	// PayloadArmAway is the payload for arm away (default: ARM_AWAY)
	// +optional
	PayloadArmAway string `json:"payloadArmAway,omitempty"`

	// PayloadArmNight is the payload for arm night (default: ARM_NIGHT)
	// +optional
	PayloadArmNight string `json:"payloadArmNight,omitempty"`

	// PayloadArmVacation is the payload for arm vacation (default: ARM_VACATION)
	// +optional
	PayloadArmVacation string `json:"payloadArmVacation,omitempty"`

	// PayloadArmCustomBypass is the payload for arm custom bypass (default: ARM_CUSTOM_BYPASS)
	// +optional
	PayloadArmCustomBypass string `json:"payloadArmCustomBypass,omitempty"`

	// PayloadDisarm is the payload for disarm (default: DISARM)
	// +optional
	PayloadDisarm string `json:"payloadDisarm,omitempty"`

	// PayloadTrigger is the payload for trigger
	// +optional
	PayloadTrigger string `json:"payloadTrigger,omitempty"`

	// CodeArmRequired indicates whether code is required to arm (default: true)
	// +optional
	CodeArmRequired *bool `json:"codeArmRequired,omitempty"`

	// CodeDisarmRequired indicates whether code is required to disarm (default: true)
	// +optional
	CodeDisarmRequired *bool `json:"codeDisarmRequired,omitempty"`

	// CodeTriggerRequired indicates whether code is required to trigger (default: true)
	// +optional
	CodeTriggerRequired *bool `json:"codeTriggerRequired,omitempty"`

	// CodeFormat is the code format
	// +kubebuilder:validation:Enum=number;text
	// +optional
	CodeFormat string `json:"codeFormat,omitempty"`

	// SupportedFeatures is the list of supported features (e.g. arm_home, arm_away, arm_night, trigger)
	// +optional
	SupportedFeatures []string `json:"supportedFeatures,omitempty"`
}

// MQTTAlarmControlPanelStatus defines the observed state of MQTTAlarmControlPanel.
type MQTTAlarmControlPanelStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTAlarmControlPanel is the Schema for the mqttalarmcontrolpanels API.
// It is an alarm control panel entity with arm/disarm modes and optional code support.
type MQTTAlarmControlPanel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTAlarmControlPanelSpec   `json:"spec,omitempty"`
	Status MQTTAlarmControlPanelStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTAlarmControlPanelList contains a list of MQTTAlarmControlPanel.
type MQTTAlarmControlPanelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTAlarmControlPanel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTAlarmControlPanel{}, &MQTTAlarmControlPanelList{})
}
