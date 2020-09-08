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

// MapContainerStatusToString will convert the container status to string
func MapContainerStatusToString(status ContainerStatus) string {
	switch status {
	case Created:
		return "Created"
	case Running:
		return "Running"
	case Completed:
		return "Completed"
	case Errored:
		return "Errored"
	default:
		return ""
	}
}

// ContainerInfo is an object that holds all important information about a container
type ContainerInfo struct {
	Name   string
	ID     string
	Image  string
	Status ContainerStatus
}
