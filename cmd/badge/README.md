# `badge`

## Table of contents

- [Summary](#summary)
- [Manual](#manual)
- [Subcommands](#subcommands)
  - [`coverage`](#coverage)
- [Usage](#usage)
  - [`coverage`](#coverage)

## Summary

A simple tool to generate badges within a file.

## Manual

``` bash
$ goforma badge --help
NAME:
   goforma badge

USAGE:
   goforma badge [command options]

DESCRIPTION:
   Document badge generator

COMMANDS:
   coverage  Generate coverage badge within a document
   help, h   Shows a list of commands or help for one command

OPTIONS:
   --document-path value  markdown file path to be updated
   --id value             placeholder identifier
   --help, -h             show help
```

## Subcommands

### `coverage`

``` bash
$ goforma badge --document-path=README.MD --id=id coverage --help
NAME:
   goforma badge coverage - Generate coverage badge within a document

USAGE:
   goforma badge coverage [command options]

OPTIONS:
   --cov-file-path value  Generate coverage file path
   --minimum value        Minimum coverage threshold (default: 0)
   --help, -h             show help
```

## Usage

### `coverage`

Use the following placeholder to indicate where the badge should be placed. Make sure `<PLACEHOLDER_ID>` is replaced
with a correct placeholder identifier (e.g.`example-badge-id`).

    ![<PLACEHOLDER_ID>]()

    ![example-badge-id]()

Use `go tool cover -func="tmp/coverage.out" -o tmp/coverage.in` to convert the coverage output before using this tool.

Pass that converted output as input to the coverage badge command.

``` bash
goforma badge \
  --document-path=README.md \
  --id=example-badge-id \
  coverage \
  --cov-file-path=tmp/coverage.in \
  --minimum=80
```
