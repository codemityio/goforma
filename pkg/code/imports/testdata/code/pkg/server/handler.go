package server

import (
	"net/http"

	"code/pkg/integration"

	"github.com/go-chi/chi/v5"
)

// Handler abstraction for the router to be used with Huma.
type Handler struct {
	config     *Config
	mux        *chi.Mux
	operations []integration.OperationWithHandlerRegisterer
}

// ServeHTTP an http handler.
func (h *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	h.mux.ServeHTTP(writer, req)
}
