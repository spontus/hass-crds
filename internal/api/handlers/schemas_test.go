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
	"net/http"
	"testing"

	"github.com/go-logr/logr"
)

func TestSchemaHandler_ListEntityTypes(t *testing.T) {
	handler := NewSchemaHandler(nil, logr.Discard())

	rr := executeRequest(handler.ListEntityTypes, http.MethodGet, "/api/v1/entity-types", nil, nil)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var response struct {
		EntityTypes []EntityType            `json:"entityTypes"`
		Categories  map[string][]EntityType `json:"categories"`
	}
	if err := parseJSONResponse(rr, &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(response.EntityTypes) == 0 {
		t.Error("expected entity types to be populated")
	}

	if len(response.Categories) == 0 {
		t.Error("expected categories to be populated")
	}

	expectedKinds := []string{"MQTTButton", "MQTTSwitch", "MQTTSensor", "MQTTLight"}
	kindSet := make(map[string]bool)
	for _, et := range response.EntityTypes {
		kindSet[et.Kind] = true
	}

	for _, kind := range expectedKinds {
		if !kindSet[kind] {
			t.Errorf("expected entity type %s to be present", kind)
		}
	}

	expectedCategories := []string{"Controls", "Sensors", "Lighting", "Climate", "Security"}
	for _, cat := range expectedCategories {
		if _, ok := response.Categories[cat]; !ok {
			t.Errorf("expected category %s to be present", cat)
		}
	}
}

func TestSchemaHandler_ListEntityTypes_ValidStructure(t *testing.T) {
	handler := NewSchemaHandler(nil, logr.Discard())

	rr := executeRequest(handler.ListEntityTypes, http.MethodGet, "/api/v1/entity-types", nil, nil)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var response struct {
		EntityTypes []EntityType            `json:"entityTypes"`
		Categories  map[string][]EntityType `json:"categories"`
	}
	if err := parseJSONResponse(rr, &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	for _, et := range response.EntityTypes {
		if et.Kind == "" {
			t.Error("entity type missing kind")
		}
		if et.Plural == "" {
			t.Error("entity type missing plural")
		}
		if et.Description == "" {
			t.Error("entity type missing description")
		}
		if et.Category == "" {
			t.Error("entity type missing category")
		}
	}

	for category, types := range response.Categories {
		if len(types) == 0 {
			t.Errorf("category %s has no types", category)
		}
		for _, et := range types {
			if et.Category != category {
				t.Errorf("entity %s in category %s has wrong category field: %s", et.Kind, category, et.Category)
			}
		}
	}
}

func TestGetEntityTypes(t *testing.T) {
	types := GetEntityTypes()

	if len(types) == 0 {
		t.Error("expected entity types to be populated")
	}

	kindCount := make(map[string]int)
	for _, et := range types {
		kindCount[et.Kind]++
	}

	for kind, count := range kindCount {
		if count > 1 {
			t.Errorf("duplicate entity type: %s appears %d times", kind, count)
		}
	}
}

func TestGetEntityTypeByKind(t *testing.T) {
	tests := []struct {
		kind        string
		expectFound bool
		expectPlural string
	}{
		{"MQTTButton", true, "mqttbuttons"},
		{"MQTTSwitch", true, "mqttswitches"},
		{"MQTTSensor", true, "mqttsensors"},
		{"MQTTLight", true, "mqttlights"},
		{"UnknownKind", false, ""},
		{"", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			et := GetEntityTypeByKind(tt.kind)

			if tt.expectFound {
				if et == nil {
					t.Errorf("expected to find entity type %s", tt.kind)
					return
				}
				if et.Kind != tt.kind {
					t.Errorf("expected kind %s, got %s", tt.kind, et.Kind)
				}
				if et.Plural != tt.expectPlural {
					t.Errorf("expected plural %s, got %s", tt.expectPlural, et.Plural)
				}
			} else {
				if et != nil {
					t.Errorf("expected not to find entity type %s", tt.kind)
				}
			}
		})
	}
}

func TestSchemaHandler_GetSchema_UnknownKind(t *testing.T) {
	handler := NewSchemaHandler(nil, logr.Discard())

	rr := executeRequest(handler.GetSchema, http.MethodGet, "/api/v1/entity-types/UnknownKind/schema", nil, map[string]string{
		"kind": "UnknownKind",
	})

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}

	var response map[string]string
	if err := parseJSONResponse(rr, &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response["error"] != "unknown entity type: UnknownKind" {
		t.Errorf("unexpected error message: %s", response["error"])
	}
}
