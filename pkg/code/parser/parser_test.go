package parser

import (
	_ "embed"
	"encoding/json"
	"os"
	"testing"

	"github.com/codemityio/goforma/pkg/code/doc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/pkg.json
	pkg string

	//go:embed testdata/pkg-custom.json
	pkgCustom string

	//go:embed testdata/pkg-app.json
	pkgApp string

	//go:embed testdata/pkg-app-recursive.json
	pkgAppRecursive string

	//go:embed testdata/pkg-recursive.json
	pkgRecursive string

	//go:embed testdata/pkg-edge.json
	pkgEdge string
)

func TestDefaultParser_Parse(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	parser := New(
		WithRootPath(wd+"/testdata/code"),
		WithDocParser(doc.New()),
	)

	tests := []struct {
		name           string
		inputPath      string
		expectedResult string
		writePath      string
	}{
		{
			name:           "pkg",
			inputPath:      "./pkg",
			expectedResult: pkg,
			writePath:      "testdata/pkg.json",
		},
		{
			name:           "pkg-app",
			inputPath:      "./pkg/app",
			expectedResult: pkgApp,
			writePath:      "testdata/pkg-app.json",
		},
		{
			name:           "pkg-custom",
			inputPath:      "./pkg/custom",
			expectedResult: pkgCustom,
			writePath:      "testdata/pkg-custom.json",
		},
		{
			name:           "pkg-app-recursive",
			inputPath:      "./pkg/app/...",
			expectedResult: pkgAppRecursive,
			writePath:      "testdata/pkg-app-recursive.json",
		},
		{
			name:           "pkg-recursive",
			inputPath:      "./pkg/...",
			expectedResult: pkgRecursive,
			writePath:      "testdata/pkg-recursive.json",
		},
		{
			name:           "pkg-edge",
			inputPath:      "./pkg/edge",
			expectedResult: pkgEdge,
			writePath:      "testdata/pkg-edge.json",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			output, err := parser.Parse(test.inputPath)
			require.NoError(t, err)

			result, err := json.MarshalIndent(output, "", "  ")
			require.NoError(t, err)

			assert.JSONEq(t, test.expectedResult, string(result))

			require.NoError(t, os.WriteFile(test.writePath, result, 0o644))
		})
	}
}
