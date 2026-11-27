package jacquard

import (
	"embed"
	"fmt"
)

//go:embed definitions/v1/*_schema.json
var schemas embed.FS

// SchemaTypes returns the supported JSON interface names in a stable order.
//
// The names are input, output, funnel_input, and funnel_output. The returned slice
// is independent of later calls and may be modified by the caller.
func SchemaTypes() []string {
	return []string{"input", "output", "funnel_input", "funnel_output"}
}

// Schema returns an embedded schema without exposing mutable global state.
//
// Parameters:
//   - version: The exact interface version, [SchemaVersion] (v1).
//   - kind: An interface name returned by [SchemaTypes].
//
// Returns:
//   - An independent byte slice containing a self-contained draft-07 JSON Schema.
//   - An error if the version or kind is unsupported. No fallback is selected.
func Schema(version, kind string) ([]byte, error) {
	if version != SchemaVersion {
		return nil, fmt.Errorf("unsupported schema version %q", version)
	}
	for _, supported := range SchemaTypes() {
		if kind == supported {
			return schemas.ReadFile("definitions/" + version + "/" + kind + "_schema.json")
		}
	}
	return nil, fmt.Errorf("unsupported schema type %q", kind)
}
