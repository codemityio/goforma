package edge

import (
	"code/pkg/app"
	cfg "code/pkg/config"
	"code/pkg/server"
	"context"
)

type TypeOfCustomType app.CustomStringWithStructReceivingMethods

type TypeOfExternalCustomType cfg.Transport

type CompositionStructWithUsageType struct {
	TypeOfExternalCustomType
	cfg.TLS
	Name string
}

func (c *CompositionStructWithUsageType) String() string {
	return ""
}

func (c *CompositionStructWithUsageType) LoadByte(context.Context, []byte) error {
	return nil
}

func (c *CompositionStructWithUsageType) Function(_, _, _ app.CustomString, d []byte) (e, f, g app.CustomString, h []byte) {
	return "", "", "", nil
}

func (c *CompositionStructWithUsageType) AnotherAdditionalFunction(_, _, _ app.CustomString, d []byte) (e, f, g app.CustomString, h []byte) {
	return "", "", "", nil
}

// CompositionStructWithUsageTypeDefinition type Definition:
//
// - Creates a new type.
// - Requires explicit conversion to assign or use values interchangeably with the base type.
// - Methods must be explicitly defined for the new type.
type CompositionStructWithUsageTypeDefinition CompositionStructWithUsageType

func (c *CompositionStructWithUsageTypeDefinition) Method() {}

// CompositionStructWithUsageTypeAlias type Alias:
//
// - No explicit conversion is needed; they are treated as the same type.
// - All methods of the original type are automatically available.
type CompositionStructWithUsageTypeAlias = CompositionStructWithUsageType

type TransportAlias = cfg.Transport

type HandlerAlias = server.Handler

type HandlerDefinition server.Handler

type CustomStructAlias = app.CustomStruct[app.Input, app.Output]

type CustomStructDefinition app.CustomStruct[app.Input, app.Output]

type CustomFunctionWithMethodDefinition app.CustomFunctionWithMethod

type CustomFunctionWithMethodAlias = app.CustomFunctionWithMethod

type StringAlias = string
