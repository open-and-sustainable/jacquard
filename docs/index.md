---
title: Home
layout: default
---

# Jacquard

Deterministic evidence checking for AI-assisted knowledge synthesis.

Jacquard is open science software for AI-assisted knowledge synthesis. It checks
whether claims extracted from a document are supported by that document. Callers
supply the original text, extraction, and a declarative contract.
The library produces reproducible findings with paths back to each claim.

The first release includes five areas: contract validation, exact quote checking,
value support, evidence reporting, and an independent filter-size audit.

The [checker contract](checker-contract.md) defines that release. The current code
provides the interface definitions and validators; the checking engines are planned.

## Start here

- [Main features](main-features.md)
- [Getting started](getting-started.md)
- [Checker contract](checker-contract.md)
- [Schemas and validation](schemas-validation.md)
- [Command line](command-line.md)
- [Use from other languages](cshared-library.md)
- [Development and synthetic tests](development.md)
- [Go API reference](https://pkg.go.dev/github.com/open-and-sustainable/jacquard)
- [Changelog](https://github.com/open-and-sustainable/jacquard/blob/main/CHANGELOG.md)

The Go API reference becomes available after the module is published and indexed.

## Open science

Rules and interfaces are versioned, findings are inspectable, and tests use invented
documents and populations. The project uses the
[GNU AGPL v3](https://github.com/open-and-sustainable/jacquard/blob/main/LICENSE).

## Citation

> Boero, R. Jacquard - Deterministic evidence checking for AI-assisted knowledge synthesis. Software under development. [GitHub](https://github.com/open-and-sustainable/jacquard).

Machine-readable [citation metadata](https://github.com/open-and-sustainable/jacquard/blob/main/CITATION.cff)
is maintained with the source. A Jacquard DOI and version-specific citation will
be added when a release is archived.
