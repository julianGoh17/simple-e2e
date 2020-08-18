package operations

import (
	"fmt"

	model "github.com/julianGoh17/simple-e2e/framework/models"
	"gopkg.in/yaml.v2"
)

// Controller is able to understand which stages and steps to run based on the test file. It is responsible for understanding if a test step has
// failed and will stop the test run prematurely if so.
type Controller struct {
	stepManager *StepManager
	procedure   *model.Procedure
}

// NewController is a constructor function which returns a pointer to the variable to work with
func NewController() *Controller {
	return &Controller{
		stepManager: NewStepManager(),
	}
}

// AddTestStep adds a Step Description and its associated function to the Controller so it knows what needs to do
func (controller *Controller) AddTestStep(description string, function func(*model.Step) error) error {
	return controller.stepManager.AddStepToManager(description, function)
}

// SetProcedure takes the read byte data from the test file and converts it to the Procedure object
func (controller *Controller) SetProcedure(procedureData []byte) error {
	procedure := &model.Procedure{}

	if err := yaml.UnmarshalStrict(procedureData, procedure); err != nil {
		return err
	}
	if procedure.Stages == nil {
		return fmt.Errorf("Test file '%s' does not have any stages to file", procedure.Name)
	}

	controller.procedure = procedure

	return nil
}

// RunTest will run a specified test and if any stages are passed in then it will only run those stages
func (controller *Controller) RunTest(test []byte, stages ...string) error {
	if err := controller.SetProcedure(test); err != nil {
		return err
	}

	set := make(map[string]bool)
	for _, value := range stages {
		set[value] = true
	}

	testPassed := true
	failedStage := ""
	for _, stage := range controller.procedure.Stages {
		if testPassed {
			if len(set) == 0 || set[stage.Name] {
				if err := controller.runStage(&stage); err != nil {
					testPassed = false
					failedStage = stage.Name
					continue
				}
			}
		} else {
			if (len(set) == 0 || set[stage.Name]) && stage.AlwaysRuns {
				if err := controller.runStage(&stage); err != nil {
					return err
				}
			}
		}
	}

	if !testPassed {
		return fmt.Errorf("Test failed at stage: %s", failedStage)
	}

	return nil
}

func (controller *Controller) runStage(stage *model.Stage) error {
	for _, step := range stage.Steps {
		function, err := controller.stepManager.GetTestMethod(step.Description)
		if err != nil {
			return err
		}
		if err := runStep(function, step); err != nil {
			return err
		}
	}
	return nil
}

func runStep(function func(*model.Step) error, step model.Step) error {
	if err := function(&step); err != nil {
		return err
	}
	if !step.HasSucceeded() {
		return fmt.Errorf("Step '%s' has failed", step.Description)
	}
	return nil
}
