package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
)

const minVersion = "1.18"


// MobyClient is wrapper of docker client
type MobyClient struct {
	instance *client.Client
}

// NewMobyClients creates instance of MobyClient
func NewMobyClients() (*MobyClient, error) {
	instance, err := client.NewClientWithOpts(client.WithVersion(minVersion))
	if err != nil {
		return nil, err

	}
	return &MobyClient{
		instance: instance,
	}, nil
}

// PullImage starts pulling image
func (c *MobyClient) PullImage(ctx context.Context, imageName string) (io.ReadCloser, error) {
	return c.instance.ImagePull(ctx, imageName, types.ImagePullOptions{})

}

// CreateContainer creates container from image
func (c *MobyClient) CreateContainer(ctx context.Context, imageName string) (string, error) {
	resp, err := c.instance.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, "")
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

// Run starts container by containerID
func (c *MobyClient) Run(ctx context.Context, containerID string) error {
	err := c.instance.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	return nil
}

// List returns a list of container
func (c *MobyClient) List(ctx context.Context) ([]types.Container, error) {
	return c.instance.ContainerList(ctx, types.ContainerListOptions{})

}

// Stop stops container by ID
func (c *MobyClient) Stop(ctx context.Context, containerID string) error {
	return c.instance.ContainerStop(context.Background(), containerID, nil)
}
