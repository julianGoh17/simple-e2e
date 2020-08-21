package client

import (
	"github.com/docker/docker/client"
)

// WrapperClient is the framework's wrapper for the Docker client which will be used to create docker images
type WrapperClient struct {
	Cli *client.Client
}

// Initialize will create a docker client for the framework to use to communicate with the host machines daemon
func (wrapper *WrapperClient) Initialize() error {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	wrapper.Cli = dockerClient

	return nil
}

// Close will close the framework's docker client
func (wrapper *WrapperClient) Close() error {
	return wrapper.Cli.Close()
}
