---
title: Getting Started
layout: default
---

# Getting started

## Requirements

Use Go 1.26 or later. A C compiler is required to build the shared library.
The library needs no model service or credentials.

From a checkout, download the pinned Go dependencies and run the synthetic suite:

```sh
go mod download
go test ./...
```

## Validate an interface

All public Go types and functions are available from the root `jacquard` package.

```go
package main

import (
    "fmt"
    "os"

    "github.com/open-and-sustainable/jacquard"
)

func main() {
    data, err := os.ReadFile("test/fixtures/check_input.json")
    if err != nil {
        panic(err)
    }
    if err := jacquard.ValidateInput(string(data)); err != nil {
        panic(err)
    }
    fmt.Println("Valid request shape")
}
```

Run the equivalent CLI:

```sh
go run ./cmd/jacquard validate -type input < test/fixtures/check_input.json
```

The example validates the request's JSON structure. The
[checker contract](checker-contract.md) defines claim and evidence-report rules.

## Next steps

Read [schemas and validation](schemas-validation.md) for the distinction between
request validation and extraction validation, and [C-shared library](cshared-library.md)
for exported functions and memory ownership.
