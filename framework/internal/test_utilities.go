package internal

/*
 * THIS PACKAGE IS MEANT FOR TEST UTILITY FUNCTIONS ONLY
 */

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/julianGoh17/simple-e2e/framework/util"
)

const (
	// ExistingImage points to an existing docker image that can be used
	ExistingImage = "docker.io/library/alpine"
	// ExistingContainerName points to an existing docker container id that can be used
	ExistingContainerName = "existing-container"
	// InvalidDockerHost will cause the Docker client to error on creation
	InvalidDockerHost = "random-host"
	// UnconnectableDockerHost will cause the client to talk to an invalid Docker daemon
	UnconnectableDockerHost = "http://localhost:9091"
	// DockerHostEnv is the env var key name
	DockerHostEnv = "DOCKER_HOST"
	// NonExistentContainerID is a containerID that doesn't exist and should cause errors when trying to interact with the container
	NonExistentContainerID = "non-existent-container"
	// NonExistentContainerName is the name of a container that doesn't exist and should cause errors when trying to interact with the container
	NonExistentContainerName = "non-existent"
	// UnconnectableContainerName is the name of a container that can't be connected to because the framework is talking to the wrong docker host
	UnconnectableContainerName = "unconnectable-container"
	// ActualDockerfile is the name of the actual Dockerfile that is present in the repository
	ActualDockerfile = "Dockerfile.simple"
	// NonExistentDockerfile is the name of the actual Dockerfile that is not present in the repository
	NonExistentDockerfile = "non-existent-Dockerfile"
)

var (
	// ErrInvalidHost is the error that occurs when the Docker client starts up with the InvalidDockerHost const
	ErrInvalidHost = fmt.Errorf("unable to parse docker host `random-host`")
	// ErrCanNotConnectToHost is the error that occurs when the Docker client tries to connect the daemon on UncconnectableDockerHost
	ErrCanNotConnectToHost = fmt.Errorf("Cannot connect to the Docker daemon at %s. Is the docker daemon running?", UnconnectableDockerHost)
	// ErrCanNotFindNonExistentContainer is the error that occurs when the Docker client tries to interact with a non-existent container ID
	ErrCanNotFindNonExistentContainer = fmt.Errorf("Error response from daemon: No such container: %s", NonExistentContainerID)
	// ErrCanNotFindNonExistentContainerInRegistry is the error that occurs when the framework tries to find a container not registered in its framework
	ErrCanNotFindNonExistentContainerInRegistry = fmt.Errorf("Could not find container '%s' in the framework's registry", NonExistentContainerName)
	// ErrClosedTarReader is the error that occurs when trying to read with a tar reader that is already closed
	ErrClosedTarReader = fmt.Errorf("archive/tar: write after close")
	// TestDuration is the standard duration a test should wait for something with a timeout to complete
	TestDuration time.Duration = 5 * time.Second
	// ForeverRunningCmd is the command used to make a Docker container run forever on start up
	ForeverRunningCmd = []string{"tail", "-f", "/dev/null"}
)

// SetDockerfilesRoot will set 'DOCKERFILES_DIR' env to the path to the 'project's root/Dockerfiles' if it's not already set
func SetDockerfilesRoot() {
	// If not in container, set as the path to the 'project's root/Dockerfiles'
	if os.Getenv(util.DockerfileDirEnv) == "" {
		_, b, _, _ := runtime.Caller(0)
		d := path.Join(path.Dir(b))
		os.Setenv(util.DockerfileDirEnv, fmt.Sprintf("%s/../Dockerfiles", filepath.Dir(d)))
	}
}

// SetTestFilesRoot will set 'TEST_DIR' env to the path to the 'project's root/tests' if it's not already set
func SetTestFilesRoot() {
	// If not in container, set as the path to the 'project's root/tests'
	if os.Getenv(util.TestDirEnv) == "" {
		_, b, _, _ := runtime.Caller(0)
		d := path.Join(path.Dir(b))
		os.Setenv(util.TestDirEnv, fmt.Sprintf("%s/../tests", filepath.Dir(d)))
	}
}

// TestCoverageReaches85Percent will ensure that test coverage passes 85%
func TestCoverageReaches85Percent(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	rc := m.Run()

	// rc 0 means we've passed,
	// and CoverMode will be non empty if run with -cover
	if rc == 0 && testing.CoverMode() != "" {
		c := testing.Coverage()
		if c < 0.85 {
			fmt.Println("Tests passed but coverage failed at", c)
			rc = -1
		}
	}
	os.Exit(rc)
}
