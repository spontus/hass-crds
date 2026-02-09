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

// MQTTCameraSpec defines the desired state of MQTTCamera.
type MQTTCameraSpec struct {
	CommonSpec `json:",inline"`

	// Topic is the MQTT topic to subscribe to for image data
	Topic string `json:"topic"`

	// ImageEncoding is the image encoding (b64 for base64-encoded images)
	// +optional
	ImageEncoding string `json:"imageEncoding,omitempty"`

	// StateClass is the state class for statistics (unusual for cameras but supported)
	// +kubebuilder:validation:Enum=measurement;total;total_increasing
	// +optional
	StateClass string `json:"stateClass,omitempty"`

	// ExpireAfter is the seconds after which the image expires
	// +optional
	ExpireAfter *int `json:"expireAfter,omitempty"`
}

// MQTTCameraStatus defines the observed state of MQTTCamera.
type MQTTCameraStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTCamera is the Schema for the mqttcameras API.
// It is a camera entity that receives images via MQTT.
type MQTTCamera struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTCameraSpec   `json:"spec,omitempty"`
	Status MQTTCameraStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTCameraList contains a list of MQTTCamera.
type MQTTCameraList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTCamera `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTCamera{}, &MQTTCameraList{})
}
