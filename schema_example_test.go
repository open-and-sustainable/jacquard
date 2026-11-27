package jacquard_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/open-and-sustainable/jacquard"
)

func ExampleSchema() {
	schema, err := jacquard.Schema("v1", "input")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(json.Valid(schema))
	// Output: true
}

func ExampleSchemaTypes() {
	fmt.Println(jacquard.SchemaTypes())
	// Output: [input output funnel_input funnel_output]
}
