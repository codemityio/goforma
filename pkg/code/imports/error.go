package imports

import (
	"errors"
	"fmt"
)

var (
	// ErrPkg a common package error.
	ErrPkg = errors.New("imports")
	// ErrPkgLoad error.
	ErrPkgLoad = fmt.Errorf("%w: unable to load packages", ErrPkg)
)
