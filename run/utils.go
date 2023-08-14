package run

import "os"

// toSwitch returns "--[value]".
func toSwitch(value string) string {
	return "--" + value
}

const TempFilePrefix = "semgrep"

// createTempFile writes the data to a temporary file and returns the handle.
// Note: We need to delete the file ourselves.
func createTempFile(data string) (string, error) {

	// Create the temp file with the prefix.
	f, err := os.CreateTemp("", TempFilePrefix)
	if err != nil {
		return "", err
	}

	// Write to file.
	_, err = f.WriteString(data)
	return f.Name(), err
}
