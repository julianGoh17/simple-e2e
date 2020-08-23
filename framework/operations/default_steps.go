package operations

import (
	"fmt"

	"github.com/julianGoh17/simple-e2e/framework/models"
)

func getDefaultSteps() map[string]func(step *models.Step) error {
	defaultSteps := map[string]func(step *models.Step) error{
		"Say hello to": SayHelloTo,
	}

	return defaultSteps
}

// SayHelloTo is just a placeholder function for testing
func SayHelloTo(step *models.Step) error {
	name, err := step.GetValueFromVariablesAsString("NAME")
	if err != nil {
		fmt.Println("Step failed!")
		step.SetFailed()
		return err
	}
	fmt.Printf("Hello there %s!\n", name)
	step.SetPassed()
	return nil
}
