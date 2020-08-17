package models

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGettingDescriptionVariables(t *testing.T) {
	tables := []struct {
		Step                         *Step
		ExpectedDescriptionVariables []string
	}{
		{
			&Step{
				Description: "This is a 'variable'",
				Variables:   make(map[string]string),
			},
			[]string{"variable"},
		},
		{
			&Step{
				Description: "This is a 'variable' and this is 'another'",
				Variables:   make(map[string]string),
			},
			[]string{"variable", "another"},
		},
		{
			&Step{
				Description: "This is should fail '",
				Variables:   make(map[string]string),
			},
			nil,
		},
		{
			&Step{
				Description: "This is has no variables",
				Variables:   map[string]string{},
			},
			[]string{},
		},
	}

	for _, table := range tables {
		descriptionVariables, err := table.Step.GetDescriptionVariables()
		if table.ExpectedDescriptionVariables == nil {
			assert.Error(t, err)
		} else {
			assert.Equal(t, table.ExpectedDescriptionVariables, descriptionVariables)
		}
	}
}

func TestHasSucceed(t *testing.T) {
	step := &Step{}
	assert.False(t, step.HasSucceeded())

	step.SetPassed()
	assert.True(t, step.HasSucceeded())

	step.SetFailed()
	assert.False(t, step.HasSucceeded())
}

func TestStepCanGetEnvVar(t *testing.T) {
	key := "key"
	value := "value"
	step := &Step{}
	os.Setenv(key, value)
	assert.Equal(t, step.GetGlobalVariable(key), value)
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
