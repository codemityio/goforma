package app

import (
	"code/pkg/config"
	cfg "code/pkg/config"
)

type CustomFloat32 float32

type CustomFloat64 float64

type CustomComplex128 complex128

type CustomIntSlice []int

type CustomIntArray [3]int

type CustomRune rune

type CustomString string

type Input string

type Output []byte

type PointerString *string

type PointerStruct *cfg.JWT

type PointerSliceOfStruct []*cfg.JWT

type MapOfPointerStruct map[string]*cfg.JWT

type MapOfStruct map[config.URL]cfg.JWT

type CustomStringWithPointerReceivingMethods string

func (c *CustomStringWithPointerReceivingMethods) String() string {
	return ""
}

type CustomStringWithStructReceivingMethods string

func (c CustomStringWithStructReceivingMethods) String() string {
	return ""
}

type CustomFunctionWithMethod func()

func (c *CustomFunctionWithMethod) String() string {
	return ""
}

type CustomFunctionWithMethodWithArgumentsAndReturns func(one string, two PointerStruct) (PointerStruct, error)

func (c *CustomFunctionWithMethodWithArgumentsAndReturns) String() string {
	return ""
}
