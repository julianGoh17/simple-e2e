package models

import (
	"fmt"
	"os"
	"strings"
)

// Step is the struct that represents that will map the human readable string to the function
type Step struct {
	Description string
	Variables   map[string]string `yaml:"variables,omitempty"`
}

// GetDescriptionVariables will get the variables from TestStep.Description. For example, "this is a 'variable'" will return ["variable"]
func (s Step) GetDescriptionVariables() ([]string, error) {
	descriptionComponents := strings.Split(s.Description, "'")
	if len(descriptionComponents)%2 != 1 {
		return nil, fmt.Errorf("Test Step '%s' is ill formatted because it contains an odd number of '", s.Description)
	}

	descriptionVariables := []string{}

	for index := 1; index < len(descriptionComponents); index += 2 {
		descriptionVariables = append(descriptionVariables, descriptionComponents[index])
	}

	return descriptionVariables, nil
}

// GetTestVariable will return the variable specific to this test from the 'variables' field
func (s Step) GetTestVariable(variableName string) string {
	return s.Variables[variableName]
}

// GetGlobalVariable will return the global variable from the Env vars. This just serves as a wrapper method to make it easier to read the
// test code
func (s Step) GetGlobalVariable(variableName string) string {
	return os.Getenv(variableName)
}
