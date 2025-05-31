package s3

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hibare/GoCommon/v2/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- Mocks ---
type MockS3Client struct{ mock.Mock }

func (m *MockS3Client) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.PutObjectOutput), args.Error(1)
}
func (m *MockS3Client) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.ListObjectsV2Output), args.Error(1)
}
func (m *MockS3Client) DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.DeleteObjectOutput), args.Error(1)
}
func (m *MockS3Client) ListObjects(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
	args := m.Called(ctx, params)
	// Ensure Get(0) is nil-checked if it can be nil for an error case, or adjust mock.
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.ListObjectsOutput), args.Error(1)
}

func TestS3(t *testing.T) {
	ctx := context.Background()
	temp := t.TempDir()

	// create  two files in the temp dir
	os.WriteFile(filepath.Join(temp, "file1"), []byte("file1"), 0644)
	os.WriteFile(filepath.Join(temp, "file2"), []byte("file2"), 0644)

	t.Run("GetPrefix and GetTimestampedPrefix", func(t *testing.T) {
		s3 := &S3{}
		prefix := s3.GetPrefix("foo", "bar")
		require.Equal(t, "foo/bar/", prefix)
		tsPrefix := s3.GetTimestampedPrefix("foo", "bar")
		require.Contains(t, tsPrefix, "foo/bar/")
		require.True(t, len(tsPrefix) > len("foo/bar/"))
	})

	t.Run("TrimPrefix", func(t *testing.T) {
		s3 := &S3{}
		keys := []string{"prefix/1", "prefix/2", "prefix/3/"}
		trimmed := s3.TrimPrefix(keys, "prefix/")
		assert.Equal(t, []string{"1", "2", "3"}, trimmed)
	})

	t.Run("UploadDir", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}

			mockClient.On("PutObject", ctx, mock.Anything).Return(&s3.PutObjectOutput{}, nil).Twice()

			resp, err := s3Client.UploadDir(ctx, "bucket", "prefix", temp, nil)
			assert.NoError(t, err)
			assert.Equal(t, 2, resp.SuccessFiles)
			assert.Equal(t, 2, resp.TotalFiles)
			assert.Equal(t, 1, resp.TotalDirs)
			assert.Empty(t, resp.FailedFiles)
		})
		t.Run("error", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}

			mockClient.On("PutObject", ctx, mock.Anything).Return(nil, errors.New("fail")).Twice()

			resp, err := s3Client.UploadDir(ctx, "bucket", "prefix", temp, nil)
			assert.NoError(t, err)
			assert.Equal(t, 0, resp.SuccessFiles)
			assert.Equal(t, 2, resp.TotalFiles)
			assert.Equal(t, 1, resp.TotalDirs)
			assert.Equal(t, 2, len(resp.FailedFiles))
			assert.Contains(t, resp.FailedFiles, filepath.Join(temp, "file1"))
			assert.Contains(t, resp.FailedFiles[filepath.Join(temp, "file1")].Error(), "fail")
		})
	})

	t.Run("UploadFile", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}
			mockClient.On("PutObject", ctx, mock.Anything).Return(&s3.PutObjectOutput{}, nil)
			key, err := s3Client.UploadFile(ctx, "bucket", "prefix", filepath.Join(temp, "file1"))
			assert.NoError(t, err)
			assert.Contains(t, key, "file1")
		})

		t.Run("upload error", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}
			mockClient.On("PutObject", ctx, mock.Anything).Return(nil, errors.New("fail"))
			key, err := s3Client.UploadFile(ctx, "bucket", "prefix", filepath.Join(temp, "file1"))
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "fail")
			assert.Empty(t, key)
		})
	})

	t.Run("ListObjectsAtPrefixRoot", func(t *testing.T) {
		t.Run("with results", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}
			mockClient.On("ListObjectsV2", ctx, mock.Anything).Return(&s3.ListObjectsV2Output{
				Contents:       []types.Object{{Key: utils.ToPtr("prefix/file1")}},
				CommonPrefixes: []types.CommonPrefix{{Prefix: utils.ToPtr("prefix/")}},
			}, nil)
			keys, err := s3Client.ListObjectsAtPrefixRoot(ctx, "bucket", "prefix")
			assert.NoError(t, err)
			assert.Contains(t, keys, "prefix/file1")
			assert.Contains(t, keys, "prefix/")
		})
		t.Run("no results", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}
			mockClient.On("ListObjectsV2", ctx, mock.Anything).Return(&s3.ListObjectsV2Output{}, nil)
			keys, err := s3Client.ListObjectsAtPrefixRoot(ctx, "bucket", "prefix")
			assert.NoError(t, err)
			assert.Empty(t, keys)
		})
		t.Run("error", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}
			mockClient.On("ListObjectsV2", ctx, mock.Anything).Return(nil, errors.New("fail"))
			keys, err := s3Client.ListObjectsAtPrefixRoot(ctx, "bucket", "prefix")
			assert.Error(t, err)
			assert.Empty(t, keys)
		})
	})

	t.Run("DeleteObjects", func(t *testing.T) {
		t.Run("non-recursive success", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}
			mockClient.On("DeleteObject", ctx, mock.Anything).Return(&s3.DeleteObjectOutput{}, nil)
			err := s3Client.DeleteObjects(ctx, "bucket", "key", false)
			assert.NoError(t, err)
		})
		t.Run("recursive success", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}
			mockClient.On("ListObjects", ctx, mock.Anything).Return(&s3.ListObjectsOutput{
				Contents: []types.Object{{Key: utils.ToPtr("key/1")}, {Key: utils.ToPtr("key/2")}},
			}, nil)
			mockClient.On("DeleteObject", ctx, mock.MatchedBy(func(input *s3.DeleteObjectInput) bool {
				return (input.Key != nil && (*input.Key == "key/1" || *input.Key == "key/2"))
			})).Return(&s3.DeleteObjectOutput{}, nil).Twice()
			mockClient.On("DeleteObject", ctx, mock.Anything).Return(&s3.DeleteObjectOutput{}, nil)
			err := s3Client.DeleteObjects(ctx, "bucket", "key", true)
			assert.NoError(t, err)
		})
		t.Run("list error", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}
			mockClient.On("ListObjects", ctx, mock.Anything).Return(nil, errors.New("fail"))
			err := s3Client.DeleteObjects(ctx, "bucket", "key", true)
			assert.Error(t, err)
		})
		t.Run("delete error", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}
			objectKey := "key/1"
			mockClient.On("ListObjects", ctx, mock.Anything).Return(&s3.ListObjectsOutput{
				Contents: []types.Object{{Key: utils.ToPtr(objectKey)}},
			}, nil)
			// Mock DeleteObject to fail for the specific key
			mockClient.On("DeleteObject", ctx, mock.MatchedBy(func(input *s3.DeleteObjectInput) bool {
				return input.Key != nil && *input.Key == objectKey
			})).Return(nil, errors.New("delete failed"))
			err := s3Client.DeleteObjects(ctx, "bucket", "key", true)
			assert.Error(t, err)
			// Optionally, check the error message if DeleteObjects wraps or returns the specific error
			assert.Contains(t, err.Error(), "delete failed")
		})
	})
}
