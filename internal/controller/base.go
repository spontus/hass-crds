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

package controller

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	mqttv1alpha1 "github.com/spontus/hass-crds/api/v1alpha1"
	"github.com/spontus/hass-crds/internal/mqtt"
	"github.com/spontus/hass-crds/internal/payload"
	"github.com/spontus/hass-crds/internal/topic"
)

const (
	// FinalizerName is the finalizer used by all MQTT entity controllers.
	FinalizerName = "mqtt.home-assistant.io/finalizer"

	// DefaultQoS is the default MQTT QoS level.
	DefaultQoS = byte(1)

	// DefaultRetain indicates whether discovery messages should be retained.
	DefaultRetain = true
)

// BaseReconciler contains common reconciliation logic for all MQTT entity controllers.
type BaseReconciler struct {
	Client     client.Client
	Log        logr.Logger
	MQTTClient mqtt.Client
}

// EntityObject is an interface for all MQTT entity types.
type EntityObject interface {
	client.Object
	GetCommonSpec() *mqttv1alpha1.CommonSpec
	GetCommonStatus() *mqttv1alpha1.CommonStatus
	SetCommonStatus(status mqttv1alpha1.CommonStatus)
}

// PayloadBuilder is a function that builds the discovery payload for an entity.
type PayloadBuilder func(obj EntityObject, uniqueID string) (*payload.Builder, error)

// PublishDiscovery publishes the MQTT discovery message for an entity.
func (r *BaseReconciler) PublishDiscovery(ctx context.Context, obj EntityObject, kind string, buildPayload PayloadBuilder) error {
	namespace := obj.GetNamespace()
	name := obj.GetName()

	// Get the common spec
	spec := obj.GetCommonSpec()

	// Generate unique ID
	uniqueID := topic.UniqueIDWithOverride(spec.UniqueId, namespace, name)

	// Build the payload
	pb, err := buildPayload(obj, uniqueID)
	if err != nil {
		return err
	}

	// Add unique_id to payload
	pb.Set("uniqueId", uniqueID)

	// Add device configuration if present
	if spec.Device != nil {
		device := payload.DeviceBlockToMap(
			spec.Device.Name,
			spec.Device.Identifiers,
			spec.Device.Connections,
			spec.Device.Manufacturer,
			spec.Device.Model,
			spec.Device.ModelId,
			spec.Device.SerialNumber,
			spec.Device.HwVersion,
			spec.Device.SwVersion,
			spec.Device.SuggestedArea,
			spec.Device.ConfigurationUrl,
			spec.Device.ViaDevice,
		)
		pb.SetDevice(device)
	}

	// Add availability configuration
	if len(spec.Availability) > 0 {
		var availList []map[string]interface{}
		for _, a := range spec.Availability {
			availList = append(availList, payload.AvailabilityToMap(
				a.Topic,
				a.PayloadAvailable,
				a.PayloadNotAvailable,
				a.ValueTemplate,
			))
		}
		pb.SetAvailability(availList)
		if spec.AvailabilityMode != "" {
			pb.Set("availabilityMode", spec.AvailabilityMode)
		}
	}

	// Build JSON payload
	jsonPayload, err := pb.Build()
	if err != nil {
		return err
	}

	// Generate discovery topic
	discoveryTopic := topic.DiscoveryTopic(kind, namespace, name)

	// Determine QoS
	qos := DefaultQoS
	if spec.Qos != nil {
		qos = byte(*spec.Qos)
	}

	// Publish to MQTT
	if err := r.MQTTClient.Publish(ctx, discoveryTopic, jsonPayload, qos, DefaultRetain); err != nil {
		return err
	}

	r.Log.Info("Published discovery message", "topic", discoveryTopic, "kind", kind, "name", name)
	return nil
}

// HandleDeletion publishes an empty payload to remove the entity from Home Assistant.
func (r *BaseReconciler) HandleDeletion(ctx context.Context, obj EntityObject, kind string) error {
	namespace := obj.GetNamespace()
	name := obj.GetName()

	// Generate discovery topic
	discoveryTopic := topic.DiscoveryTopic(kind, namespace, name)

	// Publish empty payload to remove entity
	if err := r.MQTTClient.Publish(ctx, discoveryTopic, []byte{}, DefaultQoS, DefaultRetain); err != nil {
		return err
	}

	r.Log.Info("Published deletion message", "topic", discoveryTopic, "kind", kind, "name", name)
	return nil
}

// UpdateStatusPublished updates the status to reflect a successful publish.
func (r *BaseReconciler) UpdateStatusPublished(ctx context.Context, obj EntityObject, kind string) error {
	namespace := obj.GetNamespace()
	name := obj.GetName()

	status := obj.GetCommonStatus()
	now := metav1.Now()
	status.LastPublished = &now
	status.DiscoveryTopic = topic.DiscoveryTopic(kind, namespace, name)
	status.ObservedGeneration = obj.GetGeneration()

	// Update or add Published condition
	r.SetCondition(status, mqttv1alpha1.ConditionTypePublished, mqttv1alpha1.ConditionTrue, "Success", "Discovery message published")

	obj.SetCommonStatus(*status)
	return r.Client.Status().Update(ctx, obj)
}

// UpdateStatusFailed updates the status to reflect a failed publish.
func (r *BaseReconciler) UpdateStatusFailed(ctx context.Context, obj EntityObject, reason, message string) error {
	status := obj.GetCommonStatus()

	r.SetCondition(status, mqttv1alpha1.ConditionTypePublished, mqttv1alpha1.ConditionFalse, reason, message)

	obj.SetCommonStatus(*status)
	return r.Client.Status().Update(ctx, obj)
}

// SetCondition updates or adds a condition to the status.
func (r *BaseReconciler) SetCondition(status *mqttv1alpha1.CommonStatus, condType, condStatus, reason, message string) {
	now := metav1.Now()
	newCondition := mqttv1alpha1.Condition{
		Type:               condType,
		Status:             condStatus,
		LastTransitionTime: &now,
		Reason:             reason,
		Message:            message,
	}

	for i, c := range status.Conditions {
		if c.Type == condType {
			if c.Status != condStatus {
				status.Conditions[i] = newCondition
			} else {
				// Only update reason/message, keep transition time
				status.Conditions[i].Reason = reason
				status.Conditions[i].Message = message
			}
			return
		}
	}

	// Condition not found, add it
	status.Conditions = append(status.Conditions, newCondition)
}

// EnsureFinalizer adds the finalizer if not present.
func (r *BaseReconciler) EnsureFinalizer(ctx context.Context, obj client.Object) error {
	if !controllerutil.ContainsFinalizer(obj, FinalizerName) {
		controllerutil.AddFinalizer(obj, FinalizerName)
		return r.Client.Update(ctx, obj)
	}
	return nil
}

// RemoveFinalizer removes the finalizer.
func (r *BaseReconciler) RemoveFinalizer(ctx context.Context, obj client.Object) error {
	if controllerutil.ContainsFinalizer(obj, FinalizerName) {
		controllerutil.RemoveFinalizer(obj, FinalizerName)
		return r.Client.Update(ctx, obj)
	}
	return nil
}

// IsBeingDeleted checks if the object is being deleted.
func (r *BaseReconciler) IsBeingDeleted(obj client.Object) bool {
	return !obj.GetDeletionTimestamp().IsZero()
}

// ParseRediscoverInterval parses the rediscoverInterval string into a duration.
func ParseRediscoverInterval(interval string) (time.Duration, error) {
	if interval == "" {
		return 0, nil
	}
	return time.ParseDuration(interval)
}
