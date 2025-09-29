package docker

import (
	"testing"

	containerType "github.com/docker/docker/api/types/container"
	imageType "github.com/docker/docker/api/types/image"
	networkType "github.com/docker/docker/api/types/network"
	volumeType "github.com/docker/docker/api/types/volume"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_NormalizeImageName(t *testing.T) {
	tests := []struct {
		name     string
		image    string
		expected string
		wantErr  bool
	}{
		{
			name:     "image with tag",
			image:    "nginx:1.20",
			expected: "docker.io/library/nginx:1.20",
			wantErr:  false,
		},
		{
			name:     "image without tag",
			image:    "nginx",
			expected: "docker.io/library/nginx:latest",
			wantErr:  false,
		},
		{
			name:     "image with digest",
			image:    "nginx@sha256:abc123def4567890123456789012345678901234567890123456789012345678",
			expected: "docker.io/library/nginx@sha256:abc123def4567890123456789012345678901234567890123456789012345678",
			wantErr:  false,
		},
		{
			name:     "fully qualified image name",
			image:    "docker.io/library/nginx:1.20",
			expected: "docker.io/library/nginx:1.20",
			wantErr:  false,
		},
		{
			name:     "private registry image",
			image:    "registry.example.com/myapp:v1.0",
			expected: "registry.example.com/myapp:v1.0",
			wantErr:  false,
		},
		{
			name:     "invalid image name",
			image:    "",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "malformed image name",
			image:    "invalid:",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{}
			result, err := client.NormalizeImageName(tt.image)

			if tt.wantErr {
				require.Error(t, err)
				assert.Empty(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestClient_ContainerList(t *testing.T) {
	ctx := t.Context()
	options := containerType.ListOptions{All: true}

	t.Run("successful container list", func(t *testing.T) {
		mockAPI := &mockDockerAPI{}
		client := &Client{client: mockAPI}

		expectedContainers := []containerType.Summary{
			{ID: "container1", Names: []string{"/test1"}},
			{ID: "container2", Names: []string{"/test2"}},
		}

		mockAPI.On("ContainerList", ctx, options).Return(expectedContainers, nil)

		containers, err := client.ContainerList(ctx, options)

		require.NoError(t, err)
		assert.Equal(t, expectedContainers, containers)
		mockAPI.AssertExpectations(t)
	})

	t.Run("container list error", func(t *testing.T) {
		mockAPI := &mockDockerAPI{}
		client := &Client{client: mockAPI}

		expectedError := assert.AnError
		mockAPI.On("ContainerList", ctx, options).Return([]containerType.Summary(nil), expectedError)

		containers, err := client.ContainerList(ctx, options)

		require.Error(t, err)
		assert.Nil(t, containers)
		assert.Equal(t, expectedError, err)
		mockAPI.AssertExpectations(t)
	})
}

func TestClient_ImageList(t *testing.T) {
	ctx := t.Context()
	options := imageType.ListOptions{All: true}

	t.Run("successful image list with normalization", func(t *testing.T) {
		mockAPI := &mockDockerAPI{}
		client := &Client{client: mockAPI}

		rawImages := []imageType.Summary{
			{
				ID:       "image1",
				RepoTags: []string{"nginx", "nginx:1.20"},
			},
			{
				ID:       "image2",
				RepoTags: []string{"alpine:3.18"},
			},
		}

		expectedImages := []imageType.Summary{
			{
				ID:       "image1",
				RepoTags: []string{"docker.io/library/nginx:latest", "docker.io/library/nginx:1.20"},
			},
			{
				ID:       "image2",
				RepoTags: []string{"docker.io/library/alpine:3.18"},
			},
		}

		mockAPI.On("ImageList", ctx, options).Return(rawImages, nil)

		images, err := client.ImageList(ctx, options)

		require.NoError(t, err)
		assert.Equal(t, expectedImages, images)
		mockAPI.AssertExpectations(t)
	})

	t.Run("image list with empty repo tags", func(t *testing.T) {
		mockAPI := &mockDockerAPI{}
		client := &Client{client: mockAPI}

		rawImages := []imageType.Summary{
			{
				ID:       "image1",
				RepoTags: []string{},
			},
		}

		expectedImages := []imageType.Summary{
			{
				ID:       "image1",
				RepoTags: []string{},
			},
		}

		mockAPI.On("ImageList", ctx, options).Return(rawImages, nil)

		images, err := client.ImageList(ctx, options)

		require.NoError(t, err)
		assert.Equal(t, expectedImages, images)
		mockAPI.AssertExpectations(t)
	})

	t.Run("image list with invalid tag normalization", func(t *testing.T) {
		mockAPI := &mockDockerAPI{}
		client := &Client{client: mockAPI}

		rawImages := []imageType.Summary{
			{
				ID:       "image1",
				RepoTags: []string{"invalid:"},
			},
		}

		mockAPI.On("ImageList", ctx, options).Return(rawImages, nil)

		images, err := client.ImageList(ctx, options)

		require.Error(t, err)
		assert.Nil(t, images)
		mockAPI.AssertExpectations(t)
	})

	t.Run("image list error from API", func(t *testing.T) {
		mockAPI := &mockDockerAPI{}
		client := &Client{client: mockAPI}

		expectedError := assert.AnError
		mockAPI.On("ImageList", ctx, options).Return([]imageType.Summary(nil), expectedError)

		images, err := client.ImageList(ctx, options)

		require.Error(t, err)
		assert.Nil(t, images)
		assert.Equal(t, expectedError, err)
		mockAPI.AssertExpectations(t)
	})
}

func TestClient_NetworkList(t *testing.T) {
	ctx := t.Context()
	options := networkType.ListOptions{}

	t.Run("successful network list", func(t *testing.T) {
		mockAPI := &mockDockerAPI{}
		client := &Client{client: mockAPI}

		expectedNetworks := []networkType.Summary{
			{ID: "network1", Name: "bridge"},
			{ID: "network2", Name: "host"},
		}

		mockAPI.On("NetworkList", ctx, options).Return(expectedNetworks, nil)

		networks, err := client.NetworkList(ctx, options)

		require.NoError(t, err)
		assert.Equal(t, expectedNetworks, networks)
		mockAPI.AssertExpectations(t)
	})

	t.Run("network list error", func(t *testing.T) {
		mockAPI := &mockDockerAPI{}
		client := &Client{client: mockAPI}

		expectedError := assert.AnError
		mockAPI.On("NetworkList", ctx, options).Return([]networkType.Summary(nil), expectedError)

		networks, err := client.NetworkList(ctx, options)

		require.Error(t, err)
		assert.Nil(t, networks)
		assert.Equal(t, expectedError, err)
		mockAPI.AssertExpectations(t)
	})
}

func TestClient_VolumeList(t *testing.T) {
	ctx := t.Context()
	options := volumeType.ListOptions{}

	t.Run("successful volume list", func(t *testing.T) {
		mockAPI := &mockDockerAPI{}
		client := &Client{client: mockAPI}

		expectedVolumes := volumeType.ListResponse{
			Volumes: []*volumeType.Volume{
				{Name: "volume1", Driver: "local"},
				{Name: "volume2", Driver: "local"},
			},
		}

		mockAPI.On("VolumeList", ctx, options).Return(expectedVolumes, nil)

		volumes, err := client.VolumeList(ctx, options)

		require.NoError(t, err)
		assert.Equal(t, expectedVolumes, volumes)
		mockAPI.AssertExpectations(t)
	})

	t.Run("volume list error", func(t *testing.T) {
		mockAPI := &mockDockerAPI{}
		client := &Client{client: mockAPI}

		expectedError := assert.AnError
		mockAPI.On("VolumeList", ctx, options).Return(volumeType.ListResponse{}, expectedError)

		volumes, err := client.VolumeList(ctx, options)

		require.Error(t, err)
		assert.Empty(t, volumes)
		assert.Equal(t, expectedError, err)
		mockAPI.AssertExpectations(t)
	})
}

func TestClient_Integration(t *testing.T) {
	t.Run("client with all methods", func(t *testing.T) {
		mockAPI := &mockDockerAPI{}
		client := &Client{client: mockAPI}
		ctx := t.Context()

		// Test that all methods can be called
		_, err := client.NormalizeImageName("nginx:latest")
		require.NoError(t, err)

		mockAPI.On("ContainerList", ctx, containerType.ListOptions{}).Return([]containerType.Summary{}, nil)
		_, err = client.ContainerList(ctx, containerType.ListOptions{})
		require.NoError(t, err)

		mockAPI.On("ImageList", ctx, imageType.ListOptions{}).Return([]imageType.Summary{}, nil)
		_, err = client.ImageList(ctx, imageType.ListOptions{})
		require.NoError(t, err)

		mockAPI.On("NetworkList", ctx, networkType.ListOptions{}).Return([]networkType.Summary{}, nil)
		_, err = client.NetworkList(ctx, networkType.ListOptions{})
		require.NoError(t, err)

		mockAPI.On("VolumeList", ctx, volumeType.ListOptions{}).Return(volumeType.ListResponse{}, nil)
		_, err = client.VolumeList(ctx, volumeType.ListOptions{})
		require.NoError(t, err)

		mockAPI.AssertExpectations(t)
	})
}

func TestClient_EdgeCases(t *testing.T) {
	t.Run("empty image list", func(t *testing.T) {
		mockAPI := &mockDockerAPI{}
		client := &Client{client: mockAPI}
		ctx := t.Context()

		mockAPI.On("ImageList", ctx, imageType.ListOptions{}).Return([]imageType.Summary{}, nil)

		images, err := client.ImageList(ctx, imageType.ListOptions{})
		require.NoError(t, err)
		assert.Empty(t, images)
		mockAPI.AssertExpectations(t)
	})

	t.Run("image with multiple tags", func(t *testing.T) {
		mockAPI := &mockDockerAPI{}
		client := &Client{client: mockAPI}
		ctx := t.Context()

		rawImages := []imageType.Summary{
			{
				ID:       "image1",
				RepoTags: []string{"nginx", "nginx:1.20", "nginx:latest"},
			},
		}

		expectedImages := []imageType.Summary{
			{
				ID:       "image1",
				RepoTags: []string{"docker.io/library/nginx:latest", "docker.io/library/nginx:1.20", "docker.io/library/nginx:latest"},
			},
		}

		mockAPI.On("ImageList", ctx, imageType.ListOptions{}).Return(rawImages, nil)

		images, err := client.ImageList(ctx, imageType.ListOptions{})
		require.NoError(t, err)
		assert.Equal(t, expectedImages, images)
		mockAPI.AssertExpectations(t)
	})
}
