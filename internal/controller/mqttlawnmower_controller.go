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

type MQTTLawnMowerReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Log        logr.Logger
	MQTTClient mqtt.Client
	base       BaseReconciler
}

func NewMQTTLawnMowerReconciler(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client) *MQTTLawnMowerReconciler {
	return &MQTTLawnMowerReconciler{
		Client:     c,
		Scheme:     scheme,
		Log:        log.WithName("mqttlawnmower"),
		MQTTClient: mqttClient,
		base: BaseReconciler{
			Client:     c,
			Log:        log.WithName("mqttlawnmower"),
			MQTTClient: mqttClient,
		},
	}
}

// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqttlawnmowers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqttlawnmowers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqttlawnmowers/finalizers,verbs=update

func (r *MQTTLawnMowerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("mqttlawnmower", req.NamespacedName)

	var obj mqttv1alpha1.MQTTLawnMower
	if err := r.Get(ctx, req.NamespacedName, &obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	wrapper := &mqttLawnMowerWrapper{&obj}

	if r.base.IsBeingDeleted(&obj) {
		if err := r.base.HandleDeletion(ctx, wrapper, "MQTTLawnMower"); err != nil {
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

	if err := r.base.PublishDiscovery(ctx, wrapper, "MQTTLawnMower", r.buildPayload); err != nil {
		log.Error(err, "Failed to publish discovery")
		if statusErr := r.base.UpdateStatusFailed(ctx, wrapper, "PublishFailed", err.Error()); statusErr != nil {
			log.Error(statusErr, "Failed to update status")
		}
		return ctrl.Result{RequeueAfter: 30 * time.Second}, err
	}

	if err := r.base.UpdateStatusPublished(ctx, wrapper, "MQTTLawnMower"); err != nil {
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

func (r *MQTTLawnMowerReconciler) buildPayload(obj EntityObject, uniqueID string) (*payload.Builder, error) {
	entity := obj.(*mqttLawnMowerWrapper).MQTTLawnMower
	spec := &entity.Spec

	pb := payload.New()

	pb.Set("name", spec.Name)
	pb.Set("activityStateTopic", spec.ActivityStateTopic)
	pb.Set("activityValueTemplate", spec.ActivityValueTemplate)
	pb.Set("dockCommandTopic", spec.DockCommandTopic)
	pb.Set("dockCommandTemplate", spec.DockCommandTemplate)
	pb.Set("pauseCommandTopic", spec.PauseCommandTopic)
	pb.Set("pauseCommandTemplate", spec.PauseCommandTemplate)
	pb.Set("startMowingCommandTopic", spec.StartMowingCommandTopic)
	pb.Set("startMowingCommandTemplate", spec.StartMowingCommandTemplate)
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

func (r *MQTTLawnMowerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mqttv1alpha1.MQTTLawnMower{}).
		Complete(r)
}

type mqttLawnMowerWrapper struct {
	*mqttv1alpha1.MQTTLawnMower
}

func (w *mqttLawnMowerWrapper) GetCommonSpec() *mqttv1alpha1.CommonSpec {
	return &w.Spec.CommonSpec
}

func (w *mqttLawnMowerWrapper) GetCommonStatus() *mqttv1alpha1.CommonStatus {
	return &w.Status.CommonStatus
}

func (w *mqttLawnMowerWrapper) SetCommonStatus(status mqttv1alpha1.CommonStatus) {
	w.Status.CommonStatus = status
}

func (w *mqttLawnMowerWrapper) GetObject() client.Object {
	return w.MQTTLawnMower
}
