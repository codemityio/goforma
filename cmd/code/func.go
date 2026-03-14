package code

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/codemityio/goforma/pkg/code/doc"
	"github.com/codemityio/goforma/pkg/code/gen"
	"github.com/codemityio/goforma/pkg/code/imports"
	"github.com/codemityio/goforma/pkg/code/parser"
	"github.com/urfave/cli/v2"
)

func dep(ctx *cli.Context) error {
	path := ctx.String("path")
	wd := ctx.String("workdir")
	depth := ctx.Int("depth")
	owned := ctx.StringSlice("owned")
	excludePaths := ctx.StringSlice("exclude-path")
	excludeStandard := ctx.Bool("exclude-standard")
	excludeVendor := ctx.Bool("exclude-vendor")
	excludeInternal := ctx.Bool("exclude-internal")

	parser := imports.New(
		imports.WithRootPath(wd),
		imports.WithDepth(depth),
		imports.WithOwned(owned),
		imports.WithExcludePaths(excludePaths),
		imports.WithExcludeStandard(excludeStandard),
		imports.WithExcludeVendor(excludeVendor),
		imports.WithExcludeInternal(excludeInternal),
	)

	output, err := parser.Parse(path)
	if err != nil {
		return fmt.Errorf("unable to parse `%s` path: %w", path, err)
	}

	res, err := gen.NewDepGraphGenerator().Generate(output)
	if err != nil {
		return fmt.Errorf("unable to generate diagram: %w", err)
	}

	if _, e := fmt.Fprintln(
		ctx.App.Writer,
		res,
	); e != nil {
		return fmt.Errorf("%w: %w", errWrite, e)
	}

	return nil
}

func uml(ctx *cli.Context) error {
	path := ctx.String("path")
	wd := ctx.String("workdir")
	jop := ctx.String("json-output-path")

	var docParser doc.Parser

	if !ctx.Bool("include-doc-comment-slashes") {
		docParser = doc.New()
	}

	prsr := parser.New(
		parser.WithRootPath(wd),
		parser.WithDocParser(docParser),
	)

	output, err := prsr.Parse(path)
	if err != nil {
		return fmt.Errorf("unable to parse `%s` path: %w", path, err)
	}

	generator := gen.NewDefaultUMLGraphGenerator(
		gen.WithUMLGraphGeneratorConfig(&gen.UMLGraphGeneratorConfig{
			Legend:      ctx.Bool("include-legend"),
			Primitive:   ctx.Bool("include-primitive"),
			Var:         ctx.Bool("include-var"),
			Const:       ctx.Bool("include-const"),
			Func:        ctx.Bool("include-func"),
			NotExported: ctx.Bool("include-not-exported"),
			Doc:         ctx.Bool("include-doc"),
		}),
		gen.WithUMLGraphGeneratorCodeMap(output),
	)

	if jop != "" {
		var jo []byte

		jo, err = json.MarshalIndent(output, "", "  ") //nolint:musttag
		if err != nil {
			return fmt.Errorf("unable to serialise code tree: %w", err)
		}

		if e := os.WriteFile(jop, jo, permsWrite); e != nil {
			return fmt.Errorf("unable to write json output: %w", e)
		}
	}

	res, err := generator.Generate()
	if err != nil {
		return fmt.Errorf("unable to generate diagram: %w", err)
	}

	if _, e := fmt.Fprintln(
		ctx.App.Writer,
		res,
	); e != nil {
		return fmt.Errorf("%w: %w", errWrite, e)
	}

	return nil
}
