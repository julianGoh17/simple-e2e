package docker

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/assert"
)

const (
	invalidDockerHost = "http://localhost:9090"
	dockerHostEnv     = "DOCKER_HOST"
)

var (
	canNotConnectToHostError = fmt.Sprintf("Cannot connect to the Docker daemon at %s. Is the docker daemon running?", invalidDockerHost)
)

func TestWrapperClientFailsToInitialize(t *testing.T) {
	client := WrapperClient{}
	os.Setenv(dockerHostEnv, "random-host")
	err := client.Initialize()
	assert.Error(t, err)
	assert.Equal(t, "unable to parse docker host `random-host`", err.Error())
	os.Unsetenv(dockerHostEnv)
}

func TestWrapperClientCanClose(t *testing.T) {
	client := createClient(t)
	assert.NoError(t, client.Close())
}

func TestWrapperClientBuildImageFails(t *testing.T) {
	os.Setenv(dockerHostEnv, invalidDockerHost)
	defer os.Unsetenv(dockerHostEnv)
	client := createClient(t)

	ctx := context.Background()
	err := client.BuildImage(ctx, nil, types.ImageBuildOptions{
		Tags: []string{"failed image"},
	})
	assert.Error(t, err)
	assert.Equal(t, canNotConnectToHostError, err.Error())
}

func TestWrapperClientPullImageFails(t *testing.T) {
	os.Setenv(dockerHostEnv, invalidDockerHost)
	defer os.Unsetenv(dockerHostEnv)
	client := createClient(t)

	ctx := context.Background()
	err := client.PullImage(ctx, "random-image")
	assert.Error(t, err)
	assert.Equal(t, canNotConnectToHostError, err.Error())
}

func TestWrapperClientCreateContainerFails(t *testing.T) {
	os.Setenv(dockerHostEnv, invalidDockerHost)
	defer os.Unsetenv(dockerHostEnv)
	client := createClient(t)

	ctx := context.Background()
	res, err := client.CreateContainer(ctx, &container.Config{}, "random-container")
	assert.Error(t, err)
	assert.Equal(t, canNotConnectToHostError, err.Error())
	assert.Equal(t, container.ContainerCreateCreatedBody{}, res)
}

func TestWrapperClientDeleteContainerFails(t *testing.T) {
	os.Setenv(dockerHostEnv, invalidDockerHost)
	defer os.Unsetenv(dockerHostEnv)
	client := createClient(t)

	ctx := context.Background()
	err := client.DeleteContainer(ctx, "random-id")
	assert.Error(t, err)
	assert.Equal(t, canNotConnectToHostError, err.Error())
}

func TestWrapperClientListContainersFail(t *testing.T) {
	os.Setenv(dockerHostEnv, invalidDockerHost)
	defer os.Unsetenv(dockerHostEnv)
	client := createClient(t)

	ctx := context.Background()
	containers, err := client.ListContainers(ctx)
	assert.Error(t, err)
	assert.Equal(t, canNotConnectToHostError, err.Error())
	assert.Nil(t, containers)
}

func createClient(t *testing.T) WrapperClient {
	client := WrapperClient{}
	err := client.Initialize()
	assert.NoError(t, err)
	return client
}
