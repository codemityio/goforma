package app

import (
	"code/pkg/integration"
	intg "code/pkg/integration"
)

type Five[I Input, O Output, EI integration.ExternalInput, EO integration.ExternalOutput] struct {
	Input        Input
	output       Output
	CustomMethod CustomFunction[I, O]
	// ExternalCustomMethod an example field code doc block.
	//
	// With a new line separation between line of text.
	ExternalCustomMethod intg.ExternalCustomFunction[EI, EO]
}

func (o Five[I, O, EI, EO]) FunctionOne(b bool, i int) bool {
	return true
}
