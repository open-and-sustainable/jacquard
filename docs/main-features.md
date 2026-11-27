---
title: Main Features
layout: default
---

# Main features

The [checker contract](checker-contract.md) defines five areas of evidence checking
for AI-assisted knowledge synthesis:

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

The Go library, CLI, and C exports validate the document-checking and audit request
and report interfaces. See [getting started](getting-started.md) for a complete Go
example and [schemas and validation](schemas-validation.md) for the validation APIs.
