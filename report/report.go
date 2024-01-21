package report

import (
	"fmt"
	htmltemplate "html/template"
	texttemplate "html/template"
	"strings"

	"github.com/parsiya/semgrep_go/output"
)

// Creates a generic text report based on the template and Semgrep's JSON
// output. Note, this function uses text/template which is not safe for HTML
// generation.
func GenericTextReport(tmpl string, output output.Output) (string, error) {
	return genericReport(tmpl, output, false)
}

// Creates a generic HTMl report based on the template and Semgrep's JSON
// output. This function uses html/template which is safe* for HTML generation.
func GenericHTMLReport(tmpl string, output output.Output) (string, error) {
	return genericReport(tmpl, output, true)
}

// Internal function for creating a generic report. If isHTML is false, it will
// create a text report.
func genericReport(tmpl string, output output.Output, isHTML bool) (string, error) {
	var t *htmltemplate.Template
	var err error
	if isHTML {
		t, err = htmltemplate.New("report").Parse(tmpl)
	} else {
		t, err = texttemplate.New("report").Parse(tmpl)
	}

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
