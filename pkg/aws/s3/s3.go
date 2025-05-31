// Package s3 provides utilities for working with AWS S3.
package s3

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hibare/GoCommon/v2/pkg/constants"
	commonFiles "github.com/hibare/GoCommon/v2/pkg/file"
)

// ServiceAPI is the interface for the S3 service.
type ServiceAPI interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	ListObjects(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error)
}

// Client is the interface for the S3 service.
type Client interface {
	GetPrefix(prefixes ...string) string
	GetTimestampedPrefix(prefixes ...string) string
	TrimPrefix(keys []string, prefix string) []string
	UploadDir(ctx context.Context, bucket, prefix, baseDir string, exclude []*regexp.Regexp) (UploadDirResponse, error)
	UploadFile(ctx context.Context, bucket, prefix, filePath string) (string, error)
	ListObjectsAtPrefixRoot(ctx context.Context, bucket, prefix string) ([]string, error)
}

// S3 is the implementation of the S3 service.
type S3 struct {
	Client ServiceAPI
}

// getPrefix sets the prefix for the S3 service.
func (s *S3) getPrefix(timestamped bool, parts ...string) string {
	partsSlice := []string{}

	for _, p := range parts {
		if p != "" {
			partsSlice = append(partsSlice, p)
		}
	}

	generatedPrefix := filepath.Join(partsSlice...)

	if timestamped {
		timePrefix := time.Now().Format(constants.DefaultDateTimeLayout)
		generatedPrefix = filepath.Join(generatedPrefix, timePrefix)
	}

	if !strings.HasSuffix(generatedPrefix, constants.S3PrefixSeparator) {
		generatedPrefix = fmt.Sprintf("%s%s", generatedPrefix, constants.S3PrefixSeparator)
	}

	return generatedPrefix
}

// GetTimestampedPrefix sets the timestamped prefix for the S3 service.
func (s *S3) GetTimestampedPrefix(parts ...string) string {
	return s.getPrefix(true, parts...)
}

// GetPrefix sets the prefix for the S3 service.
func (s *S3) GetPrefix(parts ...string) string {
	return s.getPrefix(false, parts...)
}

// TrimPrefix trims the prefix from the keys.
func (s *S3) TrimPrefix(keys []string, prefix string) []string {
	var trimmedKeys []string
	for _, key := range keys {
		trimmedKey := strings.TrimPrefix(key, prefix)
		trimmedKey = strings.TrimSuffix(trimmedKey, "/")
		trimmedKeys = append(trimmedKeys, trimmedKey)
	}
	return trimmedKeys
}

// UploadDirResponse holds the result of an UploadDir operation.
type UploadDirResponse struct {
	BaseKey      string
	TotalFiles   int
	TotalDirs    int
	SuccessFiles int
	FailedFiles  map[string]error
}

// UploadDir uploads a directory to the S3 service.
func (s *S3) UploadDir(ctx context.Context, bucket, prefix, baseDir string, exclude []*regexp.Regexp) (UploadDirResponse, error) {
	resp := UploadDirResponse{
		FailedFiles: make(map[string]error), // Initialize the map
	}

	baseDirParentPath := filepath.Dir(baseDir)
	files, dirs := commonFiles.ListFilesDirs(baseDir, exclude)

	resp.TotalFiles = len(files)
	resp.TotalDirs = len(dirs)

	for _, file := range files {
		fp, err := os.Open(file)
		if err != nil {
			resp.FailedFiles[file] = err
			continue
		}
		defer func() {
			_ = fp.Close()
		}()

		key := filepath.Join(prefix, strings.TrimPrefix(file, baseDirParentPath))
		_, err = s.Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: &bucket,
			Key:    &key,
			Body:   fp,
		})
		if err != nil {
			resp.FailedFiles[file] = err
			continue
		}
		resp.SuccessFiles++
	}

	if resp.SuccessFiles > 0 {
		resp.BaseKey = filepath.Join(prefix, filepath.Base(baseDir))
	}

	return resp, nil
}

// UploadFile uploads a file to the S3 service.
func (s *S3) UploadFile(ctx context.Context, bucket, prefix, filePath string) (string, error) {
	fp, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = fp.Close()
	}()

	key := filepath.Join(prefix, filepath.Base(filePath))
	_, err = s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   fp,
	})
	if err != nil {
		return "", err
	}

	return key, nil
}

// ListObjectsAtPrefixRoot lists the objects at the prefix root.
func (s *S3) ListObjectsAtPrefixRoot(ctx context.Context, bucket, prefix string) ([]string, error) {
	var keys []string
	input := &s3.ListObjectsV2Input{
		Bucket:    &bucket,
		Prefix:    &prefix,
		Delimiter: aws.String("/"),
	}

	resp, err := s.Client.ListObjectsV2(ctx, input)
	if err != nil {
		return keys, err
	}

	for _, obj := range resp.Contents {
		if *obj.Key == prefix {
			continue
		}
		keys = append(keys, *obj.Key)
	}

	if len(keys) == 0 && len(resp.CommonPrefixes) == 0 {
		return keys, nil
	}

	for _, cp := range resp.CommonPrefixes {
		keys = append(keys, *cp.Prefix)
	}

	return keys, nil
}

// DeleteObjects deletes the objects from the S3 service.
func (s *S3) DeleteObjects(ctx context.Context, bucket, key string, recursive bool) error {
	// Delete all child object recursively
	if recursive {
		resp, err := s.Client.ListObjects(ctx, &s3.ListObjectsInput{
			Bucket: &bucket,
			Prefix: &key,
		})
		if err != nil {
			return err
		}

		for _, obj := range resp.Contents {
			_, err = s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: &bucket,
				Key:    obj.Key,
			})
			if err != nil {
				return err
			}
		}
	}

	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return err
	}

	return nil
}

// Options is the options for the S3 service.
type Options struct {
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
	Bucket    string
	Prefix    string
}

// NewS3WithDeps returns a new S3 instance with injected dependencies (for testing/mocking).
func NewS3WithDeps(client ServiceAPI) Client {
	return &S3{
		Client: client,
	}
}

// NewS3 returns a new instance of S3 with the provided configuration (for production use).
func NewS3(ctx context.Context, opts Options) (Client, error) {
	// Build config options slice based on provided input
	var cfgOptions []func(*config.LoadOptions) error

	if opts.Region != "" {
		cfgOptions = append(cfgOptions, config.WithRegion(opts.Region))
	}
	if opts.AccessKey != "" && opts.SecretKey != "" {
		cfgOptions = append(cfgOptions, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(opts.AccessKey, opts.SecretKey, "")))
	}
	if opts.Endpoint != "" {
		cfgOptions = append(cfgOptions, config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(_, _ string, _ ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:               opts.Endpoint,
					HostnameImmutable: true,
				}, nil
			}),
		))
	}

	cfg, err := config.LoadDefaultConfig(ctx, cfgOptions...)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)

	return &S3{
		Client: s3Client,
	}, nil
}
