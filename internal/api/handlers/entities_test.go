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
)

func TestEntityHandler_List_AllEntities(t *testing.T) {
	button := newTestEntity("MQTTButton", "default", "test-button", true)
	sensor := newTestEntity("MQTTSensor", "kube-system", "test-sensor", false)

	handler := newTestEntityHandler(button, sensor)

	rr := executeRequest(handler.List, http.MethodGet, "/api/v1/entities", nil, nil)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var response EntityListResponse
	if err := parseJSONResponse(rr, &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response.Total != 2 {
		t.Errorf("expected total 2, got %d", response.Total)
	}

	if len(response.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(response.Items))
	}
}

func TestEntityHandler_List_FilterByKind(t *testing.T) {
	button := newTestEntity("MQTTButton", "default", "test-button", true)
	sensor := newTestEntity("MQTTSensor", "default", "test-sensor", false)

	handler := newTestEntityHandler(button, sensor)

	rr := executeRequest(handler.List, http.MethodGet, "/api/v1/entities?kind=MQTTButton", nil, nil)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var response EntityListResponse
	if err := parseJSONResponse(rr, &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response.Total != 1 {
		t.Errorf("expected total 1, got %d", response.Total)
	}

	if len(response.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(response.Items))
	}

	if response.Items[0].Kind != "MQTTButton" {
		t.Errorf("expected kind MQTTButton, got %s", response.Items[0].Kind)
	}
}

