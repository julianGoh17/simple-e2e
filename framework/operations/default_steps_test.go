package operations

import (
	"fmt"
	"testing"

	"github.com/julianGoh17/simple-e2e/framework/docker"
	"github.com/julianGoh17/simple-e2e/framework/internal"
	"github.com/julianGoh17/simple-e2e/framework/models"
	"github.com/stretchr/testify/assert"
)

func TestSayHelloToStep(t *testing.T) {
	tables := []struct {
		step      *models.Step
		willError bool
	}{
		{
			&models.Step{},
			true,
		},
		{
			&models.Step{
				Description: "Test step",
				Variables: map[string]string{
					"NAME": "me",
				},
			},
			false,
		},
	}

	for _, table := range tables {
		err := SayHelloTo(table.step)
		if table.willError {
			assert.Error(t, err)
			assert.Equal(t, table.step.HasSucceeded(), false)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, table.step.HasSucceeded(), true)
		}
	}
}

func TestPullImageStepErrors(t *testing.T) {
	docker, err := docker.NewHandler()
	assert.NoError(t, err)

	testCases := []struct {
		step *models.Step
		err  error
	}{
		{
			&models.Step{
				Variables: map[string]string{
					"IMAGE_REPOSITORY": "docker.io/library/",
				},
				Docker: docker,
			},
			fmt.Errorf("Could not find variable '%s' in step.variables", "IMAGE"),
		}, {
			&models.Step{
				Variables: map[string]string{
					"IMAGE": "alpine",
				},
				Docker: docker,
			},
			fmt.Errorf("Could not find variable '%s' in step.variables", "IMAGE_REPOSITORY"),
		},
		{
			&models.Step{
				Variables: map[string]string{
					"IMAGE_REPOSITORY": "docker.io/library",
					"IMAGE":            "blah",
					"IMAGE_TAG":        "invalid-image-tag",
				},
				Docker: docker,
			},
			fmt.Errorf("Error response from daemon: pull access denied for blah, repository does not exist or may require 'docker login': denied: requested access to the resource is denied"),
		},
	}

	for _, testCase := range testCases {
		err := PullImage(testCase.step)
		assert.Error(t, err)
		assert.Equal(t, testCase.err.Error(), err.Error())
	}
}

func TestPullImageStepPasses(t *testing.T) {
	docker, err := docker.NewHandler()
	assert.NoError(t, err)

	steps := []*models.Step{
		{
			Variables: map[string]string{
				"IMAGE_REPOSITORY": "docker.io/library",
				"IMAGE":            "alpine",
				"IMAGE_TAG":        "latest",
			},
			Docker: docker,
		},
		{
			Variables: map[string]string{
				"IMAGE_REPOSITORY": "docker.io/library",
				"IMAGE":            "alpine",
			},
			Docker: docker,
		},
	}

	for _, step := range steps {
		err := PullImage(step)
		assert.NoError(t, err)
	}
}

func TestBuildImageStepFails(t *testing.T) {
	docker, err := docker.NewHandler()
	assert.NoError(t, err)

	testCases := []struct {
		step *models.Step
		err  error
	}{
		{
			&models.Step{
				Variables: map[string]string{
					"IMAGE": "test",
				},
				Docker: docker,
			},
			fmt.Errorf("Could not find variable '%s' in step.variables", "DOCKERFILE"),
		},
		{
			&models.Step{
				Variables: map[string]string{
					"DOCKERFILE": "DockerfileThatDoesNotExist",
				},
				Docker: docker,
			},
			fmt.Errorf("Could not find variable '%s' in step.variables", "IMAGE"),
		},
		{
			&models.Step{
				Variables: map[string]string{
					"DOCKERFILE": "DockerfileThatDoesNotExist",
					"IMAGE":      "test",
				},
				Docker: docker,
			},
			fmt.Errorf("open /home/e2e/Dockerfiles/DockerfileThatDoesNotExist: no such file or directory")},
	}

	for _, testCase := range testCases {
		err := BuildImage(testCase.step)
		assert.Error(t, err)
		assert.Equal(t, testCase.err.Error(), err.Error())
	}
}

func TestCreateContainerStepFails(t *testing.T) {
	docker, err := docker.NewHandler()
	assert.NoError(t, err)

	testCases := []struct {
		step *models.Step
		err  error
	}{
		{
			&models.Step{
				Variables: map[string]string{
					"IMAGE": "test",
				},
				Docker: docker,
			},
			fmt.Errorf("Could not find variable '%s' in step.variables", "CONTAINER_NAME"),
		},
		{
			&models.Step{
				Variables: map[string]string{
					"CONTAINER_NAME": "test",
				},
				Docker: docker,
			},
			fmt.Errorf("Could not find variable '%s' in step.variables", "IMAGE"),
		},
	}

	for _, testCase := range testCases {
		err := CreateContainer(testCase.step)
		assert.Error(t, err)
		assert.Equal(t, testCase.err.Error(), err.Error())
	}
}

