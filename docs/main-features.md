---
title: Main Features
layout: default
---

# Main features

These are the five first-release capabilities. Their normative behavior is in the
[checker contract](checker-contract.md); checking engines are still to be implemented.

| Capability | Question answered | Result |
| --- | --- | --- |
| Contract validation | Does the extraction obey its schema and claim rules? | Validity and issues at JSON Pointers |
| Quote fidelity | Was each quote copied exactly, and where does it occur? | Exact occurrences or a mismatch diagnosis |
| Value support | Is each claimed text or number present in its own quote? | Supported values and parsed literals |
| Evidence reporting | Which claims have supporting evidence? | One evidence status per claim and totals |
| Filter-size audit | How many small and large items survive each stage? | A complete stage × field × bucket table |

Contracts describe caller-defined fields. The same core supports a participant
count, an inventory count, a place name, or an amount without embedding a project's
subject matter in its rules.

The current implementation validates the request and report interfaces. It does
not yet perform these five checks. See [getting started](getting-started.md).
