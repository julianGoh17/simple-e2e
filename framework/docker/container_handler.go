package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/julianGoh17/simple-e2e/framework/util"
	"github.com/rs/zerolog/log"
)

var (
	logger = util.GetStandardLogger()
	config = util.NewConfig()
)

// Handler is the framework's controller responsible for all docker related operations
type Handler struct {
	wrapper *WrapperClient
	// TODO: will need container manager to manage running containers
}

// NewHandler will create a handler object intialized and ready to use. Will error if there is any problems with setting up for docker operations
func NewHandler() (*Handler, error) {
	logger.Trace().Msg("Creating new Docker handler")
	ctx := context.Background()
	handler := &Handler{}
	handler.wrapper = &WrapperClient{}
	if err := handler.wrapper.Initialize(); err != nil {
		return nil, traceExitOfError(err, "Failed to create new Docker handler")
	}
	handler.wrapper.Cli.NegotiateAPIVersion(ctx)

	logger.Trace().Msg("Successfully created new Docker handler")
	return handler, nil
}

// PullImage will pull the image from dockerhub onto the host machine's daemon
func (handler *Handler) PullImage(image string) error {
	logger.Trace().Str("image", image).Msg("Docker handler pulling image")

	ctx := context.Background()
	return handler.wrapper.PullImage(ctx, image)
}

// BuildImage will build an image from a specified Dockrefile onto the host machine's daemon
func (handler *Handler) BuildImage(dockerfile, imageName string) error {
	logger.Trace().Str("Dockerfile", dockerfile).Str("Image Name", imageName).Msg("Docker handler building image")
	dockerfileBytes, err := readDockerfile(dockerfile)
	if err != nil {
		return traceExitDockerfileBuildingError(err, dockerfile, "Failed to read Dockerfile to create tar")
	}

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	build, err := createDockerfileBuild(dockerfile, dockerfileBytes, tw, buf)
	if err != nil {
		return traceExitOfError(err, "Failed to create tar buffer for Dockerfile")
	}
	ctx := context.Background()

	return handler.wrapper.BuildImage(ctx, build, types.ImageBuildOptions{
		Tags:       []string{imageName},
		Context:    build,
		Dockerfile: dockerfile,
	})
}

func createDockerfileBuild(dockerfile string, dockerfileBytes []byte, tw *tar.Writer, buf *bytes.Buffer) (io.Reader, error) {
	logger.Trace().Str("Dockerfile", dockerfile).Msg("Creating tar buffer for Dockerfile")

	if err := writeTarHeader(dockerfile, dockerfileBytes, tw); err != nil {
		return nil, traceExitDockerfileBuildingError(err, dockerfile, "Failed to write tar header for Dockerfile")
	}

	if err := writeTarBytes(dockerfile, dockerfileBytes, tw); err != nil {
		return nil, traceExitDockerfileBuildingError(err, dockerfile, "Failed to write tar to buffer for Dockerfile")
	}

	logger.Trace().Str("Dockerfile", dockerfile).Msg("Finished creating tar for Dockerfile")
	return bytes.NewReader(buf.Bytes()), nil
}

func writeTarHeader(dockerfile string, dockerfileBytes []byte, tw *tar.Writer) error {
	logger.Trace().Str("Dockerfile", dockerfile).Msg("Writing tar header")
	tarHeader := &tar.Header{
		Name: dockerfile,
		Size: int64(len(dockerfileBytes)),
	}

	if err := tw.WriteHeader(tarHeader); err != nil {
		return traceExitDockerfileBuildingError(err, dockerfile, "Failed to write header of tar buffer")
	}

	return traceExitDockerfileBuildingError(nil, dockerfile, "Successfully wrote header of tar buffer")
}

func writeTarBytes(dockerfile string, dockerfileBytes []byte, tw *tar.Writer) error {
	if _, err := tw.Write(dockerfileBytes); err != nil {
		return traceExitDockerfileBuildingError(err, dockerfile, "Failed to write tar bytes to buffer")
	}
	return traceExitDockerfileBuildingError(nil, dockerfile, "Successfully wrote tar bytes to buffer")
}

func readDockerfile(dockerfile string) ([]byte, error) {
	logger.Trace().Str("Dockerfile", dockerfile).Msg("Reading Dockerfile")
	dockerfileReader, err := os.Open(getDockerfilePath(dockerfile))
	if err != nil {
		return nil, traceExitDockerfileBuildingError(err, dockerfile, "Failed to open Dockerfile")
	}
	logger.Trace().Str("Dockerfile", dockerfile).Msg("Finished reading Dockerfile")
	return ioutil.ReadAll(dockerfileReader)
}

func getDockerfilePath(dockerfile string) string {
	return fmt.Sprintf("%s/%s", config.GetOrDefault(util.DockerfileDirEnv), dockerfile)
}

func traceExitDockerfileBuildingError(err error, dockerfile, msg string) error {
	logger.Trace().Err(err).Str("Dockerfile", dockerfile).Msg(msg)
	return err
}

func traceExitOfError(err error, msg string) error {
	log.Trace().
		Err(err).
		Msg(msg)
	return err
}
