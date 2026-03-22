package server

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

// Option function.
type Option func(p *Handler)

// OperationWithHandler an operation type.
type OperationWithHandler[I, O any] struct {
	huma.Operation
	Handler func(context.Context, *I) (*O, error)
}

func (op *OperationWithHandler[I, O]) Register(api huma.API) {
	huma.Register(api, op.Operation, op.Handler)
}
