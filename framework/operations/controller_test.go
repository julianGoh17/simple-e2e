package operations

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/julianGoh17/simple-e2e/framework/internal"
	models "github.com/julianGoh17/simple-e2e/framework/models"
	"github.com/julianGoh17/simple-e2e/framework/util"
	"github.com/stretchr/testify/assert"
)

const illFormatted = `
no: good
`

const correctlyFormated = `
name: example-test
description: example description
stages:
  - name: example-stage
    steps: 
      - description: example-step
`

const multiStageRun = `
name: example-test
description: example description
stages:
  - name: example-stage
    steps:
      - description: "example-step"
  - name: example-stage2
    alwaysRuns: true
    steps:
      - description: "example-step"`

const noName = `
description: example description
stages:
  - description: step
	
`

const noDescription = `
name: example-test
stages:
  - description: step
`

const noStages = `
name: example-test
description: example description
`

func TestYamlFormatting(t *testing.T) {
	controller, err := NewController()
	assert.NoError(t, err)

	testFileOutcomes := []struct {
		testFile  string
		isCorrect bool
	}{
		{illFormatted, false},
		{noName, false},
		{noDescription, false},
		{noStages, false},
		{correctlyFormated, true},
	}

	for _, outcome := range testFileOutcomes {
		if outcome.isCorrect {
			assert.NoError(t, controller.SetProcedure([]byte(outcome.testFile)))
		} else {
			assert.Error(t, controller.SetProcedure([]byte(outcome.testFile)))
		}
	}
}

func TestRunTest(t *testing.T) {
	internal.SetTestFilesRoot()
	controller, err := NewController()
	assert.NoError(t, err)

	tables := []struct {
		testFile     string
		testLocation string
		willError    bool
	}{
		{"test.yaml", os.Getenv(util.TestDirEnv), false},
		{"test.yaml", "random", true},
	}

	for _, table := range tables {
		err := controller.RunTest(fmt.Sprintf("%s/%s", table.testLocation, table.testFile))
		if table.willError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestRunStages(t *testing.T) {
	testFileOutcomes := []struct {
		testFile     string
		testFunction func(*models.Step) error
		stages       []string
		willError    bool
	}{
		{
			correctlyFormated,
			testFuncPassStep,
			[]string{},
			false,
		},
		{
			correctlyFormated,
			testFuncFailStep,
			[]string{},
			true,
		},
		{
			correctlyFormated,
			testFuncErrorStep,
			[]string{},
			true,
		},
		{
			correctlyFormated,
			testFuncPassStep,
			[]string{"example-stage"},
			false,
		},
		{
			correctlyFormated,
			testFuncFailStep,
			[]string{"example-stage"},
			true,
		},
		{
			illFormatted,
			testFuncFailStep,
			[]string{},
			true,
		},
	}

	for _, outcome := range testFileOutcomes {
		controller, err := NewController()
		assert.NoError(t, err)
		assert.NoError(t, controller.AddTestStep("example-step", outcome.testFunction))
		if outcome.willError {
			assert.Error(t, controller.runTest([]byte(outcome.testFile), outcome.stages...))
		} else {
			assert.NoError(t, controller.runTest([]byte(outcome.testFile), outcome.stages...))
		}
	}
}

func TestGetContainerInformationFails(t *testing.T) {
	os.Setenv(internal.DockerHostEnv, internal.UnconnectableDockerHost)
	defer os.Unsetenv(internal.DockerHostEnv)
	controller, err := NewController()
	assert.Error(t, err)
	assert.Equal(t, internal.ErrCanNotConnectToHost.Error(), err.Error())

	_, err = controller.GetContainerInfo(true)
	assert.Error(t, err)
	assert.Equal(t, internal.ErrCanNotConnectToHost.Error(), err.Error())
}

func TestGetContainerInformationPasses(t *testing.T) {
	controller, err := NewController()
	assert.NoError(t, err)
	assert.NoError(t, controller.docker.PullImage(internal.ExistingImage))
	assert.NoError(t, controller.docker.CreateContainer(internal.ExistingImage, internal.ExistingContainerName, []string{}))
	defer controller.docker.DeleteContainer(internal.ExistingContainerName)

	containerInfos, err := controller.GetContainerInfo(true)
	assert.NoError(t, err)

	hasCreatedContainer := false
	for _, containerInfo := range containerInfos {
		if containerInfo.Image == internal.ExistingImage {
			hasCreatedContainer = true
		}
	}

	assert.Equal(t, true, hasCreatedContainer, "Could not find created container when listing containers")

}

func TestWillRunAlwaysRunsEvenWhenFail(t *testing.T) {
	controller, err := NewController()
	assert.NoError(t, err)
	assert.NoError(t, controller.AddTestStep("example-step", testFuncFailStep))
	assert.Error(t, controller.runTest([]byte(multiStageRun)))
}

func TestFailsWhenCanNotGetStep(t *testing.T) {
	controller, err := NewController()
	assert.NoError(t, err)
	assert.Error(t, controller.runTest([]byte(multiStageRun)))
}

func testFuncPassStep(step *models.Step) error {
	step.SetPassed()
	return nil
}

func testFuncFailStep(step *models.Step) error {
	step.SetFailed()
	return nil
}

func testFuncErrorStep(step *models.Step) error {
	return errors.New("This will error")
}
