package docker

import (
	"archive/tar"
	"bytes"
	"fmt"
	"testing"

	"github.com/julianGoh17/simple-e2e/framework/internal"
	"github.com/julianGoh17/simple-e2e/framework/util"
	"github.com/stretchr/testify/assert"
)

func TestReadDockerfileFailsWhenDockerfileCanNotBeFound(t *testing.T) {
	internal.SetDockerfilesRoot()
	bytes, err := readDockerfile(nonExistentDockerfile)
	assert.Error(t, err)
	assert.Equal(t, []byte(nil), bytes)
	assert.Equal(t, fmt.Sprintf("open %s/%s: no such file or directory", config.GetOrDefault(util.DockerfileDirEnv), nonExistentDockerfile), err.Error())
}

func TestReadDockerfileFailsWhenDockerfilePasses(t *testing.T) {
	internal.SetDockerfilesRoot()
	bytes, err := readDockerfile(actualDockerfile)
	assert.NoError(t, err)
	assert.NotNil(t, bytes)
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

func TestCreateDockerfileFails(t *testing.T) {
	tw, buf := createTarWriterAndBuffer()
	tw.Close()
	reader, err := createDockerfileBuild(nonExistentDockerfile, []byte{}, tw, buf)
	assert.Error(t, err)
	assert.Equal(t, closedReaderError, err.Error())
	assert.Nil(t, reader)
}

func TestBuildImagePasses(t *testing.T) {
	internal.SetDockerfilesRoot()
	handler, err := NewHandler()
	assert.NoError(t, err)
	err = handler.BuildImage(actualDockerfile, "test")
	assert.NoError(t, err)
}

func TestBuildImageFails(t *testing.T) {
	internal.SetDockerfilesRoot()
	handler, err := NewHandler()
	assert.NoError(t, err)

	err = handler.BuildImage(nonExistentDockerfile, "test")
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("open %s/%s: no such file or directory", config.GetOrDefault(util.DockerfileDirEnv), nonExistentDockerfile), err.Error())
}

func createTarWriterAndBuffer() (*tar.Writer, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	return tar.NewWriter(buf), buf
}
