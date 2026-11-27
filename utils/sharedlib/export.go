package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"encoding/json"
	"unsafe"

	"github.com/open-and-sustainable/jacquard"
)

func validate(input, version *C.char, kind string) (result *C.char) {
	defer func() {
		if recover() != nil {
			result = C.CString(`{"valid":false,"issues":[],"error":"internal validation error"}`)
		}
	}()
	if input == nil || version == nil {
		return C.CString(`{"valid":false,"issues":[],"error":"input and version must not be NULL"}`)
	}
	err := jacquard.Validate(C.GoString(input), C.GoString(version), kind)
	data, marshalErr := json.Marshal(jacquard.ResultFor(err))
	if marshalErr != nil {
		return C.CString(`{"valid":false,"issues":[],"error":"cannot encode validation result"}`)
	}
	return C.CString(string(data))
}

// RunValidationInput validates a document checker request against an interface schema.
//
// Parameters:
//   - input: One NUL-terminated UTF-8 JSON request.
//   - version: A NUL-terminated interface version string, v1.
//
// Returns an owned UTF-8 JSON result with valid, issues, and an optional error.
// NULL parameters produce a JSON error. Free the result with [FreeCString].
// This checks the request shape; it does not check evidence.
//
//export RunValidationInput
func RunValidationInput(input, version *C.char) *C.char { return validate(input, version, "input") }

// RunValidationOutput validates a document checker report against an interface schema.
//
// Parameters:
//   - input: One NUL-terminated UTF-8 JSON report.
//   - version: A NUL-terminated interface version string, v1.
//
// Returns an owned UTF-8 JSON result with valid, issues, and an optional error.
// NULL parameters produce a JSON error. Free the result with [FreeCString].
// This checks the report shape; it does not recompute evidence findings.
//
//export RunValidationOutput
func RunValidationOutput(input, version *C.char) *C.char { return validate(input, version, "output") }

// RunValidationFunnelInput validates a size audit request against an interface schema.
//
// Parameters:
//   - input: One NUL-terminated UTF-8 JSON audit request.
//   - version: A NUL-terminated interface version string, v1.
//
// Returns an owned UTF-8 JSON result with valid, issues, and an optional error.
// NULL parameters produce a JSON error. Free the result with [FreeCString].
// This checks the request shape; it does not run the audit.
//
//export RunValidationFunnelInput
func RunValidationFunnelInput(input, version *C.char) *C.char {
	return validate(input, version, "funnel_input")
}

// RunValidationFunnelOutput validates a size audit report against an interface schema.
//
// Parameters:
//   - input: One NUL-terminated UTF-8 JSON audit report.
//   - version: A NUL-terminated interface version string, v1.
//
// Returns an owned UTF-8 JSON result with valid, issues, and an optional error.
// NULL parameters produce a JSON error. Free the result with [FreeCString].
// This checks the report shape; it does not recompute cutoffs or counts.
//
//export RunValidationFunnelOutput
func RunValidationFunnelOutput(input, version *C.char) *C.char {
	return validate(input, version, "funnel_output")
}

// FreeCString releases memory returned by a RunValidation function. NULL is safe.
//
// Call this exactly once for each non-NULL result while the library is loaded.
// Do not pass an input buffer or a pointer allocated by another library.
//
//export FreeCString
func FreeCString(value *C.char) { C.free(unsafe.Pointer(value)) }

func main() {}
