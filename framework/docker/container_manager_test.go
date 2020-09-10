package docker

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/julianGoh17/simple-e2e/framework/internal"
	"github.com/stretchr/testify/assert"
)

func TestStartContainerFails(t *testing.T) {
	client := createClient(t)
	containerManager := createNonExistentTestManager()
	ctx := context.Background()
	err := containerManager.StartContainer(ctx, &client)
	assert.Error(t, err)
	assert.Equal(t, internal.ErrCanNotFindNonExistentContainer.Error(), err.Error())
}

func TestStartContainerPasses(t *testing.T) {
	client := createClient(t)
	containerManager := createShortLifeSpanTestContainer(t, &client)
	defer deleteTestContainer(t, &client, &containerManager)

	ctx := context.Background()
	assert.NoError(t, containerManager.StartContainer(ctx, &client))
	assert.Equal(t, Running, containerManager.containerInfo.Status)
}

func TestRestartContainerFails(t *testing.T) {
	client := createClient(t)
	containerManager := createNonExistentTestManager()
	ctx := context.Background()
	err := containerManager.RestartContainer(ctx, &client)
	assert.Error(t, err)
	assert.Equal(t, internal.ErrCanNotFindNonExistentContainer.Error(), err.Error())
}

func TestRestartContainerPasses(t *testing.T) {
	client := createClient(t)
	containerManager := createForeverRunningTestContainer(t, &client)
	defer deleteTestContainer(t, &client, &containerManager)

	ctx := context.Background()
	assert.NoError(t, containerManager.StartContainer(ctx, &client))
	assert.NoError(t, containerManager.PauseContainer(ctx, &client, &internal.TestDuration))
	assert.NoError(t, containerManager.RestartContainer(ctx, &client))
	assert.Equal(t, Running, containerManager.containerInfo.Status)
}

func TestStopContainerFails(t *testing.T) {
	client := createClient(t)
	containerManager := createNonExistentTestManager()

	ctx := context.Background()
	err := containerManager.StopContainer(ctx, &client, &internal.TestDuration)
	assert.Error(t, err)
	assert.Equal(t, internal.ErrCanNotFindNonExistentContainer.Error(), err.Error())
}

func TestStopContainerPasses(t *testing.T) {
	client := createClient(t)
	containerManager := createShortLifeSpanTestContainer(t, &client)
	defer deleteTestContainer(t, &client, &containerManager)

	ctx := context.Background()
	assert.NoError(t, containerManager.StartContainer(ctx, &client))
	assert.NoError(t, containerManager.StopContainer(ctx, &client, &internal.TestDuration))
	assert.Equal(t, Exited, containerManager.containerInfo.Status)
}

func TestPauseContainerFails(t *testing.T) {
	client := createClient(t)
	containerManager := createNonExistentTestManager()
	ctx := context.Background()
	err := containerManager.PauseContainer(ctx, &client, &internal.TestDuration)
	assert.Error(t, err)
	assert.Equal(t, internal.ErrCanNotFindNonExistentContainer.Error(), err.Error())
}

func TestPauseContainerPasses(t *testing.T) {
	client := createClient(t)
	containerManager := createShortLifeSpanTestContainer(t, &client)
	defer deleteTestContainer(t, &client, &containerManager)
	ctx := context.Background()
	assert.NoError(t, containerManager.PauseContainer(ctx, &client, &internal.TestDuration))
	assert.Equal(t, Paused, containerManager.containerInfo.Status)
}

func createNonExistentTestManager() *ContainerManager {
	return &ContainerManager{
		containerInfo: &ContainerInfo{
			Name: internal.NonExistentContainerName,
			ID:   internal.NonExistentContainerID,
		},
	}
}

func createShortLifeSpanTestContainer(t *testing.T, docker *WrapperClient) ContainerManager {
	return createTestContainer(t, docker, []string{})
}

func createForeverRunningTestContainer(t *testing.T, docker *WrapperClient) ContainerManager {
	return createTestContainer(t, docker, []string{"tail", "-f", "/dev/null"})
}

func createTestContainer(t *testing.T, docker *WrapperClient, cmd []string) ContainerManager {
	ctx := context.Background()
	docker.PullImage(ctx, internal.ExistingImage)

	res, err := docker.CreateContainer(ctx, &container.Config{
		Image: internal.ExistingImage,
		Tty:   false,
		Cmd:   cmd,
	}, internal.ExistingContainerName)

	assert.NoError(t, err)

	return ContainerManager{
		image: internal.ExistingImage,
		containerInfo: &ContainerInfo{
			ID:   res.ID,
			Name: internal.ExistingContainerName,
		},
	}
}

func deleteTestContainer(t *testing.T, docker *WrapperClient, manager *ContainerManager) {
	ctx := context.Background()
	if manager.containerInfo.Status != Exited {
		assert.NoError(t, docker.StopContainer(ctx, manager.containerInfo.ID, &internal.TestDuration))
	}
	assert.NoError(t, docker.DeleteContainer(ctx, manager.containerInfo.ID))
}
