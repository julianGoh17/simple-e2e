package operations

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	regexDescription        = "This is a '${string}'"
	testDescription         = "This is a 'string'"
	literalDescription      = "This is a literal string"
	invalidRegexDescription = `'${string}'r([a-z]+)gosdf[`
	regexKey                = "This is a '('[A-Za-z]+')'"
)

func TestAddTestStep(t *testing.T) {
	tables := []struct {
		descriptions []string
		literalKeys  []string
		regexKeys    []string
		willError    bool
	}{
		{
			[]string{regexDescription},
			[]string{},
			[]string{regexKey},
			false,
		},
		{
			[]string{literalDescription},
			[]string{literalDescription},
			[]string{},
			false,
		},
		{
			[]string{literalDescription, regexDescription},
			[]string{literalDescription},
			[]string{regexKey},
			false,
		},
		{
			[]string{invalidRegexDescription},
			[]string{},
			[]string{},
			true,
		},
	}

	for _, table := range tables {
		stepManager := NewStepManager()
		errMsg := fmt.Sprintf("Failed for descriptions '%s'", table.descriptions)
		for _, description := range table.descriptions {
			if table.willError {
				assert.Error(t, stepManager.AddStepToManager(description, testFuncPassStep), errMsg)
			} else {
				assert.NoError(t, stepManager.AddStepToManager(description, testFuncPassStep), errMsg)
			}
		}
		assert.Equal(t, len(table.literalKeys), len(stepManager.literalTestMethods), errMsg)
		assert.Equal(t, len(table.regexKeys), len(stepManager.regexTestMethods), errMsg)

		for key := range stepManager.literalTestMethods {
			assert.NotNil(t, stepManager.literalTestMethods[key], errMsg)
		}

		for key := range stepManager.regexTestMethods {
			assert.NotNil(t, stepManager.regexTestMethods[key], errMsg)
		}
	}
}

func TestGetTestStep(t *testing.T) {
	tables := []struct {
		description string
		testString  string
		willError   bool
	}{
		{
			regexDescription,
			testDescription,
			false,
		},
		{
			literalDescription,
			literalDescription,
			false,
		},
		{
			regexDescription,
			literalDescription + "'",
			true,
		},
		{
			literalDescription,
			literalDescription + " random suffix",
			true,
		},
	}

	for _, table := range tables {
		stepManager := NewStepManager()
		assert.NoError(t, stepManager.AddStepToManager(table.description, testFuncPassStep))
		_, err := stepManager.GetTestMethod(table.testString)
		if table.willError {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestAddingDuplicateRegexTestSteps(t *testing.T) {
	stepManager := NewStepManager()

	assert.NoError(t, stepManager.AddStepToManager(regexDescription, testFuncPassStep))
	assert.Error(t, stepManager.AddStepToManager(regexDescription, testFuncPassStep))

	assert.Equal(t, 0, len(stepManager.literalTestMethods))
	assert.Equal(t, 1, len(stepManager.regexTestMethods))
}

func TestAddingDuplicateLiteralTestSteps(t *testing.T) {
	stepManager := NewStepManager()

	assert.NoError(t, stepManager.AddStepToManager(literalDescription, testFuncPassStep))
	assert.Error(t, stepManager.AddStepToManager(literalDescription, testFuncPassStep))

	assert.Equal(t, 1, len(stepManager.literalTestMethods))
	assert.Equal(t, 0, len(stepManager.regexTestMethods))
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	rc := m.Run()

	// rc 0 means we've passed,
	// and CoverMode will be non empty if run with -cover
	if rc == 0 && testing.CoverMode() != "" {
		c := testing.Coverage()
		if c < 0.85 {
			fmt.Println("Tests passed but coverage failed at", c)
			rc = -1
		}
	}
	os.Exit(rc)
}
