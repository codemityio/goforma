package gen

import (
	"errors"
	"fmt"
)

var (
	// ErrPkg a common package error.
	ErrPkg = errors.New("gen")
	// ErrTemplateParse error.
	ErrTemplateParse = fmt.Errorf("%w: unable to parse a template", ErrPkg)
	// ErrTemplateExecute error.
	ErrTemplateExecute = fmt.Errorf("%w: unable to execute a template renderer", ErrPkg)
)
