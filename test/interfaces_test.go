package test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/open-and-sustainable/jacquard"
)

func TestSyntheticInterfaceFixtures(t *testing.T) {
	for _, tt := range []struct {
		file, kind string
		target     any
	}{
		{"check_input", "input", &jacquard.Input{}},
		{"check_output", "output", &jacquard.Output{}},
		{"funnel_input", "funnel_input", &jacquard.FunnelInput{}},
		{"funnel_output", "funnel_output", &jacquard.FunnelOutput{}},
	} {
		t.Run(tt.file, func(t *testing.T) {
			data, err := os.ReadFile("fixtures/" + tt.file + ".json")
			if err != nil {
				t.Fatal(err)
			}
			if err := jacquard.Validate(string(data), tt.kind); err != nil {
				t.Fatalf("synthetic fixture is not valid: %v", err)
			}
			if err := json.Unmarshal(data, tt.target); err != nil {
				t.Fatal(err)
			}
			roundtrip, err := json.Marshal(tt.target)
			if err != nil {
				t.Fatal(err)
			}
			if err := jacquard.Validate(string(roundtrip), tt.kind); err != nil {
				t.Fatalf("typed round trip violates the interface: %v", err)
			}
		})
	}
}

func TestFunnelMeasurementStates(t *testing.T) {
	for _, measurement := range []string{
		`{"state":"parsed","value":0}`,
		`{"state":"parsed","value":123456789012345678901234567890}`,
		`{"state":"missing"}`,
		`{"state":"unparsed"}`,
	} {
		if err := jacquard.ValidateFunnelInput(funnelWith(measurement)); err != nil {
			t.Fatalf("valid measurement rejected: %s: %v", measurement, err)
		}
	}
	for _, measurement := range []string{
		`{"state":"parsed"}`,
		`{"state":"parsed","value":-1}`,
		`{"state":"parsed","value":1.5}`,
		`{"state":"missing","value":0}`,
		`{"state":"unparsed","value":null}`,
		`{"state":"guessed"}`,
	} {
		if err := jacquard.ValidateFunnelInput(funnelWith(measurement)); err == nil {
			t.Fatalf("invalid measurement accepted: %s", measurement)
		}
	}
}

func funnelWith(measurement string) string {
	return `{"stages":["input"],"reference_values":{"size":[1]},"items":[{"item_id":"x","values":{"size":[` + measurement + `]},"passed_stages":["input"]}]}`
}

func TestEmbeddedExtractionSchemaCanBeUsedExplicitly(t *testing.T) {
	data, err := os.ReadFile("fixtures/check_input.json")
	if err != nil {
		t.Fatal(err)
	}
	var input jacquard.Input
	if err := json.Unmarshal(data, &input); err != nil {
		t.Fatal(err)
	}
	issues, err := jacquard.ValidateJSON(input.Extraction, input.Contract.Schema)
	if err != nil || len(issues) != 0 {
		t.Fatalf("synthetic extraction doesn't follow its schema: %+v, %v", issues, err)
	}
	issues, err = jacquard.ValidateJSON([]byte(`{"records":[{"participants":{"value":"three"}}]}`), input.Contract.Schema)
	if err != nil || len(issues) == 0 {
		t.Fatalf("schema accepted a string as a count: %+v, %v", issues, err)
	}
}
