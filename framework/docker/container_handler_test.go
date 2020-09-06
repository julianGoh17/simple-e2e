package docker

import (
	"fmt"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/julianGoh17/simple-e2e/framework/internal"
	"github.com/stretchr/testify/assert"
)

const (
	actualDockerfile      = "Dockerfile.simple"
	nonExistentDockerfile = "non-existent-Dockerfile"
	closedReaderError     = "archive/tar: write after close"
	existingImage         = "docker.io/library/alpine"
)

func TestNewHandlerHasNoNils(t *testing.T) {
	handler, err := NewHandler()
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.wrapper)
	// Depending on Host Daemon's containers could have multiple containers running
	assert.GreaterOrEqual(t, len(handler.containerManagers), 0)
}

func TestNewHandlerFailsToInitialize(t *testing.T) {
	os.Setenv("DOCKER_HOST", "random-host")
	handler, err := NewHandler()
	assert.Nil(t, handler)
	assert.Error(t, err)
	assert.Equal(t, "unable to parse docker host `random-host`", err.Error())
	os.Unsetenv("DOCKER_HOST")
}

func TestHandlerPullImage(t *testing.T) {
	testCases := []struct {
		image string
		err   error
	}{
		{
			existingImage,
			nil,
		},
		{
			"non-existentImage",
			fmt.Errorf("invalid reference format: repository name must be lowercase"),
		},
		{
			"non-existent-image",
			fmt.Errorf("Error response from daemon: pull access denied for non-existent-image, repository does not exist or may require 'docker login': denied: requested access to the resource is denied"),
		},
	}

	handler, err := NewHandler()
	assert.NoError(t, err)

	for _, testCase := range testCases {
		err = handler.PullImage(testCase.image)
		if testCase.err == nil {
			assert.NoError(t, err)
		} else {
			assert.Equal(t, testCase.err.Error(), err.Error())
		}
	}
}

func TestHandlerCreateContainerFails(t *testing.T) {
	handler, err := NewHandler()
	assert.NoError(t, err)
	alreadyAddedContainer := "already-added"

	handler.containerManagers[alreadyAddedContainer] = &ContainerManager{}

	testCases := []struct {
		containerName string
		err           error
	}{
		{
			alreadyAddedContainer,
			fmt.Errorf("container with name '%s' already exists", alreadyAddedContainer),
		},
		{
			"random-client",
			fmt.Errorf("Error response from daemon: No such image: random-image:latest"),
		},
	}

	for _, testCase := range testCases {
		err := handler.CreateContainer("random-image", testCase.containerName)
		assert.Error(t, err)
		assert.Equal(t, testCase.err.Error(), err.Error())
	}
}

func TestHandlerCreateAndDeleteContainerPasses(t *testing.T) {
	handler, err := NewHandler()
	assert.NoError(t, err)
	containerName := "test"

	err = handler.CreateContainer(existingImage, containerName)
	containersBeforeDeletion := len(handler.containerManagers)
	assert.NoError(t, err)
	assert.Greater(t, containersBeforeDeletion, 0)
	assert.NotNil(t, handler.containerManagers[containerName])

	// Need to delete container for this to work, as there will be a created container that does nothing
	err = handler.DeleteContainer(containerName)
	assert.NoError(t, err)
	assert.Less(t, len(handler.containerManagers), containersBeforeDeletion)
	assert.Nil(t, handler.containerManagers[containerName])
}

func TestHandlerDeleteContainerFromHandlerFails(t *testing.T) {
	handler, err := NewHandler()
	assert.NoError(t, err)
	existingContainerName := "existing"
	nonExistentContainerID := "non-existent-id"
	nonExistentContainerName := "non-existent"
	handler.containerManagers[existingContainerName] = &ContainerManager{containerID: nonExistentContainerID}

	testCases := []struct {
		containerName string
		err           error
	}{
		{
			nonExistentContainerName,
			fmt.Errorf("Could not find container '%s' in Framework registry", nonExistentContainerName),
		},
		{
			existingContainerName,
			fmt.Errorf("Error: No such container: %s", nonExistentContainerID),
		},
	}

	for _, testCase := range testCases {
		err := handler.DeleteContainer(testCase.containerName)
		assert.Error(t, err)
		assert.Equal(t, testCase.err.Error(), err.Error())
	}
}

func TestMapContainerNamesAndIDsFails(t *testing.T) {
	os.Setenv(dockerHostEnv, invalidDockerHost)
	defer os.Unsetenv(dockerHostEnv)
	handler, err := NewHandler()
	assert.NoError(t, err)

	containers, err := handler.MapContainersNamesAndIDs()
	assert.Error(t, err)
	assert.Equal(t, canNotConnectToHostError, err.Error())
	assert.Nil(t, containers)
}

func TestMapContainerNamesAndIDsPasses(t *testing.T) {
	handler, err := NewHandler()
	assert.NoError(t, err)

	containerName := "test"

	err = handler.CreateContainer(existingImage, containerName)
	assert.NoError(t, err)
	assert.Greater(t, len(handler.containerManagers), 0)
	assert.NotNil(t, handler.containerManagers[containerName])

	containers, err := handler.MapContainersNamesAndIDs()
	assert.NoError(t, err)
	assert.Greater(t, len(containers), 0)
	assert.NotNil(t, containers[containerName])

	err = handler.DeleteContainer(containerName)
	assert.NoError(t, err)
	assert.Less(t, len(handler.containerManagers), len(containers))
	assert.Nil(t, handler.containerManagers[containerName])
}

func TestGetContainerNamesAndIDs(t *testing.T) {
	testCases := []struct {
		containers          []types.Container
		expectedNamesAndIDs map[string]string
	}{
		{
			[]types.Container{
				{
					Names: []string{"first", "second"},
					ID:    "firstID",
				},
				{
					Names: []string{"third", "fourth"},
					ID:    "secondID",
				},
			},
			map[string]string{
				"first/second": "firstID",
				"third/fourth": "secondID",
			},
		},
		{
			[]types.Container{},
			make(map[string]string),
		},
	}

	for _, testCase := range testCases {
		namesAndIDs := getContainerNamesAndIDs(testCase.containers)
		assert.Equal(t, testCase.expectedNamesAndIDs, namesAndIDs)
	}
}

func TestMain(m *testing.M) {
	internal.TestCoverageReaches85Percent(m)
}
