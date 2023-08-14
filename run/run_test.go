package run

import (
	"testing"
)

func TestOutputFormat_String(t *testing.T) {

	enums := []OutputFormat{
		Text,
		Emacs,
		JSON,
		GitLabSAST,
		GitLabSecrets,
		JUnitXML,
		SARIF,
		Vim,
	}

	formats := []string{
		"text",
		"emacs",
		"json",
		"gitlab-sast",
		"gitlab-secrets",
		"junit-xml",
		"sarif",
		"vim",
	}

	for i := range enums {
		t.Run(formats[i], func(t *testing.T) {
			want := "--" + formats[i]
			if got := enums[i].String(); got != want {
				t.Errorf("OutputFormat.String() = %v, want %v", got, want)
			}
		})
	}
}

func TestIsInstalled(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "installed",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInstalled(); got != tt.want {
				t.Errorf("IsInstalled() = %v, want %v", got, tt.want)
			}
		})
	}
}
