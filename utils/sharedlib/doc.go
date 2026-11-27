// Package main provides C-compatible exports for Jacquard's interface validators.
//
// # Exported functions
//
// [RunValidationInput] and [RunValidationOutput] validate document checker
// requests and reports. [RunValidationFunnelInput] and [RunValidationFunnelOutput]
// validate population audit requests and reports. These functions use the same
// root jacquard package as native Go callers and the CLI. They check interface
// shapes; they do not perform evidence checking or compute an audit.
//
// # Parameters and results
//
// Each validator accepts NUL-terminated UTF-8 input and schema-version strings.
// The supported version is v1. A returned JSON object contains valid and issues,
// plus error for an invalid invocation or unavailable schema. NULL inputs return
// a JSON error. The library prints no diagnostics to stdout or stderr.
//
// # Memory ownership
//
// Callers retain ownership of input buffers. Each returned string is allocated
// by the library and must be freed exactly once with [FreeCString] while the
// library remains loaded. FreeCString accepts NULL.
//
// # Build and example usage
//
//	go build -buildmode=c-shared -o dist/jacquard.so ./utils/sharedlib
//
// Go also generates the corresponding C header. A C caller can use it as follows:
//
//	#include <stdio.h>
//	#include "jacquard.h"
//
//	int main(void) {
//	    char *result = RunValidationInput("{}", "v1");
//	    puts(result);
//	    FreeCString(result);
//	    return 0;
//	}
package main
