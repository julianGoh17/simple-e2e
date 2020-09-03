package models

import (
	"fmt"
	"os"
	"testing"

	"github.com/julianGoh17/simple-e2e/framework/internal"
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

func TestGettingStringVariableFromStepVariables(t *testing.T) {
	tables := []struct {
		step     *Step
		expected string
		err      error
	}{
		{
			&Step{
				Description: "This is a step",
				Variables: map[string]string{
					"TEST": "RandomVariable",
				},
			},
			"RandomVariable",
			nil,
		},
		{
			&Step{
				Description: "This is a step",
				Variables:   map[string]string{},
			},
			"",
			fmt.Errorf("Could not find variable '%s' in step.variables", "TEST"),
		},
	}

	for _, table := range tables {
		val, err := table.step.GetValueFromVariablesAsString("TEST")
		if table.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, table.expected, val)
		} else {
			assert.Error(t, table.err, err)
			assert.Equal(t, table.err.Error(), err.Error())
		}
	}
}

func TestGettingIntegerVariableFromStepVariables(t *testing.T) {
	tables := []struct {
		step     *Step
		expected int
		err      error
	}{
		{
			&Step{
				Description: "This is a step",
				Variables: map[string]string{
					"TEST": "RandomVariable",
				},
			},
			0,
			fmt.Errorf("Could not convert 'RandomVariable' to type 'int'"),
		},
		{
			&Step{
				Description: "This is a step",
				Variables:   map[string]string{},
			},
			0,
			fmt.Errorf("Could not find variable '%s' in step.variables", "TEST"),
		},
		{
			&Step{
				Description: "This is a step",
				Variables: map[string]string{
					"TEST": "0",
				},
			},
			0,
			nil,
		},
	}

	for _, table := range tables {
		val, err := table.step.GetValueFromVariablesAsInteger("TEST")
		if table.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, table.expected, val)
		} else {
			assert.Error(t, table.err, err)
			assert.Equal(t, table.err.Error(), err.Error())
		}
	}
}

func TestGettingFloat32VariableFromStepVariables(t *testing.T) {
	tables := []struct {
		step     *Step
		expected float32
		err      error
	}{
		{
			&Step{
				Description: "This is a step",
				Variables: map[string]string{
					"TEST": "RandomVariable",
				},
			},
			0,
			fmt.Errorf("Could not convert 'RandomVariable' to type 'float32'"),
		},
		{
			&Step{
				Description: "This is a step",
				Variables:   map[string]string{},
			},
			0,
			fmt.Errorf("Could not find variable '%s' in step.variables", "TEST"),
		},
		{
			&Step{
				Description: "This is a step",
				Variables: map[string]string{
					"TEST": "0",
				},
			},
			0,
			nil,
		},
	}

	for _, table := range tables {
		val, err := table.step.GetValueFromVariablesAsFloat32("TEST")
		if table.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, table.expected, val)
		} else {
			assert.Error(t, table.err, err)
			assert.Equal(t, table.err.Error(), err.Error())
		}
	}
}

func TestGettingFloat64VariableFromStepVariables(t *testing.T) {
	tables := []struct {
		step     *Step
		expected float64
		err      error
	}{
		{
			&Step{
				Description: "This is a step",
				Variables: map[string]string{
					"TEST": "RandomVariable",
				},
			},
			0,
			fmt.Errorf("Could not convert 'RandomVariable' to type 'float64'"),
		},
		{
			&Step{
				Description: "This is a step",
				Variables:   map[string]string{},
			},
			0,
			fmt.Errorf("Could not find variable '%s' in step.variables", "TEST"),
		},
		{
			&Step{
				Description: "This is a step",
				Variables: map[string]string{
					"TEST": "0",
				},
			},
			0,
			nil,
		},
	}

	for _, table := range tables {
		val, err := table.step.GetValueFromVariablesAsFloat64("TEST")
		if table.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, table.expected, val)
		} else {
			assert.Error(t, table.err, err)
			assert.Equal(t, table.err.Error(), err.Error())
		}
	}
}

