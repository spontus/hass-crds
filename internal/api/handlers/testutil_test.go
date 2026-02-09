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

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func newTestEntityHandler(objects ...runtime.Object) *EntityHandler {
	scheme := runtime.NewScheme()

	dynamicClient := dynamicfake.NewSimpleDynamicClientWithCustomListKinds(
		scheme,
		map[schema.GroupVersionResource]string{
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttbuttons"}:         "MQTTButtonList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttswitches"}:        "MQTTSwitchList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttsensors"}:         "MQTTSensorList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttbinarysensors"}:   "MQTTBinarySensorList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttlights"}:          "MQTTLightList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttclimates"}:        "MQTTClimateList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttcovers"}:          "MQTTCoverList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttfans"}:            "MQTTFanList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttlocks"}:           "MQTTLockList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttvacuums"}:         "MQTTVacuumList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttnumbers"}:         "MQTTNumberList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttselects"}:         "MQTTSelectList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttscenes"}:          "MQTTSceneList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqtttexts"}:           "MQTTTextList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttevents"}:          "MQTTEventList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqtthumidifiers"}:     "MQTTHumidifierList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttwaterheaters"}:    "MQTTWaterHeaterList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttalarmcontrolpanels"}: "MQTTAlarmControlPanelList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttvalves"}:          "MQTTValveList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttlawnmowers"}:      "MQTTLawnMowerList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttsirens"}:          "MQTTSirenList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttcameras"}:         "MQTTCameraList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttimages"}:          "MQTTImageList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttnotifies"}:        "MQTTNotifyList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttupdates"}:         "MQTTUpdateList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttdevicetrackers"}:  "MQTTDeviceTrackerList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqtttags"}:            "MQTTTagList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttdevicetriggers"}:  "MQTTDeviceTriggerList",
			{Group: "mqtt.home-assistant.io", Version: "v1alpha1", Resource: "mqttdevices"}:         "MQTTDeviceList",
		},
		objects...,
	)

	return NewEntityHandler(dynamicClient, nil, logr.Discard())
}

func newTestNamespaceHandler(objects ...client.Object) *NamespaceHandler {
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)

	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(objects...).
		Build()

	return NewNamespaceHandler(fakeClient, logr.Discard())
}

func executeRequest(handler http.HandlerFunc, method, path string, body interface{}, urlParams map[string]string) *httptest.ResponseRecorder {
	var reqBody io.Reader
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		reqBody = bytes.NewReader(jsonBytes)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	if len(urlParams) > 0 {
		rctx := chi.NewRouteContext()
		for key, value := range urlParams {
			rctx.URLParams.Add(key, value)
		}
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

func newTestEntity(kind, namespace, name string, published bool) *unstructured.Unstructured {
	et := GetEntityTypeByKind(kind)
	if et == nil {
		return nil
	}

	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "mqtt.home-assistant.io/v1alpha1",
			"kind":       kind,
			"metadata": map[string]interface{}{
				"name":              name,
				"namespace":         namespace,
				"resourceVersion":   "1",
				"creationTimestamp": "2024-01-01T00:00:00Z",
			},
			"spec": map[string]interface{}{
				"name": "Display Name for " + name,
			},
		},
	}

	if published {
		obj.Object["status"] = map[string]interface{}{
			"conditions": []interface{}{
				map[string]interface{}{
					"type":   "Published",
					"status": "True",
				},
			},
		}
	}

	return obj
}

func newTestNamespace(name string, phase corev1.NamespacePhase) *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Status: corev1.NamespaceStatus{
			Phase: phase,
		},
	}
}

func parseJSONResponse(rr *httptest.ResponseRecorder, v interface{}) error {
	return json.Unmarshal(rr.Body.Bytes(), v)
}
