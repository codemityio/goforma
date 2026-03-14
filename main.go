package main

import (
	"log"
	"os"

	"github.com/codemityio/goforma/cmd/badge"
	"github.com/codemityio/goforma/cmd/code"
	"github.com/codemityio/goforma/internal/app"
	"github.com/urfave/cli/v2"
)

func main() {
	application := app.New(
		app.WithValues(
			name,
			`A tool to support work with Markdown.`,
			version,
			copyright,
			authorName,
			authorEmail,
			buildTime,
		),
	)

	application.Commands = []*cli.Command{
		&badge.App,
		&code.App,
	}

	if e := application.Run(os.Args); e != nil {
		log.Fatalf("error: %v", e)
	}
}
