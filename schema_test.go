package jacquard_test

import (
	"encoding/json"
	"testing"

	"github.com/open-and-sustainable/jacquard"
)

func TestEmbeddedSchemas(t *testing.T) {
	for _, kind := range jacquard.SchemaTypes() {
		t.Run(kind, func(t *testing.T) {
			data, err := jacquard.Schema(kind)
			if err != nil || !json.Valid(data) {
				t.Fatalf("invalid embedded schema: %s, %v", data, err)
			}
			data[0] = '!'
			fresh, err := jacquard.Schema(kind)
			if err != nil || !json.Valid(fresh) {
				t.Fatal("mutating a returned schema changed embedded data")
			}
		})
	}
}

func TestUnsupportedSchemaTypes(t *testing.T) {
	for _, kind := range []string{"unknown", "../input", "input/other"} {
		if _, err := jacquard.Schema(kind); err == nil {
			t.Fatalf("accepted unsupported schema type %q", kind)
		}
	}
}
