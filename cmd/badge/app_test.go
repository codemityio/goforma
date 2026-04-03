package badge

import (
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestBadgeColour(t *testing.T) {
	tests := []struct {
		input int64
		want  string
	}{
		{input: 0, want: "red"},
		{input: colourRed, want: "yellow"},
		{input: colourYellow, want: "yellowgreen"},
		{input: colourYellowGreen, want: "green"},
		{input: colourGreen, want: "brightgreen"},
		{input: colourGreen + 10, want: "brightgreen"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, badgeColour(tt.input))
	}
}

func TestSanitizePath(t *testing.T) {
	baseDir := t.TempDir()

	tests := []struct {
		name      string
		inputPath string
		wantErr   error
		wantPath  string
	}{
		{
			name:      "valid relative path",
			inputPath: "file.txt",
			wantErr:   nil,
			wantPath:  filepath.Join(baseDir, "file.txt"),
		},
		{
			name:      "relative path outside base",
			inputPath: "../outside.txt",
			wantErr:   errPathOutsideBase,
		},
		{
			name:      "absolute path outside base",
			inputPath: "/etc/passwd",
			wantErr:   errPathOutsideBase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sanitizePath(baseDir, tt.inputPath)

			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.wantPath, got)
		})
	}
}

func TestCoverageFunction(t *testing.T) {
	setupFiles := func(covContent, docContent string) (string, string) {
		covPath := "testdata/coverage.txt"
		docPath := "testdata/README.md"

		require.NoError(t, os.WriteFile(covPath, []byte(covContent), 0o600))
		require.NoError(t, os.WriteFile(docPath, []byte(docContent), 0o600))

		documentPath = "testdata/README.md"
		id = "coverage-badge-do-not-edit"

		return covPath, docPath
	}

	tests := []struct {
		name         string
		covContent   string
		minimum      int64
		wantErr      error
		expectUpdate string
	}{
		{
			name:       "valid coverage above minimum",
			covContent: "total:\t\t\t\t\t\t\t(statements)\t75.0%\n",
			minimum:    70,
			wantErr:    nil,
			expectUpdate: `# Title

![coverage-badge-do-not-edit](https://img.shields.io/badge/Coverage-75%25-yellowgreen.svg?longCache=true&style=flat)

## Table of contents

## Summary

Test file
`,
		},
		{
			name:       "coverage below minimum",
			covContent: "total:\t\t\t\t\t\t\t(statements)\t75.0%\n",
			minimum:    80,
			wantErr:    errPkg,
			expectUpdate: `# Title

![coverage-badge-do-not-edit](https://img.shields.io/badge/Coverage-75%25-yellowgreen.svg?longCache=true&style=flat)

## Table of contents

## Summary

Test file
`,
		},
		{
			name:       "malformed coverage",
			covContent: "total:\t\t\t\t\t\t\t(statements)\tabc%\n",
			minimum:    0,
			wantErr:    errPkg,
			expectUpdate: `# Title

![coverage-badge-do-not-edit](https://img.shields.io/badge/Coverage-75%00-red.svg?longCache=true&style=flat)

## Table of contents

## Summary

Test file
`,
		},
		{
			name:       "missing total",
			covContent: "foo bar\n",
			minimum:    0,
			wantErr:    errPkg,
			expectUpdate: `# Title

![coverage-badge-do-not-edit](https://img.shields.io/badge/Coverage-75%00-red.svg?longCache=true&style=flat)

## Table of contents

## Summary

Test file
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, docPath := setupFiles(tt.covContent, `# Title

![coverage-badge-do-not-edit](https://img.shields.io/badge/Coverage-75%00-red.svg?longCache=true&style=flat)

## Table of contents

## Summary

Test file
`)

			app := &cli.App{
				Commands: App.Subcommands,
			}

			coverageCmd := app.Commands[0]

			fs := flag.NewFlagSet("coverage", flag.ContinueOnError)

			for _, f := range coverageCmd.Flags {
				switch v := f.(type) {
				case *cli.StringFlag:
					fs.String(v.Name, v.Value, v.Usage)
				case *cli.Int64Flag:
					fs.Int64(v.Name, v.Value, v.Usage)
				}
			}

			c := cli.NewContext(app, fs, nil)

			require.NoError(t, c.Set("cov-file-path", "testdata/coverage.txt"))
			require.NoError(t, c.Set("minimum", strconv.FormatInt(tt.minimum, 10)))

			require.ErrorIs(t, coverage(c), tt.wantErr)

			// Check if document updated
			updated, rErr := os.ReadFile(docPath) // #nosec G304
			require.NoError(t, rErr)

			assert.Equal(t, tt.expectUpdate, string(updated))
		})
	}
}
