# `jacquard`

Open Science Software for Deterministic Evidence Checking of Extracted Claims.

## About

`jacquard` checks an extraction against the exact text supplied to its extractor.
The first release covers contract validation, quote fidelity, value support,
evidence reporting, and a separate audit of filtering by item size.

**Development status:** the first-release behavior is defined in the
[checker contract](docs/checker-contract.md). This repository currently implements
the JSON interface definitions, schema validation, a validation CLI, and C exports.
The document checker and size audit engines are not implemented yet.

## Installation (Go)

Go 1.26 or later and a C compiler for the optional shared library are required.
From this checkout:

```sh
go test ./...
go build -o bin/jacquard ./cmd/jacquard
```

Use the current interface validator:

```sh
go run ./cmd/jacquard validate -type input < test/fixtures/check_input.json
```

## Go library

Import the root package to access the public types, schema helpers, and validators:

```go
import "github.com/open-and-sustainable/jacquard"
```

For example, call `jacquard.ValidateInput(inputJSON, "v1")` or
`jacquard.ValidateJSON(instance, schema)`. See the
[Go usage example](docs/getting-started.md) for a complete program.

## Documentation

- [User guide](docs/index.md)
- [First-release checker contract](docs/checker-contract.md)
- [Schemas and validation](docs/schemas-validation.md)
- [C-shared library](docs/cshared-library.md)
- [Development and testing](docs/development.md)
- [API reference](https://pkg.go.dev/github.com/open-and-sustainable/jacquard)
- [Changelog](CHANGELOG.md)

The documentation site is configured for
[open-and-sustainable.github.io/jacquard](https://open-and-sustainable.github.io/jacquard/).
The site and pkg.go.dev reference become available after publication and indexing.

The documentation and development patterns follow
[Alembica](https://github.com/open-and-sustainable/alembica): versioned JSON schemas,
documented Go APIs, C exports, a themed documentation site, package tests, and
integration tests. Jacquard's public Go API lives in the root package.

## Testing

All documents, extractions, populations, and expected findings in tests are synthetic.
Tests require no models, credentials, external corpora, or network calls. Go downloads
pinned build dependencies during setup; the test suite itself runs offline.

## Authors and contributions

Riccardo Boero — [ribo@nilu.no](mailto:ribo@nilu.no).
See [contribution guidelines](docs/CONTRIBUTING.md) and the
[code of conduct](docs/CODE_OF_CONDUCT.md).

## License

GNU Affero General Public License, version 3. See [LICENSE](LICENSE).

## Citation

> Boero, R. `jacquard` - Open Science Software for Deterministic Evidence Checking of Extracted Claims. Software under development. [GitHub](https://github.com/open-and-sustainable/jacquard).

Machine-readable citation metadata is in [CITATION.cff](CITATION.cff). A Jacquard
DOI and version-specific citation will be added when a release is archived.
