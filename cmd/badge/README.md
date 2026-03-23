# `badge`

## Table of contents

- [Summary](#summary)
- [Manual](#manual)
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
   --document value  markdown file path to be updated
   --id value        placeholder identifier
   --help, -h        show help
```

## Usage

Use the following placeholder to indicate where the badge should be placed. Make sure `<PLACEHOLDER_ID>` is replaced
with a correct placeholder identifier (e.g. `example-badge-id`).

    ![<PLACEHOLDER_ID>]()

    ![example-badge-id]()

### `coverage`

Use `go tool cover -func="tmp/coverage.out" -o tmp/coverage.in` to cover the coverage output.

Pass that converted output as input to the coverage badge command.
