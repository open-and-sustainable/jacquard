---
title: Checker Contract
layout: default
---

# Checker contract

This contract defines document-checking and filter-size audit data, rules, and
report semantics.

## Operations and boundaries

The root `jacquard` package exports the request/report types and JSON validators:

- `Input` and `Output` represent a document, its extraction, and checking findings.
- `FunnelInput` and `FunnelOutput` represent a population and its size audit.

The [Go validation APIs](schemas-validation.md), [CLI](command-line.md), and
[C exports](cshared-library.md) validate these structures using the same schemas.

The library is deterministic. It performs no model calls or network requests,
does not modify source text, and checks documentary support rather than truth in
the world.

## Document checker input

A request contains `document`, `extraction`, and `contract`.

- `document`: a nonempty `document_id`, a BCP-47 `language`, and the exact UTF-8 `text`
  seen by the extractor. Empty document text is valid.
- `extraction`: any parsed JSON value; the caller's schema determines its structure.
- `contract.id`: a nonempty caller-defined rule identifier.
- `contract.schema`: a self-contained draft-07 JSON Schema. Boolean schemas are valid.
  References within that schema are supported. File and network references are rejected.
- `contract.language_table`: an explicit rule table, `en` for English.
  Unknown tables or an incompatible document language are request errors; there is no
  silent fallback. `en`, `en-US`, and other English subtags use the English table.
- `contract.claims`: selectors and claim kinds. An empty list permits schema-only checks.
- `contract.record_rules`: optional requirements for evidence within a whole record.
- `contract.quotes_path`: optional pointer to a quote dictionary, default `#/quotes`.

