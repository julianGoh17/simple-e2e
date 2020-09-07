package operations

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/julianGoh17/simple-e2e/framework/docker"
	model "github.com/julianGoh17/simple-e2e/framework/models"
	"github.com/julianGoh17/simple-e2e/framework/util"
	"gopkg.in/yaml.v2"
)

var (
	logger = util.GetStandardLogger()
	config = util.NewConfig()
)

// Controller is able to understand which stages and steps to run based on the test file. It is responsible for understanding if a test step has
// failed and will stop the test run prematurely if so.
type Controller struct {
	stepManager *StepManager
	procedure   *model.Procedure
	docker      *docker.Handler
}

// NewController is a constructor function which returns a pointer to the variable to work with
func NewController() (*Controller, error) {
	docker, err := docker.NewHandler()
	if err != nil {
		return nil, err
	}
	return &Controller{
		stepManager: NewStepManager(),
		docker:      docker,
	}, nil
}

// AddTestStep adds a Step Description and its associated function to the Controller so it knows what needs to do
func (controller *Controller) AddTestStep(description string, function func(*model.Step) error) error {
	logger.Trace().
		Str("step", description).
		Interface("func", function).
		Msg("Adding step to controller")
	return controller.stepManager.AddStepToManager(description, function)
}

// SetProcedure takes the read byte data from the test file and converts it to the Procedure object
func (controller *Controller) SetProcedure(procedureData []byte) error {
	logger.Trace().
		Msg("Unmarshalling test file into object")
	procedure := &model.Procedure{}

	if err := yaml.UnmarshalStrict(procedureData, procedure); err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to unmarshall object")
		return err
	}
	if procedure.Stages == nil {
		err := fmt.Errorf("Test file '%s' does not have any stages to file", procedure.Name)
		logger.Error().
			Err(err).
			Msg("Test did not contain any stages")
		return err
	}
	// Need to pass in snapshot manager/docker/etc into each step so they access same instance
	for stage := range procedure.Stages {
		for step := range procedure.Stages[stage].Steps {
			procedure.Stages[stage].Steps[step].Docker = controller.docker
		}
	}
	controller.procedure = procedure
	logger.Trace().
		Msg("Succesfully unmarshalled test file into object")
	return nil
}

// RunTest will run a specified test and if any stages are passed in then it will only run those stages
func (controller *Controller) RunTest(testPath string, stages ...string) error {
	logger.Info().
		Str("testPath", testPath).
		Str("stages", strings.Join(stages, ",")).
		Msg("Running test")

	body, err := ioutil.ReadFile(testPath)
	if err != nil {
		return fmt.Errorf("unable to read file: %v", err)
	}

	return controller.runTest(body, stages...)
}

func (controller *Controller) runTest(test []byte, stages ...string) error {
	logger.Trace().
		Str("stages", strings.Join(stages, ",")).
		Msg("Mapping test to object and then running test")
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
			logger.Debug().
				Str("stage", stage.Name).
				Bool("failed", false).
				Msg("Test has not failed, continuing to run stage.")
			if len(set) == 0 || set[stage.Name] {
				if err := controller.runStage(&stage); err != nil {
					testPassed = false
					failedStage = stage.Name
					continue
				}
			}
		} else {
			logger.Debug().
				Str("stage", stage.Name).
				Bool("failed", false).
				Bool("alwaysRun", stage.AlwaysRuns).
				Msg("Test has failed, continuing to run stages with 'alwaysRun' is true.")
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

func (controller *Controller) runStage(stagePointer *model.Stage) error {
	stage := *stagePointer
	logger.Info().
		Str("stage", stage.Name).
		Msg("Beginning to run through steps in stage")
	for _, step := range stage.Steps {
		function, err := controller.stepManager.GetTestMethod(step.Description)
		if err != nil {
			logger.Error().
				Err(err).
				Str("stage", stage.Name).
				Str("step", step.Description).
				Bool("hasFailed", true).
				Msg("Could not find step in stage manager")
			return err
		}
		if err := runStep(function, step); err != nil {
			logger.Error().
				Err(err).
				Str("stage", stage.Name).
				Bool("hasFailed", true).
				Msg("Stage has failed at step")
			return err
		}
	}
	logger.Info().
		Str("stage", stage.Name).
		Msg("Completed running through steps in stage")
	return nil
}

func runStep(function func(*model.Step) error, step model.Step) error {
	logger.Info().
		Str("step", step.Description).
		Msg("Beginning to run step")
	if err := function(&step); err != nil {
		return err
	}
	if !step.HasSucceeded() {
		err := fmt.Errorf("Step '%s' has failed", step.Description)
		logger.Error().
			Err(err).
			Str("step", step.Description).
			Msg("Step has errored")
		return err
	}
	logger.Info().
		Str("step", step.Description).
		Msg("Finished running to run step")
	return nil
}

// GetContainerNamesAndIDs will return a map of container names to container ids
func (controller *Controller) GetContainerNamesAndIDs() (map[string]string, error) {
	logger.Trace().Msg("Getting container names and ids")
	namesAndIds, err := controller.docker.MapContainersNamesAndIDs()
	if err != nil {
		logger.Trace().
			Err(err).
			Msg("Failed to get container names and ids")
		return nil, err
	}
	logger.Trace().Msg("Successfully retrieved container names and ids")
	return namesAndIds, nil
}
