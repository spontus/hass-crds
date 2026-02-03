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

type MQTTWaterHeaterReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Log        logr.Logger
	MQTTClient mqtt.Client
	base       BaseReconciler
}

func NewMQTTWaterHeaterReconciler(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client) *MQTTWaterHeaterReconciler {
	return &MQTTWaterHeaterReconciler{
		Client:     c,
		Scheme:     scheme,
		Log:        log.WithName("mqttwaterheater"),
		MQTTClient: mqttClient,
		base: BaseReconciler{
			Client:     c,
			Log:        log.WithName("mqttwaterheater"),
			MQTTClient: mqttClient,
		},
	}
}

// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqttwaterheaters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqttwaterheaters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqttwaterheaters/finalizers,verbs=update

func (r *MQTTWaterHeaterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("mqttwaterheater", req.NamespacedName)

	var obj mqttv1alpha1.MQTTWaterHeater
	if err := r.Get(ctx, req.NamespacedName, &obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	wrapper := &mqttWaterHeaterWrapper{&obj}

	if r.base.IsBeingDeleted(&obj) {
		if err := r.base.HandleDeletion(ctx, wrapper, "MQTTWaterHeater"); err != nil {
			log.Error(err, "Failed to handle deletion")
			return ctrl.Result{RequeueAfter: 30 * time.Second}, err
		}
		if err := r.base.RemoveFinalizer(ctx, &obj); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if err := r.base.EnsureFinalizer(ctx, &obj); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.base.PublishDiscovery(ctx, wrapper, "MQTTWaterHeater", r.buildPayload); err != nil {
		log.Error(err, "Failed to publish discovery")
		if statusErr := r.base.UpdateStatusFailed(ctx, wrapper, "PublishFailed", err.Error()); statusErr != nil {
			log.Error(statusErr, "Failed to update status")
		}
		return ctrl.Result{RequeueAfter: 30 * time.Second}, err
	}

	if err := r.base.UpdateStatusPublished(ctx, wrapper, "MQTTWaterHeater"); err != nil {
		log.Error(err, "Failed to update status")
		return ctrl.Result{}, err
	}

	requeueAfter := time.Duration(0)
	if obj.Spec.RediscoverInterval != "" {
		d, err := ParseRediscoverInterval(obj.Spec.RediscoverInterval)
		if err == nil && d > 0 {
			requeueAfter = d
		}
	}

	if requeueAfter > 0 {
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}
	return ctrl.Result{}, nil
}

func (r *MQTTWaterHeaterReconciler) buildPayload(obj EntityObject, uniqueID string) (*payload.Builder, error) {
	entity := obj.(*mqttWaterHeaterWrapper).MQTTWaterHeater
	spec := &entity.Spec

	pb := payload.New()

	pb.Set("name", spec.Name)
	pb.Set("temperatureCommandTopic", spec.TemperatureCommandTopic)
	pb.Set("temperatureStateTopic", spec.TemperatureStateTopic)
	pb.Set("temperatureCommandTemplate", spec.TemperatureCommandTemplate)
	pb.Set("temperatureStateTemplate", spec.TemperatureStateTemplate)
	pb.Set("currentTemperatureTopic", spec.CurrentTemperatureTopic)
	pb.Set("currentTemperatureTemplate", spec.CurrentTemperatureTemplate)
	pb.Set("modeCommandTopic", spec.ModeCommandTopic)
	pb.Set("modeStateTopic", spec.ModeStateTopic)
	pb.Set("modeCommandTemplate", spec.ModeCommandTemplate)
	pb.Set("modeStateTemplate", spec.ModeStateTemplate)
	pb.Set("modes", spec.Modes)
	pb.Set("powerCommandTopic", spec.PowerCommandTopic)
	pb.Set("payloadOn", spec.PayloadOn)
	pb.Set("payloadOff", spec.PayloadOff)
	pb.Set("minTemp", spec.MinTemp)
	pb.Set("maxTemp", spec.MaxTemp)
	pb.Set("temperatureUnit", spec.TemperatureUnit)
	pb.Set("precision", spec.Precision)
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

func (r *MQTTWaterHeaterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mqttv1alpha1.MQTTWaterHeater{}).
		Complete(r)
}

type mqttWaterHeaterWrapper struct {
	*mqttv1alpha1.MQTTWaterHeater
}

func (w *mqttWaterHeaterWrapper) GetCommonSpec() *mqttv1alpha1.CommonSpec {
	return &w.Spec.CommonSpec
}

func (w *mqttWaterHeaterWrapper) GetCommonStatus() *mqttv1alpha1.CommonStatus {
	return &w.Status.CommonStatus
}

func (w *mqttWaterHeaterWrapper) SetCommonStatus(status mqttv1alpha1.CommonStatus) {
	w.Status.CommonStatus = status
}

func (w *mqttWaterHeaterWrapper) GetObject() client.Object {
	return w.MQTTWaterHeater
}
