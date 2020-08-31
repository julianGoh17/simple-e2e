package docker

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// WrapperClient is the framework's wrapper for the Docker client which will be used to create docker images
type WrapperClient struct {
	Cli *client.Client
}

// Initialize will create a docker client for the framework to use to communicate with the host machines daemon
func (wrapper *WrapperClient) Initialize() error {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	wrapper.Cli = dockerClient

	return nil
}

// Close will close the framework's docker client
func (wrapper *WrapperClient) Close() error {
	return wrapper.Cli.Close()
}

// PullImage will pull the image from DockerHub onto the host machine's daemon.
// TODO: Add ability to pass in ImagePullOpions configured through container handler
func (wrapper *WrapperClient) PullImage(ctx context.Context, image string) error {
	reader, err := wrapper.Cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	return readOutputAndCloseReader(reader)
}

// BuildImage will build the specified image in the specified location
// TODO: Add ability to pass in ImageBuildOptions configured through container handler
func (wrapper *WrapperClient) BuildImage(ctx context.Context, buildContext io.Reader, dockerfile string) error {
	res, err := wrapper.Cli.ImageBuild(ctx, buildContext, types.ImageBuildOptions{
		Context:    buildContext,
		Dockerfile: dockerfile,
	})
	if err != nil {
		return err
	}

	return readOutputAndCloseReader(res.Body)
}

func readOutputAndCloseReader(reader io.ReadCloser) error {
	io.Copy(os.Stdout, reader)
	return reader.Close()
}
