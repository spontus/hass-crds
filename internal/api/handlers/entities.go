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
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

const (
	apiGroup   = "mqtt.home-assistant.io"
	apiVersion = "v1alpha1"
)

type EntityHandler struct {
	client     dynamic.Interface
	restConfig *rest.Config
	log        logr.Logger
}

func NewEntityHandler(client dynamic.Interface, restConfig *rest.Config, log logr.Logger) *EntityHandler {
	return &EntityHandler{
		client:     client,
		restConfig: restConfig,
		log:        log.WithName("entities"),
	}
}

type EntitySummary struct {
	Kind        string            `json:"kind"`
	APIVersion  string            `json:"apiVersion"`
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	DisplayName string            `json:"displayName,omitempty"`
	Published   bool              `json:"published"`
	CreatedAt   string            `json:"createdAt"`
	Labels      map[string]string `json:"labels,omitempty"`
}

type EntityListResponse struct {
	Items []EntitySummary `json:"items"`
	Total int             `json:"total"`
}

func (h *EntityHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	kindFilter := r.URL.Query().Get("kind")
	namespaceFilter := r.URL.Query().Get("namespace")

	var results []EntitySummary

	entityTypes := GetEntityTypes()
	for _, et := range entityTypes {
		if kindFilter != "" && et.Kind != kindFilter {
			continue
		}

		gvr := schema.GroupVersionResource{
			Group:    apiGroup,
			Version:  apiVersion,
			Resource: et.Plural,
		}

		var list *unstructured.UnstructuredList
		var err error

		if namespaceFilter != "" {
			list, err = h.client.Resource(gvr).Namespace(namespaceFilter).List(ctx, metav1.ListOptions{})
		} else {
			list, err = h.client.Resource(gvr).List(ctx, metav1.ListOptions{})
		}

		if err != nil {
			h.log.Error(err, "failed to list entities", "kind", et.Kind)
			continue
		}

		for _, item := range list.Items {
			summary := h.toSummary(&item)
			results = append(results, summary)
		}
	}

	writeJSON(w, http.StatusOK, EntityListResponse{
		Items: results,
		Total: len(results),
	})
}

func (h *EntityHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	kind := chi.URLParam(r, "kind")
	namespace := chi.URLParam(r, "namespace")
	name := chi.URLParam(r, "name")

	et := GetEntityTypeByKind(kind)
	if et == nil {
		writeError(w, http.StatusNotFound, "unknown entity type: "+kind)
		return
	}

	gvr := schema.GroupVersionResource{
		Group:    apiGroup,
		Version:  apiVersion,
		Resource: et.Plural,
	}

	obj, err := h.client.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		h.log.Error(err, "failed to get entity", "kind", kind, "namespace", namespace, "name", name)
		writeError(w, http.StatusNotFound, "entity not found")
		return
	}

	writeJSON(w, http.StatusOK, obj.Object)
}

func (h *EntityHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	kind := chi.URLParam(r, "kind")
	namespace := chi.URLParam(r, "namespace")

	et := GetEntityTypeByKind(kind)
	if et == nil {
		writeError(w, http.StatusNotFound, "unknown entity type: "+kind)
		return
	}

	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	obj := &unstructured.Unstructured{Object: body}
	obj.SetAPIVersion(apiGroup + "/" + apiVersion)
	obj.SetKind(kind)
	obj.SetNamespace(namespace)

	gvr := schema.GroupVersionResource{
		Group:    apiGroup,
		Version:  apiVersion,
		Resource: et.Plural,
	}

	created, err := h.client.Resource(gvr).Namespace(namespace).Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		h.log.Error(err, "failed to create entity", "kind", kind, "namespace", namespace)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, created.Object)
}

func (h *EntityHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	kind := chi.URLParam(r, "kind")
	namespace := chi.URLParam(r, "namespace")
	name := chi.URLParam(r, "name")

	et := GetEntityTypeByKind(kind)
	if et == nil {
		writeError(w, http.StatusNotFound, "unknown entity type: "+kind)
		return
	}

	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	gvr := schema.GroupVersionResource{
		Group:    apiGroup,
		Version:  apiVersion,
		Resource: et.Plural,
	}

	existing, err := h.client.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		writeError(w, http.StatusNotFound, "entity not found")
		return
	}

	obj := &unstructured.Unstructured{Object: body}
	obj.SetAPIVersion(apiGroup + "/" + apiVersion)
	obj.SetKind(kind)
	obj.SetNamespace(namespace)
	obj.SetName(name)
	obj.SetResourceVersion(existing.GetResourceVersion())

	updated, err := h.client.Resource(gvr).Namespace(namespace).Update(ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		h.log.Error(err, "failed to update entity", "kind", kind, "namespace", namespace, "name", name)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, updated.Object)
}

func (h *EntityHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	kind := chi.URLParam(r, "kind")
	namespace := chi.URLParam(r, "namespace")
	name := chi.URLParam(r, "name")

	et := GetEntityTypeByKind(kind)
	if et == nil {
		writeError(w, http.StatusNotFound, "unknown entity type: "+kind)
		return
	}

	gvr := schema.GroupVersionResource{
		Group:    apiGroup,
		Version:  apiVersion,
		Resource: et.Plural,
	}

	err := h.client.Resource(gvr).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		h.log.Error(err, "failed to delete entity", "kind", kind, "namespace", namespace, "name", name)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *EntityHandler) toSummary(obj *unstructured.Unstructured) EntitySummary {
	spec, _, _ := unstructured.NestedMap(obj.Object, "spec")
	displayName, _, _ := unstructured.NestedString(spec, "name")

	published := false
	conditions, found, _ := unstructured.NestedSlice(obj.Object, "status", "conditions")
	if found {
		for _, c := range conditions {
			cond, ok := c.(map[string]interface{})
			if !ok {
				continue
			}
			if cond["type"] == "Published" && cond["status"] == "True" {
				published = true
				break
			}
		}
	}

	return EntitySummary{
		Kind:        obj.GetKind(),
		APIVersion:  obj.GetAPIVersion(),
		Name:        obj.GetName(),
		Namespace:   obj.GetNamespace(),
		DisplayName: displayName,
		Published:   published,
		CreatedAt:   obj.GetCreationTimestamp().Format("2006-01-02T15:04:05Z"),
		Labels:      obj.GetLabels(),
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
