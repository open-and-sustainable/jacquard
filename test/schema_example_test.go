package test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/open-and-sustainable/jacquard"
)

func Example_schema() {
	schema, err := jacquard.Schema("input")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(json.Valid(schema))
	// Output: true
}

func Example_schemaTypes() {
	fmt.Println(jacquard.SchemaTypes())
	// Output: [input output funnel_input funnel_output]
}
