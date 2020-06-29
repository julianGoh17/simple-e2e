package models

import "os"

// TestProcedure is a struct which represents the entire
type TestProcedure struct {
	Name            string
	GlobalVariables map[string]string `yaml:"globalVariables,omitempty"`
	Stages          []TestStage
}

// SetGlobalVariables will set all the variables in 'GlobalVariables' in the terminal as an environmental variable so that it can be
// used in all test steps.
func (tp TestProcedure) SetGlobalVariables() error {
	for key, value := range tp.GlobalVariables {
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}
	return nil
}
