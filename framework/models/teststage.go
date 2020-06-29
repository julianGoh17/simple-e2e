package models

// TestStage is a struct which represents the associated test steps in a stage
type TestStage struct {
	Name  string
	Steps []TestStep 
}
