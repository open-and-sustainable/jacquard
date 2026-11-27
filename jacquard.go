// Package jacquard provides contracts and JSON interface validation for
// deterministic evidence checking in AI-assisted knowledge synthesis.
//
// # Supported interfaces
//
//   - [Input] and [Output]: document checker requests and reports.
//   - [FunnelInput] and [FunnelOutput]: population size audit requests and reports.
//
// [Contract] declares a caller's extraction schema, claim selectors, and record
// rules. The types represent document-checking and population audit data.
//
// # Validation functions
//
// [ValidateInput], [ValidateOutput], [ValidateFunnelInput], and
// [ValidateFunnelOutput] validate the four JSON interface shapes. [Validate]
// selects a schema by interface name. [ValidateJSON] accepts an
// explicit, self-contained draft-07 schema. [ResultFor] converts validation
// outcomes to the JSON result used by the CLI and C exports.
//
// # Schemas and errors
//
// [Schema] returns an independent copy of an embedded schema. [SchemaTypes]
// lists interface names. Counts and amounts that can exceed
// machine-integer precision use encoding/json.Number.
//
// Validation preserves arbitrary-precision JSON numbers and sorts findings by
// JSON Pointer, code, then message. The validators return nil for a
// valid instance, *[Error] for instance findings, or an ordinary error for an
// unavailable or unusable schema. Callers can use errors.As to distinguish them.
// [ValidateJSON] returns findings and schema errors separately. File and network
// loading are disabled; references must resolve inside the supplied schema.
//
// # Example usage
//
//	issues, err := jacquard.ValidateJSON(
//		[]byte(`{"participants":3}`),
//		[]byte(`{"type":"object","properties":{"participants":{"type":"integer"}}}`),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(len(issues))
//
// All public Go types and functions are available from the root import path,
// github.com/open-and-sustainable/jacquard. Synthetic examples exercise the public
// API in the test package.
package jacquard

import (
	"github.com/open-and-sustainable/jacquard/internal/contract"
	"github.com/open-and-sustainable/jacquard/internal/validation"
)

// Document contains exactly the UTF-8 text supplied to an extractor.
type Document = contract.Document

// ClaimRule selects claim objects using an array-wildcard JSON Pointer.
// Parts maps standard part names to relative JSON Pointers within a claim.
type ClaimRule = contract.ClaimRule

// RecordRule requires evidence in at least one selected claim within each record.
// AnyEvidence contains relative selectors, including selectors for note arrays.
type RecordRule = contract.RecordRule

// Contract declares extraction structure and the checks to apply to claim objects.
// Schema is the caller's self-contained draft-07 JSON Schema.
type Contract = contract.Contract

// Input is one document, its parsed extraction, and a checking contract.
type Input = contract.Input

// Issue is a deterministic finding at a JSON Pointer in the checked value.
type Issue = contract.Issue

// FormatResult separates extraction schema failures from semantic rule findings.
type FormatResult = contract.FormatResult

// Occurrence is a zero-based, half-open span in the original document.
// Character offsets count Unicode code points, not grapheme clusters.
type Occurrence = contract.Occurrence

// QuoteDiagnosis explains a failed exact match without changing its outcome.
type QuoteDiagnosis = contract.QuoteDiagnosis

// QuoteResult locates every exact occurrence or diagnoses a failed exact match.
type QuoteResult = contract.QuoteResult

// ParsedAmount describes an optional amount reading in integer minor units.
type ParsedAmount = contract.ParsedAmount

// ClaimResult separates evidence support from contract validity.
// Status is supported, not_checked, quote_not_in_text, or value_not_in_quote.
type ClaimResult = contract.ClaimResult

// RuleIdentity records reproducible identities for the contract and language data.
type RuleIdentity = contract.RuleIdentity

// Summary counts mutually exclusive claim evidence statuses.
type Summary = contract.Summary

// Output represents a document-checking report.
type Output = contract.Output

// Measurement distinguishes a parsed nonnegative integer from missing/unreadable data.
// State is parsed, missing, or unparsed. Value is present only for parsed values.
type Measurement = contract.Measurement

// FunnelItem can contribute to several distinct size buckets within a field.
// PassedStages must be a nonempty prefix of the declared stage order.
type FunnelItem = contract.FunnelItem

// FunnelInput declares the reference population and all audit stages up front.
type FunnelInput = contract.FunnelInput

// FunnelCell counts distinct items for one stage, field, and magnitude bucket.
type FunnelCell = contract.FunnelCell

