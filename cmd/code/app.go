package code

import (
	"strconv"

	"github.com/codemityio/goforma/pkg/code/imports"
	"github.com/urfave/cli/v2"
)

// App main application.
var App = cli.Command{ //nolint:exhaustruct,gochecknoglobals
	Name:         "code",
	Aliases:      nil,
	Usage:        "",
	UsageText:    "",
	Description:  "Go code tools",
	ArgsUsage:    "",
	Category:     "",
	BashComplete: nil,
	Before: func(c *cli.Context) error {
		return nil
	},
	After:        nil,
	Action:       nil,
	OnUsageError: nil,
	Flags:        []cli.Flag{},
	Subcommands: []*cli.Command{
		{
			Name:  "dep",
			Usage: `Generate dependency graph`,
			Flags: []cli.Flag{
				&cli.StringFlag{ //nolint:exhaustruct
					Name:     "path",
					Usage:    "Path to be scanned (e.g. ./pkg/example, ./pkg/example/...)",
					Required: true,
					Value:    ".",
				},
				&cli.StringFlag{ //nolint:exhaustruct
					Name:     "workdir",
					Usage:    "Working directory (e.g. absolute path of the current project root directory)",
					Required: false,
					Value:    ".",
				},
				&cli.IntFlag{ //nolint:exhaustruct
					Name: "depth",
					Usage: "Depth to scan for dependencies (default " + strconv.Itoa(
						imports.DefaultDepth,
					) + ")",
					Required: false,
					Value:    imports.DefaultDepth,
				},
				&cli.StringSliceFlag{ //nolint:exhaustruct
					Name:     "owned",
					Usage:    "Indicate owned packages/prefixes (e.g. --owned=github.com/example-one --owned=github.com/example-two)",
					Required: false,
				},
				&cli.StringSliceFlag{ //nolint:exhaustruct
					Name:     "exclude-path",
					Usage:    "Indicate packages/prefixes to be excluded (e.g. --exclude-path=github.com/example-one --exclude-path=github.com/example-two)", //nolint:lll
					Required: false,
				},
				&cli.BoolFlag{ //nolint:exhaustruct
					Name:     "exclude-standard",
					Usage:    "Exclude standard library packages",
					Required: false,
				},
				&cli.BoolFlag{ //nolint:exhaustruct
					Name:     "exclude-vendor",
					Usage:    "Exclude vendor packages",
					Required: false,
				},
				&cli.BoolFlag{ //nolint:exhaustruct
					Name:     "exclude-internal",
					Usage:    "Exclude internal packages",
					Required: false,
				},
			},
			Action: dep,
		},
		{
			Name:  "uml",
			Usage: `Generate UML graph`,
			Flags: []cli.Flag{
				&cli.StringFlag{ //nolint:exhaustruct
					Name:     "path",
					Usage:    "Path to be scanned (e.g. ./pkg/example, ./pkg/example/...)",
					Required: true,
					Value:    ".",
				},
				&cli.StringFlag{ //nolint:exhaustruct
					Name:     "workdir",
					Usage:    "Working directory (e.g. absolute path of the current project root directory)",
					Required: false,
					Value:    ".",
				},
				&cli.StringFlag{ //nolint:exhaustruct
					Name:     "json-output-path",
					Usage:    "Path to output code tree in json format",
					Required: false,
				},
				&cli.BoolFlag{ //nolint:exhaustruct
					Name:     "include-legend",
					Usage:    "Include diagram legend",
					Required: false,
					Value:    false,
				},
				&cli.BoolFlag{ //nolint:exhaustruct
					Name:     "include-primitive",
					Usage:    "Include primitive types",
					Required: false,
					Value:    false,
				},
				&cli.BoolFlag{ //nolint:exhaustruct
					Name:     "include-var",
					Usage:    "Include declared variables",
					Required: false,
					Value:    false,
				},
				&cli.BoolFlag{ //nolint:exhaustruct
					Name:     "include-const",
					Usage:    "include declared constants",
					Required: false,
					Value:    false,
				},
				&cli.BoolFlag{ //nolint:exhaustruct
					Name:     "include-func",
					Usage:    "Include declared functions",
					Required: false,
					Value:    false,
				},
				&cli.BoolFlag{ //nolint:exhaustruct
					Name:     "include-not-exported",
					Usage:    "Include all not exported variables, types, fields, methods, etc...",
					Required: false,
					Value:    false,
				},
				&cli.BoolFlag{ //nolint:exhaustruct
					Name:     "include-doc",
					Usage:    "Include all code doc blocks",
					Required: false,
					Value:    false,
				},
				&cli.BoolFlag{ //nolint:exhaustruct
					Name:     "include-doc-comment-slashes",
					Usage:    "Include all code doc blocks comment slashes",
					Required: false,
					Value:    false,
				},
			},
			Action: uml,
		},
	},
}
