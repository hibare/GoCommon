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

// S3APIIface is the interface for the S3 service.
type S3APIIface interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	ListObjects(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error)
}

// ClientIface is the interface for the S3 service.
type ClientIface interface {
	BuildKey(prefixes ...string) string
	BuildTimestampedKey(prefixes ...string) string
	TrimPrefix(keys []string, prefix string) []string

	UploadDir(ctx context.Context, bucket, prefix, baseDir string, exclude []*regexp.Regexp) (UploadDirResponse, error)
	UploadFile(ctx context.Context, bucket, prefix, filePath string) (string, error)
	ListObjectsAtPrefix(ctx context.Context, bucket, prefix string) ([]string, error)
	DeleteObjects(ctx context.Context, bucket, key string, recursive bool) error
}

// client is the implementation of the client service.
type client struct {
	Client S3APIIface
}

// BuildKey builds a key from the parts.
func (s *client) BuildKey(parts ...string) string {
	partsSlice := []string{}

	for _, p := range parts {
		if p != "" {
			partsSlice = append(partsSlice, p)
		}
	}

	generatedKey := filepath.Join(partsSlice...)

	if !strings.HasSuffix(generatedKey, S3PrefixSeparator) {
		generatedKey = fmt.Sprintf("%s%s", generatedKey, S3PrefixSeparator)
	}

	return generatedKey
}

// BuildTimestampedKey builds a timestamped key from the parts.
func (s *client) BuildTimestampedKey(parts ...string) string {
	// append the timestamp to the parts and reuse BuildKey
	parts = append(parts, time.Now().Format(constants.DefaultDateTimeLayout))
	return s.BuildKey(parts...)
}

// TrimPrefix trims the prefix from the keys.
func (s *client) TrimPrefix(keys []string, prefix string) []string {
	trimmedKeys := make([]string, 0, len(keys))
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
func (s *client) UploadDir(ctx context.Context, bucket, prefix, baseDir string, exclude []*regexp.Regexp) (UploadDirResponse, error) {
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
func (s *client) UploadFile(ctx context.Context, bucket, prefix, filePath string) (string, error) {
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

// ListObjectsAtPrefix lists the objects at the prefix root.
func (s *client) ListObjectsAtPrefix(ctx context.Context, bucket, prefix string) ([]string, error) {
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
func (s *client) DeleteObjects(ctx context.Context, bucket, key string, recursive bool) error {
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
}

func newClient(ctx context.Context, opts Options) (ClientIface, error) {
	// Build config options slice based on provided input
	var cfgOptions []func(*s3.Options)

	if opts.Region != "" {
		cfgOptions = append(cfgOptions, func(o *s3.Options) {
			o.Region = opts.Region
		})
	}
	if opts.AccessKey != "" && opts.SecretKey != "" {
		cfgOptions = append(cfgOptions, func(o *s3.Options) {
			o.Credentials = credentials.NewStaticCredentialsProvider(opts.AccessKey, opts.SecretKey, "")
		})
	}

	if opts.Endpoint != "" {
		cfgOptions = append(cfgOptions, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(opts.Endpoint)
		})
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg, cfgOptions...)

	return &client{
		Client: s3Client,
	}, nil
}

// NewClient returns a new instance of S3 with the provided configuration (for production use).
var NewClient = newClient
