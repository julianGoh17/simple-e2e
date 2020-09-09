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
		{
			Exited,
			"Exited",
		},
		{
			Paused,
			"Paused",
		},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expectedString, MapContainerStatusToString(testCase.status))
	}
}

func TestMapStateToStatus(t *testing.T) {
	testCases := []struct {
		str            string
		expectedStatus ContainerStatus
	}{
		{
			"created",
			Created,
		},
		{
			"running",
			Running,
		},
		{
			"completed",
			Completed,
		},
		{
			"errored",
			Errored,
		},
		{
			"exited",
			Exited,
		},
		{
			"paused",
			Paused,
		},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expectedStatus, MapStateToStatus(testCase.str))
	}
}
