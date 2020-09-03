package models

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/go-yaml/yaml"
	"github.com/julianGoh17/simple-e2e/framework/internal"
	"github.com/julianGoh17/simple-e2e/framework/util"
	"github.com/stretchr/testify/assert"
)

func TestMultiStageUnmarshal(t *testing.T) {
	procedure := unmarshalYaml("multi-stage-test", t)

	assert.Equal(t, "Multi-stage Test", procedure.Name)
	assert.Equal(t, "DO NOT ALTER OR DELETE. This is an example of a multi stage test where each stage has multiple steps.", procedure.Description)

	assert.Equal(t, 2, len(procedure.GlobalVariables))
	assert.Equal(t, "true", procedure.GlobalVariables["MOCK_API"])
	assert.Equal(t, "true", procedure.GlobalVariables["DOES_SUCCEED"])

	assert.Equal(t, 2, len(procedure.Stages))

	firstStage := procedure.Stages[0]
	assert.Equal(t, "setup", firstStage.Name)
	assert.Equal(t, false, firstStage.AlwaysRuns)
	assert.Equal(t, 2, len(firstStage.Steps))
	assert.Equal(t, "Create mock API endpoint", firstStage.Steps[0].Description)
	assert.Equal(t, 2, len(firstStage.Steps[0].Variables))
	assert.Equal(t, "localhost:8080/v1/weather", firstStage.Steps[0].Variables["URL"])
	assert.Equal(t, "1", firstStage.Steps[0].Variables["RETRIES"])

	assert.Equal(t, "Load fake data to mock API endpoint", firstStage.Steps[1].Description)
	assert.Equal(t, 2, len(firstStage.Steps[1].Variables))
	assert.Equal(t, "localhost:8080/v1/weather", firstStage.Steps[1].Variables["URL"])
	assert.Equal(t, "{\"some\":\"data\"}", firstStage.Steps[1].Variables["PAYLOAD"])

	secondStage := procedure.Stages[1]
	assert.Equal(t, "test", secondStage.Name)
	assert.Equal(t, true, secondStage.AlwaysRuns)
	assert.Equal(t, 1, len(secondStage.Steps))
	assert.Equal(t, "Get data from API Endpoint", secondStage.Steps[0].Description)
	assert.Equal(t, 2, len(secondStage.Steps[0].Variables))
	assert.Equal(t, "localhost:8080/v1/weather", secondStage.Steps[0].Variables["URL"])
	assert.Equal(t, "{\"some\":\"data\"}", secondStage.Steps[0].Variables["PAYLOAD"])
}

func TestSimpleUnmarshal(t *testing.T) {
	procedure := unmarshalYaml("simple-test", t)

	assert.Equal(t, "Simple Test", procedure.Name)
	assert.Equal(t, "DO NOT ALTER OR DELETE. This is an example of a simple test with a single stage.", procedure.Description)

	assert.Equal(t, 0, len(procedure.GlobalVariables))

	assert.Equal(t, 1, len(procedure.Stages))

	firstStage := procedure.Stages[0]
	assert.Equal(t, "test", firstStage.Name)
	assert.Equal(t, 2, len(firstStage.Steps))
	assert.Equal(t, "Stand Up UI", firstStage.Steps[0].Description)
	assert.Equal(t, 1, len(firstStage.Steps[0].Variables))
	assert.Equal(t, "localhost:8080", firstStage.Steps[0].Variables["URL"])

	assert.Equal(t, "Click some buttons", firstStage.Steps[1].Description)
	assert.Equal(t, 0, len(firstStage.Steps[1].Variables))
}

func TestGlobalVariablesMultiStage(t *testing.T) {
	procedure := unmarshalYaml("multi-stage-test", t)
	assert.NoError(t, procedure.SetGlobalVariables())

	for key, value := range procedure.GlobalVariables {
		assert.Equal(t, value, os.Getenv(key))
	}
}

func unmarshalYaml(fileName string, t *testing.T) Procedure {
	internal.SetTestFilesRoot()
	var procedure Procedure

	body, err := ioutil.ReadFile(fmt.Sprintf("%s/examples/%s.yaml", os.Getenv(util.TestDirEnv), fileName))
	if err != nil {
		t.Errorf("unable to read file: %v", err)
	}

	err = yaml.Unmarshal([]byte(body), &procedure)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	return procedure
}
