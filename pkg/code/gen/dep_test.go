package gen

import (
	_ "embed"
	"encoding/json"
	"os"
	"testing"

	"github.com/codemityio/goforma/pkg/code/imports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/dep-pkg-recursive.json
var depPkgRecursiveJSON []byte

//go:embed testdata/dep-pkg-recursive.dot
var depPkgRecursiveDOT string

//go:embed testdata/dep-pkg-recursive-exclude-standard.json
var depPkgRecursiveExcludeStandardJSON []byte

//go:embed testdata/dep-pkg-recursive-exclude-standard.dot
var depPkgRecursiveExcludeStandardDOT string

func TestDefaultDepGraphGenerator_Generate(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		result    string
		writePath string
	}{
		{
			name:      "success-dep-pkg-recursive",
			input:     depPkgRecursiveJSON,
			result:    depPkgRecursiveDOT,
			writePath: "testdata/dep-pkg-recursive.dot",
		},
		{
			name:      "success-pkg-recursive-exclude-standard",
			input:     depPkgRecursiveExcludeStandardJSON,
			result:    depPkgRecursiveExcludeStandardDOT,
			writePath: "testdata/dep-pkg-recursive-exclude-standard.dot",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := NewDepGraphGenerator()

			var input imports.Packages

			require.NoError(t, json.Unmarshal(tt.input, &input))

			result, err := generator.Generate(&input)
			require.NoError(t, err)

			assert.Equal(t, tt.result, result)

			require.NoError(t, os.WriteFile(tt.writePath, []byte(result), 0o644)) // #nosec G306
		})
	}
}

func TestStrokeColour(t *testing.T) {
	tests := []struct {
		name     string
		pkg      *imports.Package
		expected string
	}{
		{
			name:     "local-package",
			pkg:      &imports.Package{IsLocal: true},
			expected: "gray58",
		},
		{
			name:     "owned-package",
			pkg:      &imports.Package{IsOwned: true},
			expected: "gray69",
		},
		{
			name:     "external-package",
			pkg:      &imports.Package{IsExternal: true},
			expected: "gray65",
		},
		{
			name:     "standard-library-package",
			pkg:      &imports.Package{IsStandard: true},
			expected: "gray80",
		},
		{
			name:     "internal-package",
			pkg:      &imports.Package{IsInternal: true},
			expected: "grey25",
		},
		{
			name:     "vendor-package",
			pkg:      &imports.Package{IsVendor: true},
			expected: "grey36",
		},
		{
			name:     "unknown-package-type",
			pkg:      &imports.Package{},
			expected: "invis",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &DefaultDepGraphGenerator{}

			assert.Equal(t, tt.expected, g.strokeColour(tt.pkg))
		})
	}
}

func TestFillColour(t *testing.T) {
	tests := []struct {
		name     string
		pkg      *imports.Package
		expected string
	}{
		{
			name:     "local-package",
			pkg:      &imports.Package{IsLocal: true},
			expected: "gray99",
		},
		{
			name:     "owned-package",
			pkg:      &imports.Package{IsOwned: true},
			expected: "gray99",
		},
		{
			name:     "external-package",
			pkg:      &imports.Package{IsExternal: true},
			expected: "gray83",
		},
		{
			name:     "standard-library-package",
			pkg:      &imports.Package{IsStandard: true},
			expected: "gray100",
		},
		{
			name:     "internal-package",
			pkg:      &imports.Package{IsInternal: true},
			expected: "grey44",
		},
		{
			name:     "vendor-package",
			pkg:      &imports.Package{IsVendor: true},
			expected: "grey55",
		},
		{
			name:     "unknown-package-type",
			pkg:      &imports.Package{},
			expected: "invis",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &DefaultDepGraphGenerator{}

			assert.Equal(t, tt.expected, g.fillColour(tt.pkg))
		})
	}
}
