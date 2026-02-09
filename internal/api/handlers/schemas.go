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
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-logr/logr"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type EntityType struct {
	Kind        string `json:"kind"`
	Plural      string `json:"plural"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

var entityTypes = []EntityType{
	{Kind: "MQTTButton", Plural: "mqttbuttons", Description: "Stateless button that publishes when pressed", Category: "Controls"},
	{Kind: "MQTTSwitch", Plural: "mqttswitches", Description: "On/off switch with state", Category: "Controls"},
	{Kind: "MQTTScene", Plural: "mqttscenes", Description: "Scene activation", Category: "Controls"},
	{Kind: "MQTTSelect", Plural: "mqttselects", Description: "Dropdown selector from options", Category: "Controls"},
	{Kind: "MQTTNumber", Plural: "mqttnumbers", Description: "Numeric input with min/max", Category: "Controls"},
	{Kind: "MQTTText", Plural: "mqtttexts", Description: "Text input field", Category: "Controls"},
	{Kind: "MQTTSensor", Plural: "mqttsensors", Description: "Read-only sensor value", Category: "Sensors"},
	{Kind: "MQTTBinarySensor", Plural: "mqttbinarysensors", Description: "On/off sensor state", Category: "Sensors"},
	{Kind: "MQTTEvent", Plural: "mqttevents", Description: "Event trigger entity", Category: "Sensors"},
	{Kind: "MQTTLight", Plural: "mqttlights", Description: "Light with brightness/color", Category: "Lighting"},
	{Kind: "MQTTClimate", Plural: "mqttclimates", Description: "HVAC/thermostat control", Category: "Climate"},
	{Kind: "MQTTHumidifier", Plural: "mqtthumidifiers", Description: "Humidifier/dehumidifier", Category: "Climate"},
	{Kind: "MQTTWaterHeater", Plural: "mqttwaterheaters", Description: "Water heater control", Category: "Climate"},
	{Kind: "MQTTFan", Plural: "mqttfans", Description: "Fan with speed control", Category: "Climate"},
	{Kind: "MQTTLock", Plural: "mqttlocks", Description: "Lock/unlock control", Category: "Security"},
	{Kind: "MQTTAlarmControlPanel", Plural: "mqttalarmcontrolpanels", Description: "Alarm system control", Category: "Security"},
	{Kind: "MQTTCover", Plural: "mqttcovers", Description: "Blinds/garage doors", Category: "Covers"},
	{Kind: "MQTTValve", Plural: "mqttvalves", Description: "Water/gas valve control", Category: "Covers"},
	{Kind: "MQTTVacuum", Plural: "mqttvacuums", Description: "Robot vacuum control", Category: "Devices"},
	{Kind: "MQTTLawnMower", Plural: "mqttlawnmowers", Description: "Robot lawn mower", Category: "Devices"},
	{Kind: "MQTTSiren", Plural: "mqttsirens", Description: "Siren/alarm device", Category: "Devices"},
	{Kind: "MQTTCamera", Plural: "mqttcameras", Description: "Camera image entity", Category: "Media"},
	{Kind: "MQTTImage", Plural: "mqttimages", Description: "Static image entity", Category: "Media"},
	{Kind: "MQTTNotify", Plural: "mqttnotifies", Description: "Notification service", Category: "Media"},
	{Kind: "MQTTUpdate", Plural: "mqttupdates", Description: "Firmware update entity", Category: "Media"},
	{Kind: "MQTTDeviceTracker", Plural: "mqttdevicetrackers", Description: "Device location tracking", Category: "Tracking"},
	{Kind: "MQTTTag", Plural: "mqtttags", Description: "NFC/RFID tag scanner", Category: "Tracking"},
	{Kind: "MQTTDeviceTrigger", Plural: "mqttdevicetriggers", Description: "Device automation trigger", Category: "Tracking"},
	{Kind: "MQTTDevice", Plural: "mqttdevices", Description: "Shared device configuration", Category: "Utility"},
}

func GetEntityTypes() []EntityType {
	return entityTypes
}

func GetEntityTypeByKind(kind string) *EntityType {
	for _, et := range entityTypes {
		if et.Kind == kind {
			return &et
		}
	}
	return nil
}

type SchemaHandler struct {
	restConfig *rest.Config
	log        logr.Logger
}

func NewSchemaHandler(restConfig *rest.Config, log logr.Logger) *SchemaHandler {
	return &SchemaHandler{
		restConfig: restConfig,
		log:        log.WithName("schemas"),
	}
}

func (h *SchemaHandler) ListEntityTypes(w http.ResponseWriter, r *http.Request) {
	categories := make(map[string][]EntityType)
	for _, et := range entityTypes {
		categories[et.Category] = append(categories[et.Category], et)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"entityTypes": entityTypes,
		"categories":  categories,
	})
}

func (h *SchemaHandler) GetSchema(w http.ResponseWriter, r *http.Request) {
	kind := chi.URLParam(r, "kind")

	et := GetEntityTypeByKind(kind)
	if et == nil {
		writeError(w, http.StatusNotFound, "unknown entity type: "+kind)
		return
	}

	crdName := et.Plural + ".mqtt.home-assistant.io"

	_ = apiextensionsv1.AddToScheme(scheme.Scheme)
	crdClient, err := client.New(h.restConfig, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		h.log.Error(err, "failed to create client")
		writeError(w, http.StatusInternalServerError, "failed to create client")
		return
	}

	crd := &apiextensionsv1.CustomResourceDefinition{}
	if err := crdClient.Get(context.Background(), client.ObjectKey{Name: crdName}, crd); err != nil {
		h.log.Error(err, "failed to get CRD", "name", crdName)
		writeError(w, http.StatusNotFound, "CRD not found: "+crdName)
		return
	}

	var schemaProps map[string]interface{}
	for _, version := range crd.Spec.Versions {
		if version.Name == apiVersion && version.Schema != nil && version.Schema.OpenAPIV3Schema != nil {
			schemaProps = h.convertJSONSchemaProps(version.Schema.OpenAPIV3Schema)
			break
		}
	}

	if schemaProps == nil {
		writeError(w, http.StatusNotFound, "schema not found for "+kind)
		return
	}

	specSchema, ok := schemaProps["properties"].(map[string]interface{})["spec"]
	if !ok {
		writeError(w, http.StatusNotFound, "spec schema not found for "+kind)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"kind":        kind,
		"apiVersion":  apiGroup + "/" + apiVersion,
		"description": et.Description,
		"schema":      specSchema,
	})
}

func (h *SchemaHandler) convertJSONSchemaProps(props *apiextensionsv1.JSONSchemaProps) map[string]interface{} {
	if props == nil {
		return nil
	}

	result := make(map[string]interface{})

	if props.Type != "" {
		result["type"] = props.Type
	}
	if props.Description != "" {
		result["description"] = props.Description
	}
	if props.Default != nil {
		result["default"] = props.Default
	}
	if len(props.Enum) > 0 {
		enumVals := make([]interface{}, len(props.Enum))
		for i, e := range props.Enum {
			var val interface{}
			if err := e.UnmarshalJSON(e.Raw); err == nil {
				val = strings.Trim(string(e.Raw), "\"")
			}
			enumVals[i] = val
		}
		result["enum"] = enumVals
	}
	if props.Minimum != nil {
		result["minimum"] = *props.Minimum
	}
	if props.Maximum != nil {
		result["maximum"] = *props.Maximum
	}
	if props.Pattern != "" {
		result["pattern"] = props.Pattern
	}
	if props.Format != "" {
		result["format"] = props.Format
	}
	if len(props.Required) > 0 {
		result["required"] = props.Required
	}

	if len(props.Properties) > 0 {
		propsMap := make(map[string]interface{})
		for name, prop := range props.Properties {
			propsMap[name] = h.convertJSONSchemaProps(&prop)
		}
		result["properties"] = propsMap
	}

	if props.Items != nil && props.Items.Schema != nil {
		result["items"] = h.convertJSONSchemaProps(props.Items.Schema)
	}

	if props.AdditionalProperties != nil && props.AdditionalProperties.Schema != nil {
		result["additionalProperties"] = h.convertJSONSchemaProps(props.AdditionalProperties.Schema)
	}

	return result
}
