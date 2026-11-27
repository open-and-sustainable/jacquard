// Package contract defines document, claim, evidence, and population audit data.
package contract

import "encoding/json"

// Document contains exactly the UTF-8 text supplied to an extractor.
type Document struct {
	ID       string `json:"document_id"`
	Language string `json:"language"`
	Text     string `json:"text"`
}

// ClaimRule selects claim objects using an array-wildcard JSON Pointer.
// Parts maps standard part names to relative JSON Pointers within a claim.
type ClaimRule struct {
	Path          string            `json:"path"`
	Kind          string            `json:"kind"`
	Parts         map[string]string `json:"parts,omitempty"`
	RequiredParts []string          `json:"required_parts,omitempty"`
	FieldWords    []string          `json:"field_words,omitempty"`
	CountSupport  string            `json:"count_support,omitempty"`
}

// RecordRule requires evidence in at least one selected claim within each record.
// AnyEvidence contains relative selectors, including selectors for note arrays.
type RecordRule struct {
	Path        string   `json:"path"`
	AnyEvidence []string `json:"any_evidence"`
}

// Contract declares extraction structure and the checks to apply to claim objects.
// Schema is the caller's self-contained draft-07 JSON Schema.
type Contract struct {
	ID            string          `json:"id"`
	LanguageTable string          `json:"language_table"`
	Schema        json.RawMessage `json:"schema"`
	Claims        []ClaimRule     `json:"claims"`
	RecordRules   []RecordRule    `json:"record_rules,omitempty"`
	QuotesPath    string          `json:"quotes_path,omitempty"`
}

// Input is one document, its parsed extraction, and a checking contract.
type Input struct {
	Document   Document        `json:"document"`
	Extraction json.RawMessage `json:"extraction"`
	Contract   Contract        `json:"contract"`
}

// Issue is a deterministic finding at a JSON Pointer in the checked value.
type Issue struct {
	Code    string `json:"code"`
	Path    string `json:"path"`
	Message string `json:"message"`
}

// FormatResult separates extraction schema failures from semantic rule findings.
type FormatResult struct {
	SchemaValid bool    `json:"schema_valid"`
	Valid       bool    `json:"valid"`
	Issues      []Issue `json:"issues"`
}

// Occurrence is a zero-based, half-open span in the original document.
// Character offsets count Unicode code points, not grapheme clusters.
type Occurrence struct {
	ByteStart int `json:"byte_start"`
	ByteEnd   int `json:"byte_end"`
	CharStart int `json:"char_start"`
	CharEnd   int `json:"char_end"`
}

// QuoteDiagnosis explains a failed exact match without changing its outcome.
type QuoteDiagnosis struct {
	Kind   string       `json:"kind"`
	Pieces []Occurrence `json:"pieces,omitempty"`
	Gaps   []int        `json:"gaps,omitempty"`
	Prefix string       `json:"prefix,omitempty"`
	Suffix string       `json:"suffix,omitempty"`
}

// QuoteResult locates every exact occurrence or diagnoses a failed exact match.
type QuoteResult struct {
	Exact       bool            `json:"exact"`
	Occurrences []Occurrence    `json:"occurrences"`
	Diagnosis   *QuoteDiagnosis `json:"diagnosis,omitempty"`
}

// ParsedAmount describes an optional amount reading in integer minor units.
type ParsedAmount struct {
	Value    json.Number `json:"value"`
	Currency string      `json:"currency"`
	Unit     string      `json:"unit"`
}

// ClaimResult separates evidence support from contract validity.
// Status is supported, not_checked, quote_not_in_text, or value_not_in_quote.
type ClaimResult struct {
	Path             string        `json:"path"`
	Kind             string        `json:"kind"`
	Status           string        `json:"status"`
	Reason           string        `json:"reason,omitempty"`
	Quote            *QuoteResult  `json:"quote,omitempty"`
	LiteralCounts    []json.Number `json:"literal_counts,omitempty"`
	ImpliedCounts    []json.Number `json:"implied_counts,omitempty"`
	ParsedAmount     *ParsedAmount `json:"parsed_amount,omitempty"`
	AmountParseState string        `json:"amount_parse_state,omitempty"`
	ValueDiagnosis   string        `json:"value_diagnosis,omitempty"`
}

// RuleIdentity records reproducible identities for the contract and language data.
type RuleIdentity struct {
	ContractID          string `json:"contract_id"`
	ContractSHA256      string `json:"contract_sha256"`
	LanguageTable       string `json:"language_table"`
	LanguageTableSHA256 string `json:"language_table_sha256"`
}

// Summary counts mutually exclusive claim evidence statuses.
type Summary struct {
	Claims          int `json:"claims"`
	Supported       int `json:"supported"`
	NotChecked      int `json:"not_checked"`
	QuoteNotInText  int `json:"quote_not_in_text"`
	ValueNotInQuote int `json:"value_not_in_quote"`
}

// Output represents a document-checking report.
type Output struct {
	JacquardVersion string        `json:"jacquard_version"`
	DocumentID      string        `json:"document_id"`
	Rules           RuleIdentity  `json:"rules"`
	Format          FormatResult  `json:"format"`
	Claims          []ClaimResult `json:"claims"`
	Summary         Summary       `json:"summary"`
}
