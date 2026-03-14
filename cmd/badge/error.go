package badge

import (
	"errors"
	"fmt"
)

var (
	errPkg = errors.New("badge")

	errPathOutsideBase = fmt.Errorf("%w: path outside base directory", errPkg)
	errRelPath         = fmt.Errorf("%w: rel path", errPkg)
)
