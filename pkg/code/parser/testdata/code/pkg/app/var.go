package app

import (
	"context"

	"code/pkg/config"
	cfg "code/pkg/config"

	"github.com/urfave/cli/v2"
)

var (
	VarFuncGenerics AnyFunction[int, string] = func(ctx context.Context, input *int) (*string, error) {
		return nil, nil
	}
	VarFuncGenericsCustomTypes AnyFunction[CustomFloat64, CustomComplex128] = func(ctx context.Context, input *CustomFloat64) (*CustomComplex128, error) {
		return nil, nil
	}
	/*
		varNotExported is a value with extensive document.
	*/
	varNotExported string
	IntStringSwap  = func(first int, second string) (string, int) {
		return Swap(first, second)
	}
	StringValueWithoutExplicitType           = "value"
	BooleanValueWithoutExplicitType          = true
	CustomFloat64Value                       CustomFloat64
	CustomComplex128Value                    CustomComplex128 = 1 + 2i
	CustomComplex128ValueWithoutExplicitType                  = 1 + 2i
	Config                                   config.JWT
	ConfigWithAlias                          cfg.JWT
	App                                      cli.Command = cli.Command{}
	AppUntyped                                           = cli.Command{}
)
