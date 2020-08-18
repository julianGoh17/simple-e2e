package models

import (
	"fmt"
	"os"
	"strings"
)

// Step is the struct that represents that will map the human readable string to the function
type Step struct {
	Description  string
	Variables    map[string]string `yaml:"variables,omitempty"`
	converter    TypeConverter
	isSuccessful bool
}

// GetDescriptionVariables will get the variables from TestStep.Description. For example, "this is a 'variable'" will return ["variable"]
func (s *Step) GetDescriptionVariables() ([]string, error) {
	descriptionComponents := strings.Split(s.Description, "'")
	if len(descriptionComponents)%2 != 1 {
		return nil, fmt.Errorf("Test Step '%s' is ill formatted because it contains an odd number of ','", s.Description)
	}

	descriptionVariables := []string{}

	for index := 1; index < len(descriptionComponents); index += 2 {
		descriptionVariables = append(descriptionVariables, descriptionComponents[index])
	}

	return descriptionVariables, nil
}

// GetValueFromVariablesAsString will return the variable specific to this step from the step.variables if it exists otherwise it will return an error
func (s *Step) GetValueFromVariablesAsString(variableName string) (string, error) {
	if val, ok := s.Variables[variableName]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Could not find variable '%s' in step variables", variableName)
}

// GetValueFromVariablesAsInteger will return the variable specific to this step from the step.variables as an integer if it exists otherwise it will return an error
func (s *Step) GetValueFromVariablesAsInteger(variableName string) (int, error) {
	if val, ok := s.Variables[variableName]; ok {
		return s.converter.GetInteger(val)
	}
	return 0, fmt.Errorf("Could not find variable '%s' in step variables", variableName)
}

// GetValueFromVariablesAsFloat32 will return the variable specific to this step from the step.variables as a float32 if it exists otherwise it will return an error
func (s *Step) GetValueFromVariablesAsFloat32(variableName string) (float32, error) {
	if val, ok := s.Variables[variableName]; ok {
		return s.converter.GetFloat32(val)
	}
	return float32(0), fmt.Errorf("Could not find variable '%s' in step variables", variableName)
}

// GetValueFromVariablesAsFloat64 will return the variable specific to this step from the step.variables as a float64 if it exists otherwise it will return an error
func (s *Step) GetValueFromVariablesAsFloat64(variableName string) (float64, error) {
	if val, ok := s.Variables[variableName]; ok {
		return s.converter.GetFloat64(val)
	}
	return float64(0), fmt.Errorf("Could not find variable '%s' in step variables", variableName)
}

// GetValueFromVariablesAsBoolean will return the variable specific to this step from the step.variables as a boolean if it exists otherwise it will return an error
func (s *Step) GetValueFromVariablesAsBoolean(variableName string) (bool, error) {
	if val, ok := s.Variables[variableName]; ok {
		return s.converter.GetBoolean(val)
	}
	return false, fmt.Errorf("Could not find variable '%s' in step variables", variableName)
}

// GetValueFromVariablesAsStringArray will return the variable specific to this step from the step.variables as a string array (separated by commas)
// if it exists otherwise it will return an error
func (s *Step) GetValueFromVariablesAsStringArray(variableName string) ([]string, error) {
	if val, ok := s.Variables[variableName]; ok {
		return strings.Split(val, ","), nil
	}
	return []string{}, fmt.Errorf("Could not find variable '%s' in step variables", variableName)
}

// GetValueFromVariablesAsIntegerArray will return the variable specific to this step from the step.variables as a integer array (separated by commas)
// if it exists otherwise it will return an error
func (s *Step) GetValueFromVariablesAsIntegerArray(variableName string) ([]int, error) {
	if val, ok := s.Variables[variableName]; ok {
		return s.converter.GetIntegerArray(val)
	}
	return []int{}, fmt.Errorf("Could not find variable '%s' in step variables", variableName)
}

// GetValueFromVariablesAsFloat32Array will return the variable specific to this step from the step.variables as a float32 array (separated by commas)
// if it exists otherwise it will return an error
func (s *Step) GetValueFromVariablesAsFloat32Array(variableName string) ([]float32, error) {
	if val, ok := s.Variables[variableName]; ok {
		return s.converter.GetFloat32Array(val)
	}
	return []float32{}, fmt.Errorf("Could not find variable '%s' in step variables", variableName)
}

// GetValueFromVariablesAsFloat64Array will return the variable specific to this step from the step.variables as a float64 array (separated by commas)
// if it exists otherwise it will return an error
func (s *Step) GetValueFromVariablesAsFloat64Array(variableName string) ([]float64, error) {
	if val, ok := s.Variables[variableName]; ok {
		return s.converter.GetFloat64Array(val)
	}
	return []float64{}, fmt.Errorf("Could not find variable '%s' in step variables", variableName)
}

// GetValueFromVariablesAsBooleanArray will return the variable specific to this step from the step.variables as a boolean array (separated by commas)
// if it exists otherwise it will return an error
func (s *Step) GetValueFromVariablesAsBooleanArray(variableName string) ([]bool, error) {
	if val, ok := s.Variables[variableName]; ok {
		return s.converter.GetBooleanArray(val)
	}
	return []bool{}, fmt.Errorf("Could not find variable '%s' in step variables", variableName)
}

// GetGlobalVariable will return the global variable from the Env vars. This just serves as a wrapper method to make it easier to read the
// test code
func (s *Step) GetGlobalVariable(variableName string) string {
	return os.Getenv(variableName)
}

// HasSucceeded returns whether or not the TestStep has succeeded
func (s *Step) HasSucceeded() bool {
	return s.isSuccessful
}

// SetPassed sets the value of TestStep to passed
func (s *Step) SetPassed() {
	s.isSuccessful = true
}

// SetFailed sets the value of TestStep to failed
func (s *Step) SetFailed() {
	s.isSuccessful = false
}
