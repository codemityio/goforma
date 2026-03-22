package app

import (
	cfg "code/pkg/config"
)

const (
	ValueStringOne                  = "one"
	ValueStringTwo   string         = "two"
	ValueType        cfg.URL        = "http://localhost"
	AdditionalType   cfg.Additional = "http://localhost"
	ValueBoolTrue                   = true
	ValueBoolFalse   bool           = true
	valueNotExported string         = "not exported"
	// ValueWithDoc is a value with extensive document.
	// e.g. use it as you wish...
	//      ...perhaps with a little bit of indentation.
	//
	// /*
	//  * example block comment
	//  */
	//
	// Well-defined doc block should help with understanding the code.
	ValueWithDoc = "value with doc"
)
