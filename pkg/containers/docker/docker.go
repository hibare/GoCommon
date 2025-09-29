package docker

import (
	"context"

	"github.com/distribution/reference"
	containerType "github.com/docker/docker/api/types/container"
	imageType "github.com/docker/docker/api/types/image"
	networkType "github.com/docker/docker/api/types/network"
	volumeType "github.com/docker/docker/api/types/volume"
)

// DockerAPIIface is the interface for the Docker API.
type DockerAPIIface interface {
	ContainerList(ctx context.Context, options containerType.ListOptions) ([]containerType.Summary, error)
	ImageList(ctx context.Context, options imageType.ListOptions) ([]imageType.Summary, error)
	NetworkList(ctx context.Context, options networkType.ListOptions) ([]networkType.Summary, error)
	VolumeList(ctx context.Context, options volumeType.ListOptions) (volumeType.ListResponse, error)
}

// ClientIface is the interface for the Docker client.
type ClientIface interface {
	NormalizeImageName(image string) (string, error)
	ContainerList(ctx context.Context, options containerType.ListOptions) ([]containerType.Summary, error)
	ImageList(ctx context.Context, options imageType.ListOptions) ([]imageType.Summary, error)
	NetworkList(ctx context.Context, options networkType.ListOptions) ([]networkType.Summary, error)
	VolumeList(ctx context.Context, options volumeType.ListOptions) (volumeType.ListResponse, error)
}

// Client is the implementation of the ClientIface.
type Client struct {
	client DockerAPIIface
}

// NormalizeImageName normalizes the image name to a fully qualified name & tag. Add latest tag if not present.
func (c *Client) NormalizeImageName(image string) (string, error) {
	named, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return "", err
	}

	// Ensure tag is present, default to "latest".
	named = reference.TagNameOnly(named)

	return named.String(), nil
}

// ContainerList lists all containers.
func (c *Client) ContainerList(ctx context.Context, options containerType.ListOptions) ([]containerType.Summary, error) {
	return c.client.ContainerList(ctx, options)
}

// ImageList lists all images.
func (c *Client) ImageList(ctx context.Context, options imageType.ListOptions) ([]imageType.Summary, error) {
	images, err := c.client.ImageList(ctx, options)
	if err != nil {
		return nil, err
	}

	normalizedImages := make([]imageType.Summary, 0, len(images))
	for _, img := range images {
		normalizedTags := make([]string, 0, len(img.RepoTags))
		for _, tag := range img.RepoTags {
			normalizedTag, err := c.NormalizeImageName(tag)
			if err != nil {
				return nil, err
			}
			normalizedTags = append(normalizedTags, normalizedTag)
		}
		img.RepoTags = normalizedTags
		normalizedImages = append(normalizedImages, img)
	}

	return normalizedImages, nil
}

// NetworkList lists all networks.
func (c *Client) NetworkList(ctx context.Context, options networkType.ListOptions) ([]networkType.Summary, error) {
	return c.client.NetworkList(ctx, options)
}

// VolumeList lists all volumes.
func (c *Client) VolumeList(ctx context.Context, options volumeType.ListOptions) (volumeType.ListResponse, error) {
	return c.client.VolumeList(ctx, options)
}
