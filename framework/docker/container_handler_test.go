package docker

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/julianGoh17/simple-e2e/framework/internal"
	"github.com/stretchr/testify/assert"
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
	os.Setenv(internal.DockerHostEnv, internal.InvalidDockerHost)
	defer os.Unsetenv(internal.DockerHostEnv)
	handler, err := NewHandler()
	assert.Nil(t, handler)
	assert.Error(t, err)
	assert.Equal(t, internal.ErrInvalidHost.Error(), err.Error())
}

func TestHandlerPullImage(t *testing.T) {
	testCases := []struct {
		image string
		err   error
	}{
		{
			internal.ExistingImage,
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
		err := handler.CreateContainer("random-image", testCase.containerName, []string{})
		assert.Error(t, err)
		assert.Equal(t, testCase.err.Error(), err.Error())
	}
}

func TestHandlerCreateAndDeleteContainerPasses(t *testing.T) {
	handler, err := NewHandler()
	assert.NoError(t, err)
	containerName := "test"

	err = handler.CreateContainer(internal.ExistingImage, containerName, []string{})
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
	handler.containerManagers[existingContainerName] = &ContainerManager{containerInfo: &ContainerInfo{ID: nonExistentContainerID}}

	testCases := []struct {
		containerName string
		err           error
	}{
		{
			nonExistentContainerName,
			internal.ErrCanNotFindNonExistentContainerInRegistry,
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
	os.Setenv(internal.DockerHostEnv, internal.UnconnectableDockerHost)
	defer os.Unsetenv(internal.DockerHostEnv)
	handler, _ := NewHandler()

	containers, err := handler.GetContainerInfo(true)
	assert.Error(t, err)
	assert.Equal(t, internal.ErrCanNotConnectToHost.Error(), err.Error())
	assert.Nil(t, containers)
}

func TestMapContainerNamesAndIDsPasses(t *testing.T) {
	handler, err := NewHandler()
	assert.NoError(t, err)

	containerName := "test"

	err = handler.CreateContainer(internal.ExistingImage, containerName, []string{})
	assert.NoError(t, err)
	assert.Greater(t, len(handler.containerManagers), 0)
	assert.NotNil(t, handler.containerManagers[containerName])

	containers, err := handler.GetContainerInfo(true)
	assert.NoError(t, err)
	assert.Greater(t, len(containers), 0)

	hasListedCreatedContainer := false
	for _, container := range containers {
		hasListedCreatedContainer = container.Name == "/"+containerName
		if hasListedCreatedContainer {
			break
		}
	}

	assert.Equal(t, true, hasListedCreatedContainer, "Could not find created container in the listed containers")

	err = handler.DeleteContainer(containerName)
	assert.NoError(t, err)
	assert.Less(t, len(handler.containerManagers), len(containers))
	assert.Nil(t, handler.containerManagers[containerName])
}

func TestGetContainerNamesAndIDs(t *testing.T) {
	testCases := []struct {
		containers []types.Container
		expected   []*ContainerInfo
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
			[]*ContainerInfo{
				{
					Name: "first/second",
					ID:   "firstID",
				},
				{
					Name: "third/fourth",
					ID:   "secondID",
				},
			},
		},
		{
			[]types.Container{},
			[]*ContainerInfo{},
		},
	}

	for _, testCase := range testCases {
		containerInfo := convertToContainerInfo(testCase.containers)
		assert.Equal(t, testCase.expected, containerInfo)
	}
}

/*
 * ONLY TESTING ERROR CASES, SUCCESS CASES SHOULD BE TESTED IN CONTAINER MANAGER TO REDUCE DUPLICATION
 */

func TestFrameworkErrorsWhenItCanNotFindContainer(t *testing.T) {
	emptyHandler := createEmptyHandler(t)

	for _, function := range []func(string) error{
		emptyHandler.DeleteContainer,
		emptyHandler.RestartContainer,
		emptyHandler.StartContainer,
	} {
		err := function(internal.NonExistentContainerName)
		assert.Error(t, err)
		assert.Equal(t, internal.ErrCanNotFindNonExistentContainerInRegistry.Error(), err.Error())
	}

	for _, function := range []func(string, *time.Duration) error{
		emptyHandler.PauseContainer,
		emptyHandler.StopContainer,
	} {
		err := function(internal.NonExistentContainerName, &internal.TestDuration)
		assert.Error(t, err)
		assert.Equal(t, internal.ErrCanNotFindNonExistentContainerInRegistry.Error(), err.Error())
	}
}

func TestFrameworkErrorsWhenItCanNotTalkToDaemon(t *testing.T) {
	emptyHandler := createHandlerThatWillTalkToWrongDaemonHandler(t)

	for _, function := range []func(string) error{
		emptyHandler.DeleteContainer,
		emptyHandler.RestartContainer,
		emptyHandler.StartContainer,
	} {
		err := function(internal.UnconnectableContainerName)
		assert.Error(t, err)
		assert.Equal(t, internal.ErrCanNotConnectToHost.Error(), err.Error())
	}

	for _, function := range []func(string, *time.Duration) error{
		emptyHandler.PauseContainer,
		emptyHandler.StopContainer,
	} {
		err := function(internal.UnconnectableContainerName, &internal.TestDuration)
		assert.Error(t, err)
		assert.Equal(t, internal.ErrCanNotConnectToHost.Error(), err.Error())
	}
}

func createEmptyHandler(t *testing.T) Handler {
	handler, err := NewHandler()
	assert.NoError(t, err)
	return *handler
}

func createHandlerThatWillTalkToWrongDaemonHandler(t *testing.T) Handler {
	os.Setenv(internal.DockerHostEnv, internal.UnconnectableDockerHost)
	defer os.Unsetenv(internal.DockerHostEnv)
	handler, err := NewHandler()
	assert.Error(t, err)
	assert.Equal(t, internal.ErrCanNotConnectToHost.Error(), err.Error())
	handler.containerManagers[internal.UnconnectableContainerName] = &ContainerManager{
		containerInfo: &ContainerInfo{},
	}
	return *handler
}

func TestMain(m *testing.M) {
	internal.TestCoverageReaches85Percent(m)
}
