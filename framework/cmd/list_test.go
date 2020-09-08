package cmd

import (
	"os"
	"testing"

	"github.com/julianGoh17/simple-e2e/framework/docker"
	"github.com/julianGoh17/simple-e2e/framework/internal"
	"github.com/stretchr/testify/assert"
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
			internal.InvalidDockerHost,
			internal.ErrInvalidHost,
		},
		{
			internal.UnconnectableDockerHost,
			internal.ErrCanNotConnectToHost,
		},
	}

	for _, testCase := range testCases {
		os.Setenv(internal.DockerHostEnv, testCase.host)
		rootCmd := NewRootCmd()
		InitRootCmd(rootCmd)

		rootCmd.SetArgs([]string{"list"})
		err := rootCmd.Execute()
		assert.Error(t, err)
		assert.Equal(t, testCase.err.Error(), err.Error())
	}
	os.Unsetenv(internal.DockerHostEnv)
}

func TestListCommandListsAllRunningDockerContainers(t *testing.T) {
	handler, err := docker.NewHandler()
	assert.NoError(t, err)

	containerName := "test"

	assert.NoError(t, handler.PullImage(internal.ExistingImage))
	assert.NoError(t, handler.CreateContainer(internal.ExistingImage, containerName))

	containers, err := handler.GetContainerInfo(true)
	assert.NoError(t, err)

	defer handler.DeleteContainer(containerName)

	rootCmd := NewRootCmd()
	InitRootCmd(rootCmd)
	read, written, rescue := beginCaptureOfTerminalOutput()

	rootCmd.SetArgs([]string{"list", "-a", "true"})
	assert.NoError(t, rootCmd.Execute())

	output := endCaptureOfTerminalOutput(read, written, rescue)

	for _, container := range containers {
		assert.Contains(t, output, container.Name)
		assert.Contains(t, output, container.ID)
		assert.Contains(t, output, docker.MapContainerStatusToString(container.Status))
	}
}

func TestListCommandListsCertainDockerContainers(t *testing.T) {
	handler, err := docker.NewHandler()
	assert.NoError(t, err)

	containerName := "test"

	assert.NoError(t, handler.PullImage(internal.ExistingImage))
	assert.NoError(t, handler.CreateContainer(internal.ExistingImage, containerName))

	containers, err := handler.GetContainerInfo(false)
	assert.NoError(t, err)

	defer handler.DeleteContainer(containerName)

	rootCmd := NewRootCmd()
	InitRootCmd(rootCmd)
	read, written, rescue := beginCaptureOfTerminalOutput()

	rootCmd.SetArgs([]string{"list", "-a", "false"})
	assert.NoError(t, rootCmd.Execute())

	output := endCaptureOfTerminalOutput(read, written, rescue)

	for _, container := range containers {
		assert.NotEqual(t, docker.Completed, container.Status)
		assert.NotEqual(t, docker.Exited, container.Status)
		assert.Contains(t, output, container.Name)
		assert.Contains(t, output, container.ID)
		assert.NotContains(t, output, docker.MapContainerStatusToString(docker.Completed))
		assert.Contains(t, output, docker.MapContainerStatusToString(container.Status))
	}
}

func TestGetTable(t *testing.T) {
	containers := []*docker.ContainerInfo{
		{
			Name:   "test",
			ID:     "id",
			Status: docker.Completed,
		},
	}

	table := getTable(containers)
	assert.Equal(t, table.NumLines(), 1)
}
