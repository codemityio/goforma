package server

import (
	"errors"
	"fmt"
)

var (
	// ErrPkg a common package error.
	ErrPkg = errors.New("huma")
	// ErrConfigUnableToLoad error.
	ErrConfigUnableToLoad = fmt.Errorf("%w: unable to load config", ErrPkg)
)
