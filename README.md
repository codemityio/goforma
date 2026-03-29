# ![GoForma](logo.jpg)

![coverage-badge-do-not-edit](https://img.shields.io/badge/Coverage-87%25-green.svg?longCache=true&style=flat)

## Table of contents

- [Summary](#summary)
- [Development](#development)
  - [`make`](#make)
- [Installation](#installation)
- [Usage](#usage)
  - [Manual](#manual)
  - [Subcommands](#subcommands)
  - [Docker](#docker)
- [Packages](#packages)
- [License](#license)

## Summary

Tool to support work with Go code.

## Development

To work with the codebase, use `make` command as the primary entry point for all project tools.

Navigate the available options using the arrow keys: `↓ ↑ → ←`. Use `/` to toggle search.

### `make`

``` bash
$ make help
build                          Build container image
buildx                         Build container multi platform images and push
check                          Run all CI required targets
cleanup                        Cleanup project
cmd                            Run a command passed as COMMAND= value (e.g. make cmd COMMAND="make check")
cov-open                       Inspect coverage in the browser
cov-report                     Check coverage report
cov                            Check coverage
diff                           Check diff to ensure this project consistency
docs-cmd                       Generate pkg docs
docs-depgraph                  Generate dependency graph
docs-main                      Generate main docs
docs-pkg                       Generate pkg docs
docs-render                    Render diagrams
docs-uml                       Generate UML documentation
docs                           Generate all docs
exec                           Execute built bin (use FLAGS= and COMMAND= environment variables to pass main command flags and subcommand with flags when needed)
fmt                            Format code
gen                            Go generate
go                             Build Go
help                           Prints help for targets with comments
install                        Install binary locally
next                           Create a new version (bump prerelease or patch)
prep                           Prepare dev tools
push                           Push image
reset                          Stop and remove project containers, remove project volumes, remove project images
run-container                  Run container (use FLAGS= and COMMAND= environment variables to pass main command flags and subcommand with flags when needed)
run-go                         Run go (use FLAGS= and COMMAND= environment variables to pass main command flags and subcommand with flags when needed)
statan-fix                     Analyze code and fix
statan                         Analyze code
test-race                      Run race tests
test                           Run tests
update                         Update all dependencies
vendor                         Run go mod vendor
version                        Print the most recent version
```

## Installation

To install the tool use `make install` (directly from the repository clone) or use
`go install github.com/codemityio/goforma@latest`.

## Usage

Once you have the tool installed, just use the `goforma` command to get started.

### Manual

``` bash
$ goforma --help
NAME:
   goforma - A new cli application

USAGE:
   goforma [global options] command [command options]

VERSION:
   latest

DESCRIPTION:
   A tool to support work with Markdown.

AUTHOR:
   codemityio

COMMANDS:
   badge    
   code     
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

COPYRIGHT:
   codemityio
```

### Subcommands

- [`badge`](cmd/badge/README.md) - A simple tool to generate badges within a file.
- [`code`](cmd/code/README.md) - A simple tool to perform tasks on **Go** code.

### Docker

``` bash
$ docker run codemityio/goforma
```

## Packages

- [`code`](pkg/code/README.md) - A package containing tools to perform code analysis, generate documentation and so on.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
