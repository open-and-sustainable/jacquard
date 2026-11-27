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
- Go module with versioned JSON interfaces, arbitrary-precision interface validators,
  a validation CLI, and C exports.
- Synthetic contract fixtures, package and integration tests, and shared-library smoke tests.
- GreenDoc user, integrator, and contributor guides with an indigo palette.
- Citation metadata and citation sections in the README and documentation home page.
- Go package overviews, API documentation, and executable usage examples.
- GitHub Actions build and deployment of the documentation from `docs/` on `main`.
- Documentation asset directories and conditional logo/favicon rendering.

### Changed

- Collect all public Go types, schema helpers, and validation functions in the
  root `github.com/open-and-sustainable/jacquard` package, so callers use one import.
- Describe Jacquard as deterministic evidence checking for AI-assisted knowledge
  synthesis across the README, documentation, Go package overview, and citations,
  with a plain project name in headings and citations.
- Use direct capability descriptions throughout the documentation and keep
  synthetic-data statements specific to the test suite.

### Fixed

- Quote the coverage flag in CI so Windows PowerShell passes the complete filename
  to Go instead of treating `.out` as a separate package argument.
