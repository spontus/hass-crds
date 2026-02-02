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

// MQTTTagSpec defines the desired state of MQTTTag.
type MQTTTagSpec struct {
	CommonSpec `json:",inline"`

	// Topic is the topic to subscribe to for tag scans
	Topic string `json:"topic"`

	// ValueTemplate is the template to extract tag ID from payload
	// +optional
	ValueTemplate string `json:"valueTemplate,omitempty"`
}

// MQTTTagStatus defines the observed state of MQTTTag.
type MQTTTagStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTTag is the Schema for the mqtttags API.
// It is a tag scanner entity for NFC, RFID, or QR code scanning.
type MQTTTag struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTTagSpec   `json:"spec,omitempty"`
	Status MQTTTagStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTTagList contains a list of MQTTTag.
type MQTTTagList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTTag `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTTag{}, &MQTTTagList{})
}
