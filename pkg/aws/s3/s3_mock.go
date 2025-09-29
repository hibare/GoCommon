// Package s3 provides a mock implementation of the * interface.
package s3

import (
	"context"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/mock"
)

// mockS3API is a mock implementation of the S3Client interface.
type mockS3API struct{ mock.Mock }

// PutObject is a mock implementation of the PutObject method.
func (m *mockS3API) PutObject(ctx context.Context, params *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.PutObjectOutput), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}

// ListObjectsV2 is a mock implementation of the ListObjectsV2 method.
func (m *mockS3API) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, _ ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.ListObjectsV2Output), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}

// DeleteObject is a mock implementation of the DeleteObject method.
func (m *mockS3API) DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, _ ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.DeleteObjectOutput), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}

// ListObjects is a mock implementation of the ListObjects method.
func (m *mockS3API) ListObjects(ctx context.Context, params *s3.ListObjectsInput, _ ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
	args := m.Called(ctx, params)
	// Ensure Get(0) is nil-checked if it can be nil for an error case, or adjust mock.
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.ListObjectsOutput), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}

// MockClient is a mock implementation of the Client interface.
type MockClient struct {
	mock.Mock
}

// BuildKey is a mock implementation of the BuildKey method.
func (m *MockClient) BuildKey(parts ...string) string {
	args := m.Called(parts)
	return args.Get(0).(string) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}

// BuildTimestampedKey is a mock implementation of the BuildTimestampedKey method.
func (m *MockClient) BuildTimestampedKey(parts ...string) string {
	args := m.Called(parts)
	return args.Get(0).(string) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}

// TrimPrefix is a mock implementation of the TrimPrefix method.
func (m *MockClient) TrimPrefix(keys []string, prefix string) []string {
	args := m.Called(keys, prefix)
	return args.Get(0).([]string) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}

// UploadDir is a mock implementation of the UploadDir method.
func (m *MockClient) UploadDir(ctx context.Context, bucket, prefix, baseDir string, exclude []*regexp.Regexp) (UploadDirResponse, error) {
	args := m.Called(ctx, bucket, prefix, baseDir, exclude)
	return args.Get(0).(UploadDirResponse), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}

// UploadFile is a mock implementation of the UploadFile method.
func (m *MockClient) UploadFile(ctx context.Context, bucket, prefix, filePath string) (string, error) {
	args := m.Called(ctx, bucket, prefix, filePath)
	return args.Get(0).(string), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}

// ListObjectsAtPrefix is a mock implementation of the ListObjectsAtPrefix method.
func (m *MockClient) ListObjectsAtPrefix(ctx context.Context, bucket, prefix string) ([]string, error) {
	args := m.Called(ctx, bucket, prefix)
	return args.Get(0).([]string), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}

// DeleteObjects is a mock implementation of the DeleteObjects method.
func (m *MockClient) DeleteObjects(ctx context.Context, bucket, key string, recursive bool) error {
	args := m.Called(ctx, bucket, key, recursive)
	return args.Error(0)
}

// SetMockClient sets the mock client for the S3 package.
func SetMockClient(t *testing.T) *MockClient {
	mockClient := new(MockClient)
	NewClient = func(_ context.Context, _ Options) (ClientIface, error) {
		return mockClient, nil
	}
	t.Cleanup(func() {
		NewClient = newClient
		mockClient.AssertExpectations(t)
	})
	return mockClient
}
