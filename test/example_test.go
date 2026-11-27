package test

import (
	"errors"
	"fmt"
	"log"

	"github.com/open-and-sustainable/jacquard"
)

func Example_validateInput() {
	input := `{
		"document": {
			"document_id": "synthetic-example",
			"language": "en",
			"text": "Three volunteers joined."
		},
		"extraction": {"participants": 3},
		"contract": {
			"id": "synthetic-example",
			"language_table": "en",
			"schema": {"type": "object"},
			"claims": []
		}
	}`
	if err := jacquard.ValidateInput(input); err != nil {
		log.Fatal(err)
	}
	// This confirms the request shape; it does not check the claim's evidence.
	fmt.Println("valid request shape")
	// Output: valid request shape
}

func Example_validateJSON() {
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

func Example_error() {
	err := jacquard.ValidateInput(`{}`)
	var findings *jacquard.Error
	if errors.As(err, &findings) {
		fmt.Println("instance findings:", len(findings.Issues))
	}
	// Output: instance findings: 3
}
