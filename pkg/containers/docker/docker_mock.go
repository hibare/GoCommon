package docker

import (
	"context"
	"testing"

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

// MockClient is a mock implementation of the Client interface.
type MockClient struct {
	mock.Mock
}

// NormalizeImageName is a mock implementation of the NormalizeImageName method.
func (m *MockClient) NormalizeImageName(image string) (string, error) {
	args := m.Called(image)
	return args.Get(0).(string), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}

// ContainerList is a mock implementation of the ContainerList method.
func (m *MockClient) ContainerList(ctx context.Context, options containerType.ListOptions) ([]containerType.Summary, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]containerType.Summary), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}

// ImageList is a mock implementation of the ImageList method.
func (m *MockClient) ImageList(ctx context.Context, options imageType.ListOptions) ([]imageType.Summary, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]imageType.Summary), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}

// NetworkList is a mock implementation of the NetworkList method.
func (m *MockClient) NetworkList(ctx context.Context, options networkType.ListOptions) ([]networkType.Summary, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]networkType.Summary), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}

// VolumeList is a mock implementation of the VolumeList method.
func (m *MockClient) VolumeList(ctx context.Context, options volumeType.ListOptions) (volumeType.ListResponse, error) {
	args := m.Called(ctx, options)
	return args.Get(0).(volumeType.ListResponse), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}

// SetMockClient sets the mock client for the Docker package.
func SetMockClient(t *testing.T) *MockClient {
	mockClient := &MockClient{}
	NewClient = func(_ context.Context, _ Options) (ClientIface, error) {
		return mockClient, nil
	}
	t.Cleanup(func() {
		NewClient = newClient
		mockClient.AssertExpectations(t)
	})
	return mockClient
}