func TestGettingBooleanVariableFromStepVariables(t *testing.T) {
	tables := []struct {
		step     *Step
		expected bool
		err      error
	}{
		{
			&Step{
				Description: "This is a step",
				Variables: map[string]string{
					"TEST": "RandomVariable",
				},
			},
			false,
			fmt.Errorf("Could not convert 'RandomVariable' to type 'bool'"),
		},
		{
			&Step{
				Description: "This is a step",
				Variables:   map[string]string{},
			},
			false,
			fmt.Errorf("Could not find variable '%s' in step.variables", "TEST"),
		},
		{
			&Step{
				Description: "This is a step",
				Variables: map[string]string{
					"TEST": "True",
				},
			},
			true,
			nil,
		},
	}

	for _, table := range tables {
		val, err := table.step.GetValueFromVariablesAsBoolean("TEST")
		if table.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, table.expected, val)
		} else {
			assert.Error(t, table.err, err)
			assert.Equal(t, table.err.Error(), err.Error())
		}
	}
}

func TestGettingStringArrayFromStepVariables(t *testing.T) {
	tables := []struct {
		step     *Step
		expected []string
		err      error
	}{
		{
			&Step{
				Description: "This is a step",
				Variables: map[string]string{
					"TEST": "RandomVariable",
				},
			},
			[]string{"RandomVariable"},
			nil,
		},
		{
			&Step{
				Description: "This is a step",
				Variables:   map[string]string{},
			},
			[]string{},
			fmt.Errorf("Could not find variable '%s' in step.variables", "TEST"),
		},
	}

	for _, table := range tables {
		val, err := table.step.GetValueFromVariablesAsStringArray("TEST")
		if table.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, table.expected, val)
		} else {
			assert.Error(t, table.err, err)
			assert.Equal(t, table.err.Error(), err.Error())
		}
	}
}

func TestGettingIntegerArrayFromStepVariables(t *testing.T) {
	tables := []struct {
		step     *Step
		expected []int
		err      error
	}{
		{
			&Step{
				Description: "This is a step",
				Variables: map[string]string{
					"TEST": "RandomVariable",
				},
			},
			[]int{},
			fmt.Errorf("Could not convert '%s' to type '[]int'", "RandomVariable"),
		},
		{
			&Step{
				Description: "This is a step",
				Variables:   map[string]string{},
			},
			[]int{},
			fmt.Errorf("Could not find variable '%s' in step.variables", "TEST"),
		},
		{
			&Step{
				Description: "This is a step",
				Variables: map[string]string{
					"TEST": "1,2,3,4",
				},
			},
			[]int{1, 2, 3, 4},
			nil,
		},
	}

	for _, table := range tables {
		val, err := table.step.GetValueFromVariablesAsIntegerArray("TEST")
		if table.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, table.expected, val)
		} else {
			assert.Error(t, table.err, err)
			assert.Equal(t, table.err.Error(), err.Error())
		}
	}
}

func TestGettingFloat32ArrayFromStepVariables(t *testing.T) {
	tables := []struct {
		step     *Step
		expected []float32
		err      error
	}{
		{
			&Step{
				Description: "This is a step",
				Variables: map[string]string{
					"TEST": "RandomVariable",
				},
			},
			[]float32{},
			fmt.Errorf("Could not convert '%s' to type '[]float32'", "RandomVariable"),
		},
		{
			&Step{
				Description: "This is a step",
				Variables:   map[string]string{},
			},
			[]float32{},
			fmt.Errorf("Could not find variable '%s' in step.variables", "TEST"),
		},
		{
			&Step{
				Description: "This is a step",
				Variables: map[string]string{
					"TEST": "1.111,2,3,4",
				},
			},
			[]float32{1.111, 2, 3, 4},
			nil,
		},
	}

	for _, table := range tables {
		val, err := table.step.GetValueFromVariablesAsFloat32Array("TEST")
		if table.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, table.expected, val)
		} else {
			assert.Error(t, table.err, err)
			assert.Equal(t, table.err.Error(), err.Error())
		}
	}
}

