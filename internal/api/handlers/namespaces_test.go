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

	corev1 "k8s.io/api/core/v1"
)

func TestNamespaceHandler_List_Success(t *testing.T) {
	ns1 := newTestNamespace("default", corev1.NamespaceActive)
	ns2 := newTestNamespace("production", corev1.NamespaceActive)
	ns3 := newTestNamespace("staging", corev1.NamespaceTerminating)

	handler := newTestNamespaceHandler(ns1, ns2, ns3)

	rr := executeRequest(handler.List, http.MethodGet, "/api/v1/namespaces", nil, nil)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var response struct {
		Namespaces []NamespaceSummary `json:"namespaces"`
		Total      int                `json:"total"`
	}
	if err := parseJSONResponse(rr, &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response.Total != 3 {
		t.Errorf("expected total 3, got %d", response.Total)
	}

	if len(response.Namespaces) != 3 {
		t.Errorf("expected 3 namespaces, got %d", len(response.Namespaces))
	}

	nsMap := make(map[string]NamespaceSummary)
	for _, ns := range response.Namespaces {
		nsMap[ns.Name] = ns
	}

	if ns, ok := nsMap["default"]; !ok {
		t.Error("expected default namespace")
	} else if ns.Status != "Active" {
		t.Errorf("expected default namespace status Active, got %s", ns.Status)
	}

	if ns, ok := nsMap["staging"]; !ok {
		t.Error("expected staging namespace")
	} else if ns.Status != "Terminating" {
		t.Errorf("expected staging namespace status Terminating, got %s", ns.Status)
	}
}

func TestNamespaceHandler_List_Empty(t *testing.T) {
	handler := newTestNamespaceHandler()

	rr := executeRequest(handler.List, http.MethodGet, "/api/v1/namespaces", nil, nil)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var response struct {
		Namespaces []NamespaceSummary `json:"namespaces"`
		Total      int                `json:"total"`
	}
	if err := parseJSONResponse(rr, &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response.Total != 0 {
		t.Errorf("expected total 0, got %d", response.Total)
	}

	if response.Namespaces == nil {
		t.Error("expected namespaces to be empty array, not nil")
	}

	if len(response.Namespaces) != 0 {
		t.Errorf("expected 0 namespaces, got %d", len(response.Namespaces))
	}
}

func TestNamespaceHandler_List_WithLabels(t *testing.T) {
	ns := newTestNamespace("labeled-ns", corev1.NamespaceActive)
	ns.Labels = map[string]string{
		"environment": "production",
		"team":        "platform",
	}

	handler := newTestNamespaceHandler(ns)

	rr := executeRequest(handler.List, http.MethodGet, "/api/v1/namespaces", nil, nil)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var response struct {
		Namespaces []NamespaceSummary `json:"namespaces"`
		Total      int                `json:"total"`
	}
	if err := parseJSONResponse(rr, &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(response.Namespaces) != 1 {
		t.Fatalf("expected 1 namespace, got %d", len(response.Namespaces))
	}

	nsResult := response.Namespaces[0]
	if nsResult.Labels == nil {
		t.Fatal("expected labels to be present")
	}

	if nsResult.Labels["environment"] != "production" {
		t.Errorf("expected label environment=production, got %s", nsResult.Labels["environment"])
	}

	if nsResult.Labels["team"] != "platform" {
		t.Errorf("expected label team=platform, got %s", nsResult.Labels["team"])
	}
}
