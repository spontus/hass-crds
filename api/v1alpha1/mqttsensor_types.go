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

// MQTTSensorSpec defines the desired state of MQTTSensor.
type MQTTSensorSpec struct {
	CommonSpec `json:",inline"`

	// StateTopic is the topic to read sensor value
	StateTopic string `json:"stateTopic"`

	// ValueTemplate is the template to extract value from payload
	// +optional
	ValueTemplate string `json:"valueTemplate,omitempty"`

	// UnitOfMeasurement is the unit displayed in HA (e.g. °C, %, W)
	// +optional
	UnitOfMeasurement string `json:"unitOfMeasurement,omitempty"`

	// DeviceClass is the HA device class (e.g. temperature, humidity, power, energy, battery)
	// +optional
	DeviceClass string `json:"deviceClass,omitempty"`

	// StateClass is the state class for statistics
	// +kubebuilder:validation:Enum=measurement;total;total_increasing
	// +optional
	StateClass string `json:"stateClass,omitempty"`

	// ExpireAfter is the seconds after which the sensor value expires
	// +optional
	ExpireAfter *int `json:"expireAfter,omitempty"`

	// ForceUpdate indicates whether to update HA state even if the value hasn't changed
	// +optional
	ForceUpdate *bool `json:"forceUpdate,omitempty"`

	// LastResetValueTemplate is the template for the last reset timestamp
	// +optional
	LastResetValueTemplate string `json:"lastResetValueTemplate,omitempty"`

	// SuggestedDisplayPrecision is the number of decimal places to display
	// +optional
	SuggestedDisplayPrecision *int `json:"suggestedDisplayPrecision,omitempty"`
}

// MQTTSensorStatus defines the observed state of MQTTSensor.
type MQTTSensorStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTSensor is the Schema for the mqttsensors API.
// It is a read-only sensor that reports a value from an MQTT topic.
type MQTTSensor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTSensorSpec   `json:"spec,omitempty"`
	Status MQTTSensorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTSensorList contains a list of MQTTSensor.
type MQTTSensorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTSensor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTSensor{}, &MQTTSensorList{})
}
