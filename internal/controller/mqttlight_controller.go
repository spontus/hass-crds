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

type MQTTLightReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Log        logr.Logger
	MQTTClient mqtt.Client
	base       BaseReconciler
}

func NewMQTTLightReconciler(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client) *MQTTLightReconciler {
	return &MQTTLightReconciler{
		Client:     c,
		Scheme:     scheme,
		Log:        log.WithName("mqttlight"),
		MQTTClient: mqttClient,
		base: BaseReconciler{
			Client:     c,
			Log:        log.WithName("mqttlight"),
			MQTTClient: mqttClient,
		},
	}
}

// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqttlights,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqttlights/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqttlights/finalizers,verbs=update

func (r *MQTTLightReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("mqttlight", req.NamespacedName)

	var obj mqttv1alpha1.MQTTLight
	if err := r.Get(ctx, req.NamespacedName, &obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	wrapper := &mqttLightWrapper{&obj}

	if r.base.IsBeingDeleted(&obj) {
		if err := r.base.HandleDeletion(ctx, wrapper, "MQTTLight"); err != nil {
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

	if err := r.base.PublishDiscovery(ctx, wrapper, "MQTTLight", r.buildPayload); err != nil {
		log.Error(err, "Failed to publish discovery")
		if statusErr := r.base.UpdateStatusFailed(ctx, wrapper, "PublishFailed", err.Error()); statusErr != nil {
			log.Error(statusErr, "Failed to update status")
		}
		return ctrl.Result{RequeueAfter: 30 * time.Second}, err
	}

	if err := r.base.UpdateStatusPublished(ctx, wrapper, "MQTTLight"); err != nil {
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

func (r *MQTTLightReconciler) buildPayload(obj EntityObject, uniqueID string) (*payload.Builder, error) {
	entity := obj.(*mqttLightWrapper).MQTTLight
	spec := &entity.Spec

	pb := payload.New()

	pb.Set("commandTopic", spec.CommandTopic)
	pb.Set("name", spec.Name)
	pb.Set("schema", spec.Schema)
	pb.Set("stateTopic", spec.StateTopic)
	pb.Set("payloadOn", spec.PayloadOn)
	pb.Set("payloadOff", spec.PayloadOff)
	pb.Set("brightnessCommandTopic", spec.BrightnessCommandTopic)
	pb.Set("brightnessStateTopic", spec.BrightnessStateTopic)
	pb.Set("brightnessScale", spec.BrightnessScale)
	pb.Set("brightnessValueTemplate", spec.BrightnessValueTemplate)
	pb.Set("colorTempCommandTopic", spec.ColorTempCommandTopic)
	pb.Set("colorTempStateTopic", spec.ColorTempStateTopic)
	pb.Set("colorTempValueTemplate", spec.ColorTempValueTemplate)
	pb.Set("rgbCommandTopic", spec.RgbCommandTopic)
	pb.Set("rgbStateTopic", spec.RgbStateTopic)
	pb.Set("rgbCommandTemplate", spec.RgbCommandTemplate)
	pb.Set("rgbValueTemplate", spec.RgbValueTemplate)
	pb.Set("effectCommandTopic", spec.EffectCommandTopic)
	pb.Set("effectStateTopic", spec.EffectStateTopic)
	pb.Set("effectList", spec.EffectList)
	pb.Set("effectValueTemplate", spec.EffectValueTemplate)
	pb.Set("minMireds", spec.MinMireds)
	pb.Set("maxMireds", spec.MaxMireds)
	pb.Set("optimistic", spec.Optimistic)
	pb.Set("onCommandType", spec.OnCommandType)
	pb.Set("brightness", spec.Brightness)
	pb.Set("colorTemp", spec.ColorTemp)
	pb.Set("effect", spec.Effect)
	pb.Set("supportedColorModes", spec.SupportedColorModes)
	pb.Set("commandOnTemplate", spec.CommandOnTemplate)
	pb.Set("commandOffTemplate", spec.CommandOffTemplate)
	pb.Set("stateTemplate", spec.StateTemplate)
	pb.Set("brightnessTemplate", spec.BrightnessTemplate)
	pb.Set("colorTempTemplate", spec.ColorTempTemplate)
	pb.Set("redTemplate", spec.RedTemplate)
	pb.Set("greenTemplate", spec.GreenTemplate)
	pb.Set("blueTemplate", spec.BlueTemplate)
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

func (r *MQTTLightReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mqttv1alpha1.MQTTLight{}).
		Complete(r)
}

type mqttLightWrapper struct {
	*mqttv1alpha1.MQTTLight
}

func (w *mqttLightWrapper) GetCommonSpec() *mqttv1alpha1.CommonSpec {
	return &w.Spec.CommonSpec
}

func (w *mqttLightWrapper) GetCommonStatus() *mqttv1alpha1.CommonStatus {
	return &w.Status.CommonStatus
}

func (w *mqttLightWrapper) SetCommonStatus(status mqttv1alpha1.CommonStatus) {
	w.Status.CommonStatus = status
}

func (w *mqttLightWrapper) GetObject() client.Object {
	return w.MQTTLight
}
