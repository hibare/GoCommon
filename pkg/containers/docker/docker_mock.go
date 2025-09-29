package docker

import (
	"context"

	containerType "github.com/docker/docker/api/types/container"
	imageType "github.com/docker/docker/api/types/image"
	networkType "github.com/docker/docker/api/types/network"
	volumeType "github.com/docker/docker/api/types/volume"
	"github.com/stretchr/testify/mock"
)

type mockDockerAPI struct {
	mock.Mock
}

func (m *mockDockerAPI) ContainerList(ctx context.Context, options containerType.ListOptions) ([]containerType.Summary, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]containerType.Summary), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}

func (m *mockDockerAPI) ImageList(ctx context.Context, options imageType.ListOptions) ([]imageType.Summary, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]imageType.Summary), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}

func (m *mockDockerAPI) NetworkList(ctx context.Context, options networkType.ListOptions) ([]networkType.Summary, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]networkType.Summary), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}

func (m *mockDockerAPI) VolumeList(ctx context.Context, options volumeType.ListOptions) (volumeType.ListResponse, error) {
	args := m.Called(ctx, options)
	return args.Get(0).(volumeType.ListResponse), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}
