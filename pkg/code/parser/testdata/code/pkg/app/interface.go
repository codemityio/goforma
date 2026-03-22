package app

import (
	"code/pkg/integration"
	intg "code/pkg/integration"
)

type EmbeddableInterface interface {
	String() string
}

type CompositionInterfaceWithoutUsage interface {
	EmbeddableInterface
	integration.BytesLoader
	[]intg.EnvLoader
	AdditionalFunction(a, b, c CustomString, d []byte) (e, f, g CustomString, h []byte)
	// AnotherAdditionalFunction field doc.
	AnotherAdditionalFunction(a, b, c CustomString, d []byte) (e, f, g CustomString, h []byte)
}

type CompositionInterfaceWithUsage interface {
	EmbeddableInterface
	integration.BytesLoader
	Function(a, b, c CustomString, d []byte) (e, f, g CustomString, h []byte)
	// AnotherAdditionalFunction field doc.
	AnotherAdditionalFunction(a, b, c CustomString, d []byte) (e, f, g CustomString, h []byte)
}
