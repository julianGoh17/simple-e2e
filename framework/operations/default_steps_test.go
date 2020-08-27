package operations

import (
	"testing"

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
