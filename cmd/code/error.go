package code

import (
	"errors"
	"fmt"
)

var (
	errPkg   = errors.New("code")
	errWrite = fmt.Errorf("%w: unable to write", errPkg)
)
