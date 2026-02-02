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

// MQTTDeviceSpec defines the desired state of MQTTDevice.
// MQTTDevice is a utility resource for shared device definitions.
type MQTTDeviceSpec struct {
	// Name is the device display name
	// +optional
	Name string `json:"name,omitempty"`

	// Identifiers is a list of identifiers (at least one of identifiers or connections is needed)
	// +optional
	Identifiers []string `json:"identifiers,omitempty"`

	// Connections is a list of [type, value] pairs (e.g. [["mac", "aa:bb:cc:dd:ee:ff"]])
	// +optional
	Connections [][]string `json:"connections,omitempty"`

	// Manufacturer is the device manufacturer
	// +optional
	Manufacturer string `json:"manufacturer,omitempty"`

	// Model is the device model
	// +optional
	Model string `json:"model,omitempty"`

	// ModelId is the device model identifier
	// +optional
	ModelId string `json:"modelId,omitempty"`

	// SerialNumber is the device serial number
	// +optional
	SerialNumber string `json:"serialNumber,omitempty"`

	// HwVersion is the hardware version
	// +optional
	HwVersion string `json:"hwVersion,omitempty"`

	// SwVersion is the software version
	// +optional
	SwVersion string `json:"swVersion,omitempty"`

	// SuggestedArea is the suggested area in Home Assistant (e.g. Living Room)
	// +optional
	SuggestedArea string `json:"suggestedArea,omitempty"`

	// ConfigurationUrl is the URL for device configuration
	// +optional
	ConfigurationUrl string `json:"configurationUrl,omitempty"`

	// ViaDevice is the identifier of device that routes messages
	// +optional
	ViaDevice string `json:"viaDevice,omitempty"`
}

// MQTTDeviceStatus defines the observed state of MQTTDevice.
type MQTTDeviceStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTDevice is the Schema for the mqttdevices API.
// It is a shared device definition for multiple MQTT entities.
type MQTTDevice struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTDeviceSpec   `json:"spec,omitempty"`
	Status MQTTDeviceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTDeviceList contains a list of MQTTDevice.
type MQTTDeviceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTDevice `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTDevice{}, &MQTTDeviceList{})
}

// ToDeviceBlock converts MQTTDeviceSpec to DeviceBlock for use in discovery payloads.
func (s *MQTTDeviceSpec) ToDeviceBlock() DeviceBlock {
	return DeviceBlock{
		Name:             s.Name,
		Identifiers:      s.Identifiers,
		Connections:      s.Connections,
		Manufacturer:     s.Manufacturer,
		Model:            s.Model,
		ModelId:          s.ModelId,
		SerialNumber:     s.SerialNumber,
		HwVersion:        s.HwVersion,
		SwVersion:        s.SwVersion,
		SuggestedArea:    s.SuggestedArea,
		ConfigurationUrl: s.ConfigurationUrl,
		ViaDevice:        s.ViaDevice,
	}
}
