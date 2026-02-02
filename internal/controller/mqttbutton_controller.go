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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mqttv1alpha1 "github.com/spontus/hass-crds/api/v1alpha1"
	"github.com/spontus/hass-crds/internal/mqtt"
	"github.com/spontus/hass-crds/internal/payload"
)

// MQTTButtonReconciler reconciles a MQTTButton object.
type MQTTButtonReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Log        logr.Logger
	MQTTClient mqtt.Client
	base       BaseReconciler
}

// NewMQTTButtonReconciler creates a new MQTTButtonReconciler.
func NewMQTTButtonReconciler(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client) *MQTTButtonReconciler {
	return &MQTTButtonReconciler{
		Client:     c,
		Scheme:     scheme,
		Log:        log.WithName("mqttbutton"),
		MQTTClient: mqttClient,
		base: BaseReconciler{
			Client:     c,
			Log:        log.WithName("mqttbutton"),
			MQTTClient: mqttClient,
		},
	}
}

// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqttbuttons,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqttbuttons/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqttbuttons/finalizers,verbs=update

func (r *MQTTButtonReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("mqttbutton", req.NamespacedName)

	// Fetch the MQTTButton instance
	var button mqttv1alpha1.MQTTButton
	if err := r.Get(ctx, req.NamespacedName, &button); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Create a wrapper that implements EntityObject
	wrapper := &mqttButtonWrapper{&button}

	// Check if being deleted
	if r.base.IsBeingDeleted(&button) {
		if err := r.base.HandleDeletion(ctx, wrapper, "MQTTButton"); err != nil {
			log.Error(err, "Failed to handle deletion")
			return ctrl.Result{RequeueAfter: 30 * time.Second}, err
		}
		if err := r.base.RemoveFinalizer(ctx, &button); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Ensure finalizer
	if err := r.base.EnsureFinalizer(ctx, &button); err != nil {
		return ctrl.Result{}, err
	}

	// Publish discovery message
	if err := r.base.PublishDiscovery(ctx, wrapper, "MQTTButton", r.buildPayload); err != nil {
		log.Error(err, "Failed to publish discovery")
		if statusErr := r.base.UpdateStatusFailed(ctx, wrapper, "PublishFailed", err.Error()); statusErr != nil {
			log.Error(statusErr, "Failed to update status")
		}
		return ctrl.Result{RequeueAfter: 30 * time.Second}, err
	}

	// Update status
	if err := r.base.UpdateStatusPublished(ctx, wrapper, "MQTTButton"); err != nil {
		log.Error(err, "Failed to update status")
		return ctrl.Result{}, err
	}

	// Calculate requeue interval
	requeueAfter := time.Duration(0)
	if button.Spec.RediscoverInterval != "" {
		d, err := ParseRediscoverInterval(button.Spec.RediscoverInterval)
		if err == nil && d > 0 {
			requeueAfter = d
		}
	}

	if requeueAfter > 0 {
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}
	return ctrl.Result{}, nil
}

func (r *MQTTButtonReconciler) buildPayload(obj EntityObject, uniqueID string) (*payload.Builder, error) {
	button := obj.(*mqttButtonWrapper).MQTTButton
	spec := &button.Spec

	pb := payload.New()

	// Required field
	pb.Set("commandTopic", spec.CommandTopic)

	// Optional fields
	pb.Set("name", spec.Name)
	pb.Set("commandTemplate", spec.CommandTemplate)
	pb.Set("payloadPress", spec.PayloadPress)
	pb.Set("deviceClass", spec.DeviceClass)
	pb.Set("icon", spec.Icon)
	pb.Set("entityCategory", spec.EntityCategory)
	pb.Set("enabledByDefault", spec.EnabledByDefault)
	pb.Set("objectId", spec.ObjectId)
	pb.Set("qos", spec.Qos)
	pb.Set("retain", spec.Retain)
	pb.Set("encoding", spec.Encoding)
	pb.Set("jsonAttributesTopic", spec.JsonAttributesTopic)
	pb.Set("jsonAttributesTemplate", spec.JsonAttributesTemplate)

	return pb, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MQTTButtonReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mqttv1alpha1.MQTTButton{}).
		Complete(r)
}

// mqttButtonWrapper wraps MQTTButton to implement EntityObject.
type mqttButtonWrapper struct {
	*mqttv1alpha1.MQTTButton
}

func (w *mqttButtonWrapper) GetCommonSpec() *mqttv1alpha1.CommonSpec {
	return &w.Spec.CommonSpec
}

func (w *mqttButtonWrapper) GetCommonStatus() *mqttv1alpha1.CommonStatus {
	return &w.Status.CommonStatus
}

func (w *mqttButtonWrapper) SetCommonStatus(status mqttv1alpha1.CommonStatus) {
	w.Status.CommonStatus = status
}
