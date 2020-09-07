package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/julianGoh17/simple-e2e/framework/docker"
	"github.com/stretchr/testify/assert"
)

const (
	// TODO: move this into test tools
	existingImage           = "docker.io/library/alpine"
	invalidHost             = "random-host"
	unconnectableDockerHost = "http://localhost:9091"
	dockerHostEnv           = "DOCKER_HOST"
	invalidHostError        = "unable to parse docker host `random-host`"
)

var (
	canNotConnectToHostError = fmt.Sprintf("Cannot connect to the Docker daemon at %s. Is the docker daemon running?", unconnectableDockerHost)
)

func TestListCmd(t *testing.T) {
	listCmd := NewListCmd()

	assert.Equal(t, listCmd.Use, "list")
	assert.Equal(t, listCmd.Short, "Lists the container names and IDs")
	assert.Equal(t, listCmd.Long, `Lists the container names and IDs running on the host's daemon`)
}

func TestListCommandFails(t *testing.T) {
	testCases := []struct {
		host string
		err  error
	}{
		{
			invalidHost,
			fmt.Errorf(invalidHostError),
		},
		{
			unconnectableDockerHost,
			fmt.Errorf(canNotConnectToHostError),
		},
	}

	for _, testCase := range testCases {
		os.Setenv(dockerHostEnv, testCase.host)
		rootCmd := NewRootCmd()
		InitRootCmd(rootCmd)

		rootCmd.SetArgs([]string{"list"})
		err := rootCmd.Execute()
		assert.Error(t, err)
		assert.Equal(t, testCase.err.Error(), err.Error())
	}
	os.Unsetenv(dockerHostEnv)
}

func TestListCommandListsRunningDockerContainers(t *testing.T) {
	handler, err := docker.NewHandler()
	assert.NoError(t, err)

	containerName := "test"

	err = handler.CreateContainer(existingImage, containerName)
	assert.NoError(t, err)

	namesAndIDs, err := handler.MapContainersNamesAndIDs()
	assert.NoError(t, err)

	defer handler.DeleteContainer(containerName)

	rootCmd := NewRootCmd()
	InitRootCmd(rootCmd)
	read, written, rescue := beginCaptureOfTerminalOutput()

	rootCmd.SetArgs([]string{"list"})
	assert.NoError(t, rootCmd.Execute())

	output := endCaptureOfTerminalOutput(read, written, rescue)

	for name, ids := range namesAndIDs {
		assert.Contains(t, output, ids)
		assert.Contains(t, output, name)
	}
}

func TestGetTable(t *testing.T) {
	namesAndIDs := map[string]string{
		"name": "id",
	}

	table := getTable(namesAndIDs)
	assert.Equal(t, table.NumLines(), 1)
}
