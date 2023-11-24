package run

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/parsiya/semgrep_go/output"
)

// The Semgrep command and switches. See all with `semgrep scan --help`.
const (
	// The Semgrep command.
	Semgrep = "semgrep"

	// The version switch.
	VersionSwitch = "--version"

	// The config switch.
	ConfigSwitch = "--config"

	// Metrics on and off.
	MetricsOn  = "--metrics=on"
	MetricsOff = "--metrics=off"
)

// OutputFormat enums.
type OutputFormat string

const (
	Text          OutputFormat = "text"
	Emacs         OutputFormat = "emacs"
	JSON          OutputFormat = "json"
	GitLabSAST    OutputFormat = "gitlab-sast"
	GitLabSecrets OutputFormat = "gitlab-secrets"
	JUnitXML      OutputFormat = "junit-xml"
	SARIF         OutputFormat = "sarif"
	Vim           OutputFormat = "vim"
)

// The string representation of the enum is "--[value]".
func (f OutputFormat) String() string {
	return toSwitch(string(f))
}

// ----------

// Verbosity enums.
type Verbosity string

const (
	Quiet   Verbosity = "quiet"
	Verbose Verbosity = "verbose"
	Debug   Verbosity = "debug" // All of verbose with additional debugging.
)

// The string representation of the enum is "--[value]".
func (v Verbosity) String() string {
	return toSwitch(string(v))
}

// ----------

// Semgrep CLI switches.
type Options struct {

	// Default is false, but it doesn't disable metrics for every run. E.g.,
	// metrics are collected when pulling rules from the registry. See more at
	// https://semgrep.dev/docs/metrics/.
	//
	// You can enable metrics by calling EnableMetrics().
	metrics bool

	// Output format, becomes the `--[value]` switch.
	Output OutputFormat

	// Verbosity. Print extra output to the command line or in the JSON file.
	Verbosity Verbosity

	// The path(s) to scan.
	Paths []string

	// The Semgrep rules.
	Rules []string

	// If false (default), the strings in Rules are passed to Semgrep as-is with
	// one or more `--config` switches. This is useful for local paths, registry
	// rulesets (e.g., auto), or a URI with the ruleset.
	//
	// If true, it's assumed the strings in Rules contains one or more Semgrep
	// rules. The contents will be concatenated and stored in a temp file and
	// the name will be passed to Semgrep with `--config`. This is useful when
	// we are using dynamicly created rules.
	//
	// You can set it to true by calling StringRule().
	stringRule bool

	// Extra switches. The user is responsible for their validity.
	Extra []string
}

// Return a new Options struct.
func DefaultOptions(rules []string, paths []string) *Options {
	return &Options{
		Output:    JSON,
		Verbosity: Debug,
		Rules:     rules,
		Paths:     paths,
	}
}

// Enable Metrics
func (s *Options) EnableMetrics() {
	s.metrics = true
}

// Concat all the rules together in a temp file and pass to Semgrep.
func (s *Options) StringRule() {
	s.stringRule = true
}

// Return the options as a string array that can be passed to os/exec.Command.
func (o *Options) string() ([]string, error) {
	var optStr []string

	// If stringRule is true, store the rules in a temp file.
	if o.stringRule {
		var err error
		// Concatenate all strings in rules together.
		mergedRules := ""
		for _, r := range o.Rules {
			mergedRules += r
		}
		// Store them all in a temp file.
		ruleFile, err := createTempFile(mergedRules)
		if err != nil {
			return nil, err
		}
		// Add the rules string to the options string.
		optStr = append(optStr, ConfigSwitch, ruleFile)
	} else {
		// Else, each item in o.Rules must be passed with `--config`.
		for _, r := range o.Rules {
			optStr = append(optStr, ConfigSwitch, r)
		}
	}

	// Add metrics.
	metrics := MetricsOff
	if o.metrics {
		metrics = MetricsOn
	}
	optStr = append(optStr, metrics)

	// Add output format.
	// If output format is empty, use JSON.
	if o.Output.String() == "" {
		o.Output = JSON
	}
	optStr = append(optStr, o.Output.String())

	// Add verbosity.
	// If verbosity is empty, use Debug.
	if o.Verbosity.String() == "" {
		o.Verbosity = Debug
	}
	optStr = append(optStr, o.Verbosity.String())

	// Add any extra switches.
	optStr = append(optStr, o.Extra...)

	// Add the paths.
	optStr = append(optStr, o.Paths...)

	return optStr, nil
}

// ----------

// Run Semgrep and return the results without deserializing.
func internalRun(o *Options) ([]byte, error) {
	// Convert the switches to a string.
	opts, err := o.string()
	if err != nil {
		return nil, err
	}
	// Add them to the command.
	cmd := exec.Command(Semgrep, opts...)

	log.Printf("Running Semgrep as: %s", cmd.String())

	// Set stdout and stderr.
	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	// Run Semgrep.
	if err := cmd.Run(); err != nil {
		// Return the error.
		return nil, errors.New(stdErr.String())
	}

	// If Semgrep started successfully but has an internal error, stderr will be
	// empty, the error will be in the JSON output. We will return it as-is.
	return stdOut.Bytes(), nil
}

// Run Semgrep and return output as-is.
//
// Note that just because the Semgrep command was executed, doesn't mean there
// were no errors. Semgrep will store internal errors in Output.Errors which is
// a []CliError.
//
// We will return them as-is and let the user decide what they want to do with
// the errors.
func (o *Options) Run() ([]byte, error) {
	// Run Semgrep and return any errors if Semgrep did not run.
	return internalRun(o)
}

// Run Semgrep and return the deserialized JSON output.
func (o *Options) RunJSON() (out output.Output, err error) {
	// Change the output format to JSON in case we made a mistake when creating
	// the Options object.
	o.Output = JSON

	// Run Semgrep and return any errors if Semgrep did not run.
	data, err := internalRun(o)
	if err != nil {
		return out, err
	}
	// Deserialize the output.
	return output.Deserialize(data)
}

// Checks if Semgrep is installed.
func IsInstalled() bool {
	_, err := Version()
	return err == nil
}

// Returns the Semgrep version.
func Version() (string, error) {
	cmd := exec.Command(Semgrep, VersionSwitch)
	var stdout strings.Builder
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error running semgrep: %v", err)
	}
	return strings.TrimSpace(stdout.String()), nil
}
