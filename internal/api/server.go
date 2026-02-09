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

package api

import (
	"context"
	"embed"
	"io/fs"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-logr/logr"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/spontus/hass-crds/internal/api/handlers"
)

//go:embed static/*
var staticFiles embed.FS

type Server struct {
	addr          string
	client        client.Client
	dynamicClient dynamic.Interface
	restConfig    *rest.Config
	log           logr.Logger
	server        *http.Server
}

func NewServer(addr string, client client.Client, restConfig *rest.Config, log logr.Logger) (*Server, error) {
	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return &Server{
		addr:          addr,
		client:        client,
		dynamicClient: dynamicClient,
		restConfig:    restConfig,
		log:           log.WithName("api-server"),
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	entityHandler := handlers.NewEntityHandler(s.dynamicClient, s.restConfig, s.log)
	schemaHandler := handlers.NewSchemaHandler(s.restConfig, s.log)
	namespaceHandler := handlers.NewNamespaceHandler(s.client, s.log)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/entity-types", schemaHandler.ListEntityTypes)
		r.Get("/entity-types/{kind}/schema", schemaHandler.GetSchema)

		r.Get("/namespaces", namespaceHandler.List)

		r.Get("/entities", entityHandler.List)
		r.Get("/entities/{kind}/{namespace}/{name}", entityHandler.Get)
		r.Post("/entities/{kind}/{namespace}", entityHandler.Create)
		r.Put("/entities/{kind}/{namespace}/{name}", entityHandler.Update)
		r.Delete("/entities/{kind}/{namespace}/{name}", entityHandler.Delete)
	})

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return err
	}
	fileServer := http.FileServer(http.FS(staticFS))

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			f, err := staticFS.Open(r.URL.Path[1:])
			if err != nil {
				r.URL.Path = "/"
			} else {
				_ = f.Close()
			}
		}
		fileServer.ServeHTTP(w, r)
	})

	s.server = &http.Server{
		Addr:    s.addr,
		Handler: r,
	}

	s.log.Info("starting API server", "addr", s.addr)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.server.Shutdown(shutdownCtx)
	}()

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
