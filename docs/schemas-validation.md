---
title: Schemas & Validation
layout: default
---

# Schemas and validation

`internal/validation/schemas/` contains embedded draft-07 schemas for four interfaces:

| Type | Go type | Validator |
| --- | --- | --- |
| `input` | `jacquard.Input` | `jacquard.ValidateInput` |
| `output` | `jacquard.Output` | `jacquard.ValidateOutput` |
| `funnel_input` | `jacquard.FunnelInput` | `jacquard.ValidateFunnelInput` |
| `funnel_output` | `jacquard.FunnelOutput` | `jacquard.ValidateFunnelOutput` |

Import `github.com/open-and-sustainable/jacquard` for every type and function above.
The schema files are internal data exposed through `jacquard.Schema`.

Each validator takes one JSON string. `jacquard.Validate(instance, kind)` selects
the interface explicitly; an unknown kind returns an error. Retrieve a fresh
schema copy with `jacquard.Schema("input")`.

## Two validation layers

Interface validation checks a request or report's JSON shape. It does not execute
the caller's embedded extraction schema, resolve claim selectors, check evidence,
verify summary arithmetic, or compute an audit. These are engine responsibilities
defined by the [checker contract](checker-contract.md).

The checker contract distinguishes extraction-schema findings from request errors.
Use `jacquard.ValidateJSON(input.Extraction, input.Contract.Schema)` to validate an
extraction against the schema declared by its caller.

## Explicit schema validation

`jacquard.ValidateJSON(instance, schema)` returns sorted `jacquard.Issue`
values and an error for an unusable schema. Malformed instance JSON is an
`invalid_json` finding. Valid instances return a non-null empty issue slice.

Schema compilation is restricted to the supplied schema. Local `$ref` references
work; attempts to load a file or URL are rejected by the schema loader. Numbers are
decoded without a floating-point conversion. Issue paths preserve property names
containing dots and escape slashes and tildes as JSON Pointer tokens.