func TestDeleteContainerStepFails(t *testing.T) {
	docker, err := docker.NewHandler()
	assert.NoError(t, err)

	testCases := []struct {
		step *models.Step
		err  error
	}{
		{
			&models.Step{
				Variables: map[string]string{},
				Docker:    docker,
			},
			fmt.Errorf("Could not find variable '%s' in step.variables", "CONTAINER_NAME"),
		},
	}

	for _, testCase := range testCases {
		err := DeleteContainer(testCase.step)
		assert.Error(t, err)
		assert.Equal(t, testCase.err.Error(), err.Error())
	}
}

func TestStartContainerStepFails(t *testing.T) {
	docker, err := docker.NewHandler()
	assert.NoError(t, err)

	testCases := []struct {
		step *models.Step
		err  error
	}{
		{
			&models.Step{
				Variables: map[string]string{},
				Docker:    docker,
			},
			fmt.Errorf("Could not find variable '%s' in step.variables", "CONTAINER_NAME"),
		},
		{
			&models.Step{
				Variables: map[string]string{
					"CONTAINER_NAME": internal.NonExistentContainerName,
				},
				Docker: docker,
			},
			internal.ErrCanNotFindNonExistentContainerInRegistry,
		},
	}

	for _, testCase := range testCases {
		err := StartContainer(testCase.step)
		assert.Error(t, err)
		assert.Equal(t, testCase.err.Error(), err.Error())
	}
}

func TestRestartContainerFails(t *testing.T) {
	docker, err := docker.NewHandler()
	assert.NoError(t, err)

	testCases := []struct {
		step *models.Step
		err  error
	}{
		{
			&models.Step{
				Variables: map[string]string{},
				Docker:    docker,
			},
			fmt.Errorf("Could not find variable '%s' in step.variables", "CONTAINER_NAME"),
		},
		{
			&models.Step{
				Variables: map[string]string{
					"CONTAINER_NAME": internal.NonExistentContainerName,
				},
				Docker: docker,
			},
			internal.ErrCanNotFindNonExistentContainerInRegistry,
		},
	}

	for _, testCase := range testCases {
		err := RestartContainer(testCase.step)
		assert.Error(t, err)
		assert.Equal(t, testCase.err.Error(), err.Error())
	}
}

func TestStopContainerStepFails(t *testing.T) {
	docker, err := docker.NewHandler()
	assert.NoError(t, err)

	testCases := []struct {
		step *models.Step
		err  error
	}{
		{
			&models.Step{
				Variables: map[string]string{
					"CONTAINER_NAME": "TEST",
				},
				Docker: docker,
			},
			fmt.Errorf("Could not find variable '%s' in step.variables", "TIME_DURATION"),
		},
		{
			&models.Step{
				Variables: map[string]string{
					"TIME_DURATION": "TEST",
				},
				Docker: docker,
			},
			fmt.Errorf("Could not find variable '%s' in step.variables", "CONTAINER_NAME"),
		},
		{
			&models.Step{
				Variables: map[string]string{
					"CONTAINER_NAME": internal.NonExistentContainerName,
					"TIME_DURATION":  "1m",
				},
				Docker: docker,
			},
			internal.ErrCanNotFindNonExistentContainerInRegistry,
		},
	}

	for _, testCase := range testCases {
		err := StopContainer(testCase.step)
		assert.Error(t, err)
		assert.Equal(t, testCase.err.Error(), err.Error())
	}
}

func TestPauseContainerStepFails(t *testing.T) {
	docker, err := docker.NewHandler()
	assert.NoError(t, err)

	testCases := []struct {
		step *models.Step
		err  error
	}{
		{
			&models.Step{
				Variables: map[string]string{
					"CONTAINER_NAME": "TEST",
				},
				Docker: docker,
			},
			fmt.Errorf("Could not find variable '%s' in step.variables", "TIME_DURATION"),
		},
		{
			&models.Step{
				Variables: map[string]string{
					"TIME_DURATION": "TEST",
				},
				Docker: docker,
			},
			fmt.Errorf("Could not find variable '%s' in step.variables", "CONTAINER_NAME"),
		},
		{
			&models.Step{
				Variables: map[string]string{
					"CONTAINER_NAME": internal.NonExistentContainerName,
					"TIME_DURATION":  "1m",
				},
				Docker: docker,
			},
			internal.ErrCanNotFindNonExistentContainerInRegistry,
		},
	}

	for _, testCase := range testCases {
		err := PauseContainer(testCase.step)
		assert.Error(t, err)
		assert.Equal(t, testCase.err.Error(), err.Error())
	}
}

func TestCreateAndDeleteContainerStepPasses(t *testing.T) {
	docker, err := docker.NewHandler()
	assert.NoError(t, err)

	step := &models.Step{
		Variables: map[string]string{
			"CONTAINER_NAME": "test-container",
			"IMAGE":          internal.ExistingImage,
		},
		Docker: docker,
	}

	assert.NoError(t, CreateContainer(step))
	assert.NoError(t, DeleteContainer(step))
}

func TestBuildImageStepPasses(t *testing.T) {
	internal.SetDockerfilesRoot()
	docker, err := docker.NewHandler()
	assert.NoError(t, err)

	step := &models.Step{
		Variables: map[string]string{
			"DOCKERFILE": "Dockerfile.simple",
			"IMAGE":      "test:e2e-test",
		},
		Docker: docker,
	}

	err = BuildImage(step)
	assert.NoError(t, err)
}
