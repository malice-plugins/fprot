package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

const resultString = `
`

const versionString = `
`

func extractVirusName(line string) string {
	r := regexp.MustCompile(`<(.+)>`)
	res := r.FindStringSubmatch(line)
	if len(res) != 2 {
		return ""
	}
	return res[1]
}

func parseUpdatedDate(date string) string {
	layout := "200601021504"
	t, _ := time.Parse(layout, date)
	return fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func parseVersion(versionOut string) (version string, database string) {

	lines := strings.Split(versionOut, "\n")
	colonSeparated := []string{}

	for _, line := range lines {
		if len(line) != 0 {
			if strings.Contains(line, ":") {
				colonSeparated = append(colonSeparated, line)
			}
		}
	}
	// fmt.Println(lines)

	// Extract FPROT Details from scan output
	if len(colonSeparated) != 0 {
		for _, line := range colonSeparated {
			if len(line) != 0 {
				keyvalue := strings.Split(line, ":")
				if len(keyvalue) != 0 {
					switch {
					case strings.Contains(keyvalue[0], "Virus signatures"):
						version = parseUpdatedDate(strings.TrimSpace(keyvalue[1]))
					case strings.Contains(line, "Engine version"):
						database = strings.TrimSpace(keyvalue[1])
					}
				}
			}
		}
	} else {
		fmt.Println("[ERROR] colonSeparated was empty: ", colonSeparated)
		os.Exit(2)
	}
	return
}

func parseOutput(output string) (string, error) {
	colonSeparated := []string{}

	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if len(line) != 0 {
			if strings.Contains(line, ":") {
				colonSeparated = append(colonSeparated, line)
			}
			if strings.Contains(line, "[Found virus]") {
				result := extractVirusName(line)
				if len(result) != 0 {
					return result, nil
				} else {
					fmt.Println("[ERROR] Virus name extracted was empty: ", result)
					os.Exit(2)
				}
			}
		}
	}
	return "", nil
}

// TestParseResult tests the ParseFSecureOutput function.
func TestParseResult(t *testing.T) {
	b, err := ioutil.ReadFile("av_scan.out") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	results, err := parseOutput(string(b))

	if err != nil {
		t.Log(err)
	}

	if true {
		t.Log("results: ", results)
	}

}

// TestParseVersion tests the GetFSecureVersion function.
func TestParseVersion(t *testing.T) {
	b, err := ioutil.ReadFile("av_scan.out") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	version, database := parseVersion(string(b))

	if true {
		t.Log("version: ", version)
		t.Log("database: ", database)
	}

}
