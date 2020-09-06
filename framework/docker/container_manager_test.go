package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContainerManager(t *testing.T) {
	image := "random-image"
	containerName := "random-name"
	containerID := "random-ID"
	manager := NewContainerManager(image, containerName, containerID)

	assert.NotNil(t, manager)
	assert.Equal(t, manager.image, image)
	assert.Equal(t, manager.containerName, containerName)
	assert.Equal(t, manager.containerID, containerID)
	assert.Equal(t, manager.containerStatus, Created)
}
