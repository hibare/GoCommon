package s3

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hibare/GoCommon/pkg/testhelper"
	"github.com/stretchr/testify/assert"
)

var (
	endpoint  = "http://localhost:9000"
	region    = "us-east-1"
	accessKey = "admin"
	secretKey = "5ee4392a-cb32-4f9d-8c19-d91e19e30834"
	bucket    = "test-bucket"
	prefix    = "test-prefix"
)

func TestSetPrefix(t *testing.T) {
	s := S3{}
	s.SetPrefix("test-prefix", "test-hostname", false)
	assert.NotEmpty(t, s.Prefix)
	assert.Equal(t, s.Prefix, "test-prefix/test-hostname/")
}

func TestSetPrefixTimeStamped(t *testing.T) {
	s := S3{}
	s.SetPrefix("test-prefix", "test-hostname", true)
	assert.NotEmpty(t, s.Prefix)
	assert.Contains(t, s.Prefix, "test-prefix/test-hostname/")
}

func TestTrimPrefix(t *testing.T) {
	keysToTrim := []string{
		"test-prefix/test-hostname/1",
		"test-prefix/test-hostname/2",
		"test-prefix/test-hostname/3",
	}
	s := S3{}
	s.SetPrefix("test-prefix", "test-hostname", false)
	trimmedKeys := s.TrimPrefix(keysToTrim)
	assert.Equal(t, []string{"1", "2", "3"}, trimmedKeys)
}

func TestNewSession(t *testing.T) {
	s := S3{
		Endpoint:  endpoint,
		Region:    region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
		Prefix:    prefix,
	}
	err := s.NewSession()
	assert.NoError(t, err)
}

func TestUploadFile(t *testing.T) {
	s := S3{
		Endpoint:  endpoint,
		Region:    region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
		Prefix:    prefix,
	}
	err := s.NewSession()
	assert.NoError(t, err)

	_, filepath, err := testhelper.CreateTestFile("")
	assert.NoError(t, err)
	defer os.Remove(filepath)

	key, err := s.UploadFile(filepath)
	assert.NoError(t, err)
	assert.NotEmpty(t, key)
}

func TestUploadNoBucket(t *testing.T) {
	s := S3{
		Endpoint:  endpoint,
		Region:    region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    "non-existing-bucket",
		Prefix:    prefix,
	}
	err := s.NewSession()
	assert.NoError(t, err)

	_, filepath, err := testhelper.CreateTestFile("")
	assert.NoError(t, err)
	defer os.Remove(filepath)

	key, err := s.UploadFile(filepath)
	assert.Error(t, err)
	if aerr, ok := err.(awserr.Error); ok {
		assert.Equal(t, s3.ErrCodeNoSuchBucket, aerr.Code())
	}
	assert.Empty(t, key)
}

func TestUploadNoFile(t *testing.T) {
	s := S3{
		Endpoint:  endpoint,
		Region:    region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
		Prefix:    prefix,
	}
	err := s.NewSession()
	assert.NoError(t, err)

	key, err := s.UploadFile("/tmp/no-such-file.txt")
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
	assert.Empty(t, key)
}

func TestUploadDir(t *testing.T) {
	s := S3{
		Endpoint:  endpoint,
		Region:    region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
		Prefix:    prefix,
	}
	err := s.NewSession()
	assert.NoError(t, err)

	testDir, err := testhelper.CreateTestDir()
	assert.NoError(t, err)
	defer os.RemoveAll(testDir)

	key, totalFiles, totalDirs, successFiles := s.UploadDir(testDir)
	assert.Greater(t, totalFiles, 0)
	assert.Greater(t, totalDirs, 0)
	assert.Greater(t, successFiles, 0)
	assert.NotEmpty(t, key)
}

func TestUploadDirNoDir(t *testing.T) {
	s := S3{
		Endpoint:  endpoint,
		Region:    region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
		Prefix:    prefix,
	}
	err := s.NewSession()
	assert.NoError(t, err)

	key, totalFiles, totalDirs, successFiles := s.UploadDir("does-not-exists")
	assert.Equal(t, 0, totalFiles)
	assert.Equal(t, 0, totalDirs)
	assert.Equal(t, 0, successFiles)
	assert.Empty(t, key)
}

func TestListObjectsAtPrefixRoot(t *testing.T) {
	s := S3{
		Endpoint:  endpoint,
		Region:    region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
		Prefix:    "test-prefix/lTQIZP",
	}
	err := s.NewSession()
	assert.NoError(t, err)

	_, filepath, err := testhelper.CreateTestFile("")
	assert.NoError(t, err)
	defer os.Remove(filepath)

	key, err := s.UploadFile(filepath)
	assert.NoError(t, err)
	assert.NotEmpty(t, key)

	// List
	keys, err := s.ListObjectsAtPrefixRoot()
	assert.NoError(t, err)
	assert.Len(t, keys, 1)
}

func TestDeleteObjects(t *testing.T) {
	s := S3{
		Endpoint:  endpoint,
		Region:    region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
		Prefix:    "test-prefix/lTQIZP",
	}
	err := s.NewSession()
	assert.NoError(t, err)

	err = s.DeleteObjects(s.Prefix, true)
	assert.NoError(t, err)
}
