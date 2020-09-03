package docker

import (
	"archive/tar"
	"bytes"
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
	closedReaderError     = "archive/tar: write after close"
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

func TestReadDockerfileFailsWhenDockerfileCanNotBeFound(t *testing.T) {
	SetDockerfilesRoot()
	bytes, err := readDockerfile(nonExistentDockerfile)
	assert.Error(t, err)
	assert.Equal(t, []byte(nil), bytes)
	assert.Equal(t, fmt.Sprintf("open %s/%s: no such file or directory", config.GetOrDefault(util.DockerfileDirEnv), nonExistentDockerfile), err.Error())
}

func TestWriteTarHeaderFailsAndPasses(t *testing.T) {
	tw, _ := createTarWriterAndBuffer()
	errors := []error{
		nil,
		fmt.Errorf(closedReaderError),
	}

	for _, err := range errors {
		if err != nil {
			tw.Close()
		}
		err := writeTarHeader(nonExistentDockerfile, []byte{}, tw)
		if err != nil {
			assert.Error(t, err)
			assert.Equal(t, err.Error(), closedReaderError)
		} else {
			assert.NoError(t, err)
			assert.Nil(t, err)
		}
	}
}
func TestWriteTarBytesFails(t *testing.T) {
	tw, _ := createTarWriterAndBuffer()
	errors := []error{
		nil,
		fmt.Errorf(closedReaderError),
	}

	for _, err := range errors {
		if err != nil {
			tw.Close()
		}
		err := writeTarBytes(nonExistentDockerfile, []byte{}, tw)
		if err != nil {
			assert.Error(t, err)
			assert.Equal(t, err.Error(), closedReaderError)
		} else {
			assert.NoError(t, err)
			assert.Nil(t, err)
		}
	}
}

func TestBuildImageFails(t *testing.T) {
	SetDockerfilesRoot()
	handler, err := NewHandler()
	assert.NoError(t, err)

	err = handler.BuildImage(nonExistentDockerfile, "test")
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("open %s/%s: no such file or directory", config.GetOrDefault(util.DockerfileDirEnv), nonExistentDockerfile), err.Error())
}

func TestBuildDockerfilePasses(t *testing.T) {
	SetDockerfilesRoot()
	handler, err := NewHandler()
	assert.NoError(t, err)
	err = handler.BuildImage(actualDockerfile, "test")
	assert.NoError(t, err)
}

func TestCreateDockerfileFails(t *testing.T) {
	tw, buf := createTarWriterAndBuffer()
	tw.Close()
	reader, err := createDockerfileBuild(nonExistentDockerfile, []byte{}, tw, buf)
	assert.Error(t, err)
	assert.Equal(t, closedReaderError, err.Error())
	assert.Nil(t, reader)
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

func createTarWriterAndBuffer() (*tar.Writer, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	return tar.NewWriter(buf), buf
}

// TODO: figure out a way to have this imported as a function in all test packages to prevent copying and pasting this method and SetTestfilesRoot
func SetDockerfilesRoot() {
	// If not in container, set as the path to the 'project's root/Dockerfiles'
	if os.Getenv(util.DockerfileDirEnv) == "" {
		_, b, _, _ := runtime.Caller(0)
		d := path.Join(path.Dir(b))
		os.Setenv(util.DockerfileDirEnv, fmt.Sprintf("%s/../Dockerfiles", filepath.Dir(d)))
	}
}
