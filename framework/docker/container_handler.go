package docker

import (
	"context"
)

// Handler is the framework's controller responsible for all docker related operations
type Handler struct {
	wrapper *WrapperClient
	// TODO: will need container manager to manage running containers
}

// NewHandler will create a handler object intialized and ready to use. Will error if there is any problems with setting up for docker operations
func NewHandler() (*Handler, error) {
	ctx := context.Background()
	handler := &Handler{}
	handler.wrapper = &WrapperClient{}
	if err := handler.wrapper.Initialize(); err != nil {
		return nil, err
	}
	handler.wrapper.Cli.NegotiateAPIVersion(ctx)

	return handler, nil
}

// PullImage will pull the image from dockerhub onto the host machine's daemon
func (handler *Handler) PullImage(image string) error {
	ctx := context.Background()
	return handler.wrapper.PullImage(ctx, image)
}

// BuildImage will build an image from a specified Dockrefile onto the host machine's daemon
func (handler *Handler) BuildImage() error {
	return nil
}
