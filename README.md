# semgrep-go
Go package to interact with the Semgrep CLI programmatically.

Similar package for Rust: https://github.com/parsiya/semgrep-rs

## Usage
`go get github.com/parsiya/semgrep_go`.

```go
// Setup Semgrep switches.
opts := run.Options{
    Output:    run.JSON,       // Output format is JSON.
    Paths:     []string{"code/juice-shop"},
    Rules:     []string{"p/default"},
    Verbosity: run.Debug,
    Extra:     []string{"--no-rewrite-rule-ids"},
}

log.Print("Running Semgrep, this might take a minute.")
// Run Semgrep and get the deserialized output.
out, err := opts.Run()
if err != nil {
    // handle error
}

for _, hit := range out.Results {
    // Print the rule ID and message.
    fmt.Println("RuleID: ", hit.RuleID())
    fmt.Println("Message: ", hit.Message())
}
```

For more examples, please see the following blog and code:

* https://github.com/parsiya/semgrep-fun
* 

## Semgrep Output Structs
The structure of the output is defined in [semgrep/semgrep-interfaces][int-gh].
The source of truth is the atd file, but it's an OCaml thing and we cannot parse
it, so we rely on the automatically generated JSON schema in
[semgrep_output_v1.jsonschema][schema].

I use [omissis/go-jsonschema][go-schema] (formerly at
`atombender/go-jsonschema`) to generate the Go structs from the JSON schema.
From time to time, the schema might break backwards compatibility. I've seen it
happen twice in the last few months, but I also do not upgrade Semgrep and run
tests in every single version.

Keep this in mind before upgrading your Semgrep version. Generate the structs
and a quick compare to see if anything major has changed.

### Generating Structs
You can generate structs like this:

```
git clone --recurse-submodules https://github.com/parsiya/semgrep_go
# optional: to generate the structs for a specific version checkout that tag
cd semgrep_go/output/semgrep-interfaces
git checkout v1.52.0

# install go-jsonschema
go install github.com/omissis/go-jsonschema/cmd/gojsonschema@latest

# generate the output
# -p output: package name is output
# -o output.go: write the structs to output.go
gojsonschema -p output -o ../output.go --verbose semgrep_output_v1.jsonschema
```

What I tried when generating Go structs from JSON schemas and didn't work:
https://parsiya.io/abandoned-research/semgrep-output-json/.

Similar experiment for Rust
https://parsiya.net/blog/2022-10-16-yaml-wrangling-with-rust/.

[int-gh]: https://github.com/semgrep/semgrep-interfaces
[schema]: https://github.com/semgrep/semgrep-interfaces/blob/main/semgrep_output_v1.jsonschema
[go-schema]: https://github.com/omissis/go-jsonschema

### Compatibility
Currently, the package is using the v1.48.0 structs. Somewhere around v1.45.0
the type of `CliError.Type` was changed from string to `ErrorType`.

The current output struct was tested with these Semgrep versions:

* 1.48.0
* 1.49.0
* 1.50.0
* 1.52.0

[si]: https://github.com/returntocorp/semgrep-interfaces