---
title: Development
layout: default
---

# Development

## Repository layout

The public Go API is exposed by the root `jacquard` package:

| Path | Purpose |
| --- | --- |
| `jacquard.go` | Public entry point: exported types and functions |
| `internal/contract/` | Document, claim, report, and audit data types |
| `internal/validation/` | Validation implementation, schema loading, and unit tests |
| `internal/validation/schemas/` | Embedded JSON schema files |
| `cmd/jacquard/` | JSON validation and schema-display CLI |
| `utils/sharedlib/` | C exports and string ownership boundary |
| `test/` | Public API examples, integration fixtures, and FFI smoke test |
| `docs/` | User/integrator/contributor guide using the GreenDoc theme |
| `.github/workflows/` | Go checks and native shared-library release builds |

The root contains one Go source file, `jacquard.go`. It exposes the public API and
delegates to the internal packages. Unit tests sit beside the implementation;
public API examples and cross-package tests live in `test/`. The [checker contract](checker-contract.md)
defines the document-checking and audit data formats and rules.

## Go documentation

The entry point's package comment covers purpose, supported interfaces, behavior,
and usage. Public declarations document their inputs, outputs, errors, and ownership
where applicable. Internal packages include package comments. Go doc links connect
related types and functions.

Usage examples live under `test/` so `go test` compiles and executes them against
the public API. The package comment also includes a usage example. All test data
is synthetic.
Inspect the rendered documentation locally without publishing:

```sh
go doc .
go doc ValidateJSON
go test -run Example ./...
```

## Setup and checks

Use Go 1.26+ and a C compiler. Download dependencies once with `go mod download`.
Then run:

```sh
gofmt -w .
go vet ./...
go test -race ./...
go build -o bin/jacquard ./cmd/jacquard
go build -buildmode=c-shared -o dist/jacquard.so ./utils/sharedlib
python3 test/shared_library_test.py dist/jacquard.so
```

Create `bin/` and `dist/` first when building manually. Build outputs are ignored.
CI uses `go-version-file: go.mod` so its toolchain follows the module's declared
minimum rather than a separately maintained version number.

## Synthetic testing

The test suite uses synthetic documents and populations stored in this repository.
It runs without credentials, model answers, private data, or an external corpus.
The [testing guide](testing.md) lists executable coverage and contract scenarios.

## Documentation

Configure the site title, description, navigation, colors, and artwork paths in
`docs/_config.yml`. The site uses the `open-and-sustainable/greendoc_theme` theme
with a local layout in `docs/_layouts/default.html`.

The [Documentation Pages workflow](publishing.md) builds `docs/` from `main`
and publishes the generated artifact to GitHub Pages.

## Changelog and citation metadata

Use these categories in `CHANGELOG.md`: Added, Changed, Deprecated, Removed, Fixed,
and Security. Keep pending changes under `[Unreleased]`. At release time, move
them to `[MAJOR.MINOR.PATCH] - YYYY-MM-DD`,
record the actual release date, and start a new unreleased entry.

Keep `CITATION.cff`, the README citation, and the documentation home-page citation
consistent. The CFF records the title, author, ORCID, affiliation, repository,
documentation URL, abstract, keywords, and license. Add version, release date,
DOI, and the DOI identifier entry when those values exist. The Go API reference
is generated from source comments and executable examples when the module is
published and indexed by pkg.go.dev.