// FunnelOutput includes zero-count cells and the cutoffs used to assign deciles.
type FunnelOutput = contract.FunnelOutput

// Error contains stable, sorted schema findings for an invalid JSON instance.
// Use errors.As to distinguish these findings from schema/configuration errors.
type Error = validation.Error

// Result is the JSON result of interface validation, distinct from a checking report.
// A schema/configuration error populates Error; instance findings populate Issues.
type Result = validation.Result

// SchemaTypes returns the supported JSON interface names in a stable order.
//
// The names are input, output, funnel_input, and funnel_output. The returned slice
// is independent of later calls and may be modified by the caller.
func SchemaTypes() []string { return validation.SchemaTypes() }

// Schema returns an embedded schema without exposing mutable global state.
//
// Parameters:
//   - kind: An interface name returned by [SchemaTypes].
//
// Returns:
//   - An independent byte slice containing a self-contained draft-07 JSON Schema.
//   - An error if the kind is unsupported.
func Schema(kind string) ([]byte, error) { return validation.Schema(kind) }

// ValidateInput validates the shape of a document checker request.
//
// Parameters:
//   - input: One UTF-8 JSON request containing document, extraction, and contract.
//
// Returns nil for a valid request shape, *[Error] for instance findings, or an
// ordinary error for an unavailable schema. It does not check evidence or execute
// the extraction schema embedded in the contract.
func ValidateInput(input string) error { return validation.ValidateInput(input) }

// ValidateOutput validates the shape of a document checker report.
//
// Parameters:
//   - output: One UTF-8 JSON report, including format findings, claims, and summary.
//
// Returns nil for a valid report shape, *[Error] for instance findings, or an
// ordinary error for an unavailable schema. It does not recompute findings or
// verify summary arithmetic.
func ValidateOutput(output string) error { return validation.ValidateOutput(output) }

// ValidateFunnelInput validates the shape of a size audit request.
//
// Parameters:
//   - input: One UTF-8 JSON audit request with stages, reference values, and items.
//
// Returns nil for a valid request shape, *[Error] for instance findings, or an
// ordinary error for an unavailable schema. Population rules such as unique item
// ids and stage prefixes are the audit engine's responsibility.
func ValidateFunnelInput(input string) error { return validation.ValidateFunnelInput(input) }

// ValidateFunnelOutput validates the shape of a size audit report.
//
// Parameters:
//   - output: One UTF-8 JSON audit report with cutoffs and stage/field/bucket cells.
//
// Returns nil for a valid report shape, *[Error] for instance findings, or an
// ordinary error for an unavailable schema. It does not recompute cutoffs or counts.
func ValidateFunnelOutput(output string) error { return validation.ValidateFunnelOutput(output) }

// Validate checks an instance against one of the embedded schemas.
//
// Parameters:
//   - instance: One UTF-8 JSON value to validate.
//   - kind: input, output, funnel_input, or funnel_output.
//
// Returns nil for a valid instance, *[Error] for malformed/invalid instance JSON,
// or an ordinary error if the schema is unavailable or cannot be compiled.
// Findings are sorted by JSON Pointer, code, then message. This is interface
// shape validation, not execution of the document checker or population audit.
func Validate(instance, kind string) error { return validation.Validate(instance, kind) }

// ValidateJSON validates one UTF-8 JSON value against a self-contained draft-07
// schema, preserving arbitrary-precision numbers.
//
// Parameters:
//   - instance: Exactly one UTF-8 JSON value; malformed JSON produces an invalid_json finding.
//   - schema: A draft-07 schema with references confined to the supplied resource.
//
// Returns:
//   - A non-nil slice of findings sorted by JSON Pointer, code, then message.
//     A valid instance has an empty slice.
//   - An error if the schema cannot be decoded/compiled, names another draft,
//     or attempts to load an external resource. Findings are nil on these errors.
//
// The instance is decoded first; malformed instance JSON returns a finding before
// schema compilation. File/network loading is disabled. This function checks only
// constraints expressed in the supplied schema, not claim-rule semantics.
func ValidateJSON(instance, schema []byte) ([]Issue, error) {
	return validation.ValidateJSON(instance, schema)
}

// ResultFor converts a validation error to a CLI/FFI result without logging.
//
// A nil error produces Valid=true and an empty, non-nil Issues slice. A *[Error]
// supplies instance findings; any other error supplies the Error message.
// ResultFor does not run validation or evidence checking itself.
func ResultFor(err error) Result { return validation.ResultFor(err) }
