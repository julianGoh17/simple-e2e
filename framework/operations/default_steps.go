package operations

import (
	"fmt"

	"github.com/julianGoh17/simple-e2e/framework/models"
)

func getDefaultSteps() map[string]func(step *models.Step) error {
	defaultSteps := map[string]func(step *models.Step) error{
		"Say hello to":     SayHelloTo,
		"Pull image":       PullImage,
		"Build image":      BuildImage,
		"Create container": CreateContainer,
		"Delete container": DeleteContainer,
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
// 	- IMAGE: The name of the actual image to pull from
// 	- IMAGE_TAG: The image tag which is used to specify which version of the image to pull (if not set then will not do anything)
func PullImage(step *models.Step) error {
	if err := step.CheckIfStepVariablesExists("IMAGE_REPOSITORY", "IMAGE"); err != nil {
		return traceStepExit(step, err)
	}

	imageLocation, _ := step.GetValueFromVariablesAsString("IMAGE_REPOSITORY")
	imageName, _ := step.GetValueFromVariablesAsString("IMAGE")
	imageTag, _ := step.GetValueFromVariablesAsString("IMAGE_TAG")

	image := fmt.Sprintf("%s/%s", imageLocation, imageName)
	if imageTag != "" {
		image = fmt.Sprintf("%s:%s", image, imageTag)
	}

	return traceStepExit(step, step.Docker.PullImage(image))
}

// BuildImage will build the specified image from the specified Dockerfile located in the 'Dockerfiles' directory
// Environmental Variables:
// 	- DOCKERFILE: The name of the Dockerfile to be built
//  - IMAGE: The name to give the built image
func BuildImage(step *models.Step) error {
	traceStepEntrance(step)
	if err := step.CheckIfStepVariablesExists("DOCKERFILE", "IMAGE"); err != nil {
		return traceStepExit(step, err)
	}

	dockerfile, _ := step.GetValueFromVariablesAsString("DOCKERFILE")
	buildTag, _ := step.GetValueFromVariablesAsString("IMAGE")

	return traceStepExit(step, step.Docker.BuildImage(dockerfile, buildTag))
}

// CreateContainer will create a container (but will not run the container) from an image and create a ContainerManager to manage that Container
// Environmental Variables:
// 	- IMAGE: The name of the image to create the container with
//  - CONTAINER_NAME: The name to give to the created container
func CreateContainer(step *models.Step) error {
	traceStepEntrance(step)
	if err := step.CheckIfStepVariablesExists("IMAGE", "CONTAINER_NAME"); err != nil {
		return traceStepExit(step, err)
	}

	image, _ := step.GetValueFromVariablesAsString("IMAGE")
	containerName, _ := step.GetValueFromVariablesAsString("CONTAINER_NAME")

	return traceStepExit(step, step.Docker.CreateContainer(image, containerName, []string{}))
}

// DeleteContainer will delete a container (that has been registered with the framework) based on the container name given.
// Environmental Variables:
//  - CONTAINER_NAME: The name of the container to delete
func DeleteContainer(step *models.Step) error {
	traceStepEntrance(step)

	containerName, err := step.GetValueFromVariablesAsString("CONTAINER_NAME")
	if err != nil {
		return traceStepExit(step, err)
	}
	return traceStepExit(step, step.Docker.DeleteContainer(containerName))
}

func traceStepEntrance(step *models.Step) {
	trace := logger.Trace().Str("description", step.Description)
	for key, val := range step.Variables {
		trace.Str(key, val)
	}
	trace.Msg("Step.variables")
	logger.Info().Str("description", step.Description).Msg("Beginning of step")
}

func traceStepExit(step *models.Step, err error) error {
	logger.Info().Bool("hasStepPassed", step.HasSucceeded()).Err(err).Msg("End of step")
	step.SetErrored(err)
	return err
}
