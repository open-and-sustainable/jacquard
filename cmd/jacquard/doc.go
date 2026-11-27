// Command jacquard validates and displays Jacquard's versioned JSON interfaces.
//
// # Commands
//
// The validate command reads one JSON value from stdin and writes a JSON result
// containing valid, issues, and an optional configuration error. The schema
// command writes an embedded draft-07 schema. Both accept -type and -version;
// the defaults are input and v1. Supported types are input, output, funnel_input,
// and funnel_output.
//
// # Exit status
//
// Exit status 0 indicates successful validation or schema display, 1 indicates an
// invalid JSON instance, and 2 indicates an invocation, configuration, or I/O error.
// Usage diagnostics go to stderr. Validation checks request and report structures.
//
// # Example usage
//
//	jacquard validate -type input < test/fixtures/check_input.json
//	jacquard schema -type funnel_output -version v1
package main
