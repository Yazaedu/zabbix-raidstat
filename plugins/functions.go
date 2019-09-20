package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// GetInputData - get input data from RAID tool
func GetInputData(execPath string, args ...string) []byte {
	timeout := 5
	execContext, contextCancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer contextCancel()

	cmd := exec.CommandContext(execContext, execPath, args...)
	data, err := cmd.Output()
	if err != nil {
		if os.Getenv("RAIDSTAT_DEBUG") == "y" {
			fmt.Printf("Command output is:\n'''\n%s\n'''\n", string(data))
		}

		if execContext.Err() == context.DeadlineExceeded {
			fmt.Printf("Command '%s' timed out.\n", cmd)
		} else {
			fmt.Println(err)
		}

		os.Exit(1)
	}

	return data
}

// GetRegexpSubmatch - returns string from 1st capture group
func GetRegexpSubmatch(buf []byte, re string) (data string) {
	result := regexp.MustCompile(re).FindStringSubmatch(string(buf))

	if os.Getenv("RAIDSTAT_DEBUG") == "y" {
		fmt.Printf("Regexp is '%s'\n", re)
		fmt.Printf("Result is '%s'\n", result)
		fmt.Printf("Input data is:\n'''\n%s\n'''\n", string(buf))
	}

	if len(result) > 0 {
		data = result[1]
	}

	return
}

// GetRegexpAllSubmatch - returns strings from all capture groups
func GetRegexpAllSubmatch(buf []byte, re string) (data []string) {
	result := regexp.MustCompile(re).FindAllStringSubmatch(string(buf), -1)

	if os.Getenv("RAIDSTAT_DEBUG") == "y" {
		fmt.Printf("Regexp is '%s'\n", re)
		fmt.Printf("Result is '%s'\n", result)
		fmt.Printf("Input data is:\n'''\n%s\n'''\n", string(buf))
	}

	if len(result) > 0 {
		for _, v := range result {
			data = append(data, v[1])
		}
	}

	return
}

// MarshallJSON - returns json object
func MarshallJSON(data interface{}, indent int) []byte {
	var JSON []byte
	var jErr error

	if indent > 0 {
		JSON, jErr = json.MarshalIndent(data, "", strings.Repeat(" ", indent))
	} else {
		JSON, jErr = json.Marshal(data)
	}

	if jErr != nil {
		fmt.Println(jErr)
		os.Exit(1)
	}

	return JSON
}