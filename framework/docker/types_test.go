package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapContainerStatusToString(t *testing.T) {
	testCases := []struct {
		status         ContainerStatus
		expectedString string
	}{
		{
			Created,
			"Created",
		},
		{
			Completed,
			"Completed",
		},
		{
			Running,
			"Running",
		},
		{
			Errored,
			"Errored",
		},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expectedString, MapContainerStatusToString(testCase.status))
	}
}
