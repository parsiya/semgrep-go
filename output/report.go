package output

import (
	"fmt"
	htmltemplate "html/template"
	texttemplate "html/template"
	"strings"
)

// Create reports from the Semgrep JSON output.

// RuleIDHitMap returns a slice of sorted HitMapRows. The key in each row is the
// ruleID and the value is the number of hits for that rule. if sortByCount is
// true, the slice is sorted by hits in descending order (e.g., ruleIDs with
// more hits first). Otherwise, the slice is sorted by ruleID alphabetically.
func (o Output) RuleIDHitMap(sortByCount bool) []HitMapRow {
	// Create the HitMap from results.
	hm := ruleHitMap(o)
	// Get the sorted data.
	return hm.SortedData(sortByCount)
}

// FilePathHitMap returns a slice of sorted HitMapRows. The key in each row is
// the file path and the value is the number of hits for that fule. if
// sortByCount is true, the slice is sorted by hits in descending order (e.g.,
// ruleIDs with more hits first). Otherwise, the slice is sorted by file path
// alphabetically.
func (o Output) FilePathHitMap(sortByCount bool) []HitMapRow {
	// Create the HitMap from results.
	hm := fileHitMap(o)
	// Get the sorted data.
	return hm.SortedData(sortByCount)
}

// RuleIDStringTable returns a string table of sorted ruleIDs and hits.
func (o Output) RuleIDTextReport(sortByCount bool) string {
	// Create the hitmap from results.
	hm := ruleHitMap(o)
	// Get the sorted data as a table.
	return hm.ToStringTable([]string{"Rule ID", "Hits"}, true)
}

// FilePathStringTable returns a string table of sorted file paths and hits.
func (o Output) FilePathTextReport(sortByCount bool) string {
	// Create the hitmap from results.
	hm := fileHitMap(o)
	// Get the sorted data as a table.
	return hm.ToStringTable([]string{"File Path", "Hits"}, true)
}

// -----

// Creates a generic text report based on the template and Semgrep's JSON
// output. Note, this function uses text/template which is not safe for HTML
// generation.
func (o Output) GenericTextReport(tmpl string, output Output) (string, error) {
	return genericReport(tmpl, output, false)
}

// Creates a generic HTMl report based on the template and Semgrep's JSON
// output. This function uses html/template which is safe* for HTML generation.
func (o Output) GenericHTMLReport(tmpl string, output Output) (string, error) {
	return genericReport(tmpl, output, true)
}

// Internal function for creating a generic report. If isHTML is false, it will
// create a text report.
func genericReport(tmpl string, output Output, isHTML bool) (string, error) {
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
