package integration

import (
	"github.com/danielgtaylor/huma/v2"
)

// OperationWithHandlerRegisterer an operation with handler registerer.
type OperationWithHandlerRegisterer interface {
	Register(api huma.API)
}
