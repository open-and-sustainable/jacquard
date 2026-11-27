---
title: Research & Development
layout: default
---

# Research and development

Jacquard supports reproducible, AI-assisted knowledge synthesis through explicit,
deterministic checks of documentary support, independent of the extraction method.

The [checker contract](checker-contract.md) records the rules for claim validation,
quote location and diagnosis, value support, evidence reports, and size audits.
Go, CLI, and C interfaces share the same JSON schemas and validation rules.

Versioned contracts make behavior changes explicit. Language rules, currency
readings, and quote-matching policies belong to their declared rule versions.
Changes preserve reproducibility and include a changelog entry.

The [test suite](testing.md) uses synthetic fixtures to exercise interface shapes,
numeric boundaries, error reporting, and cross-language bindings.

Contributions should include a minimal invented case demonstrating the behavior,
the applicable contract section, and an expected result. See
[contribution guidelines](CONTRIBUTING.md).
