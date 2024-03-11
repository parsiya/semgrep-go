package output

import (
	_ "embed"
	"fmt"
	"reflect"
	"testing"
)

// Embed `test/juice-shop-1.42.0.json` as a string
//
//go:embed test/juice-shop-1.42.0.json
var juiceShopJSON []byte

func TestDeserialize(t *testing.T) {
	out, err := Deserialize(juiceShopJSON)
	if err != nil {
		t.Errorf("Couldn't deserialize test JSON: %s", err.Error())
	}

	// Do some quick sanity checks on the output.

	// Check if version is 1.42.0.
	if *out.Version != "1.42.0" {
		t.Errorf("Got version: %s - wanted: %s", *out.Version, "1.42.0")
	}
}

//go:embed test/test-juice-shop-1.42.0.json
var testBytes []byte

// Deserialize testBytes and return the results.
func deser() (res []CliMatch, err error) {
	out, err := Deserialize(testBytes)
	if err != nil {
		return res, fmt.Errorf("Couldn't deserialize test JSON: %w", err)
	}
	// Return the results.
	return out.Results, nil
}

func TestCliMatch_Metavar(t *testing.T) {

	// Get the results from the test file.
	res, err := deser()
	if err != nil {
		t.Errorf(err.Error())
	}

	tests := []struct {
		name        string
		result      CliMatch
		metavarName string
		want        string
		wantErr     bool
	}{
		// result 0
		{
			name:        "result 0 - uses - exists",
			result:      res[0],
			metavarName: "uses",
			want:        "github/codeql-action/init@v2",
			wantErr:     false,
		},
		{
			name:        "result 0 - $USES - exists",
			result:      res[0],
			metavarName: "$USES",
			want:        "github/codeql-action/init@v2",
			wantErr:     false,
		},
		{
			name:        "result 0 - $RANDOM - doesn't exist",
			result:      res[0],
			metavarName: "$RANDOM",
			want:        "",
			wantErr:     true,
		},
		{
			name:        "result 0 - random - doesn't exist",
			result:      res[0],
			metavarName: "random",
			want:        "",
			wantErr:     true,
		},
		// result 1
		{
			name:        "result 1 - 1 - exists",
			result:      res[1],
			metavarName: "1",
			want:        "post",
			wantErr:     false,
		},
		{
			name:        "result 1 - app - propagated",
			result:      res[1],
			metavarName: "app",
			want:        "express.Router()",
			wantErr:     false,
		},
		{
			name:        "result 1 - $app - propagated",
			result:      res[1],
			metavarName: "$app",
			want:        "express.Router()",
			wantErr:     false,
		},
		// result 2
		{
			name:        "result 2 - sink - propagated value",
			result:      res[2],
			metavarName: "sink",
			want:        "params.file",
			wantErr:     false,
		},
		{
			name:        "result 2 - next - exists",
			result:      res[2],
			metavarName: "$nExT",
			want:        "next",
			wantErr:     false,
		},
		{
			name:        "result 2 - req - exists",
			result:      res[2],
			metavarName: "rEq",
			want:        "params",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.result.Metavar(tt.metavarName)
			if (err != nil) != tt.wantErr {
				t.Errorf("CliMatch.Metavar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CliMatch.Metavar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTesty(t *testing.T) {

	// Get the results from the test file.
	res, err := deser()
	if err != nil {
		t.Errorf(err.Error())
	}

	for _, r := range res {
		nem := r.Extra.Metadata

		// try to cast it to map[string]interface{}
		cast, ok := nem.(map[string]interface{})
		if !ok {
			t.Errorf("Couldn't cast metadata. Current type is: %T", nem)
		}

		n := fmt.Sprintf("%T", nem)
		fmt.Println(n)

		for k, v := range cast {
			fmt.Printf("k: %s - v: %v\n", k, v)
		}

	}

}

func TestCliMatch_Metadata(t *testing.T) {

	// Get the results from the test file.
	res, err := deser()
	if err != nil {
		t.Errorf(err.Error())
	}

	tests := []struct {
		name         string
		result       CliMatch
		metadataName string
		want         interface{}
		wantErr      bool
	}{
		{
			name:         "result 0 - category - exists",
			result:       res[0],
			metadataName: "category",
			want:         "security",
			wantErr:      false,
		},
		{
			name:         "result 0 - confidence - exists",
			result:       res[0],
			metadataName: "confidence",
			want:         "HIGH",
			wantErr:      false,
		},
		{
			name:         "result 0 - random - doesn't exist",
			result:       res[0],
			metadataName: "random",
			want:         nil,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.result.Metadata(tt.metadataName)
			if (err != nil) != tt.wantErr {
				t.Errorf("CliMatch.Metadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CliMatch.Metadata() = %v, want %v", got, tt.want)
			}
		})
	}
}
