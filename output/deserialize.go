package output

import (
	"fmt"

	"github.com/parsiya/semgrep_go/output/v1_35_0"
)

// Create an alias for the latest version of the generated Go structs for the
// Semgrep JSON output.
type Output = v1_35_0.SemgrepOutputV1Jsonschema

// Deserialize the Semgrep's JSON output to a Go struct
func Deserialize(data []byte) (Output, error) {
	var out Output
	if err := out.UnmarshalJSON(data); err != nil {
		return out, fmt.Errorf("failed to deserialize Semgrep's output: %w", err)
	}
	return out, nil
}
