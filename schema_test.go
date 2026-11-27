package jacquard_test

import (
	"encoding/json"
	"testing"

	"github.com/open-and-sustainable/jacquard"
)

func TestEmbeddedSchemas(t *testing.T) {
	for _, kind := range jacquard.SchemaTypes() {
		t.Run(kind, func(t *testing.T) {
			data, err := jacquard.Schema("v1", kind)
			if err != nil || !json.Valid(data) {
				t.Fatalf("invalid embedded schema: %s, %v", data, err)
			}
			data[0] = '!'
			fresh, err := jacquard.Schema("v1", kind)
			if err != nil || !json.Valid(fresh) {
				t.Fatal("mutating a returned schema changed embedded data")
			}
		})
	}
}

func TestSchemaVersionAndTypeAreExact(t *testing.T) {
	for _, pair := range [][2]string{{"v2", "input"}, {"1.0", "input"}, {"v1", "unknown"}, {"../v1", "input"}} {
		if _, err := jacquard.Schema(pair[0], pair[1]); err == nil {
			t.Fatalf("accepted unsupported schema %v", pair)
		}
	}
}
