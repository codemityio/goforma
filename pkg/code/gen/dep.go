package gen

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/codemityio/goforma/pkg/code/imports"
)

//go:embed dep.dot.tpl
var depDotTpl string

// DefaultDepGraphGenerator default generator implementation.
type DefaultDepGraphGenerator struct{}

// Generate generate dependency graph.
func (g *DefaultDepGraphGenerator) Generate(packages *imports.Packages) (string, error) {
	// create a new template and parse the letter into it
	tmpl, err := template.New("graphviz").Funcs(template.FuncMap{
		"fillColour":   g.fillColour,
		"strokeColour": g.strokeColour,
	}).Parse(depDotTpl)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrTemplateParse, err)
	}

	var buf bytes.Buffer

	// execute the template, passing in the data structure
	err = tmpl.Execute(&buf, packages)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrTemplateExecute, err)
	}

	return buf.String(), nil
}

// Helper function to decide what colour of the box to return.
func (g *DefaultDepGraphGenerator) fillColour(pkg *imports.Package) string {
	if pkg.IsLocal {
		return "gray99"
	}

	if pkg.IsOwned {
		return "gray99"
	}

	if pkg.IsExternal {
		return "gray83"
	}

	switch {
	case pkg.IsStandard:
		return "gray100"
	case pkg.IsInternal:
		return "grey44"
	case pkg.IsVendor:
		return "grey55"
	}

	return "invis"
}

// Helper function to decide what colour of the box to return.
func (g *DefaultDepGraphGenerator) strokeColour(pkg *imports.Package) string {
	if pkg.IsLocal {
		return "gray58"
	}

	if pkg.IsOwned {
		return "gray69"
	}

	if pkg.IsExternal {
		return "gray65"
	}

	switch {
	case pkg.IsStandard:
		return "gray80"
	case pkg.IsInternal:
		return "grey25"
	case pkg.IsVendor:
		return "grey36"
	}

	return "invis"
}
