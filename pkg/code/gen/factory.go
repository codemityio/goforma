package gen

import "github.com/codemityio/goforma/pkg/code/parser"

// NewDepGraphGenerator factory function.
func NewDepGraphGenerator() *DefaultDepGraphGenerator {
	dcp := &DefaultDepGraphGenerator{}

	return dcp
}

// NewDefaultUMLGraphGenerator factory function.
func NewDefaultUMLGraphGenerator(options ...UMLGraphGeneratorOption) *DefaultUMLGraphGenerator {
	dumlgg := &DefaultUMLGraphGenerator{
		codeMaps: map[string]*parser.CodeMap[*parser.Var, *parser.Type, *parser.Func, *parser.Const]{},
		config: &UMLGraphGeneratorConfig{
			Legend:      false,
			Primitive:   false,
			Var:         false,
			Const:       false,
			Func:        false,
			NotExported: false,
			Doc:         false,
		},
		types:      map[string]struct{}{},
		typesCache: map[string]struct{}{},
		linksCache: map[string]struct{}{},
		primitiveTypes: map[string]bool{
			"bool":       true,
			"int":        true,
			"int8":       true,
			"int16":      true,
			"int32":      true,
			"int64":      true,
			"uint":       true,
			"uint8":      true,
			"uint16":     true,
			"uint32":     true,
			"uint64":     true,
			"uintptr":    true,
			"float32":    true,
			"float64":    true,
			"complex64":  true,
			"complex128": true,
			"string":     true,
			"byte":       true, // alias for uint8
			"rune":       true, // alias for int32
		},
	}

	for _, option := range options {
		option(dumlgg)
	}

	return dumlgg
}
