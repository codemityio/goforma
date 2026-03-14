package app

import (
	"context"

	"code/pkg/integration"
	intg "code/pkg/integration"
)

type AnyFunction[I, O any] func(context.Context, *I) (*O, error)

type CustomFunction[I Input, O Output] func(context.Context, I) (O, error)

type CustomStruct[I Input, O Output] struct{}

type CustomFunctionGenericsSlices[I []*Input, O []Output] func(context.Context, I) (O, error)

type CustomFunctionGenericsSlicesExternalTypes[I []*integration.ExternalInput, O []intg.ExternalOutput] func(context.Context, I) (O, error)

type CustomFunctionGenericsMaps[I []*Input, O map[string]Output] func(context.Context, I) (O, error)

func Swap[T any, U any](first T, second U) (U, T) {
	return second, first
}
