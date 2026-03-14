package gen

import (
	"github.com/codemityio/goforma/pkg/code/imports"
)

// DepGraphGenerator generate dependency graph.
type DepGraphGenerator interface {
	Generate(packages *imports.Packages) (string, error)
}

// UMLGraphGenerator generate UML diagram.
type UMLGraphGenerator interface {
	Generate() (string, error)
}
