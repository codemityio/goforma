package cmd

import (
	"github.com/urfave/cli/v2"
)

// App main application.
var App = cli.Command{
	Before: func(c *cli.Context) error {
		var err error

		return err
	},
	Subcommands: []*cli.Command{
		{
			Name:   "cp",
			Usage:  "",
			Action: cp,
		},
		{
			Name:   "rm",
			Usage:  "",
			Action: rm,
		},
	},
}

func cp(c *cli.Context) error {
	var err error

	return err
}

func rm(c *cli.Context) error {
	var err error

	return err
}
