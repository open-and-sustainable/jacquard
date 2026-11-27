---
title: Contributing
layout: default
---

# Contributing to Jacquard

Use Go 1.26 or later, clone the repository, and run `go mod download` followed by
`go test ./...`. A C compiler and Python 3 are needed for shared-library smoke tests.

Follow the [checker contract](checker-contract.md), document exported APIs, format
Go code with `gofmt`, and include focused tests for changed behavior. Unit tests live
beside their packages; integration fixtures and cross-interface tests live in `test/`.
All test data must be synthetic. Use fixed seeds for generated cases and never require
models, credentials, downloaded test data, or a private dataset.

Follow Alembica's documentation pattern: package overviews in `doc.go`, function
comments describing parameters, results, and errors, and executable synthetic
examples in `example_test.go`. Document C string ownership at each exported entry
point. Preview Go documentation locally with `go doc`.

Before submitting, run the checks in [development](development.md). Describe the
problem, resulting behavior, and validation in a pull request. Interface/rule changes
need explicit version decisions and a categorized `[Unreleased]` changelog entry.
Keep the README citation, documentation home-page citation, and `CITATION.cff`
consistent. Add release dates, versions, and Jacquard archive identifiers when
they are assigned.

Participants follow the [code of conduct](CODE_OF_CONDUCT.md).
