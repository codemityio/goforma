package app

type Three struct {
	Input  Input
	output Output
}

func (o Three) FunctionOne(b bool, i int) bool {
	return true
}

func (o Three) functionTwo(a, b, c Input, i int) bool {
	return true
}

func (o Three) FunctionThree(a, b, c CustomString, i int) bool {
	return true
}

func (o Three) FunctionWithReturnList() (Output, *Output, error) {
	return Output{}, &Output{}, nil
}

func (o Three) functionWithReturnListParametrized() (a, b, c Output, err error) {
	return
}
