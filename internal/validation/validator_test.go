package validation_test

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/open-and-sustainable/jacquard/internal/validation"
)

func TestExplicitSchemaBoundaries(t *testing.T) {
	tests := []struct {
		name, input, schema, code, path string
	}{
		{"valid", `{"n":3}`, `{"type":"object","properties":{"n":{"type":"integer"}}}`, "", ""},
		{"fractional integer", `2.5`, `{"type":"integer"}`, "schema_type", "#"},
		{"large integer preserved", `9007199254740993`, `{"enum":[9007199254740993]}`, "", ""},
		{"nearby integer distinguished", `9007199254740992`, `{"enum":[9007199254740993]}`, "schema_enum", "#"},
		{"large const distinguished", `9007199254740992`, `{"const":9007199254740993}`, "schema_const", "#"},
		{"large minimum distinguished", `9007199254740992`, `{"minimum":9007199254740993}`, "schema_minimum", "#"},
		{"arbitrary precision", `1000000000000000000000000000000000000`, `{"type":"integer","minimum":0}`, "", ""},
		{"missing property", `{}`, `{"required":["value"]}`, "schema_required", "#/value"},
		{"escaped pointer", `{"a.b/~/":{"count":"x"}}`, `{"properties":{"a.b/~/":{"properties":{"count":{"type":"integer"}}}}}`, "schema_type", "#/a.b~1~0~1/count"},
		{"null character key", `{"a\u0000b":"x"}`, `{"properties":{"a\u0000b":{"type":"integer"}}}`, "schema_type", "#/a\x00b"},
		{"local reference", `3`, `{"definitions":{"n":{"type":"integer"}},"$ref":"#/definitions/n"}`, "", ""},
		{"boolean false schema", `{}`, `false`, "schema_false", "#"},
		{"negated schema", `{}`, `{"not":{}}`, "schema_not", "#"},
		{"draft URI without fragment", `3`, `{"$schema":"https://json-schema.org/draft-07/schema","type":"integer"}`, "", ""},
		{"trailing JSON", `{} {}`, `{}`, "invalid_json", "#"},
		{"empty input", ``, `{}`, "invalid_json", "#"},
		{"malformed input", `{"value":`, `{}`, "invalid_json", "#"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := validation.ValidateJSON([]byte(tt.input), []byte(tt.schema))
			if err != nil {
				t.Fatal(err)
			}
			if tt.code == "" {
				if issues == nil || len(issues) != 0 {
					t.Fatalf("want empty non-nil issues, got %#v", issues)
				}
				return
			}
			if len(issues) != 1 || issues[0].Code != tt.code || issues[0].Path != tt.path {
				t.Fatalf("want %s at %q, got %+v", tt.code, tt.path, issues)
			}
		})
	}
}

func TestInvalidUTF8IsRejected(t *testing.T) {
	issues, err := validation.ValidateJSON([]byte{'"', 0xff, '"'}, []byte(`{}`))
	if err != nil || len(issues) != 1 || issues[0].Code != "invalid_json" {
		t.Fatalf("invalid UTF-8 accepted: %+v, %v", issues, err)
	}
}

func TestSchemaCannotLoadExternalResources(t *testing.T) {
	for _, reference := range []string{"https://invalid.example/schema.json", "file:///etc/passwd"} {
		schema, err := json.Marshal(map[string]string{"$ref": reference})
		if err != nil {
			t.Fatal(err)
		}
		_, err = validation.ValidateJSON([]byte(`{}`), schema)
		if err == nil || !strings.Contains(err.Error(), "external schema references are disabled") {
			t.Fatalf("expected loader refusal for %s, got %v", reference, err)
		}
	}
}

func TestMalformedSchemaIsAnError(t *testing.T) {
	for _, schema := range []string{`{"type":`, `{"type":"not-a-type"}`, `{"$schema":"https://json-schema.org/draft/2020-12/schema"}`} {
		if _, err := validation.ValidateJSON([]byte(`{}`), []byte(schema)); err == nil {
			t.Fatalf("accepted invalid schema %s", schema)
		}
	}
}

func TestFindingOrderIsDeterministic(t *testing.T) {
	schema := []byte(`{"properties":{"z":{"type":"integer"},"a":{"type":"integer"}},"required":["missing"]}`)
	var previous string
	for i := 0; i < 30; i++ {
		issues, err := validation.ValidateJSON([]byte(`{"z":"bad","a":"bad"}`), schema)
		if err != nil {
			t.Fatal(err)
		}
		paths := []string{issues[0].Path, issues[1].Path, issues[2].Path}
		if !reflect.DeepEqual(paths, []string{"#/a", "#/missing", "#/z"}) {
			t.Fatalf("unsorted findings: %+v", issues)
		}
		encoded, _ := json.Marshal(issues)
		if i > 0 && string(encoded) != previous {
			t.Fatal("findings changed between identical requests")
		}
		previous = string(encoded)
	}
}

func TestTypedErrorsAndJSONResults(t *testing.T) {
	err := validation.ValidateInput(`{}`)
	var findings *validation.Error
	if !errors.As(err, &findings) || len(findings.Issues) == 0 || !strings.Contains(err.Error(), "validation errors") {
		t.Fatalf("want validation findings, got %v", err)
	}
	result := validation.ResultFor(err)
	if result.Valid || len(result.Issues) == 0 || result.Error != "" {
		t.Fatalf("invalid finding result: %+v", result)
	}
	result = validation.ResultFor(validation.Validate(`{}`, "unknown"))
	if result.Valid || result.Error == "" || len(result.Issues) != 0 {
		t.Fatalf("invalid configuration error result: %+v", result)
	}
	encoded, _ := json.Marshal(validation.ResultFor(nil))
	if string(encoded) != `{"valid":true,"issues":[]}` {
		t.Fatalf("unexpected success shape: %s", encoded)
	}
}
