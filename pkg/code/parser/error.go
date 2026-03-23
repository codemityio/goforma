package parser

import (
	"errors"
	"fmt"
)

var (
	// ErrPkg a common package error.
	ErrPkg = errors.New("parser")
	// ErrPkgLoad error.
	ErrPkgLoad = fmt.Errorf("%w: unable to load packages", ErrPkg)
	// ErrGetRelPath error.
	ErrGetRelPath = fmt.Errorf("%w: unable to get relative path", ErrPkg)
)
