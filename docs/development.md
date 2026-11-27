---
title: Development
layout: default
---

# Development

## Repository layout

The documentation, schema definitions, bindings, and testing follow Alembica's
patterns. The public Go API is collected in the root `jacquard` package:

| Path | Purpose |
| --- | --- |
| Repository root | Public Go types, schema helpers, and validation functions |
| `definitions/` | Embedded JSON schema data; no separate Go package |
| `cmd/jacquard/` | JSON validation and schema-display CLI |
| `utils/sharedlib/` | C exports and string ownership boundary |
| `test/` | Synthetic integration fixtures and FFI smoke test |
| `docs/` | User/integrator/contributor guide using the GreenDoc theme |
| `.github/workflows/` | Go checks and native shared-library release builds |

Packages include `doc.go` or a package comment. Unit tests sit beside their package;
cross-package tests live in `test/`. The [checker contract](checker-contract.md)
defines the document-checking and audit data formats and rules.

## Go documentation

Following Alembica's documentation pattern, each package has a `doc.go` overview
covering its purpose, supported interfaces, main components, behavior, and usage.
Public declarations document their inputs, outputs, errors, and ownership where
applicable. Go doc links connect related types and functions.

Usage examples live in `example_test.go` so `go test` compiles and executes them
and pkg.go.dev can display them beside the relevant API. All examples are synthetic.
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

`docs/_config.yml` uses the same `open-and-sustainable/greendoc_theme` theme and
user/integrator/contributor navigation as Alembica. Jacquard uses indigo navigation
and indigo interaction colors. The site title, description, navigation, colors, and
artwork paths are configured in that file.

The [Documentation Pages workflow](publishing.md) builds `docs/` from `main`
and publishes the generated artifact to GitHub Pages. The local layout makes
logo and favicon rendering conditional on the presence of their asset files.

## Changelog and citation metadata

The root `CHANGELOG.md` follows Alembica's categorized semantic-versioning pattern:
Added, Changed, Deprecated, Removed, Fixed, and Security. Keep pending changes under
`[Unreleased]`. At release time, move them to `[MAJOR.MINOR.PATCH] - YYYY-MM-DD`,
record the actual release date, and start a new unreleased entry.

Keep `CITATION.cff`, the README citation, and the documentation home-page citation
consistent. The CFF records the title, author, ORCID, affiliation, repository,
documentation URL, abstract, keywords, and license. Add version, release date,
DOI, and the DOI identifier entry when those values exist. Use Jacquard's own
archive identifiers. The Go API reference is generated from source comments and
executable examples when the module is published and indexed by pkg.go.dev.
