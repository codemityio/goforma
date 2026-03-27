# `code`

## Table of contents

- [Summary](#summary)
- [Manual](#manual)
- [Subcommands](#subcommands)
  - [`uml`](#uml)
  - [`dep`](#dep)
- [Usage](#usage)

## Summary

A simple tool to perform tasks on **Go** code.

## Manual

``` bash
$ goforma code --help
NAME:
   goforma code

USAGE:
   goforma code [command options]

DESCRIPTION:
   Go code tools

COMMANDS:
   dep      Generate dependency graph
   uml      Generate UML graph
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

## Subcommands

### `uml`

``` bash
$ goforma code uml --help
NAME:
   goforma code uml - Generate UML graph

USAGE:
   goforma code uml [command options]

OPTIONS:
   --path value                   Path to be scanned (e.g. ./pkg/example, ./pkg/example/...) (default: ".")
   --workdir value                Working directory (e.g. absolute path of the current project root directory) (default: ".")
   --json-output-path value       Path to output code tree in json format
   --include-legend               Include diagram legend (default: false)
   --include-primitive            Include primitive types (default: false)
   --include-var                  Include declared variables (default: false)
   --include-const                include declared constants (default: false)
   --include-func                 Include declared functions (default: false)
   --include-not-exported         Include all not exported variables, types, fields, methods, etc... (default: false)
   --include-doc                  Include all code doc blocks (default: false)
   --include-doc-comment-slashes  Include all code doc blocks comment slashes (default: false)
   --help, -h                     show help
```

### `dep`

``` bash
$ goforma code dep --help
NAME:
   goforma code dep - Generate dependency graph

USAGE:
   goforma code dep [command options]

OPTIONS:
   --path value                                   Path to be scanned (e.g. ./pkg/example, ./pkg/example/...) (default: ".")
   --workdir value                                Working directory (e.g. absolute path of the current project root directory) (default: ".")
   --depth value                                  Depth to scan for dependencies (default 256) (default: 256)
   --owned value [ --owned value ]                Indicate owned packages/prefixes (e.g. --owned=github.com/example-one --owned=github.com/example-two)
   --exclude-path value [ --exclude-path value ]  Indicate packages/prefixes to be excluded (e.g. --exclude-path=github.com/example-one --exclude-path=github.com/example-two)
   --exclude-standard                             Exclude standard library packages (default: false)
   --exclude-vendor                               Exclude vendor packages (default: false)
   --exclude-internal                             Exclude internal packages (default: false)
   --help, -h                                     show help
```

## Usage

``` bash
goforma code uml \
  --workdir=${PWD} \
  --include-var \
  --include-const \
  --include-func \
  --include-not-exported \
  --path "./..." >"graph.puml"
docker run --rm \
  -w "${PWD}" \
  -v "${PWD}:${PWD}" \
  codemityio/notatio:latest plantuml --input-path="graph.puml" --output-format=svg
```

``` bash
goforma code dep \
  --workdir=${PWD} \
  --exclude-standard \
  --exclude-vendor \
  --path "./..." >"depgraph.dot"
docker run --rm \
  -w "${PWD}" \
  -v "${PWD}:${PWD}" \
  codemityio/notatio:latest graphviz --input-path="depgraph.dot" --output-format=svg
```
