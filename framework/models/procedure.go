package models

import "os"

// Procedure is a struct which represents the entire
type Procedure struct {
	Name            string
	Description     string
	GlobalVariables map[string]string `yaml:"globalVariables,omitempty"`
	Stages          []Stage
}

// SetGlobalVariables will set all the variables in 'GlobalVariables' in the terminal as an environmental variable so that it can be
// used in all test steps.
func (p Procedure) SetGlobalVariables() error {
	for key, value := range p.GlobalVariables {
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}
	return nil
}
