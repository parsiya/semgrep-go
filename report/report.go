package report

import (
	"fmt"
	htmltemplate "html/template"
	texttemplate "html/template"
	"strings"

	"github.com/parsiya/semgrep_go/output"
)

// Creates a text report based on the template and Semgrep's JSON output. Note,
// this function uses text/template which is not safe for HTML generation.
func TextReport(tmpl string, output output.Output) (string, error) {
	t, err := texttemplate.New("report").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse the report template: %w", err)
	}

	var report strings.Builder
	err = t.Execute(&report, output)
	if err != nil {
		return "", fmt.Errorf("error generating the report: %w", err)
	}
	return report.String(), nil
}

// Creates an HTMl report based on the template and Semgrep's JSON output. This
// function uses html/template which is safe* for HTML generation.
func HTMLReport(tmpl string, output output.Output) (string, error) {
	t, err := htmltemplate.New("report").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse the report template: %w", err)
	}

	var report strings.Builder
	err = t.Execute(&report, output)
	if err != nil {
		return "", fmt.Errorf("error generating the report: %w", err)
	}
	return report.String(), nil
}
