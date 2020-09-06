package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog"
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
		return traceExitOfError(err, "Failed to close Framework's Docker client")
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
		Msg("Successfully built Docker image")
	return nil
}

// CreateContainer will create a container with a specified configuration (but this does not start any processes in the container)
func (wrapper *WrapperClient) CreateContainer(ctx context.Context, config *container.Config, containerName string) (container.ContainerCreateCreatedBody, error) {
	beginningLog := traceCreateContainer(config)
	beginningLog.Msg("Creating Docker container")

	// Note for the future, to set up container to container communication may need to pass in host config
	resp, err := wrapper.Cli.ContainerCreate(ctx, config, nil, nil, nil, containerName)
	if err != nil {
		return resp, traceExitCreateContainerError(err, config, "Failed to create Docker container")
	}

	return resp, traceExitCreateContainerError(err, config, "Successfully created Docker container")
}

// DeleteContainer will kill and remove a container (started or running) from the host's docker daemon
func (wrapper *WrapperClient) DeleteContainer(ctx context.Context, containerID string) error {
	logger.Trace().
		Str("containerID", containerID).
		Msg("Beginning to delete container")

	if err := wrapper.Cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{}); err != nil {
		logger.Trace().Err(err).Str("containerID", containerID).Msg("Failed to delete container")
		return err
	}

	logger.Trace().
		Str("containerID", containerID).
		Msg("Successfully deleted container")
	return nil
}

// ListContainers will list all the existing containers on the host daemon
func (wrapper *WrapperClient) ListContainers(ctx context.Context) ([]types.Container, error) {
	logger.Trace().
		Msg("Beginning to list containers on host daemon")

	containers, err := wrapper.Cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		logger.Trace().Err(err).Msg("Failed to list containers on host daemon")
		return nil, err
	}

	logger.Trace().
		Strs("containerIDs", getContainerIDs(containers)).
		Strs("containerNames", getContainerNames(containers)).
		Msg("Successfully listed container")
	return containers, nil
}

func readOutputAndCloseReader(reader io.ReadCloser) error {
	io.Copy(os.Stdout, reader)
	return reader.Close()
}

func getContainerIDs(containers []types.Container) []string {
	ids := []string{}
	for _, container := range containers {
		ids = append(ids, container.ID)
	}
	return ids
}

func getContainerNames(containers []types.Container) []string {
	ids := []string{}
	for _, container := range containers {
		ids = append(ids, strings.Join(container.Names, "/"))
	}

	return ids
}

func traceExitOfBuildingImageForError(err error, buildOptions types.ImageBuildOptions, msg string) error {
	logger.Trace().
		Err(err).
		Str("Image", buildOptions.Tags[0]).
		Bool("hasSuccessfullyBuiltDockerfile", false).
		Msg(msg)
	return err
}

func traceExitCreateContainerError(err error, config *container.Config, msg string) error {
	event := traceCreateContainer(config)
	event.Err(err).Msg(msg)
	return err
}

func traceCreateContainer(config *container.Config) *zerolog.Event {
	return logger.Trace().
		Str("image", config.Image).
		Strs("environmentalVariables", config.Env)
}
