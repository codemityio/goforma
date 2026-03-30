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
