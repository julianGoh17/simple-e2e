package operations

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/julianGoh17/simple-e2e/framework/models"
	"github.com/julianGoh17/simple-e2e/framework/util"
)

func getDefaultSteps() map[string]func(step *models.Step) error {
	defaultSteps := map[string]func(step *models.Step) error{
		"Say hello to": SayHelloTo,
		"Pull image":   PullImage,
		"Build image":  BuildImage,
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

	return traceStepExit(step, step.Docker.PullImage(image))
}

// BuildImage will build the specified image from the specified Dockerfile located in the 'Dockerfiles' directory
// Environmental Variables:
// 	- DOCKERFILE: The name of the Dockerfile to be built
//  - IMAGE_NAME: The build tag of the Docker image
func BuildImage(step *models.Step) error {
	traceStepEntrance(step)
	if err := step.CheckIfStepVariablesExists("DOCKERFILE", "IMAGE_NAME"); err != nil {
		return err
	}

	dockerfile, _ := step.GetValueFromVariablesAsString("DOCKERFILE")
	buildTag, _ := step.GetValueFromVariablesAsString("IMAGE_NAME")

	build, err := createDockerfileBuild(dockerfile)

	if err != nil {
		return err
	}
	return traceStepExit(step, step.Docker.BuildImage(build, dockerfile, buildTag))
}

func createDockerfileBuild(dockerfile string) (io.Reader, error) {
	logger.Trace().Str("Dockerfile", dockerfile).Msg("Creating tar for dockerfile")
	dockerfileBytes, err := readDockerfile(dockerfile)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()
	tarHeader := &tar.Header{
		Name: dockerfile,
		Size: int64(len(dockerfileBytes)),
	}
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		return nil, err
	}
	_, err = tw.Write(dockerfileBytes)
	if err != nil {
		return nil, err
	}
	logger.Trace().Str("Dockerfile", dockerfile).Msg("Finished creating tar for dockerfile")
	return bytes.NewReader(buf.Bytes()), nil
}

func readDockerfile(dockerfile string) ([]byte, error) {
	logger.Trace().Str("Dockerfile", dockerfile).Msg("Reading dockerfile")
	dockerfileReader, err := os.Open(getDockerfilePath(dockerfile))
	if err != nil {
		return nil, err
	}
	logger.Trace().Str("Dockerfile", dockerfile).Msg("Finished reading dockerfile")
	return ioutil.ReadAll(dockerfileReader)
}

func getDockerfilePath(dockerfile string) string {
	return fmt.Sprintf("%s/%s", config.GetOrDefault(util.DockerfileDirEnv), dockerfile)
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