func TestEntityHandler_List_FilterByNamespace(t *testing.T) {
	button1 := newTestEntity("MQTTButton", "default", "button1", true)
	button2 := newTestEntity("MQTTButton", "production", "button2", true)

	handler := newTestEntityHandler(button1, button2)

	rr := executeRequest(handler.List, http.MethodGet, "/api/v1/entities?namespace=default", nil, nil)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var response EntityListResponse
	if err := parseJSONResponse(rr, &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response.Total != 1 {
		t.Errorf("expected total 1, got %d", response.Total)
	}

	if response.Items[0].Namespace != "default" {
		t.Errorf("expected namespace default, got %s", response.Items[0].Namespace)
	}
}

func TestEntityHandler_List_EmptyResult(t *testing.T) {
	handler := newTestEntityHandler()

	rr := executeRequest(handler.List, http.MethodGet, "/api/v1/entities", nil, nil)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var response EntityListResponse
	if err := parseJSONResponse(rr, &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response.Total != 0 {
		t.Errorf("expected total 0, got %d", response.Total)
	}

	if response.Items != nil && len(response.Items) != 0 {
		t.Errorf("expected empty or nil items, got %d items", len(response.Items))
	}
}

func TestEntityHandler_Get_Success(t *testing.T) {
	button := newTestEntity("MQTTButton", "default", "test-button", true)

	handler := newTestEntityHandler(button)

	rr := executeRequest(handler.Get, http.MethodGet, "/api/v1/entities/MQTTButton/default/test-button", nil, map[string]string{
		"kind":      "MQTTButton",
		"namespace": "default",
		"name":      "test-button",
	})

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var response map[string]interface{}
	if err := parseJSONResponse(rr, &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response["kind"] != "MQTTButton" {
		t.Errorf("expected kind MQTTButton, got %v", response["kind"])
	}

	metadata := response["metadata"].(map[string]interface{})
	if metadata["name"] != "test-button" {
		t.Errorf("expected name test-button, got %v", metadata["name"])
	}
}

func TestEntityHandler_Get_UnknownKind(t *testing.T) {
	handler := newTestEntityHandler()

	rr := executeRequest(handler.Get, http.MethodGet, "/api/v1/entities/UnknownKind/default/test", nil, map[string]string{
		"kind":      "UnknownKind",
		"namespace": "default",
		"name":      "test",
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

func TestEntityHandler_Get_NotFound(t *testing.T) {
	handler := newTestEntityHandler()

	rr := executeRequest(handler.Get, http.MethodGet, "/api/v1/entities/MQTTButton/default/nonexistent", nil, map[string]string{
		"kind":      "MQTTButton",
		"namespace": "default",
		"name":      "nonexistent",
	})

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestEntityHandler_Create_Success(t *testing.T) {
	handler := newTestEntityHandler()

	body := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name": "new-button",
		},
		"spec": map[string]interface{}{
			"name":         "My Button",
			"commandTopic": "homeassistant/button/test/command",
		},
	}

	rr := executeRequest(handler.Create, http.MethodPost, "/api/v1/entities/MQTTButton/default", body, map[string]string{
		"kind":      "MQTTButton",
		"namespace": "default",
	})

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	var response map[string]interface{}
	if err := parseJSONResponse(rr, &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response["kind"] != "MQTTButton" {
		t.Errorf("expected kind MQTTButton, got %v", response["kind"])
	}
}

func TestEntityHandler_Create_InvalidJSON(t *testing.T) {
	handler := newTestEntityHandler()

	req := executeRequest(handler.Create, http.MethodPost, "/api/v1/entities/MQTTButton/default", nil, map[string]string{
		"kind":      "MQTTButton",
		"namespace": "default",
	})

	// With nil body, the request should fail to parse JSON
	// Note: Actually with nil body, json.Decode will fail
	// Let's test with invalid JSON string instead
	if req.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, req.Code)
	}
}

func TestEntityHandler_Create_UnknownKind(t *testing.T) {
	handler := newTestEntityHandler()

	body := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name": "test",
		},
		"spec": map[string]interface{}{},
	}

	rr := executeRequest(handler.Create, http.MethodPost, "/api/v1/entities/UnknownKind/default", body, map[string]string{
		"kind":      "UnknownKind",
		"namespace": "default",
	})

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestEntityHandler_Update_Success(t *testing.T) {
	button := newTestEntity("MQTTButton", "default", "test-button", true)

	handler := newTestEntityHandler(button)

	body := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name": "test-button",
		},
		"spec": map[string]interface{}{
			"name":         "Updated Button Name",
			"commandTopic": "homeassistant/button/test/command",
		},
	}

	rr := executeRequest(handler.Update, http.MethodPut, "/api/v1/entities/MQTTButton/default/test-button", body, map[string]string{
		"kind":      "MQTTButton",
		"namespace": "default",
		"name":      "test-button",
	})

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var response map[string]interface{}
	if err := parseJSONResponse(rr, &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	spec := response["spec"].(map[string]interface{})
	if spec["name"] != "Updated Button Name" {
		t.Errorf("expected updated name, got %v", spec["name"])
	}
}

func TestEntityHandler_Update_NotFound(t *testing.T) {
	handler := newTestEntityHandler()

	body := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name": "nonexistent",
		},
		"spec": map[string]interface{}{},
	}

	rr := executeRequest(handler.Update, http.MethodPut, "/api/v1/entities/MQTTButton/default/nonexistent", body, map[string]string{
		"kind":      "MQTTButton",
		"namespace": "default",
		"name":      "nonexistent",
	})

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestEntityHandler_Update_InvalidJSON(t *testing.T) {
	button := newTestEntity("MQTTButton", "default", "test-button", true)

	handler := newTestEntityHandler(button)

	rr := executeRequest(handler.Update, http.MethodPut, "/api/v1/entities/MQTTButton/default/test-button", nil, map[string]string{
		"kind":      "MQTTButton",
		"namespace": "default",
		"name":      "test-button",
	})

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestEntityHandler_Delete_Success(t *testing.T) {
	button := newTestEntity("MQTTButton", "default", "test-button", true)

	handler := newTestEntityHandler(button)

	rr := executeRequest(handler.Delete, http.MethodDelete, "/api/v1/entities/MQTTButton/default/test-button", nil, map[string]string{
		"kind":      "MQTTButton",
		"namespace": "default",
		"name":      "test-button",
	})

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d: %s", http.StatusNoContent, rr.Code, rr.Body.String())
	}
}

func TestEntityHandler_Delete_UnknownKind(t *testing.T) {
	handler := newTestEntityHandler()

	rr := executeRequest(handler.Delete, http.MethodDelete, "/api/v1/entities/UnknownKind/default/test", nil, map[string]string{
		"kind":      "UnknownKind",
		"namespace": "default",
		"name":      "test",
	})

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestEntityHandler_ToSummary_Published(t *testing.T) {
	button := newTestEntity("MQTTButton", "default", "test-button", true)

	handler := newTestEntityHandler()
	summary := handler.toSummary(button)

	if !summary.Published {
		t.Error("expected published to be true")
	}

	if summary.Kind != "MQTTButton" {
		t.Errorf("expected kind MQTTButton, got %s", summary.Kind)
	}

	if summary.Name != "test-button" {
		t.Errorf("expected name test-button, got %s", summary.Name)
	}

	if summary.DisplayName != "Display Name for test-button" {
		t.Errorf("expected display name, got %s", summary.DisplayName)
	}
}

func TestEntityHandler_ToSummary_Unpublished(t *testing.T) {
	sensor := newTestEntity("MQTTSensor", "default", "test-sensor", false)

	handler := newTestEntityHandler()
	summary := handler.toSummary(sensor)

	if summary.Published {
		t.Error("expected published to be false")
	}
}
