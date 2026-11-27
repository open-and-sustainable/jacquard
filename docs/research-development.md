---
title: Research & Development
layout: default
---

# Research and development

Jacquard supports reproducible, AI-assisted knowledge synthesis through explicit,
deterministic checks of documentary support, independent of the extraction method.

The first release must implement all five capabilities in the
[checker contract](checker-contract.md). Development proceeds from the interface
foundation to schema/claim rules, quote location and diagnosis, value readers,
evidence reports, the size audit, and parity between Go/CLI/C interfaces.
Acceptance comes entirely from the repository's [synthetic suite](testing.md).

Future work can add language tables, additional currencies, explicit quote
tolerances, and comparison against outside values. Those features require their
own versioned contracts; they must not silently alter first-release results.

Contributions should include a minimal invented case demonstrating the behavior,
the applicable contract section, and an expected result. See
[contribution guidelines](CONTRIBUTING.md).
