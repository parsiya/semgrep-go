package run

import (
	"bytes"
	"errors"
	"fmt"
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
	Rules string

	// If set to true, it's assumed the string in Rules contains one or more
	// Semgrep rules. The contents will be stored in a temp file and the name
	// will be passed to Semgrep with `--config`. This is useful when we are
	// using dynamicly created rules.
	//
	// If set to false (default), the string in Rules is passed to Semgrep
	// as-is. This is useful for local paths, registry rulesets (e.g., auto), or
	// a URI with the ruleset.
	//
	// You can set it to true by calling StringRule().
	stringRule bool

	// Extra switches. The user is responsible for their validity.
	Extra []string
}

// Return a new Options struct.
func DefaultOptions(rules string, paths []string) *Options {
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

// Do not store the rules in a text file. Pass the rule string direclty as-is.
func (s *Options) StringRule() {
	s.stringRule = true
}

// Return the options as a string array that can be passed to os/exec.Command.
func (o *Options) string() ([]string, error) {
	var optStr []string

	ruleFile := o.Rules
	// If stringRule is true, store the rules in a temp file.
	if o.stringRule {
		var err error
		ruleFile, err = createTempFile(o.Rules)
		if err != nil {
			return nil, err
		}
	}
	// Add the rules string to the options string.
	optStr = append(optStr, ConfigSwitch, ruleFile)

	// Add metrics.
	metrics := MetricsOff
	if o.metrics {
		metrics = MetricsOn
	}
	optStr = append(optStr, metrics)

	// Add output format.
	optStr = append(optStr, o.Output.String())

	// Add verbosity.
	optStr = append(optStr, o.Verbosity.String())

	// Add the extra switches.
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

// Runs Semgrep. Return the deserialized output and errors (if any).
func Run(o *Options) (output.Output, error) {
	var out output.Output

	// Run Semgrep and return any errors if Semgrep did not run.
	data, err := internalRun(o)
	if err != nil {
		return out, err
	}

	// If Semgrep was executed. Deserialize the output and return it.
	//
	// Note that just because the Semgrep command was executed, it doesn't mean
	// there were no errors. Semgrep will store internal errors in Output.Errors
	// which is a []CliError.
	//
	// We will return them as-is and let the user decide what they want to do
	// with the errors.
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
