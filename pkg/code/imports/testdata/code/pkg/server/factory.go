package server

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

// New a factory function to create a new server.
func New(options ...Option) *Handler {
	handler := &Handler{
		config:     nil,
		mux:        chi.NewMux(),
		operations: nil,
	}

	for _, option := range options {
		option(handler)
	}

	config := huma.DefaultConfig("API", "0.0.0")

	config.OpenAPI.Info.Title = handler.config.Title
	config.OpenAPI.Info.Version = handler.config.Version
	config.OpenAPIPath = handler.config.OpenAPIPath
	config.DocsPath = handler.config.DocsPath
	config.SchemasPath = handler.config.SchemasPath
	config.DefaultFormat = handler.config.DefaultFormat

	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{}

	api := humachi.New(handler.mux, config)

	for _, v := range handler.operations {
		v.Register(api)
	}

	return handler
}
