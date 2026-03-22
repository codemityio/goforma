package server

import "code/pkg/integration"

// WithConfig configuration option.
func WithConfig(config *Config) Option {
	return func(handler *Handler) {
		handler.config = config
	}
}

// WithOperationWithHandler configuration option.
func WithOperationWithHandler(operations []integration.OperationWithHandlerRegisterer) Option {
	return func(handler *Handler) {
		handler.operations = operations
	}
}
