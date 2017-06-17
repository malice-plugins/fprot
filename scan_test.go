package main_test

import (
	"fmt"
	"strings"
	"testing"
)

const resultString = `
`

const versionString = `
`

func parseVersion(versionOut string) (version string, database string) {

	lines := strings.Split(versionOut, "\n")
	fmt.Println(lines)
	return
}

func parseOutput(output string) (string, error) {

	lines := strings.Split(output, "\n")
	fmt.Println(lines)
	return "", nil
}

// TestParseResult tests the ParseFSecureOutput function.
func TestParseResult(t *testing.T) {

	results, err := parseOutput(resultString)

	if err != nil {
		t.Log(err)
	}

	if true {
		t.Log("results: ", results)
	}

}

// TestParseVersion tests the GetFSecureVersion function.
func TestParseVersion(t *testing.T) {

	version, database := parseVersion(versionString)

	if true {
		t.Log("version: ", version)
		t.Log("database: ", database)
	}

}
