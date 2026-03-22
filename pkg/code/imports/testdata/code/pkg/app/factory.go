package app

import (
	"code/pkg/integration"
	intg "code/pkg/integration"
)

func NewFour[I Input, O Output, EI integration.ExternalInput, EO integration.ExternalOutput](
	input Input, output Output,
	customMethod CustomFunction[I, O], externalCustomMethod integration.ExternalCustomFunction[EI, EO],
) Four[I, O, EI, EO] {
	return Four[I, O, EI, EO]{
		Input:                input,
		output:               output,
		CustomMethod:         customMethod,
		ExternalCustomMethod: externalCustomMethod,
	}
}

func NewFive[I Input, O Output, EI intg.ExternalInput, EO integration.ExternalOutput](
	input Input, output Output,
	customMethod CustomFunction[I, O], externalCustomMethod intg.ExternalCustomFunction[EI, EO],
) Five[I, O, EI, EO] {
	return Five[I, O, EI, EO]{
		Input:                input,
		output:               output,
		CustomMethod:         customMethod,
		ExternalCustomMethod: externalCustomMethod,
	}
}
