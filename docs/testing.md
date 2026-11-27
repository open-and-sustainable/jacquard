---
title: Testing
layout: default
---

# Testing

The test suite uses synthetic fixtures written in this repository. Generated test
data uses fixed seeds. Tests run offline after dependency setup, without model
calls, credentials, or an external document collection.

## Executable test coverage

- Embedded schemas: supported types, exact versions, independent returned copies.
- Interface validation: valid/rejected requests and reports, integer boundaries,
  local schema references, blocked file/network references, JSON Pointer escaping,
  deterministic finding order, malformed JSON, and invalid UTF-8.
- CLI: validation exit codes, schema display, help, invalid invocation, and I/O errors.
- Go integration: synthetic request/report fixture shapes and typed JSON round trips.
- C interface: all validators, malformed input, NULL arguments, repeated allocation/free.

The report fixtures are hand-authored examples. Interface tests verify their
schemas and typed JSON round trips. Checking a report's structure and checking
its evidence are separate assertions.

## Checker contract scenarios

These scenarios describe the boundary rules in the
[checker contract](checker-contract.md).

| Area | Scenarios |
| --- | --- |
| Contract validation | All five claim kinds; absent vs null; blank strings; invalid types; unknown reference; inline/reference conflict; custom part paths; escaped property names; array selectors; overlapping selectors; record rules; schema failure short-circuit |
| Quote fidelity | Exact/repeated/overlapping matches; empty source; multibyte Unicode offsets; combining characters; case/ellipsis/typography/whitespace diagnoses in precedence order; unmatched prefix/suffix; diagnoses never pass |
| Count support | Zero; digits; commas; decimals rejected; ordinals rejected; hyphenated words; hundreds/scales; conjunction boundaries; ranges; arbitrary precision; explicit implied-count policy; unsupported values |
| Amount support | As-written parts; dollar markers; decimal scales; cent rounding; multiple numbers/ranges; unsupported/conflicting currencies; unquantified claims |
| Evidence reports | Every status/reason; nonexact quote precedence; text elsewhere; exact note; invalid format with supported evidence; schema-invalid report; correct summaries; stable hashes/serialization |
| Size audit | Missing/unparsed/zero; nearest-rank cutoffs; ties; empty and small reference sets; out-of-range positives; several linked buckets; per-item deduplication; zero cells; stage-prefix errors; duplicate ids |
| Interface parity | The same synthetic cases through Go, CLI single JSON/JSON Lines, and C; identical reports/errors and correct memory ownership |

Test fixtures live in `test/fixtures/`; focused cases also live in package table
tests. Evidence assertions compare computed findings against independently stated
expected results and source text. Schema assertions check JSON structure.

For example, `three volunteers` must yield 3; `between three and five` must yield
3 and 5, while `one hundred and five` must yield 105. A quote `2.5 volunteers` must
not yield either 2 or 5. A repeated exact quote must report each source occurrence.
These cases make boundary decisions visible without depending on external data.
