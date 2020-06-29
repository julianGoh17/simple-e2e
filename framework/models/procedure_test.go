package models

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/go-yaml/yaml"
)

func TestMultiStageUnmarshal(t *testing.T) {
	procedure := unmarshalYaml("multi-stage-test", t)

	AssertEqual(t, procedure.Name, "Multi-stage Test")
	AssertEqual(t, procedure.Description, "DO NOT ALTER OR DELETE. This is an example of a multi stage test where each stage has multiple steps.")

	AssertEqual(t, len(procedure.GlobalVariables), 2)
	AssertEqual(t, procedure.GlobalVariables["MOCK_API"], "true")
	AssertEqual(t, procedure.GlobalVariables["DOES_SUCCEED"], "true")

	AssertEqual(t, len(procedure.Stages), 2)

	firstStage := procedure.Stages[0]
	AssertEqual(t, firstStage.Name, "setup")
	AssertEqual(t, len(firstStage.Steps), 2)
	AssertEqual(t, firstStage.Steps[0].Description, "Create mock API endpoint")
	AssertEqual(t, len(firstStage.Steps[0].Variables), 2)
	AssertEqual(t, firstStage.Steps[0].Variables["URL"], "localhost:8080/v1/weather")
	AssertEqual(t, firstStage.Steps[0].Variables["RETRIES"], "1")

	AssertEqual(t, firstStage.Steps[1].Description, "Load fake data to mock API endpoint")
	AssertEqual(t, len(firstStage.Steps[1].Variables), 2)
	AssertEqual(t, firstStage.Steps[1].Variables["URL"], "localhost:8080/v1/weather")
	AssertEqual(t, firstStage.Steps[1].Variables["PAYLOAD"], "{\"some\":\"data\"}")

	secondStage := procedure.Stages[1]
	AssertEqual(t, secondStage.Name, "test")
	AssertEqual(t, len(secondStage.Steps), 1)
	AssertEqual(t, secondStage.Steps[0].Description, "Get data from API Endpoint")
	AssertEqual(t, len(secondStage.Steps[0].Variables), 2)
	AssertEqual(t, secondStage.Steps[0].Variables["URL"], "localhost:8080/v1/weather")
	AssertEqual(t, secondStage.Steps[0].Variables["PAYLOAD"], "{\"some\":\"data\"}")
}

func TestSimpleUnmarshal(t *testing.T) {
	procedure := unmarshalYaml("simple-test", t)

	AssertEqual(t, procedure.Name, "Simple Test")
	AssertEqual(t, procedure.Description, "DO NOT ALTER OR DELETE. This is an example of a simple test with a single stage.")

	AssertEqual(t, len(procedure.GlobalVariables), 0)

	AssertEqual(t, len(procedure.Stages), 1)

	firstStage := procedure.Stages[0]
	AssertEqual(t, firstStage.Name, "test")
	AssertEqual(t, len(firstStage.Steps), 2)
	AssertEqual(t, firstStage.Steps[0].Description, "Stand Up UI")
	AssertEqual(t, len(firstStage.Steps[0].Variables), 1)
	AssertEqual(t, firstStage.Steps[0].Variables["URL"], "localhost:8080")

	AssertEqual(t, firstStage.Steps[1].Description, "Click some buttons")
	AssertEqual(t, len(firstStage.Steps[1].Variables), 0)
}

func unmarshalYaml(fileName string, t *testing.T) Procedure {
	var procedure Procedure
	filePath := getTestFileDirectory(fileName, t)

	body, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Errorf("unable to read file: %v", err)
	}

	err = yaml.Unmarshal([]byte(body), &procedure)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	return procedure
}

func getTestFileDirectory(fileName string, t *testing.T) string {
	dir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting Test YAML because: %s", err)
	}

	return fmt.Sprintf("%s/../../tests/examples/%s.yaml", dir, fileName)
}

// AssertEqual checks if values are equal
func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}

	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}
