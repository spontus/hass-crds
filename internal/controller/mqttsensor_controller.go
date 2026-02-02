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

// MQTTSensorReconciler reconciles a MQTTSensor object.
type MQTTSensorReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Log        logr.Logger
	MQTTClient mqtt.Client
	base       BaseReconciler
}

// NewMQTTSensorReconciler creates a new MQTTSensorReconciler.
func NewMQTTSensorReconciler(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client) *MQTTSensorReconciler {
	return &MQTTSensorReconciler{
		Client:     c,
		Scheme:     scheme,
		Log:        log.WithName("mqttsensor"),
		MQTTClient: mqttClient,
		base: BaseReconciler{
			Client:     c,
			Log:        log.WithName("mqttsensor"),
			MQTTClient: mqttClient,
		},
	}
}

// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqttsensors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqttsensors/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqttsensors/finalizers,verbs=update

func (r *MQTTSensorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("mqttsensor", req.NamespacedName)

	var sensor mqttv1alpha1.MQTTSensor
	if err := r.Get(ctx, req.NamespacedName, &sensor); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	wrapper := &mqttSensorWrapper{&sensor}

	if r.base.IsBeingDeleted(&sensor) {
		if err := r.base.HandleDeletion(ctx, wrapper, "MQTTSensor"); err != nil {
			log.Error(err, "Failed to handle deletion")
			return ctrl.Result{RequeueAfter: 30 * time.Second}, err
		}
		if err := r.base.RemoveFinalizer(ctx, &sensor); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if err := r.base.EnsureFinalizer(ctx, &sensor); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.base.PublishDiscovery(ctx, wrapper, "MQTTSensor", r.buildPayload); err != nil {
		log.Error(err, "Failed to publish discovery")
		if statusErr := r.base.UpdateStatusFailed(ctx, wrapper, "PublishFailed", err.Error()); statusErr != nil {
			log.Error(statusErr, "Failed to update status")
		}
		return ctrl.Result{RequeueAfter: 30 * time.Second}, err
	}

	if err := r.base.UpdateStatusPublished(ctx, wrapper, "MQTTSensor"); err != nil {
		log.Error(err, "Failed to update status")
		return ctrl.Result{}, err
	}

	requeueAfter := time.Duration(0)
	if sensor.Spec.RediscoverInterval != "" {
		d, err := ParseRediscoverInterval(sensor.Spec.RediscoverInterval)
		if err == nil && d > 0 {
			requeueAfter = d
		}
	}

	if requeueAfter > 0 {
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}
	return ctrl.Result{}, nil
}

func (r *MQTTSensorReconciler) buildPayload(obj EntityObject, uniqueID string) (*payload.Builder, error) {
	sensor := obj.(*mqttSensorWrapper).MQTTSensor
	spec := &sensor.Spec

	pb := payload.New()

	pb.Set("stateTopic", spec.StateTopic)
	pb.Set("name", spec.Name)
	pb.Set("valueTemplate", spec.ValueTemplate)
	pb.Set("unitOfMeasurement", spec.UnitOfMeasurement)
	pb.Set("deviceClass", spec.DeviceClass)
	pb.Set("stateClass", spec.StateClass)
	pb.Set("expireAfter", spec.ExpireAfter)
	pb.Set("forceUpdate", spec.ForceUpdate)
	pb.Set("lastResetValueTemplate", spec.LastResetValueTemplate)
	pb.Set("suggestedDisplayPrecision", spec.SuggestedDisplayPrecision)
	pb.Set("icon", spec.Icon)
	pb.Set("entityCategory", spec.EntityCategory)
	pb.Set("enabledByDefault", spec.EnabledByDefault)
	pb.Set("objectId", spec.ObjectId)
	pb.Set("qos", spec.Qos)
	pb.Set("encoding", spec.Encoding)
	pb.Set("jsonAttributesTopic", spec.JsonAttributesTopic)
	pb.Set("jsonAttributesTemplate", spec.JsonAttributesTemplate)

	return pb, nil
}

func (r *MQTTSensorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mqttv1alpha1.MQTTSensor{}).
		Complete(r)
}

type mqttSensorWrapper struct {
	*mqttv1alpha1.MQTTSensor
}

func (w *mqttSensorWrapper) GetCommonSpec() *mqttv1alpha1.CommonSpec {
	return &w.Spec.CommonSpec
}

func (w *mqttSensorWrapper) GetCommonStatus() *mqttv1alpha1.CommonStatus {
	return &w.Status.CommonStatus
}

func (w *mqttSensorWrapper) SetCommonStatus(status mqttv1alpha1.CommonStatus) {
	w.Status.CommonStatus = status
}
