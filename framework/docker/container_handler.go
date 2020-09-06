package docker

import (
	"context"
	"fmt"
	"strings"

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
		return nil, traceExitOfError(err, "Failed to initialize container managers")
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
func (handler *Handler) CreateContainer(image, containerName string) error {
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
	}, containerName)
	if err != nil {
		return traceExitCreateContainerAndContainerManagerError(err, image, containerName, "Failed to create container")
	}

	handler.containerManagers[containerName] = NewContainerManager(image, containerName, resp.ID)

	logger.Trace().
		Str("image", image).
		Str("containerName", containerName).
		Str("containerID", resp.ID).
		Msg("Successfully created container and manager")
	return nil
}

// DeleteContainer will delete a specified container and its corresponding ContainerManager
func (handler *Handler) DeleteContainer(containerName string) error {
	logger.Trace().
		Str("containerName", containerName).
		Msg("Attempting to delete container and corresponding container manager")

	if _, ok := handler.containerManagers[containerName]; !ok {
		return traceExitDeleteContainerAndContainerManagerError(fmt.Errorf("Could not find container '%s' in Framework registry", containerName),
			containerName, "", "Attempted to delete unregistered container")
	}

	manager := handler.containerManagers[containerName]
	ctx := context.Background()
	if err := handler.wrapper.DeleteContainer(ctx, manager.containerID); err != nil {
		return traceExitDeleteContainerAndContainerManagerError(err, containerName, manager.containerID, "Failed to delete container")
	}

	delete(handler.containerManagers, containerName)

	logger.Trace().
		Str("containerName", containerName).
		Msg("Successfully deleted container and corresponding container manager")
	return nil
}

// MapContainersNamesAndIDs will create a map of container name to container ID of all containers currently on the host daemon
func (handler *Handler) MapContainersNamesAndIDs() (map[string]string, error) {
	logger.Trace().
		Msg("Attemping to list containers")

	ctx := context.Background()
	containers, err := handler.wrapper.ListContainers(ctx)
	if err != nil {
		logger.Trace().
			Err(err).
			Msg("Failed to list containers")
		return nil, err
	}

	logger.Trace().Strs("containers", getContainerNames(containers)).Msg("Successfully listed containers")

	return getContainerNamesAndIDs(containers), nil
}

func getContainerNamesAndIDs(containers []types.Container) map[string]string {
	namesAndIDs := make(map[string]string)
	for _, container := range containers {
		namesAndIDs[strings.Join(container.Names, "/")] = container.ID
	}

	return namesAndIDs
}

func (handler *Handler) initializeContainerManagers() error {
	logger.Trace().Msg("Attempting to initialize container managers")

	containerNamesAndIDs, err := handler.MapContainersNamesAndIDs()
	if err != nil {
		logger.Trace().Err(err).Msg("Failed to initialize container managers")
	}

	for name, id := range containerNamesAndIDs {
		handler.containerManagers[name] = NewSimpleContainerManager(name, id)
	}

	logger.Trace().Msg("Successfully initialized contianer managers")
	return nil
}

func traceExitCreateContainerAndContainerManagerError(err error, image, containerName, msg string) error {
	logger.Trace().
		Str("image", image).
		Str("containerName", containerName).
		Err(err).
		Msg(msg)

	return err
}

func traceExitDeleteContainerAndContainerManagerError(err error, containerName, containerID, msg string) error {
	logger.Trace().
		Err(err).
		Str("containerID", containerID).
		Str("containerName", containerName).
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
