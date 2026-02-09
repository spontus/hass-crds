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

// EntityMetadata contains common fields for all MQTT entity types.
type EntityMetadata struct {
	// Name is the display name in Home Assistant
	// +optional
	Name string `json:"name,omitempty"`

	// UniqueId is the unique identifier for HA entity registry (defaults to <namespace>-<name>)
	// +optional
	UniqueId string `json:"uniqueId,omitempty"`

	// Icon is the MDI icon (e.g. mdi:thermometer)
	// +optional
	Icon string `json:"icon,omitempty"`

	// EntityCategory is the entity category
	// +kubebuilder:validation:Enum=config;diagnostic
	// +optional
	EntityCategory string `json:"entityCategory,omitempty"`

	// EnabledByDefault indicates whether the entity is enabled when first discovered
	// +optional
	EnabledByDefault *bool `json:"enabledByDefault,omitempty"`

	// ObjectId is the override for HA entity ID generation
	// +optional
	ObjectId string `json:"objectId,omitempty"`
}

// DeviceBlock contains the device configuration for Home Assistant device registry.
type DeviceBlock struct {
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

// DeviceRef is a reference to an MQTTDevice resource instead of inline device block.
type DeviceRef struct {
	// Name is the name of an MQTTDevice resource in the same namespace
	Name string `json:"name"`
}

// AvailabilityConfig defines an availability topic configuration.
type AvailabilityConfig struct {
	// Topic is the MQTT topic for availability
	Topic string `json:"topic"`

	// PayloadAvailable is the payload indicating available (default: online)
	// +optional
	PayloadAvailable string `json:"payloadAvailable,omitempty"`

	// PayloadNotAvailable is the payload indicating unavailable (default: offline)
	// +optional
	PayloadNotAvailable string `json:"payloadNotAvailable,omitempty"`

	// ValueTemplate is the template to extract availability from payload
	// +optional
	ValueTemplate string `json:"valueTemplate,omitempty"`
}

// CommonSpec contains fields common to all MQTT entity specs.
type CommonSpec struct {
	EntityMetadata `json:",inline"`

	// Device is the device configuration for Home Assistant device registry
	// +optional
	Device *DeviceBlock `json:"device,omitempty"`

	// DeviceRef is a reference to an MQTTDevice resource instead of inline device block
	// +optional
	DeviceRef *DeviceRef `json:"deviceRef,omitempty"`

	// Availability is a list of availability topics
	// +optional
	Availability []AvailabilityConfig `json:"availability,omitempty"`

	// AvailabilityTopic is a simple availability topic (shorthand for single availability)
	// +optional
	AvailabilityTopic string `json:"availabilityTopic,omitempty"`

	// AvailabilityMode is how to combine multiple availability topics
	// +kubebuilder:validation:Enum=all;any;latest
	// +optional
	AvailabilityMode string `json:"availabilityMode,omitempty"`

	// Qos is the MQTT QoS level
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=2
	// +optional
	Qos *int `json:"qos,omitempty"`

	// Retain indicates whether to retain messages on command/state topics
	// +optional
	Retain *bool `json:"retain,omitempty"`

	// Encoding is the payload encoding (default: utf-8)
	// +optional
	Encoding string `json:"encoding,omitempty"`

	// JsonAttributesTopic is the MQTT topic for JSON attributes
	// +optional
	JsonAttributesTopic string `json:"jsonAttributesTopic,omitempty"`

	// JsonAttributesTemplate is the template to extract attributes from payload
	// +optional
	JsonAttributesTemplate string `json:"jsonAttributesTemplate,omitempty"`

	// RediscoverInterval is how often to re-publish the discovery config payload (e.g. 5m, 1h)
	// +optional
	RediscoverInterval string `json:"rediscoverInterval,omitempty"`
}

// Condition contains details for the current condition of this resource.
type Condition struct {
	// Type is the condition type (Published, MQTTConnected)
	Type string `json:"type"`

	// Status is the status of the condition
	// +kubebuilder:validation:Enum=True;False;Unknown
	Status string `json:"status"`

	// LastTransitionTime is the last time the condition transitioned
	// +optional
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`

	// Reason is a brief reason for the condition's last transition
	// +optional
	Reason string `json:"reason,omitempty"`

	// Message is a human-readable message indicating details
	// +optional
	Message string `json:"message,omitempty"`
}

// CommonStatus contains fields common to all MQTT entity statuses.
type CommonStatus struct {
	// ObservedGeneration is the generation observed by the controller
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// LastPublished is the timestamp of last discovery publish
	// +optional
	LastPublished *metav1.Time `json:"lastPublished,omitempty"`

	// DiscoveryTopic is the MQTT discovery topic path
	// +optional
	DiscoveryTopic string `json:"discoveryTopic,omitempty"`

	// Conditions is the list of conditions for this resource
	// +optional
	Conditions []Condition `json:"conditions,omitempty"`
}

// ConditionType constants for status conditions.
const (
	ConditionTypePublished     = "Published"
	ConditionTypeMQTTConnected = "MQTTConnected"
)

// ConditionStatus constants.
const (
	ConditionTrue    = "True"
	ConditionFalse   = "False"
	ConditionUnknown = "Unknown"
)
