package jacquard

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Error contains stable, sorted schema findings for an invalid JSON instance.
// Use errors.As to distinguish these findings from schema/configuration errors.
type Error struct {
	// Issues are sorted by path, code, then message.
	Issues []Issue
}

// Error returns the findings as a human-readable message in their sorted order.
func (e *Error) Error() string {
	messages := make([]string, len(e.Issues))
	for i, issue := range e.Issues {
		messages[i] = issue.Path + ": " + issue.Message
	}
	return "validation errors: " + strings.Join(messages, "; ")
}

// ValidateInput validates the shape of a document checker request.
//
// Parameters:
//   - input: One UTF-8 JSON request containing metadata, document, extraction, and contract.
//   - version: The exact interface version, v1.
//
// Returns nil for a valid request shape, *[Error] for instance findings, or an
// ordinary error for an unavailable schema. It does not check evidence or execute
// the extraction schema embedded in the contract.
func ValidateInput(input, version string) error {
	return Validate(input, version, "input")
}

// ValidateOutput validates the shape of a document checker report.
//
// Parameters:
//   - output: One UTF-8 JSON report, including format findings, claims, and summary.
//   - version: The exact interface version, v1.
//
// Returns nil for a valid report shape, *[Error] for instance findings, or an
// ordinary error for an unavailable schema. It does not recompute findings or
// verify summary arithmetic.
func ValidateOutput(output, version string) error {
	return Validate(output, version, "output")
}

// ValidateFunnelInput validates the shape of a size audit request.
//
// Parameters:
//   - input: One UTF-8 JSON audit request with stages, reference values, and items.
//   - version: The exact interface version, v1.
//
// Returns nil for a valid request shape, *[Error] for instance findings, or an
// ordinary error for an unavailable schema. Population rules such as unique item
// ids and stage prefixes are the audit engine's responsibility.
func ValidateFunnelInput(input, version string) error {
	return Validate(input, version, "funnel_input")
}

// ValidateFunnelOutput validates the shape of a size audit report.
//
// Parameters:
//   - output: One UTF-8 JSON audit report with cutoffs and stage/field/bucket cells.
//   - version: The exact interface version, v1.
//
// Returns nil for a valid report shape, *[Error] for instance findings, or an
// ordinary error for an unavailable schema. It does not recompute cutoffs or counts.
func ValidateFunnelOutput(output, version string) error {
	return Validate(output, version, "funnel_output")
}

// Validate checks an instance against one of the embedded versioned schemas.
//
// Parameters:
//   - instance: One UTF-8 JSON value to validate.
//   - version: The exact interface version, v1.
//   - kind: input, output, funnel_input, or funnel_output.
//
// Returns nil for a valid instance, *[Error] for malformed/invalid instance JSON,
// or an ordinary error if the schema is unavailable or cannot be compiled.
// Findings are sorted by JSON Pointer, code, then message. This is interface
// shape validation, not execution of the document checker or population audit.
func Validate(instance, version, kind string) error {
	schema, err := Schema(version, kind)
	if err != nil {
		return err
	}
	issues, err := ValidateJSON([]byte(instance), schema)
	if err != nil {
		return err
	}
	if len(issues) != 0 {
		return &Error{Issues: issues}
	}
	return nil
}

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
	value, err := decode(instance)
	if err != nil {
		return []Issue{{Code: "invalid_json", Path: "#", Message: "input must be one UTF-8 JSON value"}}, nil
	}
	schemaValue, err := decode(schema)
	if err != nil {
		return nil, fmt.Errorf("invalid schema JSON: %w", err)
	}
	if object, ok := schemaValue.(map[string]any); ok {
		if declared, exists := object["$schema"]; exists {
			draft, ok := declared.(string)
			draft = strings.TrimSuffix(draft, "#")
			if !ok || (draft != "http://json-schema.org/draft-07/schema" && draft != "https://json-schema.org/draft-07/schema") {
				return nil, errors.New("only draft-07 schemas are supported")
			}
		}
	}
	compiler := jsonschema.NewCompiler()
	compiler.DefaultDraft(jsonschema.Draft7)
	compiler.UseLoader(blockedLoader{})
	const resource = "urn:jacquard:validation"
	if err := compiler.AddResource(resource, schemaValue); err != nil {
		return nil, fmt.Errorf("load schema: %w", err)
	}
	compiled, err := compiler.Compile(resource)
	if err != nil {
		return nil, fmt.Errorf("compile schema: %w", err)
	}
	err = compiled.Validate(value)
	if err == nil {
		return []Issue{}, nil
	}
	var invalid *jsonschema.ValidationError
	if !errors.As(err, &invalid) {
		return nil, fmt.Errorf("validate JSON: %w", err)
	}
	issues := []Issue{}
	printer := message.NewPrinter(language.English)
	var collect func(*jsonschema.ValidationError)
	collect = func(finding *jsonschema.ValidationError) {
		if len(finding.Causes) != 0 {
			for _, cause := range finding.Causes {
				collect(cause)
			}
			return
		}
		path := pointer(finding.InstanceLocation)
		code := "schema_validation"
		if keyword := finding.ErrorKind.KeywordPath(); len(keyword) > 0 {
			code = "schema_" + keyword[0]
		}
		switch detail := finding.ErrorKind.(type) {
		case *kind.FalseSchema:
			issues = append(issues, Issue{Code: "schema_false", Path: path, Message: "false schema"})
		case *kind.Not:
			issues = append(issues, Issue{Code: "schema_not", Path: path, Message: "not constraint failed"})
		case *kind.Required:
			for _, property := range detail.Missing {
				issues = append(issues, Issue{Code: code, Path: path + "/" + escape(property), Message: fmt.Sprintf("missing property %q", property)})
			}
		case *kind.AdditionalProperties:
			for _, property := range detail.Properties {
				issues = append(issues, Issue{Code: code, Path: path + "/" + escape(property), Message: fmt.Sprintf("additional property %q is not allowed", property)})
			}
		default:
			issues = append(issues, Issue{Code: code, Path: path, Message: finding.ErrorKind.LocalizedString(printer)})
		}
	}
	collect(invalid)
	sort.Slice(issues, func(i, j int) bool {
		a, b := issues[i], issues[j]
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		if a.Code != b.Code {
			return a.Code < b.Code
		}
		return a.Message < b.Message
	})
	return issues, nil
}

func decode(data []byte) (any, error) {
	if !utf8.Valid(data) {
		return nil, errors.New("invalid UTF-8")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, errors.New("expected one JSON value")
	}
	return value, nil
}

func pointer(tokens []string) string {
	path := "#"
	for _, token := range tokens {
		path += "/" + escape(token)
	}
	return path
}

func escape(token string) string {
	return strings.ReplaceAll(strings.ReplaceAll(token, "~", "~0"), "/", "~1")
}

type blockedLoader struct{}

func (blockedLoader) Load(source string) (any, error) {
	return nil, fmt.Errorf("external schema references are disabled: %s", source)
}
