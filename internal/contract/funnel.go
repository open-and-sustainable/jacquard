package contract

import "encoding/json"

// Measurement distinguishes a parsed nonnegative integer from missing/unreadable data.
// State is parsed, missing, or unparsed. Value is present only for parsed values.
type Measurement struct {
	State string       `json:"state"`
	Value *json.Number `json:"value,omitempty"`
}

// FunnelItem can contribute to several distinct size buckets within a field.
// PassedStages must be a nonempty prefix of the declared stage order.
type FunnelItem struct {
	ID           string                   `json:"item_id"`
	Values       map[string][]Measurement `json:"values"`
	PassedStages []string                 `json:"passed_stages"`
}

// FunnelInput declares the reference population and all audit stages up front.
type FunnelInput struct {
	Stages          []string                 `json:"stages"`
	ReferenceValues map[string][]json.Number `json:"reference_values"`
	Items           []FunnelItem             `json:"items"`
}

// FunnelCell counts distinct items for one stage, field, and magnitude bucket.
type FunnelCell struct {
	Stage  string `json:"stage"`
	Field  string `json:"field"`
	Bucket string `json:"bucket"`
	Count  int    `json:"count"`
}

// FunnelOutput includes zero-count cells and the cutoffs used to assign deciles.
type FunnelOutput struct {
	JacquardVersion string                   `json:"jacquard_version"`
	Cutoffs         map[string][]json.Number `json:"cutoffs"`
	Cells           []FunnelCell             `json:"cells"`
}
