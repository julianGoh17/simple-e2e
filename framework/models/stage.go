package models

// Stage is a struct which represents the associated test steps in a stage
type Stage struct {
	Name       string `yaml:"name"`
	AlwaysRuns bool   `yaml:"alwaysRuns"`
	Steps      []Step `yaml:",flow"`
}
