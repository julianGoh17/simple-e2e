package operations

import (
	"fmt"

	"github.com/julianGoh17/simple-e2e/framework/models"
)

func getDefaultSteps() map[string]func(step *models.Step) error {
	defaultSteps := map[string]func(step *models.Step) error{
		"Say hello to": SayHelloTo,
		"Pull image":   PullImage,
	}

	return defaultSteps
}

// SayHelloTo is just a placeholder function for testing
// Environmental Variables:
//   - NAME: Describes who to say hello to
func SayHelloTo(step *models.Step) error {
	traceStepEntrance(step)

	name, err := step.GetValueFromVariablesAsString("NAME")
	if err != nil {
		fmt.Println("Step failed!")
		step.SetFailed()
		return traceStepExit(step, err)
	}
	fmt.Printf("Hello there %s!\n", name)
	step.SetPassed()
	return traceStepExit(step, nil)
}

// PullImage will pull an image from a specified location onto the host machines daemon
// Environmental Variables:
// 	- IMAGE_REPOSITORY: The docker image repository to pull from
// 	- IMAGE_NAME: The name of the actual image to pull from
// 	- IMAGE_TAG: The image tag which is used to specify which version of the image to pull (if not set then will not do anything)
func PullImage(step *models.Step) error {
	if err := step.CheckIfStepVariablesExists("IMAGE_REPOSITORY", "IMAGE_NAME"); err != nil {
		return err
	}

	imageLocation, _ := step.GetValueFromVariablesAsString("IMAGE_REPOSITORY")
	imageName, _ := step.GetValueFromVariablesAsString("IMAGE_NAME")
	imageTag, _ := step.GetValueFromVariablesAsString("IMAGE_TAG")

	image := fmt.Sprintf("%s/%s", imageLocation, imageName)
	if imageTag != "" {
		image = fmt.Sprintf("%s:%s", image, imageTag)
	}

	err := step.Docker.PullImage(image)
	step.SetErrored(err)
	return traceStepExit(step, err)
}

func traceStepEntrance(step *models.Step) {
	logger.Trace().Str("step", step.Description)
	for key, val := range step.Variables {
		logger.Trace().Str(key, val)
	}
	logger.Trace().Msg("Beginning of step")
}

func traceStepExit(step *models.Step, err error) error {
	logger.Trace().Bool("hasStepPassed", step.HasSucceeded()).Err(err).Msg("End of step")
	return err
}
