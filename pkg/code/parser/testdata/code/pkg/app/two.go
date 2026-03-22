package app

type Two struct {
}

func (o Two) FunctionOne(b bool, i int) bool {
	return true
}

func (o Two) FunctionTwo(a, b, c Input, i int) bool {
	return true
}

func (o Two) FunctionThree(a, b, c CustomString, i int) bool {
	return true
}

func (o Two) FunctionWithReturnList() (Output, *Output, error) {
	return Output{}, &Output{}, nil
}

func (o Two) FunctionWithReturnListParametrized() (a, b, c Output, err error) {
	return
}