The [input fixture](https://github.com/open-and-sustainable/jacquard/blob/main/test/fixtures/check_input.json)
is a complete synthetic example.

### Claim selection and mapping

`path` is a JSON Pointer starting with `#`. A token consisting of `*` selects each
element of an array, in array order; it does not select object keys. For example,
`#/records/*/participants` selects the participant claim in each record. `~0` and
`~1` escape a literal tilde and slash. Other glob operators are not supported.
`*` is reserved in selectors; a property literally named `*` cannot be selected.

Each selected claim is an object or null. Missing/null claims produce `not_checked`
when their parent record exists. An absent/null intermediate array produces no
matches; require that array in the extraction schema when it must exist.
Conflicting selectors that resolve to the same claim path are invalid contracts.

The default part names are below. `parts` maps a standard part to a relative JSON
Pointer within the selected object, allowing another project's existing format:
`{"value":"#/total","evidence":"#/source/quote"}`. An unmapped part uses its default
name; an absent part has the same semantic meaning as null. A caller's schema can
still require explicit parts. `required_parts` adds non-null requirements when
evidence is present. Unknown parts or wildcard part pointers are contract errors.

| Kind | Standard parts | Evidence support test |
| --- | --- | --- |
| `text` | `as_written`, `evidence` | Text occurs inside its quote |
| `category` | `type`, `as_written`, `evidence` | Label as written occurs inside its quote |
| `count` | `value`, `upper_value`, `qualification`, optional `causal_scope`, `evidence` | Every stated bound is recognized in its quote |
| `amount` | `amount_as_written`, `currency_as_written`, `scope_as_written`, `qualification`, `evidence` | Each supplied as-written part occurs inside its quote |
| `note` | `kind`, `evidence` | Its quote is exact |

Every kind also supports `evidence_ref`. Categories and note kinds are closed lists
declared by the caller's extraction schema. Field names and category vocabularies
are never built into the core.

### Shared evidence and absent values

A claim supplies either a nonblank `evidence` string or an `evidence_ref` into the
quote dictionary. Both non-null at once are invalid, even if they resolve to equal
strings. A reference must resolve to a nonblank string; an unknown reference is an
issue, never an empty quote. The resolved quote goes through exactly the same checks
as inline evidence. Reusing a quote does not merge claim results.

Absent/null parts are equivalent for semantic checks. An explicitly blank string
is invalid. Omitting empty fields never overrides the caller's JSON Schema.

## 1. Contract validation

First validate the extraction against `contract.schema`. If it fails, return
`format.schema_valid=false`, `format.valid=false`, sorted schema issues, an empty
`claims` array, and zero summary counts. No quote or value checking runs.

If the schema passes, apply semantic rules to each selected claim:

| Kind | Required consistency |
| --- | --- |
| `text` | Null text requires null evidence/reference; nonblank text requires evidence |
| `category` | Null type requires null as-written text and evidence; a type requires nonblank as-written text and evidence |
| `count` | No evidence requires null value/bounds/qualification/scope; evidence requires qualification; quantified claims require a nonnegative integer value |
| `amount` | No evidence requires null remaining parts; evidence requires qualification; quantified claims require nonblank amount text |
| `note` | A populated note requires a nonblank kind and evidence |

Count qualifications are `exact`, `approximate`, `at_least`, `more_than`, `at_most`,
`less_than`, `range`, and `unquantified`. `range` requires `upper_value >= value`;
other qualifications prohibit `upper_value`. `unquantified` requires null bounds.
Amounts use the same qualifications except `range`; unquantified
amounts have null amount text, and quoted prose can describe the uncertainty.
Qualification labels are checked for consistency, not inferred from natural language.
`causal_scope` has no built-in vocabulary or requirement; callers declare those.

Every evidence string must be exact. Every non-null as-written part must occur
inside its own quote, using case-sensitive substring matching. A count's stated
bounds must be supported under its declared count policy (see §3).

A record rule selects records by `path` and requires nonempty evidence in at least
one relative selector in `any_evidence`. Selectors can include notes, for example
`["#/participants", "#/funding", "#/notes/*"]`. Evidence existence meets this rule;
the separate quote/value checks determine its support.

Stable semantic issue codes are:

| Code | Trigger |
| --- | --- |
| `blank_part` | A supplied string contains only whitespace |
| `missing_evidence` | A populated claim has no evidence |
| `unexpected_part` | A part must be absent/null for this claim state |
| `missing_part` | A required claim part is absent/null |
| `invalid_part` | A part has an invalid type or a qualification outside the allowed set |
| `conflicting_evidence` | Both inline evidence and a reference are supplied |
| `evidence_reference_not_found` | A reference has no valid dictionary entry |
| `invalid_range` | A range is missing its upper bound or has descending bounds |
| `quote_not_exact` | Resolved evidence is not an exact substring |
| `value_not_in_quote` | A supplied value is not supported by its quote |
| `record_without_evidence` | A record rule has no evidenced target |

Schema issues use `schema_` followed by the failed JSON Schema keyword, for example
`schema_type`, `schema_enum`, and `schema_required`; a false schema uses `schema_false`.
Claim-part type and qualification rules apply even if the caller's extraction
schema is permissive. Counts are nonnegative integers; as-written parts, labels,
scope, and evidence/reference parts are strings or null. A nonobject/non-null
selected claim has an `invalid_part` issue at its claim path and is `not_checked`
with reason `invalid_claim`.
Each issue has `code`, `path`, and a human-readable `message`. Paths identify the
offending part (or the whole record for a record rule). Codes and paths are the
compatibility interface; message wording is informative. Semantic findings do not
stop checking other claims.

## 2. Quote fidelity

An exact quote has all occurrences reported as zero-based, half-open byte spans
and Unicode code-point spans. Include overlapping occurrences and preserve source
order. Characters are code points, not grapheme clusters. The document is never
trimmed, case-folded, Unicode-normalized, or otherwise changed.

A nonexact quote has no exact occurrences and one diagnosis, chosen in this order:

1. `first_letter_case`: flipping only the first code point's case gives an exact match.
2. `ellipsis_gap`: splitting on `...` or `…`, trimming the pieces, and discarding empty
   pieces leaves at least two nonempty pieces found in order without overlap. Choose
   the earliest occurrence of each successive piece. Report source spans and skipped
   byte lengths between them.
3. `typography`: matching becomes possible after applying the language table's explicit
   curly-quote, dash, and Unicode-space mappings to quote and source.
4. `whitespace`: matching becomes possible after collapsing Unicode whitespace runs
   in both quote and source to one ASCII space, without stripping ends.
5. `other`: report the longest prefix and suffix of the quote found in the document,
   measured on code-point boundaries. Either may be empty.

Diagnosis never changes `exact=false`. Quote checking requires exact matching.

## 3. Value support

Text and category checks use the supplied `as_written` string literally. A missing
match is diagnosed as `as_written_elsewhere_in_text` if it occurs elsewhere in the
source, otherwise `as_written_not_in_text`.

Count reading reports sorted, unique arbitrary-precision nonnegative integers.
The English rules recognize digit literals with valid thousands grouping and number words
from zero through billions, including hundreds, hyphens, and conjunctions inside a
number. Decimals and ordinal tokens do not become integer counts. Count reading
uses integer arithmetic throughout, without floating-point conversion.

`count_support` defaults to `literal`. `literal_or_implied` additionally enables the
table's zero phrases and indefinite-person rules, scoped by the rule's `field_words`.
Both literal and implied readings are reported separately. Implied readings must
not be silently enabled. Count checks establish presence, not a causal relationship
between a number and a named outcome. A quote containing several numbers can support
several literal candidates; this limitation is explicit.

For amounts, all supplied as-written parts must occur inside the quote. An optional
amount reading converts a single USD amount to integer cents using decimal arithmetic,
with half-cent ties rounded away from zero. The English rules recognize `$`, `USD`, and
`dollar(s)` and the scales thousand/k, million/m/mln, billion/b/bn. Conflicting or
non-USD currency markers give `unsupported_currency`; ranges, multiple numeric
amounts, or unreadable amounts give `unparsed`. Parsing does not invent a conversion
or change evidence support for the supplied as-written text. Report its parse state.

## 4. Evidence reporting

Each selected claim has exactly one status, applied in this order:

| Status | Meaning |
| --- | --- |
| `not_checked` | Claim absent/null; no resolvable evidence; or no quantified value to test |
| `quote_not_in_text` | A resolved quote is not exact |
| `value_not_in_quote` | Quote exact, but at least one supplied value is unsupported |
| `supported` | Quote exact and every required evidence support test passes |

A populated note with an exact quote is `supported`. An unquantified count or amount
with an exact quote is `not_checked` with reason `unquantified`; a nonexact quote is
still `quote_not_in_text`. Other reasons are `absent_claim`, `missing_evidence`,
`unresolved_evidence`, `missing_value`, and `invalid_claim`. Do not count
unquantified claims as zero.

Contract validity and evidence support are separate results. A count with an exact,
supporting quote but a missing qualification can be `supported` while `format.valid`
is false. Callers can require both. A failed extraction schema produces no claim statuses.

The report contains software version, document id, contract id/hash, language-table
id/hash, format findings, claim results, and summary counts. Summary counts must sum
to the number of claim results. Required arrays are empty arrays rather than null.

Canonical output uses UTF-8 JSON with recursively sorted object keys, numbers written
in decimal without exponent notation or redundant trailing fractional zeros (negative
zero becomes zero), no insignificant whitespace, and no timestamps. Preserve
array order: claims by contract declaration then selected array index, issues by path
then code then message, and occurrences by source position. Contract and language-table
SHA-256 values hash this canonical representation of their parsed JSON, including all
declared contract fields. Identical requests and rule data produce identical bytes.

## 5. Filter-size audit

A size audit uses declared stage names in order, `reference_values` per caller-defined
field, and items with `item_id`, per-field measurements, and `passed_stages`.
Reference values are nonnegative integers; zero values are excluded from cutoff
calculation. Each measurement is explicitly `parsed` with an integer `value`, `missing`,
or `unparsed`. Multiple measurements allow one item to link to several source values.

Stages and item ids must be unique. An item's passed stages must be a nonempty prefix
of the declared order. Unknown stages/fields, negative/fractional measurements, and
inconsistent measurement states are request errors. A missing field or empty measurement
list contributes to `missing`.

For each field, sort all positive reference values. For n positives, cutoff k is
the value at one-based position `ceil(k*n/10)`, for k=1…10. A positive measurement
belongs to the first cutoff greater than or equal to it; values beyond the last
cutoff belong to decile 10. Ties are never split. If there are no positive references,
return an empty cutoff list; a positive item measurement in that field is a request
error because no decile can be assigned.

The 13 buckets, in order, are `missing`, `unparsed`, `zero`, and
`positive_decile_01` through `positive_decile_10`. Each item contributes at most once
per distinct stage/field/bucket, even if several linked values fall in that bucket.
An item can contribute to several buckets, so bucket counts need not sum to item count.

Return every stage × field × bucket cell, including zeros, together with cutoffs.
Cell order is declared stage order, lexicographically sorted fields, then the bucket
order above. The audit reports counts; it does not declare statistical bias or impose
a pass/fail threshold.

## Errors, versioning, and batch format

Malformed requests, invalid contracts, unsupported table versions, and invalid audit
populations are operation errors. An extraction that violates its schema or has
unsupported evidence is a successful operation returning a negative report. Do not
mix these two error channels. JSON interface validation alone does not verify semantic
request rules such as stage prefixes or selector overlaps.

The batch format uses one JSON Lines envelope per input record, in order:
`{"ok":true,"report":...}`
or `{"ok":false,"error":{"code":...,"message":...}}`. Blank lines are invalid records.
An invalid line does not discard other records. Operation error codes are
`invalid_json`, `invalid_request`, `invalid_contract`, `unsupported_language_table`,
and `internal_error`. The command and C-library guides describe their validation
result formats, invocation errors, and memory ownership.

Behavior changes update the contract, schemas, tests, and changelog. Content hashes
identify the contract and language rules. Go, CLI, and C validation use the same schemas and finding rules.
The [testing guide](testing.md) describes the test suite and contract scenarios.
