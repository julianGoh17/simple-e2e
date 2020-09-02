package docker

import (
	"context"
	"fmt"
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
	logger.Trace().Msg("Initializing Framework's Docker client")
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return traceExitOfError(err, "Failed to initialize Framework's Docker client")
	}

	wrapper.Cli = dockerClient

	return traceExitOfError(nil, "Successfully initilaized Framework's Docker client")
}

// Close will close the framework's docker client
func (wrapper *WrapperClient) Close() error {
	logger.Trace().Msg("Closing Framework's Docker client")
	if err := wrapper.Cli.Close(); err != nil {
		return traceExitOfError(wrapper.Cli.Close(), "Failed to close Framework's Docker client")
	}
	return traceExitOfError(nil, "Successfully closed Framework's Docker client")
}

// PullImage will pull the image from DockerHub onto the host machine's daemon.
// TODO: Add ability to pass in ImagePullOpions configured through container handler
func (wrapper *WrapperClient) PullImage(ctx context.Context, image string) error {
	logger.Trace().Str("Image", image).Msg("Wrapper client beginning to pull docker image")
	reader, err := wrapper.Cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return traceExitOfError(err, fmt.Sprintf("Wrapper client failed to pull docker image '%s'", image))
	}
	if err := readOutputAndCloseReader(reader); err != nil {
		return traceExitOfError(err, "Failed to close reader of standard output")
	}
	return traceExitOfError(nil, "Wrapper client successfully pulled docker image")
}

// BuildImage will build the specified image in the specified location
func (wrapper *WrapperClient) BuildImage(ctx context.Context, buildContext io.Reader, buildOptions types.ImageBuildOptions) error {
	logger.Trace().Str("Image", buildOptions.Tags[0]).Msg("Wrapper client beginning to build docker image")
	res, err := wrapper.Cli.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		return traceExitOfBuildingImageForError(err, buildOptions, "Wrapper client failed to build docker image")
	}

	if err := readOutputAndCloseReader(res.Body); err != nil {
		return traceExitOfBuildingImageForError(err, buildOptions, "Failed to close reader of standard output")
	}
	logger.Trace().
		Str("Image", buildOptions.Tags[0]).
		Bool("hasSuccessfullyBuiltDockerfile", true).
		Msg("Successfully built docker image")
	return nil
}

func readOutputAndCloseReader(reader io.ReadCloser) error {
	io.Copy(os.Stdout, reader)
	return reader.Close()
}

func traceExitOfBuildingImageForError(err error, buildOptions types.ImageBuildOptions, msg string) error {
	logger.Trace().
		Err(err).
		Str("Image", buildOptions.Tags[0]).
		Bool("hasSuccessfullyBuiltDockerfile", false).
		Msg(msg)
	return err
}
