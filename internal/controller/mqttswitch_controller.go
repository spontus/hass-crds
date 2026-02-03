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

// MQTTSwitchReconciler reconciles a MQTTSwitch object.
type MQTTSwitchReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Log        logr.Logger
	MQTTClient mqtt.Client
	base       BaseReconciler
}

// NewMQTTSwitchReconciler creates a new MQTTSwitchReconciler.
func NewMQTTSwitchReconciler(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client) *MQTTSwitchReconciler {
	return &MQTTSwitchReconciler{
		Client:     c,
		Scheme:     scheme,
		Log:        log.WithName("mqttswitch"),
		MQTTClient: mqttClient,
		base: BaseReconciler{
			Client:     c,
			Log:        log.WithName("mqttswitch"),
			MQTTClient: mqttClient,
		},
	}
}

// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqttswitches,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqttswitches/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqttswitches/finalizers,verbs=update

func (r *MQTTSwitchReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("mqttswitch", req.NamespacedName)

	var sw mqttv1alpha1.MQTTSwitch
	if err := r.Get(ctx, req.NamespacedName, &sw); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	wrapper := &mqttSwitchWrapper{&sw}

	if r.base.IsBeingDeleted(&sw) {
		if err := r.base.HandleDeletion(ctx, wrapper, "MQTTSwitch"); err != nil {
			log.Error(err, "Failed to handle deletion")
			return ctrl.Result{RequeueAfter: 30 * time.Second}, err
		}
		if err := r.base.RemoveFinalizer(ctx, &sw); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if err := r.base.EnsureFinalizer(ctx, &sw); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.base.PublishDiscovery(ctx, wrapper, "MQTTSwitch", r.buildPayload); err != nil {
		log.Error(err, "Failed to publish discovery")
		if statusErr := r.base.UpdateStatusFailed(ctx, wrapper, "PublishFailed", err.Error()); statusErr != nil {
			log.Error(statusErr, "Failed to update status")
		}
		return ctrl.Result{RequeueAfter: 30 * time.Second}, err
	}

	if err := r.base.UpdateStatusPublished(ctx, wrapper, "MQTTSwitch"); err != nil {
		log.Error(err, "Failed to update status")
		return ctrl.Result{}, err
	}

	requeueAfter := time.Duration(0)
	if sw.Spec.RediscoverInterval != "" {
		d, err := ParseRediscoverInterval(sw.Spec.RediscoverInterval)
		if err == nil && d > 0 {
			requeueAfter = d
		}
	}

	if requeueAfter > 0 {
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}
	return ctrl.Result{}, nil
}

func (r *MQTTSwitchReconciler) buildPayload(obj EntityObject, uniqueID string) (*payload.Builder, error) {
	sw := obj.(*mqttSwitchWrapper).MQTTSwitch
	spec := &sw.Spec

	pb := payload.New()

	pb.Set("commandTopic", spec.CommandTopic)
	pb.Set("stateTopic", spec.StateTopic)
	pb.Set("name", spec.Name)
	pb.Set("commandTemplate", spec.CommandTemplate)
	pb.Set("valueTemplate", spec.ValueTemplate)
	pb.Set("payloadOn", spec.PayloadOn)
	pb.Set("payloadOff", spec.PayloadOff)
	pb.Set("stateOn", spec.StateOn)
	pb.Set("stateOff", spec.StateOff)
	pb.Set("deviceClass", spec.DeviceClass)
	pb.Set("optimistic", spec.Optimistic)
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

func (r *MQTTSwitchReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mqttv1alpha1.MQTTSwitch{}).
		Complete(r)
}

type mqttSwitchWrapper struct {
	*mqttv1alpha1.MQTTSwitch
}

func (w *mqttSwitchWrapper) GetCommonSpec() *mqttv1alpha1.CommonSpec {
	return &w.Spec.CommonSpec
}

func (w *mqttSwitchWrapper) GetCommonStatus() *mqttv1alpha1.CommonStatus {
	return &w.Status.CommonStatus
}

func (w *mqttSwitchWrapper) SetCommonStatus(status mqttv1alpha1.CommonStatus) {
	w.Status.CommonStatus = status
}

func (w *mqttSwitchWrapper) GetObject() client.Object {
	return w.MQTTSwitch
}
