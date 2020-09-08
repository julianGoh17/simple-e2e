package docker

// ContainerManager is the object which is responsible for interacting with a specific container and handling any interactions with it
type ContainerManager struct {
	image         string
	containerInfo *ContainerInfo
}
