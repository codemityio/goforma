package app

import (
	"fmt"

	"code/pkg/integration"
)

type Four[I Input, O Output, EI integration.ExternalInput, EO integration.ExternalOutput] struct {
	Input        Input
	output       Output
	CustomMethod CustomFunction[I, O]
	// ExternalCustomMethod an example field code doc block.
	//
	// With a new line separation between line of text.
	ExternalCustomMethod integration.ExternalCustomFunction[EI, EO]
}

func (o Four[I, O, EI, EO]) FunctionOne(b bool, i int) bool {
	var err error

	const error string = "error"

	if !b {
		err = fmt.Errorf(error)
	}

	if err != nil {
		return false
	}

	return true
}
