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

	"github.com/julianGoh17/simple-e2e/framework/util"
)

const (
	// ExistingImage points to an existing docker image that can be used
	ExistingImage = "docker.io/library/alpine"
	// InvalidDockerHost will cause the Docker client to error on creation
	InvalidDockerHost = "random-host"
	// UnconnectableDockerHost  will cause the client to talk to an invalid Docker daemon
	UnconnectableDockerHost = "http://localhost:9091"
	// DockerHostEnv is the env var key name
	DockerHostEnv = "DOCKER_HOST"
)

var (
	// ErrInvalidHost is the error that occurs when the Docker client starts up with the InvalidDockerHost const
	ErrInvalidHost = fmt.Errorf("unable to parse docker host `random-host`")
	// ErrCanNotConnectToHost is the error that occurs when the Docker client tries to connect the daemon on UncconnectableDockerHost
	ErrCanNotConnectToHost = fmt.Errorf("Cannot connect to the Docker daemon at %s. Is the docker daemon running?", UnconnectableDockerHost)
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
