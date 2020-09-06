package docker

import (
	"context"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/assert"
)

func TestWrapperClientFailsToInitialize(t *testing.T) {
	client := WrapperClient{}
	os.Setenv("DOCKER_HOST", "random-host")
	err := client.Initialize()
	assert.Error(t, err)
	assert.Equal(t, "unable to parse docker host `random-host`", err.Error())
	os.Unsetenv("DOCKER_HOST")
}

func TestWrapperClientCanClose(t *testing.T) {
	client := createClient(t)
	assert.NoError(t, client.Close())
}

func TestWrapperClientBuildImageFails(t *testing.T) {
	client := createClient(t)

	ctx := context.Background()
	err := client.BuildImage(ctx, nil, types.ImageBuildOptions{
		Tags: []string{"failed image"},
	})
	assert.Error(t, err)
	assert.Equal(t, "Error response from daemon: client version 1.41 is too new. Maximum supported API version is 1.40", err.Error())
}

func TestWrapperClientPullImageFails(t *testing.T) {
	client := createClosedClient(t)

	ctx := context.Background()
	err := client.PullImage(ctx, "random-image")
	assert.Error(t, err)
	assert.Equal(t, "Error response from daemon: client version 1.41 is too new. Maximum supported API version is 1.40", err.Error())
}

func TestWrapperClientCreateContainerFails(t *testing.T) {
	client := createClosedClient(t)

	ctx := context.Background()
	res, err := client.CreateContainer(ctx, &container.Config{}, "random-container")
	assert.Error(t, err)
	assert.Equal(t, container.ContainerCreateCreatedBody{}, res)
}

func TestWrapperClientDeleteContainerFails(t *testing.T) {
	client := createClient(t)

	ctx := context.Background()
	err := client.DeleteContainer(ctx, "random-id")
	assert.Error(t, err)
	assert.Equal(t, "Error response from daemon: client version 1.41 is too new. Maximum supported API version is 1.40", err.Error())
}

func createClosedClient(t *testing.T) WrapperClient {
	client := createClient(t)
	assert.NoError(t, client.Close())
	return client
}

func createClient(t *testing.T) WrapperClient {
	client := WrapperClient{}
	err := client.Initialize()
	assert.NoError(t, err)
	return client
}
