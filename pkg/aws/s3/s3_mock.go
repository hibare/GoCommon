// Package s3 provides a mock implementation of the * interface.
package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/mock"
)

// mockServiceAPI is a mock implementation of the S3Client interface.
type mockServiceAPI struct{ mock.Mock }

// PutObject is a mock implementation of the PutObject method.
func (m *mockServiceAPI) PutObject(ctx context.Context, params *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.PutObjectOutput), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}

// ListObjectsV2 is a mock implementation of the ListObjectsV2 method.
func (m *mockServiceAPI) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, _ ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.ListObjectsV2Output), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}

// DeleteObject is a mock implementation of the DeleteObject method.
func (m *mockServiceAPI) DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, _ ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.DeleteObjectOutput), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}

// ListObjects is a mock implementation of the ListObjects method.
func (m *mockServiceAPI) ListObjects(ctx context.Context, params *s3.ListObjectsInput, _ ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
	args := m.Called(ctx, params)
	// Ensure Get(0) is nil-checked if it can be nil for an error case, or adjust mock.
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.ListObjectsOutput), args.Error(1) //nolint:errcheck // reason: type assertion on mock, error not possible/needed
}
