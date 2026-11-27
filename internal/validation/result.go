package validation

import (
	"errors"
	"github.com/open-and-sustainable/jacquard/internal/contract"
)

// Result is the JSON result of interface validation, distinct from a checking report.
// A schema/configuration error populates Error; instance findings populate Issues.
type Result struct {
	Valid  bool             `json:"valid"`
	Issues []contract.Issue `json:"issues"`
	Error  string           `json:"error,omitempty"`
}

// ResultFor converts a validation error to a CLI/FFI result without logging.
//
// A nil error produces Valid=true and an empty, non-nil Issues slice. A *[Error]
// supplies instance findings; any other error supplies the Error message.
// ResultFor does not run validation or evidence checking itself.
func ResultFor(err error) Result {
	result := Result{Valid: err == nil, Issues: []contract.Issue{}}
	if err == nil {
		return result
	}
	var findings *Error
	if errors.As(err, &findings) {
		result.Issues = findings.Issues
	} else {
		result.Error = err.Error()
	}
	return result
}
