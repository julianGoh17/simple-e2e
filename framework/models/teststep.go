package models

import (
	"fmt"
	"os"
	"strings"
)

// TestStep is the struct that represents that will map the human readable string to the function
type TestStep struct {
	Description string
	Variables   map[string]string `yaml:"variables,omitempty"`
}

// GetDescriptionVariables will get the variables from TestStep.Description. For example, "this is a 'variable'" will return ["variable"]
func (ts TestStep) GetDescriptionVariables() ([]string, error) {
	descriptionComponents := strings.Split(ts.Description, "'")
	if len(descriptionComponents)%2 != 1 {
		return nil, fmt.Errorf("Test Step '%s' is ill formatted because it contains an odd number of '", ts.Description)
	}

	var descriptionVariables []string

	for index := 1; index < len(descriptionComponents); index += 2 {
		descriptionVariables = append(descriptionVariables, descriptionComponents[index])
	}

	return descriptionVariables, nil
}

// GetTestVariable will return the variable specific to this test from the 'variables' field
func (ts TestStep) GetTestVariable(variableName string) string {
	return ts.Variables[variableName]
}

// GetGlobalVariable will return the global variable from the Env vars. This just serves as a wrapper method to make it easier to read the
// test code
func (ts TestStep) GetGlobalVariable(variableName string) string {
	return os.Getenv(variableName)
}
