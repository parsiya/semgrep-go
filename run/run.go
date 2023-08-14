package run

import (
	"bytes"
	"errors"
	"os/exec"

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
	metrics bool

	// Output format, becomes the `--[value]` switch.
	Output OutputFormat

	// Verbosity. Print extra output to the command line or in the JSON file.
	Verbosity Verbosity

	// The path(s) to scan.
	Paths []string

	// The Semgrep rules.
	Rules string

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

// Return the options as a string array that can be passed to os/exec.Command.
func (o *Options) string() ([]string, error) {
	var optStr []string

	// Write the rules to a temporary file.
	ruleFile, err := createTempFile(o.Rules)
	if err != nil {
		return nil, err
	}
	// Add it to the options string.
	optStr = append(optStr, ConfigSwitch, ruleFile)

	// Add metrics.
	if o.metrics {
		optStr = append(optStr, MetricsOn)
	} else {
		optStr = append(optStr, MetricsOff)
	}

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

// Runs Semgrep. Return the output and errors (if any).
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
	// We will return des as-is and let the user decide what they want to do
	// with the errors.
	return output.Deserialize(data)
}

// Checks if Semgrep is installed.
func IsInstalled() bool {
	cmd := exec.Command(Semgrep, VersionSwitch)
	// Run Semgrep.
	if err := cmd.Run(); err != nil {
		// Return the error.
		return false
	}
	return true
}
