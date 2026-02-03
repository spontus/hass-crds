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

type MQTTTagReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Log        logr.Logger
	MQTTClient mqtt.Client
	base       BaseReconciler
}

func NewMQTTTagReconciler(c client.Client, scheme *runtime.Scheme, log logr.Logger, mqttClient mqtt.Client) *MQTTTagReconciler {
	return &MQTTTagReconciler{
		Client:     c,
		Scheme:     scheme,
		Log:        log.WithName("mqtttag"),
		MQTTClient: mqttClient,
		base: BaseReconciler{
			Client:     c,
			Log:        log.WithName("mqtttag"),
			MQTTClient: mqttClient,
		},
	}
}

// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqtttags,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqtttags/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mqtt.home-assistant.io,resources=mqtttags/finalizers,verbs=update

func (r *MQTTTagReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("mqtttag", req.NamespacedName)

	var obj mqttv1alpha1.MQTTTag
	if err := r.Get(ctx, req.NamespacedName, &obj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	wrapper := &mqttTagWrapper{&obj}

	if r.base.IsBeingDeleted(&obj) {
		if err := r.base.HandleDeletion(ctx, wrapper, "MQTTTag"); err != nil {
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

	if err := r.base.PublishDiscovery(ctx, wrapper, "MQTTTag", r.buildPayload); err != nil {
		log.Error(err, "Failed to publish discovery")
		if statusErr := r.base.UpdateStatusFailed(ctx, wrapper, "PublishFailed", err.Error()); statusErr != nil {
			log.Error(statusErr, "Failed to update status")
		}
		return ctrl.Result{RequeueAfter: 30 * time.Second}, err
	}

	if err := r.base.UpdateStatusPublished(ctx, wrapper, "MQTTTag"); err != nil {
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

func (r *MQTTTagReconciler) buildPayload(obj EntityObject, uniqueID string) (*payload.Builder, error) {
	entity := obj.(*mqttTagWrapper).MQTTTag
	spec := &entity.Spec

	pb := payload.New()

	pb.Set("topic", spec.Topic)
	pb.Set("valueTemplate", spec.ValueTemplate)
	pb.Set("qos", spec.Qos)
	pb.Set("encoding", spec.Encoding)

	return pb, nil
}

func (r *MQTTTagReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mqttv1alpha1.MQTTTag{}).
		Complete(r)
}

type mqttTagWrapper struct {
	*mqttv1alpha1.MQTTTag
}

func (w *mqttTagWrapper) GetCommonSpec() *mqttv1alpha1.CommonSpec {
	return &w.Spec.CommonSpec
}

func (w *mqttTagWrapper) GetCommonStatus() *mqttv1alpha1.CommonStatus {
	return &w.Status.CommonStatus
}

func (w *mqttTagWrapper) SetCommonStatus(status mqttv1alpha1.CommonStatus) {
	w.Status.CommonStatus = status
}

func (w *mqttTagWrapper) GetObject() client.Object {
	return w.MQTTTag
}
