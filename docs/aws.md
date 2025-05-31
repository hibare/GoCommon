# AWS Package Documentation

## Overview

The `aws` package in GoCommon provides utilities for working with AWS services, with a primary focus on Amazon S3. It offers interfaces and implementations for interacting with S3 buckets, uploading files and directories, listing objects, and deleting objects. The package is designed to be testable and mockable, supporting dependency injection for easier testing.

---

## Subpackages

- **s3**: Provides utilities for working with AWS S3, including file and directory uploads, object listing, and deletion.

---

## Key Types and Interfaces

### S3 Service

- **ServiceAPI**: Interface for the S3 service, matching the AWS SDK's S3 client methods (`PutObject`, `ListObjectsV2`, `DeleteObject`, `ListObjects`).
- **Client**: Interface for high-level S3 operations, such as uploading files/directories and listing objects.
- **S3**: Implementation of the `Client` interface, wrapping an AWS S3 client.

### Options

- **Options**: Struct for configuring the S3 client (endpoint, region, access key, secret key, bucket, prefix).

---

## Main Functions

- **NewS3WithDeps(client ServiceAPI) Client**: Returns a new S3 client with injected dependencies (for testing/mocking).
- **NewS3(ctx, opts) (Client, error)**: Returns a new S3 client for production use, using the provided configuration.
- **UploadDir(ctx, bucket, prefix, baseDir, exclude) (UploadDirResponse, error)**: Uploads a directory to S3, optionally excluding files by regex.
- **UploadFile(ctx, bucket, prefix, filePath) (string, error)**: Uploads a single file to S3.
- **ListObjectsAtPrefixRoot(ctx, bucket, prefix) ([]string, error)**: Lists objects at the root of a given prefix.
- **DeleteObjects(ctx, bucket, key, recursive) error**: Deletes an object or all objects under a prefix (if recursive).

---

## Testing and Mocking

- The package provides a `MockS3Client` for use in tests, allowing you to mock S3 operations.

---

## Example Usage

```go
import (
    "context"
    "github.com/hibare/GoCommon/v2/pkg/aws/s3"
)

func main() {
    opts := s3.Options{
        Endpoint:  "https://s3.amazonaws.com",
        Region:    "us-east-1",
        AccessKey: "your-access-key",
        SecretKey: "your-secret-key",
        Bucket:    "your-bucket",
    }
    client, err := s3.NewS3(context.Background(), opts)
    if err != nil {
        panic(err)
    }
    // Upload a file
    key, err := client.UploadFile(context.Background(), opts.Bucket, "prefix/", "/path/to/file.txt")
    if err != nil {
        panic(err)
    }
    fmt.Println("Uploaded file to:", key)
}
```

---

## Notes

- The S3 client supports both static credentials and custom endpoints (for S3-compatible services).
- Prefix management utilities are provided for organizing S3 objects.
- The package is designed for extensibility and testability.
