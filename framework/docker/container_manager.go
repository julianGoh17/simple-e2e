package docker

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ContainerManager is the object which is responsible for interacting with a specific container and handling any interactions with it
type ContainerManager struct {
	image         string
	containerInfo *ContainerInfo
}

// StartContainer will start the docker container managed by the container manager
func (manager ContainerManager) StartContainer(ctx context.Context, docker *WrapperClient) error {
	baseTrace := getTraceContainerManager(*manager.containerInfo)
	baseTrace.Msg("Container manager attempting to start container")

	if err := docker.StartContainer(ctx, manager.containerInfo.ID); err != nil {
		baseTrace.
			Err(err).
			Msg("Container manager failed to start container")

		return err
	}

	manager.containerInfo.Status = Running
	baseTrace.Msg("Container manager successfully started container")
	return nil
}

// StopContainer will attempt to stop the Docker container managed by the container manager
func (manager ContainerManager) StopContainer(ctx context.Context, docker *WrapperClient, time *time.Duration) error {
	baseTrace := getTraceContainerManager(*manager.containerInfo)
	baseTrace.
		Str("duration", time.String()).
		Msg("Container manager attempting to stop container")

	if err := docker.StopContainer(ctx, manager.containerInfo.ID, time); err != nil {
		baseTrace.
			Err(err).
			Str("duration", time.String()).
			Msg("Container manager failed to stop container")

		return err
	}

	manager.containerInfo.Status = Exited
	baseTrace.
		Str("duration", time.String()).
		Msg("Container manager successfully stopped container")
	return nil
}

// RestartContainer will attempt to stop and restart the container managed by the container manager
func (manager ContainerManager) RestartContainer(ctx context.Context, docker *WrapperClient) error {
	baseTrace := getTraceContainerManager(*manager.containerInfo)
	baseTrace.Msg("Container manager attempting to restart container")

	if err := docker.PauseContainer(ctx, manager.containerInfo.ID); err != nil {
		baseTrace.
			Err(err).
			Msg("Container manager failed to restart container")

		return err
	}

	manager.containerInfo.Status = Running
	baseTrace.Msg("Container manager attempting to successfully restarted container")
	return nil
}

// PauseContainer will pause the docker container managed by the container manager
func (manager ContainerManager) PauseContainer(ctx context.Context, docker *WrapperClient, time *time.Duration) error {
	baseTrace := getTraceContainerManager(*manager.containerInfo)
	baseTrace.
		Str("time", time.String()).
		Msg("Container manager attempting to pause container")

	if err := docker.RestartContainer(ctx, manager.containerInfo.ID, time); err != nil {
		baseTrace.
			Err(err).
			Str("time", time.String()).
			Msg("Container manager failed to start container")

		return err
	}

	manager.containerInfo.Status = Paused

	baseTrace.
		Str("time", time.String()).
		Msg("Container manager successfully started container")
	return nil
}

func getTraceContainerManager(containerInfo ContainerInfo) *zerolog.Event {
	return log.Trace().
		Str("containerID", containerInfo.ID).
		Str("containerName", containerInfo.ID).
		Str("containerStatus", MapContainerStatusToString(containerInfo.Status))
}
