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

// MQTTUpdateSpec defines the desired state of MQTTUpdate.
type MQTTUpdateSpec struct {
	CommonSpec `json:",inline"`

	// StateTopic is the topic with JSON payload containing update info
	StateTopic string `json:"stateTopic"`

	// ValueTemplate is the template to extract state from payload
	// +optional
	ValueTemplate string `json:"valueTemplate,omitempty"`

	// CommandTopic is the topic to trigger update installation
	// +optional
	CommandTopic string `json:"commandTopic,omitempty"`

	// PayloadInstall is the payload to trigger installation (default: INSTALL)
	// +optional
	PayloadInstall string `json:"payloadInstall,omitempty"`

	// LatestVersionTopic is the topic to read latest available version
	// +optional
	LatestVersionTopic string `json:"latestVersionTopic,omitempty"`

	// LatestVersionTemplate is the template to extract latest version
	// +optional
	LatestVersionTemplate string `json:"latestVersionTemplate,omitempty"`

	// DeviceClass is the update device class
	// +kubebuilder:validation:Enum=firmware
	// +optional
	DeviceClass string `json:"deviceClass,omitempty"`

	// EntityPicture is the URL to an image for the update entity
	// +optional
	EntityPicture string `json:"entityPicture,omitempty"`

	// ReleaseUrl is the URL to release notes
	// +optional
	ReleaseUrl string `json:"releaseUrl,omitempty"`

	// ReleaseSummary is the summary of the release
	// +optional
	ReleaseSummary string `json:"releaseSummary,omitempty"`

	// Title is the title of the software/firmware
	// +optional
	Title string `json:"title,omitempty"`
}

// MQTTUpdateStatus defines the observed state of MQTTUpdate.
type MQTTUpdateStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTUpdate is the Schema for the mqttupdates API.
// It is a firmware/software update entity that tracks available updates via MQTT.
type MQTTUpdate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTUpdateSpec   `json:"spec,omitempty"`
	Status MQTTUpdateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTUpdateList contains a list of MQTTUpdate.
type MQTTUpdateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTUpdate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTUpdate{}, &MQTTUpdateList{})
}
