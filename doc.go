// Package jacquard provides versioned contracts and JSON interface validation for
// deterministic evidence checking of extracted claims.
//
// # Supported interfaces
//
//   - [Input] and [Output]: document checker requests and reports.
//   - [FunnelInput] and [FunnelOutput]: population size audit requests and reports.
//
// [Contract] declares a caller's extraction schema, claim selectors, and record
// rules. These types describe the first-release interfaces; the evidence checker
// and population audit engines are not yet implemented.
//
// # Validation functions
//
// [ValidateInput], [ValidateOutput], [ValidateFunnelInput], and
// [ValidateFunnelOutput] validate the four JSON interface shapes. [Validate]
// selects a schema by version and interface name. [ValidateJSON] accepts an
// explicit, self-contained draft-07 schema. [ResultFor] converts validation
// outcomes to the JSON result used by the CLI and C exports.
//
// # Schemas and errors
//
// [Schema] returns an independent copy of an embedded schema. [SchemaTypes]
// lists interface names, and [SchemaVersion] identifies the current JSON format
// independently of software release versions. Counts and amounts that can exceed
// machine-integer precision use encoding/json.Number.
//
// Validation preserves arbitrary-precision JSON numbers and sorts findings by
// JSON Pointer, code, then message. The versioned validators return nil for a
// valid instance, *[Error] for instance findings, or an ordinary error for an
// unavailable or unusable schema. Callers can use errors.As to distinguish them.
// [ValidateJSON] returns findings and schema errors separately. File and network
// loading are disabled; references must resolve inside the supplied schema.
//
// # Example usage
//
//	issues, err := jacquard.ValidateJSON(
//		[]byte(`{"participants":3}`),
//		[]byte(`{"type":"object","properties":{"participants":{"type":"integer"}}}`),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(len(issues))
//
// All public Go types and functions are available from the root import path,
// github.com/open-and-sustainable/jacquard. Executable synthetic examples
// accompany the functions in the Go API reference.
package jacquard
