# Jacquard Changelog

This file records notable changes to Jacquard.
Releases use semantic versioning as in `MAJOR.MINOR.PATCH`.

## Change entries

- Added: New features.
- Changed: Changes in existing functionality.
- Deprecated: Once-stable features scheduled for removal.
- Removed: Removed features.
- Fixed: Bug fixes.
- Security: Vulnerability fixes.

## [Unreleased]

### Added

- Checker contract for contract validation, quote fidelity, value support,
  evidence reporting, and a separate filter-size audit.
- Go module with JSON interfaces, arbitrary-precision interface validators,
  a validation CLI, and C exports.
- Synthetic contract fixtures, package and integration tests, and shared-library smoke tests.
- GreenDoc user, integrator, and contributor guides with an indigo palette.
- Citation metadata and citation sections in the README and documentation home page.
- Go package overviews, API documentation, and executable usage examples.
- GitHub Actions build and deployment of the documentation from `docs/` on `main`.
- Documentation asset directories and conditional logo/favicon rendering.

### Changed

- Keep the public entry point in `jacquard.go`; move implementation, data types,
  schemas, and unit tests into internal packages and public API examples into `test/`.
- License Jacquard under Apache-2.0 and include the license with shared-library releases.
- Collect all public Go types, schema helpers, and validation functions in the
  root `github.com/open-and-sustainable/jacquard` package, so callers use one import.
- Describe Jacquard as deterministic evidence checking for AI-assisted knowledge
  synthesis across the README, documentation, Go package overview, and citations,
  with a plain project name in headings and citations.
- Use direct capability descriptions throughout the documentation and keep
  synthetic-data statements specific to the test suite.
- Use the uploaded logo and favicon on the documentation site and the logo in
  the README, with the logo's indigo as the documentation accent.

### Removed

- Manual schema version selectors, interface-version metadata, and versioned
  directories for static schema files.
- The assets README and unused placeholder directory.

### Fixed

- Quote the coverage flag in CI so Windows PowerShell passes the complete filename
  to Go instead of treating `.out` as a separate package argument.
