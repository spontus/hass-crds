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

// MQTTLockSpec defines the desired state of MQTTLock.
type MQTTLockSpec struct {
	CommonSpec `json:",inline"`

	// CommandTopic is the topic to publish lock/unlock commands
	CommandTopic string `json:"commandTopic"`

	// StateTopic is the topic to read current lock state
	// +optional
	StateTopic string `json:"stateTopic,omitempty"`

	// CommandTemplate is the template for the command payload
	// +optional
	CommandTemplate string `json:"commandTemplate,omitempty"`

	// ValueTemplate is the template to extract state from payload
	// +optional
	ValueTemplate string `json:"valueTemplate,omitempty"`

	// PayloadLock is the payload for lock command (default: LOCK)
	// +optional
	PayloadLock string `json:"payloadLock,omitempty"`

	// PayloadUnlock is the payload for unlock command (default: UNLOCK)
	// +optional
	PayloadUnlock string `json:"payloadUnlock,omitempty"`

	// PayloadOpen is the payload for open command (unlatch)
	// +optional
	PayloadOpen string `json:"payloadOpen,omitempty"`

	// StateLocked is the state value meaning locked (default: LOCKED)
	// +optional
	StateLocked string `json:"stateLocked,omitempty"`

	// StateUnlocked is the state value meaning unlocked (default: UNLOCKED)
	// +optional
	StateUnlocked string `json:"stateUnlocked,omitempty"`

	// StateLocking is the state value meaning locking (default: LOCKING)
	// +optional
	StateLocking string `json:"stateLocking,omitempty"`

	// StateUnlocking is the state value meaning unlocking (default: UNLOCKING)
	// +optional
	StateUnlocking string `json:"stateUnlocking,omitempty"`

	// StateJammed is the state value meaning jammed (default: JAMMED)
	// +optional
	StateJammed string `json:"stateJammed,omitempty"`

	// CodeFormat is the regex for valid codes (e.g. ^\d{4}$)
	// +optional
	CodeFormat string `json:"codeFormat,omitempty"`

	// Optimistic indicates whether to assume state changes immediately
	// +optional
	Optimistic *bool `json:"optimistic,omitempty"`
}

// MQTTLockStatus defines the observed state of MQTTLock.
type MQTTLockStatus struct {
	CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MQTTLock is the Schema for the mqttlocks API.
// It is a lock entity with optional code support.
type MQTTLock struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MQTTLockSpec   `json:"spec,omitempty"`
	Status MQTTLockStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MQTTLockList contains a list of MQTTLock.
type MQTTLockList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MQTTLock `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MQTTLock{}, &MQTTLockList{})
}
