package s3

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hibare/GoCommon/v2/pkg/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestS3(t *testing.T) {
	temp := t.TempDir()

	// create  two files in the temp dir
	err := os.WriteFile(filepath.Join(temp, "file1"), []byte("file1"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(temp, "file2"), []byte("file2"), 0644)
	require.NoError(t, err)

	t.Run("GetPrefix and GetTimestampedPrefix", func(t *testing.T) {
		s3 := &S3{}
		prefix := s3.GetPrefix("foo", "bar")
		require.Equal(t, "foo/bar/", prefix)
		tsPrefix := s3.GetTimestampedPrefix("foo", "bar")
		require.Contains(t, tsPrefix, "foo/bar/")
		require.Greater(t, len(tsPrefix), len("foo/bar/"))
	})

	t.Run("TrimPrefix", func(t *testing.T) {
		s3 := &S3{}
		keys := []string{"prefix/1", "prefix/2", "prefix/3/"}
		trimmed := s3.TrimPrefix(keys, "prefix/")
		require.Equal(t, []string{"1", "2", "3"}, trimmed)
	})

	t.Run("UploadDir", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}

			mockClient.On("PutObject", t.Context(), mock.Anything).Return(&s3.PutObjectOutput{}, nil).Twice()

			resp, err := s3Client.UploadDir(t.Context(), "bucket", "prefix", temp, nil)
			require.NoError(t, err)
			require.Equal(t, 2, resp.SuccessFiles)
			require.Equal(t, 2, resp.TotalFiles)
			require.Equal(t, 1, resp.TotalDirs)
			require.Empty(t, resp.FailedFiles)
		})
		t.Run("error", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}

			mockClient.On("PutObject", t.Context(), mock.Anything).Return(nil, errors.New("fail")).Twice()

			resp, err := s3Client.UploadDir(t.Context(), "bucket", "prefix", temp, nil)
			require.NoError(t, err)
			require.Equal(t, 0, resp.SuccessFiles)
			require.Equal(t, 2, resp.TotalFiles)
			require.Equal(t, 1, resp.TotalDirs)
			require.Len(t, resp.FailedFiles, 2)
			require.Contains(t, resp.FailedFiles[filepath.Join(temp, "file1")].Error(), "fail")
		})
	})

	t.Run("UploadFile", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}
			mockClient.On("PutObject", t.Context(), mock.Anything).Return(&s3.PutObjectOutput{}, nil)
			key, err := s3Client.UploadFile(t.Context(), "bucket", "prefix", filepath.Join(temp, "file1"))
			require.NoError(t, err)
			require.Contains(t, key, "file1")
		})

		t.Run("upload error", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}
			mockClient.On("PutObject", t.Context(), mock.Anything).Return(nil, errors.New("fail"))
			key, err := s3Client.UploadFile(t.Context(), "bucket", "prefix", filepath.Join(temp, "file1"))
			require.Error(t, err)
			require.Contains(t, err.Error(), "fail")
			require.Empty(t, key)
		})
	})

	t.Run("ListObjectsAtPrefixRoot", func(t *testing.T) {
		t.Run("with results", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}
			mockClient.On("ListObjectsV2", t.Context(), mock.Anything).Return(&s3.ListObjectsV2Output{
				Contents:       []types.Object{{Key: utils.ToPtr("prefix/file1")}},
				CommonPrefixes: []types.CommonPrefix{{Prefix: utils.ToPtr("prefix/")}},
			}, nil)
			keys, err := s3Client.ListObjectsAtPrefixRoot(t.Context(), "bucket", "prefix")
			require.NoError(t, err)
			require.Contains(t, keys, "prefix/file1")
			require.Contains(t, keys, "prefix/")
		})
		t.Run("no results", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}
			mockClient.On("ListObjectsV2", t.Context(), mock.Anything).Return(&s3.ListObjectsV2Output{}, nil)
			keys, err := s3Client.ListObjectsAtPrefixRoot(t.Context(), "bucket", "prefix")
			require.NoError(t, err)
			require.Empty(t, keys)
		})
		t.Run("error", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}
			mockClient.On("ListObjectsV2", t.Context(), mock.Anything).Return(nil, errors.New("fail"))
			keys, err := s3Client.ListObjectsAtPrefixRoot(t.Context(), "bucket", "prefix")
			require.Error(t, err)
			require.Empty(t, keys)
		})
	})

	t.Run("DeleteObjects", func(t *testing.T) {
		t.Run("non-recursive success", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}
			mockClient.On("DeleteObject", t.Context(), mock.Anything).Return(&s3.DeleteObjectOutput{}, nil)
			err := s3Client.DeleteObjects(t.Context(), "bucket", "key", false)
			require.NoError(t, err)
		})
		t.Run("recursive success", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}
			mockClient.On("ListObjects", t.Context(), mock.Anything).Return(&s3.ListObjectsOutput{
				Contents: []types.Object{{Key: utils.ToPtr("key/1")}, {Key: utils.ToPtr("key/2")}},
			}, nil)
			mockClient.On("DeleteObject", t.Context(), mock.MatchedBy(func(input *s3.DeleteObjectInput) bool {
				return (input.Key != nil && (*input.Key == "key/1" || *input.Key == "key/2"))
			})).Return(&s3.DeleteObjectOutput{}, nil).Twice()
			mockClient.On("DeleteObject", t.Context(), mock.Anything).Return(&s3.DeleteObjectOutput{}, nil)
			err := s3Client.DeleteObjects(t.Context(), "bucket", "key", true)
			require.NoError(t, err)
		})
		t.Run("list error", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}
			mockClient.On("ListObjects", t.Context(), mock.Anything).Return(nil, errors.New("fail"))
			err := s3Client.DeleteObjects(t.Context(), "bucket", "key", true)
			require.Error(t, err)
		})
		t.Run("delete error", func(t *testing.T) {
			mockClient := new(MockS3Client)
			s3Client := &S3{Client: mockClient}
			objectKey := "key/1"
			mockClient.On("ListObjects", t.Context(), mock.Anything).Return(&s3.ListObjectsOutput{
				Contents: []types.Object{{Key: utils.ToPtr(objectKey)}},
			}, nil)
			// Mock DeleteObject to fail for the specific key
			mockClient.On("DeleteObject", t.Context(), mock.MatchedBy(func(input *s3.DeleteObjectInput) bool {
				return input.Key != nil && *input.Key == objectKey
			})).Return(nil, errors.New("delete failed"))
			err := s3Client.DeleteObjects(t.Context(), "bucket", "key", true)
			require.Error(t, err)
			// Optionally, check the error message if DeleteObjects wraps or returns the specific error
			require.Contains(t, err.Error(), "delete failed")
		})
	})
}
