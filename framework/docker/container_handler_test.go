package docker

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/julianGoh17/simple-e2e/framework/util"
	"github.com/stretchr/testify/assert"
)

const (
	actualDockerfile      = "Dockerfile.simple"
	nonExistentDockerfile = "non-existent-Dockerfile"
)

func TestPullImage(t *testing.T) {
	testCases := []struct {
		image string
		err   error
	}{
		{
			"docker.io/library/alpine",
			nil,
		},
		{
			"non-existentImage",
			fmt.Errorf("invalid reference format: repository name must be lowercase"),
		},
		{
			"non-existent-image",
			fmt.Errorf("Error response from daemon: pull access denied for non-existent-image, repository does not exist or may require 'docker login': denied: requested access to the resource is denied"),
		},
	}

	handler, err := NewHandler()
	assert.NoError(t, err)

	for _, testCase := range testCases {
		err = handler.PullImage(testCase.image)
		if testCase.err == nil {
			assert.NoError(t, err)
		} else {
			assert.Equal(t, testCase.err.Error(), err.Error())
		}
	}
}

func TestBuildDockerfilePasses(t *testing.T) {
	SetDockerfilesRoot()
	handler, err := NewHandler()
	assert.NoError(t, err)
	err = handler.BuildImage(actualDockerfile, "test")
	assert.NoError(t, err)
}

func TestReadDockerfileFailsWhenDockerfileCanNotBeFound(t *testing.T) {
	SetDockerfilesRoot()
	bytes, err := readDockerfile(nonExistentDockerfile)
	assert.Error(t, err)
	assert.Equal(t, []byte(nil), bytes)
	assert.Equal(t, fmt.Sprintf("open %s/%s: no such file or directory", config.GetOrDefault(util.DockerfileDirEnv), nonExistentDockerfile), err.Error())
}

func TestCreateDockerfileBuildFails(t *testing.T) {
	SetDockerfilesRoot()
	// actualDockerfile := "Dockerfile.simple"
	testCases := []struct {
		dockerfile string
		err        error
	}{
		{
			nonExistentDockerfile,
			fmt.Errorf("open %s/%s: no such file or directory", config.GetOrDefault(util.DockerfileDirEnv), nonExistentDockerfile),
		},
		// TODO: figure out how to cause writing tar header to fail
	}

	for _, testCase := range testCases {
		_, err := createDockerfileBuild(testCase.dockerfile)
		assert.Error(t, err)
		assert.Equal(t, testCase.err.Error(), err.Error())
	}
}

func TestCreateDockerfileBuildPasses(t *testing.T) {
	SetDockerfilesRoot()
	reader, err := createDockerfileBuild(actualDockerfile)
	assert.NoError(t, err)
	assert.NotNil(t, reader)
}

func TestReadDockerfileFailsWhenDockerfilePasses(t *testing.T) {
	SetDockerfilesRoot()
	bytes, err := readDockerfile(actualDockerfile)
	assert.NoError(t, err)
	assert.NotNil(t, bytes)
}

func TestMain(m *testing.M) {
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

// func removeWritePermissionForFile(file string) {
// 	filePath := fmt.Sprintf("%s/%s", config.GetOrDefault(util.DockerfileDirEnv), file)
// 	os.Chmod(filePath, 0444)
// }

// func giveWritePermissionForFile(file string) {
// 	filePath := fmt.Sprintf("%s/%s", config.GetOrDefault(util.DockerfileDirEnv), file)
// 	os.Chmod(filePath, 0644)
// }

// TODO: figure out a way to have this imported as a function in all test packages to prevent copying and pasting this method and SetTestfilesRoot
func SetDockerfilesRoot() {
	// If not in container, set as the path to the 'project's root/Dockerfiles'
	if os.Getenv(util.DockerfileDirEnv) == "" {
		_, b, _, _ := runtime.Caller(0)
		d := path.Join(path.Dir(b))
		os.Setenv(util.DockerfileDirEnv, fmt.Sprintf("%s/../Dockerfiles", filepath.Dir(d)))
	}
}
