package app

type One struct {
}

func (o *One) FunctionOne(b bool, i int) bool {
	return true
}

func (o *One) FunctionTwo(a, b, c bool, i int) bool {
	return true
}

func (o *One) FunctionThree(a, b, c CustomString, i int) bool {
	return true
}

func (o *One) FunctionWithReturnListParametrized() (a, b, c []byte, err error) {
	return
}
