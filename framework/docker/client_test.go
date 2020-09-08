package docker

import (
	"context"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/julianGoh17/simple-e2e/framework/internal"
	"github.com/stretchr/testify/assert"
)

func TestWrapperClientFailsToInitialize(t *testing.T) {
	client := WrapperClient{}
	os.Setenv(internal.DockerHostEnv, internal.InvalidDockerHost)
	err := client.Initialize()
	assert.Error(t, err)
	assert.Equal(t, internal.ErrInvalidHost.Error(), err.Error())
	os.Unsetenv(internal.DockerHostEnv)
}

func TestWrapperClientCanClose(t *testing.T) {
	client := createClient(t)
	assert.NoError(t, client.Close())
}

func TestWrapperClientBuildImageFails(t *testing.T) {
	os.Setenv(internal.DockerHostEnv, internal.UnconnectableDockerHost)
	defer os.Unsetenv(internal.DockerHostEnv)
	client := createClient(t)

	ctx := context.Background()
	err := client.BuildImage(ctx, nil, types.ImageBuildOptions{
		Tags: []string{"failed image"},
	})
	assert.Error(t, err)
	assert.Equal(t, internal.ErrCanNotConnectToHost.Error(), err.Error())
}

func TestWrapperClientPullImageFails(t *testing.T) {
	os.Setenv(internal.DockerHostEnv, internal.UnconnectableDockerHost)
	defer os.Unsetenv(internal.DockerHostEnv)
	client := createClient(t)

	ctx := context.Background()
	err := client.PullImage(ctx, "random-image")
	assert.Error(t, err)
	assert.Equal(t, internal.ErrCanNotConnectToHost.Error(), err.Error())
}

func TestWrapperClientCreateContainerFails(t *testing.T) {
	os.Setenv(internal.DockerHostEnv, internal.UnconnectableDockerHost)
	defer os.Unsetenv(internal.DockerHostEnv)
	client := createClient(t)

	ctx := context.Background()
	res, err := client.CreateContainer(ctx, &container.Config{}, "random-container")
	assert.Error(t, err)
	assert.Equal(t, internal.ErrCanNotConnectToHost.Error(), err.Error())
	assert.Equal(t, container.ContainerCreateCreatedBody{}, res)
}

func TestWrapperClientDeleteContainerFails(t *testing.T) {
	os.Setenv(internal.DockerHostEnv, internal.UnconnectableDockerHost)
	defer os.Unsetenv(internal.DockerHostEnv)
	client := createClient(t)

	ctx := context.Background()
	err := client.DeleteContainer(ctx, "random-id")
	assert.Error(t, err)
	assert.Equal(t, internal.ErrCanNotConnectToHost.Error(), err.Error())
}

func TestWrapperClientListContainersFail(t *testing.T) {
	os.Setenv(internal.DockerHostEnv, internal.UnconnectableDockerHost)
	defer os.Unsetenv(internal.DockerHostEnv)
	client := createClient(t)

	ctx := context.Background()
	containers, err := client.ListContainers(ctx, false)
	assert.Error(t, err)
	assert.Equal(t, internal.ErrCanNotConnectToHost.Error(), err.Error())
	assert.Nil(t, containers)
}

func createClient(t *testing.T) WrapperClient {
	client := WrapperClient{}
	err := client.Initialize()
	assert.NoError(t, err)
	return client
}
