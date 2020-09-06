package docker

// ContainerStatus is an enum which represents the state of the container
type ContainerStatus int

const (
	// Created means that the container has just been created and has not started nor has it errored
	Created ContainerStatus = iota
	// Running means that the container is currently running without any errors
	Running
	// Completed means that the container has finished whatever process is running inside of it and is currently stopped
	Completed
	// Errored means that the container has errored in someway
	Errored
)

// ContainerManager is the object which is responsible for interacting with a specific container and handling any interactions with it
type ContainerManager struct {
	image           string
	containerName   string
	containerID     string
	containerStatus ContainerStatus
}

// NewContainerManager will return a Container Manager intialized for a specified name and image
func NewContainerManager(image, containerName, containerID string) *ContainerManager {
	return &ContainerManager{
		image:           image,
		containerName:   containerName,
		containerStatus: Created,
		containerID:     containerID,
	}
}
