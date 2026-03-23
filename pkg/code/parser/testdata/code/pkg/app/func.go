package app

import (
	"fmt"

	"code/pkg/integration"
	intg "code/pkg/integration"
)

func nonExportedFunction(input integration.ExternalInput) (intg.ExternalOutput, error) {
	var err error

	const error string = "error"

	if input == "" {
		err = fmt.Errorf(error)
	}

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func nonExportedFunctionWithNamedReturn(
	input integration.ExternalInput,
) (output integration.ExternalOutput, err error) {
	return nil, nil
}

func ExportedFunctionWithMultipleSameTypeArgumentsAndNamedReturn(
	one, two, three, four integration.ExternalInput,
) (output intg.ExternalOutput, err error) {
	return nil, nil
}

func ExportedFunctionWithMultipleSameTypeArguments(
	one integration.ExternalInput,
	two intg.ExternalInput,
) (output integration.ExternalOutput, err error) {
	return nil, nil
}

func ExportedVariadicFunctionWithMultipleSameTypeArgumentsAndNamedReturn(
	input ...integration.ExternalInput,
) (output integration.ExternalOutput, err error) {
	return nil, nil
}

func ExportedVariadicFunctionWithMultipleSameTypeArgumentsAndNamedReturnWithGenerics[I intg.ExternalInput, O integration.ExternalOutput](
	input ...I,
) (output O, err error) {
	return output, nil
}

func doNothingFunction() {
	return
}
