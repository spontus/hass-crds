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

type MQTTAlarmControlPanelReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Log        logr.Logger
	MQTTClient mqtt.Client
	base       BaseReconciler
}

func NewMQTTAlarmControlPanelReconciler(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client) *MQTTAlarmControlPanelReconciler {
	return &MQTTAlarmControlPanelReconciler{
		Client:     c,
		Scheme:     scheme,
		Log:        log.WithName("mqttalarmcontrolpanel"),
		MQTTClient: mqttClient,
		base: BaseReconciler{
			Client:     c,
			Log:        log.WithName("mqttalarmcontrolpanel"),
			MQTTClient: mqttClient,
		},
	}
}

// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqttalarmcontrolpanels,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqttalarmcontrolpanels/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqttalarmcontrolpanels/finalizers,verbs=update

func (r *MQTTAlarmControlPanelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("mqttalarmcontrolpanel", req.NamespacedName)

	var obj mqttv1alpha1.MQTTAlarmControlPanel
	if err := r.Get(ctx, req.NamespacedName, &obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	wrapper := &mqttAlarmControlPanelWrapper{&obj}

	if r.base.IsBeingDeleted(&obj) {
		if err := r.base.HandleDeletion(ctx, wrapper, "MQTTAlarmControlPanel"); err != nil {
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

	if err := r.base.PublishDiscovery(ctx, wrapper, "MQTTAlarmControlPanel", r.buildPayload); err != nil {
		log.Error(err, "Failed to publish discovery")
		if statusErr := r.base.UpdateStatusFailed(ctx, wrapper, "PublishFailed", err.Error()); statusErr != nil {
			log.Error(statusErr, "Failed to update status")
		}
		return ctrl.Result{RequeueAfter: 30 * time.Second}, err
	}

	if err := r.base.UpdateStatusPublished(ctx, wrapper, "MQTTAlarmControlPanel"); err != nil {
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

func (r *MQTTAlarmControlPanelReconciler) buildPayload(obj EntityObject, uniqueID string) (*payload.Builder, error) {
	entity := obj.(*mqttAlarmControlPanelWrapper).MQTTAlarmControlPanel
	spec := &entity.Spec

	pb := payload.New()

	pb.Set("commandTopic", spec.CommandTopic)
	pb.Set("stateTopic", spec.StateTopic)
	pb.Set("name", spec.Name)
	pb.Set("commandTemplate", spec.CommandTemplate)
	pb.Set("valueTemplate", spec.ValueTemplate)
	pb.Set("payloadArmHome", spec.PayloadArmHome)
	pb.Set("payloadArmAway", spec.PayloadArmAway)
	pb.Set("payloadArmNight", spec.PayloadArmNight)
	pb.Set("payloadArmVacation", spec.PayloadArmVacation)
	pb.Set("payloadArmCustomBypass", spec.PayloadArmCustomBypass)
	pb.Set("payloadDisarm", spec.PayloadDisarm)
	pb.Set("payloadTrigger", spec.PayloadTrigger)
	pb.Set("codeArmRequired", spec.CodeArmRequired)
	pb.Set("codeDisarmRequired", spec.CodeDisarmRequired)
	pb.Set("codeTriggerRequired", spec.CodeTriggerRequired)
	pb.Set("codeFormat", spec.CodeFormat)
	pb.Set("supportedFeatures", spec.SupportedFeatures)
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

func (r *MQTTAlarmControlPanelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mqttv1alpha1.MQTTAlarmControlPanel{}).
		Complete(r)
}

type mqttAlarmControlPanelWrapper struct {
	*mqttv1alpha1.MQTTAlarmControlPanel
}

func (w *mqttAlarmControlPanelWrapper) GetCommonSpec() *mqttv1alpha1.CommonSpec {
	return &w.Spec.CommonSpec
}

func (w *mqttAlarmControlPanelWrapper) GetCommonStatus() *mqttv1alpha1.CommonStatus {
	return &w.Status.CommonStatus
}

func (w *mqttAlarmControlPanelWrapper) SetCommonStatus(status mqttv1alpha1.CommonStatus) {
	w.Status.CommonStatus = status
}

func (w *mqttAlarmControlPanelWrapper) GetObject() client.Object {
	return w.MQTTAlarmControlPanel
}
