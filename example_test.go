package jacquard_test

import (
	"errors"
	"fmt"
	"log"

	"github.com/open-and-sustainable/jacquard"
)

func ExampleValidateInput() {
	input := `{
		"metadata": {"schemaVersion": "v1"},
		"document": {
			"document_id": "synthetic-example",
			"language": "en",
			"text": "Three volunteers joined."
		},
		"extraction": {"participants": 3},
		"contract": {
			"id": "synthetic-v1",
			"language_table": "en-v1",
			"schema": {"type": "object"},
			"claims": []
		}
	}`
	if err := jacquard.ValidateInput(input, "v1"); err != nil {
		log.Fatal(err)
	}
	// This confirms the request shape; it does not check the claim's evidence.
	fmt.Println("valid request shape")
	// Output: valid request shape
}

func ExampleValidateJSON() {
	issues, err := jacquard.ValidateJSON(
		[]byte(`{"participants":"three"}`),
		[]byte(`{"type":"object","properties":{"participants":{"type":"integer"}}}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	for _, issue := range issues {
		fmt.Println(issue.Code, issue.Path)
	}
	// Output: schema_type #/participants
}

func ExampleError() {
	err := jacquard.ValidateInput(`{}`, "v1")
	var findings *jacquard.Error
	if errors.As(err, &findings) {
		fmt.Println("instance findings:", len(findings.Issues))
	}
	// Output: instance findings: 4
}
