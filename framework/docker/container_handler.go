package docker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/julianGoh17/simple-e2e/framework/util"
	"github.com/rs/zerolog/log"
)

var (
	logger = util.GetStandardLogger()
	config = util.NewConfig()
)

// Handler is the framework's controller responsible for all docker related operations
type Handler struct {
	wrapper           *WrapperClient
	containerManagers map[string]*ContainerManager
}

// NewHandler will create a handler object intialized and ready to use. Will error if there is any problems with setting up for docker operations
func NewHandler() (*Handler, error) {
	logger.Trace().Msg("Creating new Docker handler")
	ctx := context.Background()
	handler := &Handler{containerManagers: make(map[string]*ContainerManager)}

	handler.wrapper = &WrapperClient{}
	if err := handler.wrapper.Initialize(); err != nil {
		return nil, traceExitOfError(err, "Failed to create new Docker handler")
	}
	handler.wrapper.Cli.NegotiateAPIVersion(ctx)

	if err := handler.initializeContainerManagers(); err != nil {
		return handler, traceExitOfError(err, "Failed to initialize container managers")
	}

	logger.Trace().Msg("Successfully created new Docker handler")
	return handler, nil
}

// PullImage will pull the image from dockerhub onto the host machine's daemon
func (handler *Handler) PullImage(image string) error {
	logger.Trace().
		Str("image", image).
		Msg("Docker handler pulling image")

	ctx := context.Background()
	return handler.wrapper.PullImage(ctx, image)
}

// CreateContainer will create a container for a specified image and name. The framework will then create a ContainerManager to manage that container
func (handler *Handler) CreateContainer(image, containerName string, cmd []string) error {
	logger.Trace().
		Str("image", image).
		Str("containerName", containerName).
		Msg("Creating container and manager")

	if _, ok := handler.containerManagers[containerName]; ok {
		return traceExitCreateContainerAndContainerManagerError(fmt.Errorf("container with name '%s' already exists", containerName),
			image,
			containerName,
			"Container with specified name already exists")
	}

	ctx := context.Background()
	resp, err := handler.wrapper.CreateContainer(ctx, &container.Config{
		Image: image,
		Tty:   false,
		Cmd:   cmd,
	}, containerName)
	if err != nil {
		return traceExitCreateContainerAndContainerManagerError(err, image, containerName, "Failed to create container")
	}

	handler.containerManagers[containerName] = &ContainerManager{image: image, containerInfo: &ContainerInfo{
		Name:  containerName,
		ID:    resp.ID,
		Image: image,
	}}

	logger.Trace().
		Str("image", image).
		Str("containerName", containerName).
		Str("containerID", resp.ID).
		Msg("Successfully created container and manager")
	return nil
}

// DeleteContainer will delete a specified container and its corresponding ContainerManager
func (handler *Handler) DeleteContainer(containerName string) error {
	manager, err := handler.findManagerForContainer(containerName)
	if err != nil {
		return err
	}
	ctx := context.Background()
	if err := handler.wrapper.DeleteContainer(ctx, manager.containerInfo.ID); err != nil {
		logger.Trace().
			Err(err).
			Str("containerID", manager.containerInfo.ID).
			Str("containerName", containerName).
			Msg("Failed to delete container")

		return err
	}

	delete(handler.containerManagers, containerName)

	logger.Trace().
		Str("containerName", containerName).
		Msg("Successfully deleted container and corresponding container manager")
	return nil
}

// GetContainerInfo will return a list of ContainerInfo objects gathered from the host machine
func (handler *Handler) GetContainerInfo(showAll bool) ([]*ContainerInfo, error) {
	logger.Trace().
		Bool("showAll", showAll).
		Msg("Attemping to list containers")

	ctx := context.Background()
	containers, err := handler.wrapper.ListContainers(ctx, showAll)
	if err != nil {
		logger.Trace().
			Err(err).
			Bool("showAll", showAll).
			Msg("Failed to list containers")
		return nil, err
	}

	logger.Trace().
		Bool("showAll", showAll).
		Strs("containers", getContainerNames(containers)).
		Msg("Successfully listed containers")

	return convertToContainerInfo(containers), nil
}

func convertToContainerInfo(containers []types.Container) []*ContainerInfo {
	infos := []*ContainerInfo{}
	for _, container := range containers {
		infos = append(infos, &ContainerInfo{Name: strings.Join(container.Names, "/"), ID: container.ID, Image: container.Image, Status: MapStateToStatus(container.State)})
	}

	return infos
}

// TODO: It may be useful to have this step run everytime before a docker operation as this is only run on start up currently
func (handler *Handler) initializeContainerManagers() error {
	logger.Trace().Msg("Attempting to initialize container managers")

	containerInfos, err := handler.GetContainerInfo(true)
	if err != nil {
		logger.Trace().
			Err(err).
			Msg("Failed to initialize container managers")
		return err
	}

	for _, containerInfo := range containerInfos {
		handler.containerManagers[containerInfo.Name] = &ContainerManager{image: containerInfo.Image, containerInfo: containerInfo}
	}

	logger.Trace().Msg("Successfully initialized contianer managers")
	return nil
}

// StartContainer will start running a container that has been created (with whatever it has been configured with on creation)
func (handler *Handler) StartContainer(containerName string) error {
	manager, err := handler.findManagerForContainer(containerName)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return manager.StartContainer(ctx, handler.wrapper)
}

// StopContainer will gracefully stop running a container within a certain period of time before it will forcibly kill a container
func (handler *Handler) StopContainer(containerName string, time *time.Duration) error {
	manager, err := handler.findManagerForContainer(containerName)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return manager.StopContainer(ctx, handler.wrapper, time)
}

// RestartContainer will restart a container that is currently running
func (handler *Handler) RestartContainer(containerName string) error {
	manager, err := handler.findManagerForContainer(containerName)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return manager.RestartContainer(ctx, handler.wrapper)
}

// PauseContainer will pause the main executing process in the container without terminating it (leaving it running)
func (handler *Handler) PauseContainer(containerName string, time *time.Duration) error {
	manager, err := handler.findManagerForContainer(containerName)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return manager.PauseContainer(ctx, handler.wrapper, time)
}

func (handler *Handler) findManagerForContainer(containerName string) (*ContainerManager, error) {
	logger.Trace().
		Str("containerName", containerName).
		Msg("Attempting to find container in registry")

	if _, ok := handler.containerManagers[containerName]; !ok {
		err := fmt.Errorf("Could not find container '%s' in the framework's registry", containerName)
		logger.Trace().
			Str("containerName", containerName).
			Err(err).
			Msg("Failed to find container in registry")
		return &ContainerManager{}, err
	}

	logger.Trace().
		Str("containerName", containerName).
		Msg("Succcesfully found container in registry")
	return handler.containerManagers[containerName], nil
}

func traceExitCreateContainerAndContainerManagerError(err error, image, containerName, msg string) error {
	logger.Trace().
		Str("image", image).
		Str("containerName", containerName).
		Err(err).
		Msg(msg)

	return err
}

func traceExitDockerfileBuildingError(err error, dockerfile, msg string) error {
	logger.Trace().Err(err).Str("Dockerfile", dockerfile).Msg(msg)
	return err
}

func traceExitOfError(err error, msg string) error {
	log.Trace().
		Err(err).
		Msg(msg)
	return err
}
