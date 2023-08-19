package s3

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/hibare/GoCommon/pkg/constants"
	commonFiles "github.com/hibare/GoCommon/pkg/file"
)

type S3 struct {
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
	Bucket    string
	Prefix    string
	Sess      *session.Session
}

func (s *S3) SetPrefix(prefix, hostname string, timestamped bool) {
	prefixSlice := []string{}

	if prefix != "" {
		prefixSlice = append(prefixSlice, prefix)
	}

	if hostname != "" {
		prefixSlice = append(prefixSlice, hostname)
	}

	generatedPrefix := filepath.Join(prefixSlice...)

	if timestamped {
		timePrefix := time.Now().Format(constants.DefaultDateTimeLayout)
		generatedPrefix = filepath.Join(generatedPrefix, timePrefix)
	}

	if !strings.HasSuffix(generatedPrefix, constants.S3PrefixSeparator) {
		generatedPrefix = fmt.Sprintf("%s%s", generatedPrefix, constants.S3PrefixSeparator)
	}

	s.Prefix = generatedPrefix
}

func (s *S3) TrimPrefix(keys []string) []string {
	var trimmedKeys []string
	for _, key := range keys {
		trimmedKey := strings.TrimPrefix(key, s.Prefix)
		trimmedKey = strings.TrimSuffix(trimmedKey, "/")
		trimmedKeys = append(trimmedKeys, trimmedKey)
	}
	return trimmedKeys
}

func (s *S3) NewSession() error {
	sess, err := session.NewSession(&aws.Config{
		Region:           &s.Region,
		Endpoint:         &s.Endpoint,
		Credentials:      credentials.NewStaticCredentials(s.AccessKey, s.SecretKey, ""),
		S3ForcePathStyle: aws.Bool(true),
	})

	if err != nil {
		return err
	}

	s.Sess = sess

	return nil
}

func (s *S3) UploadDir(baseDir string) (string, int, int, int) {
	totalFiles, totalDirs, successFiles := 0, 0, 0
	baseKey := ""

	client := s3.New(s.Sess)
	baseDirParentPath := filepath.Dir(baseDir)

	files, dirs := commonFiles.ListFilesDirs(baseDir, nil)

	totalFiles = len(files)
	totalDirs = len(dirs)

	for _, file := range files {
		fp, err := os.Open(file)
		if err != nil {
			//ToDo: Add logs
			continue
		}
		defer fp.Close()

		key := filepath.Join(s.Prefix, strings.TrimPrefix(file, baseDirParentPath))
		_, err = client.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(key),
			Body:   fp,
		})
		if err != nil {
			//ToDo: Add logs
			continue
		}
		successFiles += 1
		//ToDo: Add logs
	}

	if successFiles > 0 {
		baseKey = filepath.Join(s.Prefix, filepath.Base(baseDir))
	}

	return baseKey, totalFiles, totalDirs, successFiles
}

func (s *S3) UploadFile(filePath string) (string, error) {
	uploader := s3manager.NewUploader(s.Sess)

	f, err := os.Open(filePath)
	if err != nil {
		//ToDo: Add logs
		return "", err
	}
	defer f.Close()

	// Upload the file to S3
	key := filepath.Join(s.Prefix, filepath.Base(filePath))
	if _, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   f,
	}); err != nil {
		//ToDo: Add logs
		return "", err
	}

	return key, nil
}

func (s *S3) ListObjectsAtPrefixRoot() ([]string, error) {
	client := s3.New(s.Sess)

	var keys []string
	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(s.Bucket),
		Prefix:    aws.String(s.Prefix),
		Delimiter: aws.String("/"),
	}

	resp, err := client.ListObjectsV2(input)
	if err != nil {
		return keys, err
	}

	for _, obj := range resp.Contents {
		if *obj.Key == s.Prefix {
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

func (s *S3) DeleteObjects(key string, recursive bool) error {
	client := s3.New(s.Sess)

	// Delete all child object recursively
	if recursive {
		// List all objects in the bucket with the given key
		resp, err := client.ListObjects(&s3.ListObjectsInput{
			Bucket: aws.String(s.Bucket),
			Prefix: aws.String(key),
		})
		if err != nil {
			return err
		}

		// Delete all objects with the given key
		for _, obj := range resp.Contents {
			if _, err = client.DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String(s.Bucket),
				Key:    obj.Key,
			}); err != nil {
				return err
			}
		}
	}

	// Delete the prefix
	if _, err := client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}); err != nil {
		return err
	}

	return nil
}
