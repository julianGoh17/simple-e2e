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
	// Exited means that the container was running but a process inside closed the container
	Exited
	// Paused means that the container is paused and currently not running anything
	Paused
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
	case Exited:
		return "Exited"
	case Paused:
		return "Paused"
	default:
		return ""
	}
}

// MapStateToStatus will convert the string container state to ContainerStatus enum
func MapStateToStatus(status string) ContainerStatus {
	switch status {
	case "running":
		return Running
	case "completed":
		return Completed
	case "errored":
		return Errored
	case "exited":
		return Exited
	case "paused":
		return Paused
	default:
		return Created
	}
}

// ContainerInfo is an object that holds all important information about a container
type ContainerInfo struct {
	Name   string
	ID     string
	Image  string
	Status ContainerStatus
}
