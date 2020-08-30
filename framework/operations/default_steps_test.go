package operations

import (
	"fmt"
	"testing"

	"github.com/julianGoh17/simple-e2e/framework/docker"
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
			fmt.Errorf("Could not find variable with name '%s' in step.variables", "IMAGE_NAME"),
		}, {
			&models.Step{
				Variables: map[string]string{
					"IMAGE_NAME": "alpine",
				},
				Docker: docker,
			},
			fmt.Errorf("Could not find variable with name '%s' in step.variables", "IMAGE_REPOSITORY"),
		},
		{
			&models.Step{
				Variables: map[string]string{
					"IMAGE_REPOSITORY": "docker.io/library",
					"IMAGE_NAME":       "blah",
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
				"IMAGE_NAME":       "alpine",
				"IMAGE_TAG":        "latest",
			},
			Docker: docker,
		},
		{
			Variables: map[string]string{
				"IMAGE_REPOSITORY": "docker.io/library",
				"IMAGE_NAME":       "alpine",
			},
			Docker: docker,
		},
	}

	for _, step := range steps {
		err := PullImage(step)
		assert.NoError(t, err)
	}
}
