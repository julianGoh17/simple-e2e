package operations

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	models "github.com/julianGoh17/simple-e2e/framework/models"
	"github.com/stretchr/testify/assert"
)

const (
	regexDescription        = "This is a '${string}'"
	literalDescription      = "This is a literal string"
	invalidRegexDescription = `'${string}'r([a-z]+)gosdf[`
	regexKey                = "This is a '('[A-Za-z]+')'"
)

func TestAddTestStep(t *testing.T) {
	tables := []struct {
		descriptions []string
		literalKeys  []string
		regexKeys    []string
		controller   *Controller
		willError    bool
	}{
		{
			[]string{regexDescription},
			[]string{},
			[]string{regexKey},
			NewController(),
			false,
		},
		{
			[]string{literalDescription},
			[]string{literalDescription},
			[]string{},
			NewController(),
			false,
		},
		{
			[]string{literalDescription, regexDescription},
			[]string{literalDescription},
			[]string{regexKey},
			NewController(),
			false,
		},
		{
			[]string{invalidRegexDescription},
			[]string{},
			[]string{},
			NewController(),
			true,
		},
	}

	for _, table := range tables {
		errMsg := fmt.Sprintf("Failed for descriptions '%s'", table.descriptions)
		for _, description := range table.descriptions {
			if table.willError {
				assert.Error(t, table.controller.AddTestStep(description, testFuncPassStep), errMsg)
			} else {
				assert.NoError(t, table.controller.AddTestStep(description, testFuncPassStep), errMsg)
			}
		}
		assert.Equal(t, len(table.literalKeys), len(table.controller.literalTestMethods), errMsg)
		assert.Equal(t, len(table.regexKeys), len(table.controller.regexTestMethods), errMsg)

		for key := range table.controller.literalTestMethods {
			assert.NotNil(t, table.controller.literalTestMethods[key], errMsg)
		}

		for key := range table.controller.regexTestMethods {
			assert.NotNil(t, table.controller.regexTestMethods[key], errMsg)
		}
	}
}

func TestAddingDuplicateRegexTestSteps(t *testing.T) {
	controller := NewController()

	assert.NoError(t, controller.AddTestStep(regexDescription, testFuncPassStep))
	assert.Error(t, controller.AddTestStep(regexDescription, testFuncPassStep))

	assert.Equal(t, 0, len(controller.literalTestMethods))
	assert.Equal(t, 1, len(controller.regexTestMethods))
}

func TestAddingDuplicateLiteralTestSteps(t *testing.T) {
	controller := NewController()

	assert.NoError(t, controller.AddTestStep(literalDescription, testFuncPassStep))
	assert.Error(t, controller.AddTestStep(literalDescription, testFuncPassStep))

	assert.Equal(t, 1, len(controller.literalTestMethods))
	assert.Equal(t, 0, len(controller.regexTestMethods))
}

func TestCorrectlyFormattedYaml(t *testing.T) {
	controller := NewController()

	data := unmarshalYaml("multi-stage-test", t)
	assert.NoError(t, controller.SetProcedure(data))
}

func TestIncorrectlyFormattedYaml(t *testing.T) {
	controller := NewController()

	data := unmarshalYaml("ill-formatted", t)
	assert.Error(t, controller.SetProcedure(data))
}

func testFuncPassStep(step *models.Step) {
	step.SetPassed()
}

func testFuncFailStep(step *models.Step) {
	step.SetFailed()
}

func unmarshalYaml(fileName string, t *testing.T) []byte {
	filePath := getTestFileDirectory(fileName, t)

	body, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Errorf("unable to read file: %v", err)
	}

	return body
}

func getTestFileDirectory(fileName string, t *testing.T) string {
	dir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting Test YAML because: %s", err)
	}

	return fmt.Sprintf("%s/../../tests/examples/%s.yaml", dir, fileName)
}
