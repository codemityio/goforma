# `badge`

## Table of contents

## Summary

A simple tool to generate badges within a file.

## Usage

Use the following placeholder to indicate where the badge should be placed. Make sure `<PLACEHOLDER_ID>` is replaced
with a correct placeholder identifier (e.g. `example-badge-id`).

    ![<PLACEHOLDER_ID>]()

    ![example-badge-id]()

### `coverage`

Use `go tool cover -func="tmp/coverage.out" -o tmp/coverage.in` to cover the coverage output.

Pass that converted output as input to the coverage badge command.
