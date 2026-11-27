// Package test contains synthetic cross-package interface tests for Jacquard.
//
// # Test coverage
//
// The Go tests validate request/report fixture shapes, typed JSON round trips,
// measurement-state boundaries, and a synthetic extraction's declared schema.
// The separate shared_library_test.py script exercises the C exports, including
// NULL handling and allocation/free cycles, using the same synthetic fixtures.
//
// # Fixtures and requirements
//
// All documents, extractions, and populations are invented in this repository.
// Tests need no credentials, models, external datasets, or network calls after
// Go dependencies have been downloaded. C-interface tests additionally require
// a C compiler, a built shared library, and Python 3.
//
// Hand-authored report fixtures exercise JSON structure and typed round trips.
//
// # Example usage
//
//	go test ./...
//	go build -buildmode=c-shared -o dist/jacquard.so ./utils/sharedlib
//	python3 test/shared_library_test.py dist/jacquard.so
package test
