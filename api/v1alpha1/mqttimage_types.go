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

// MQTTImageSpec defines the desired state of MQTTImage.
type MQTTImageSpec struct {
	CommonSpec `json:",inline"`

	// ImageTopic is the topic to receive raw image data
	// +optional
	ImageTopic string `json:"imageTopic,omitempty"`

	// ImageEncoding is the image encoding (b64 for base64-encoded images)
	// +optional
	ImageEncoding string `json:"imageEncoding,omitempty"`

	// UrlTopic is the topic to receive image URL
	// +optional
	UrlTopic string `json:"urlTopic,omitempty"`

	// UrlTemplate is the template to extract URL from payload
	// +optional
	UrlTemplate string `json:"urlTemplate,omitempty"`

	// ContentType is the image MIME type (default: image/png)
	// +optional
	ContentType string `json:"contentType,omitempty"`
}

// MQTTImageStatus defines the observed state of MQTTImage.
type MQTTImageStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTImage is the Schema for the mqttimages API.
// It is an image entity that displays a static image from an MQTT topic or URL.
type MQTTImage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTImageSpec   `json:"spec,omitempty"`
	Status MQTTImageStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTImageList contains a list of MQTTImage.
type MQTTImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTImage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTImage{}, &MQTTImageList{})
}