func TestGettingFloat64ArrayFromStepVariables(t *testing.T) {
	tables := []struct {
		step     *Step
		expected []float64
		err      error
	}{
		{
			&Step{
				Description: "This is a step",
				Variables: map[string]string{
					"TEST": "RandomVariable",
				},
			},
			[]float64{},
			fmt.Errorf("Could not convert '%s' to type '[]float64'", "RandomVariable"),
		},
		{
			&Step{
				Description: "This is a step",
				Variables:   map[string]string{},
			},
			[]float64{},
			fmt.Errorf("Could not find variable '%s' in step.variables", "TEST"),
		},
		{
			&Step{
				Description: "This is a step",
				Variables: map[string]string{
					"TEST": "1.111,2,3,4",
				},
			},
			[]float64{1.111, 2, 3, 4},
			nil,
		},
	}

	for _, table := range tables {
		val, err := table.step.GetValueFromVariablesAsFloat64Array("TEST")
		if table.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, table.expected, val)
		} else {
			assert.Error(t, table.err, err)
			assert.Equal(t, table.err.Error(), err.Error())
		}
	}
}

func TestGettingBooleanArrayFromStepVariables(t *testing.T) {
	tables := []struct {
		step     *Step
		expected []bool
		err      error
	}{
		{
			&Step{
				Description: "This is a step",
				Variables: map[string]string{
					"TEST": "RandomVariable",
				},
			},
			[]bool{},
			fmt.Errorf("Could not convert '%s' to type '[]bool'", "RandomVariable"),
		},
		{
			&Step{
				Description: "This is a step",
				Variables:   map[string]string{},
			},
			[]bool{},
			fmt.Errorf("Could not find variable '%s' in step.variables", "TEST"),
		},
		{
			&Step{
				Description: "This is a step",
				Variables: map[string]string{
					"TEST": "True,false,false",
				},
			},
			[]bool{true, false, false},
			nil,
		},
	}

	for _, table := range tables {
		val, err := table.step.GetValueFromVariablesAsBooleanArray("TEST")
		if table.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, table.expected, val)
		} else {
			assert.Error(t, table.err, err)
			assert.Equal(t, table.err.Error(), err.Error())
		}
	}
}

func TestVariableExistsFunction(t *testing.T) {
	step := &Step{
		Description: "This is a step",
		Variables: map[string]string{
			"TEST":   "RandomVariable",
			"RANDOM": "RandomVariable",
		},
	}

	noExistEnvVar := "NO_EXIST"
	tables := []struct {
		variablesToCheck []string
		err              error
	}{
		{
			[]string{"TEST", noExistEnvVar},
			fmt.Errorf("Could not find variable '%s' in step.variables", noExistEnvVar),
		},
		{
			[]string{"TEST"},
			nil,
		},
		{
			[]string{"TEST", "RANDOM"},
			nil,
		},
	}

	for _, table := range tables {
		err := step.CheckIfStepVariablesExists(table.variablesToCheck...)
		assert.Equal(t, table.err, err)
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

func TestStepPassesOrFailsTestCorrectlyDueToError(t *testing.T) {
	step := &Step{}
	errors := []error{
		fmt.Errorf("Random Error"),
		nil,
	}

	for _, err := range errors {
		step.SetErrored(err)
		if err != nil {
			assert.Equal(t, false, step.HasSucceeded())
		} else {
			assert.Equal(t, true, step.HasSucceeded())
		}
	}
}

func TestMain(m *testing.M) {
	internal.TestCoverageReaches85Percent(m)
}
