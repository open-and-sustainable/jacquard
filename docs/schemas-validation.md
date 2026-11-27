---
title: Schemas & Validation
layout: default
---

# Schemas and validation

`definitions/v1/` contains embedded draft-07 schemas for four interfaces:

| Type | Go type | Validator |
| --- | --- | --- |
| `input` | `jacquard.Input` | `jacquard.ValidateInput` |
| `output` | `jacquard.Output` | `jacquard.ValidateOutput` |
| `funnel_input` | `jacquard.FunnelInput` | `jacquard.ValidateFunnelInput` |
| `funnel_output` | `jacquard.FunnelOutput` | `jacquard.ValidateFunnelOutput` |

Import `github.com/open-and-sustainable/jacquard` for every type and function above.
The `definitions/v1/` directory holds schema data and is not a Go package.

All validators take JSON and a version string, currently `v1`. Versions are exact;
unknown versions/types return errors rather than selecting a fallback. Retrieve a
fresh schema copy with `jacquard.Schema("v1", "input")`.

## Two validation layers

Interface validation checks a request or report's JSON shape. It does not execute
the caller's embedded extraction schema, resolve claim selectors, check evidence,
verify summary arithmetic, or compute an audit. These are engine responsibilities
defined by the [checker contract](checker-contract.md).

During a future `Check`, the caller's `contract.schema` validates the extraction,
then claim rules check its semantics. A failure of the extraction schema is a
finding in the checking report, while a malformed request is an operation error.

## Explicit schema validation

`jacquard.ValidateJSON(instance, schema)` returns sorted `jacquard.Issue`
values and an error for an unusable schema. Malformed instance JSON is an
`invalid_json` finding. Valid instances return a non-null empty issue slice.

Schema compilation is restricted to the supplied schema. Local `$ref` references
work; attempts to load a file or URL are rejected by the schema loader. Numbers are
decoded without a floating-point conversion. Issue paths preserve property names
containing dots and escape slashes and tildes as JSON Pointer tokens.
