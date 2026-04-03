package badge

import (
	"github.com/urfave/cli/v2"
)

var documentPath, id string //nolint:gochecknoglobals

// App main application.
var App = cli.Command{ //nolint:exhaustruct,gochecknoglobals
	Name:         "badge",
	Aliases:      nil,
	Usage:        "",
	UsageText:    "",
	Description:  "Document badge generator",
	ArgsUsage:    "",
	Category:     "",
	BashComplete: nil,
	Before: func(c *cli.Context) error {
		documentPath = c.String("document-path")
		id = c.String("id")

		return nil
	},
	After:        nil,
	Action:       nil,
	OnUsageError: nil,
	Flags: []cli.Flag{
		&cli.StringFlag{ //nolint:exhaustruct
			Name:     "document-path",
			Usage:    "markdown file path to be updated",
			Required: true,
		},
		&cli.StringFlag{ //nolint:exhaustruct
			Name:     "id",
			Usage:    "placeholder identifier",
			Required: true,
		},
	},
	Subcommands: []*cli.Command{
		{
			Name:  "coverage",
			Usage: `Generate coverage badge within a document`,
			Flags: []cli.Flag{
				&cli.StringFlag{ //nolint:exhaustruct
					Name:     "cov-file-path",
					Usage:    "Generate coverage file path",
					Required: true,
				},
				&cli.Int64Flag{ //nolint:exhaustruct
					Name:     "minimum",
					Usage:    "Minimum coverage threshold",
					Required: false,
					Value:    0,
				},
			},
			Action: coverage,
		},
	},
}
