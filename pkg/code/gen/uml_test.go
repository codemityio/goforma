package gen

import (
	_ "embed"
	"encoding/json"
	"os"
	"testing"

	"github.com/codemityio/goforma/pkg/code/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/uml-pkg-app.json
var umlPkgAppJSON []byte

//go:embed testdata/uml-pkg-app.puml
var umlPkgAppPUML string

//go:embed testdata/uml-pkg-app-without-config.puml
var umlPkgAppWithoutConfigPUML string

//go:embed testdata/uml-pkg-edge.json
var umlPkgEdgeJSON []byte

//go:embed testdata/uml-pkg-edge.puml
var umlPkgEdgePUML string

//go:embed testdata/uml-pkg-recursive.json
var umlPkgRecursiveJSON []byte

//go:embed testdata/uml-pkg-recursive.puml
var umlPkgRecursivePUML string

//go:embed testdata/uml-pkg-recursive-skip-const.puml
var umlPkgRecursiveSkipConstPUML string

//go:embed testdata/uml-pkg-recursive-skip-doc.puml
var umlPkgRecursiveSkipDocPUML string

//go:embed testdata/uml-pkg-recursive-skip-func.puml
var umlPkgRecursiveSkipFuncPUML string

//go:embed testdata/uml-pkg-recursive-skip-notexported.puml
var umlPkgRecursiveSkipNotExportedPUML string

//go:embed testdata/uml-pkg-recursive-skip-primitive.puml
var umlPkgRecursiveSkipPrimitivePUML string

//go:embed testdata/uml-pkg-recursive-skip-var.puml
var umlPkgRecursiveSkipVarPUML string

func TestDefaultUMLGraphGenerator_Generate(t *testing.T) {
	tests := []struct {
		name      string
		config    *UMLGraphGeneratorConfig
		input     []byte
		result    string
		writePath string
	}{
		{
			name: "success-app",
			config: &UMLGraphGeneratorConfig{
				Legend:      true,
				Primitive:   true,
				Var:         true,
				Const:       true,
				Func:        true,
				NotExported: true,
				Doc:         true,
			},
			input:     umlPkgAppJSON,
			result:    umlPkgAppPUML,
			writePath: "testdata/uml-pkg-app.puml",
		},
		{
			name:      "success-app-without-config",
			config:    nil,
			input:     umlPkgAppJSON,
			result:    umlPkgAppWithoutConfigPUML,
			writePath: "testdata/uml-pkg-app-without-config.puml",
		},
		{
			name: "success-edge",
			config: &UMLGraphGeneratorConfig{
				Legend:      true,
				Primitive:   true,
				Var:         true,
				Const:       true,
				Func:        true,
				NotExported: true,
				Doc:         true,
			},
			input:     umlPkgEdgeJSON,
			result:    umlPkgEdgePUML,
			writePath: "testdata/uml-pkg-edge.puml",
		},
		{
			name: "success-recursive",
			config: &UMLGraphGeneratorConfig{
				Legend:      true,
				Primitive:   true,
				Var:         true,
				Const:       true,
				Func:        true,
				NotExported: true,
				Doc:         true,
			},
			input:     umlPkgRecursiveJSON,
			result:    umlPkgRecursivePUML,
			writePath: "testdata/uml-pkg-recursive.puml",
		},
		{
			name: "success-recursive-skip-const",
			config: &UMLGraphGeneratorConfig{
				Legend:      true,
				Primitive:   true,
				Var:         true,
				Const:       false,
				Func:        true,
				NotExported: true,
				Doc:         true,
			},
			input:     umlPkgRecursiveJSON,
			result:    umlPkgRecursiveSkipConstPUML,
			writePath: "testdata/uml-pkg-recursive-skip-const.puml",
		},
		{
			name: "success-recursive-skip-doc",
			config: &UMLGraphGeneratorConfig{
				Legend:      true,
				Primitive:   true,
				Var:         true,
				Const:       true,
				Func:        true,
				NotExported: true,
				Doc:         false,
			},
			input:     umlPkgRecursiveJSON,
			result:    umlPkgRecursiveSkipDocPUML,
			writePath: "testdata/uml-pkg-recursive-skip-doc.puml",
		},
		{
			name: "success-recursive-skip-func",
			config: &UMLGraphGeneratorConfig{
				Legend:      true,
				Primitive:   true,
				Var:         true,
				Const:       true,
				Func:        false,
				NotExported: true,
				Doc:         true,
			},
			input:     umlPkgRecursiveJSON,
			result:    umlPkgRecursiveSkipFuncPUML,
			writePath: "testdata/uml-pkg-recursive-skip-func.puml",
		},
		{
			name: "success-recursive-skip-notexported",
			config: &UMLGraphGeneratorConfig{
				Legend:      true,
				Primitive:   true,
				Var:         true,
				Const:       true,
				Func:        true,
				NotExported: false,
				Doc:         true,
			},
			input:     umlPkgRecursiveJSON,
			result:    umlPkgRecursiveSkipNotExportedPUML,
			writePath: "testdata/uml-pkg-recursive-skip-notexported.puml",
		},
		{
			name: "success-recursive-skip-primitive",
			config: &UMLGraphGeneratorConfig{
				Legend:      true,
				Primitive:   false,
				Var:         true,
				Const:       true,
				Func:        true,
				NotExported: true,
				Doc:         true,
			},
			input:     umlPkgRecursiveJSON,
			result:    umlPkgRecursiveSkipPrimitivePUML,
			writePath: "testdata/uml-pkg-recursive-skip-primitive.puml",
		},
		{
			name: "success-recursive-skip-var",
			config: &UMLGraphGeneratorConfig{
				Legend:      true,
				Primitive:   true,
				Var:         false,
				Const:       true,
				Func:        true,
				NotExported: true,
				Doc:         true,
			},
			input:     umlPkgRecursiveJSON,
			result:    umlPkgRecursiveSkipVarPUML,
			writePath: "testdata/uml-pkg-recursive-skip-var.puml",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var input parser.CodeMap[*parser.Var, *parser.Type, *parser.Func, *parser.Const]

			require.NoError(t, json.Unmarshal(test.input, &input))

			generator := NewDefaultUMLGraphGenerator(
				WithUMLGraphGeneratorConfig(test.config),
				WithUMLGraphGeneratorCodeMap(&input),
			)

			result, err := generator.Generate()
			require.NoError(t, err)

			assert.Equal(t, test.result, result)

			err = os.WriteFile(test.writePath, []byte(result), 0o644)
			require.NoError(t, err)
		})
	}
}
