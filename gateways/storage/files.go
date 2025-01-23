package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"registry-backend/config"
	"registry-backend/tracing"

	"cloud.google.com/go/storage"
	"github.com/rs/zerolog/log"
)

// StorageService defines the interface for interacting with cloud storage.
type StorageService interface {
	UploadFile(ctx context.Context, bucket, object, filePath string) (string, error)
	StreamFileUpload(w io.Writer, objectName, blob string) (string, string, error)
	GetFileUrl(ctx context.Context, bucketName, objectPath string) (string, error)
	GenerateSignedURL(bucketName, objectName string) (string, error)
}

// Ensure storageService struct implements StorageService interface
var _ StorageService = (*storageService)(nil)

// storageService struct holds the GCP storage client and configuration.
type storageService struct {
	client *storage.Client
	config *config.Config
}

// NewStorageService creates a new storage service using the provided config or returns a noop implementation if the config is missing.
func NewStorageService(cfg *config.Config) (StorageService, error) {
	if cfg == nil {
		// Return a noop implementation if config is nil or storage is not enabled
		log.Info().Msg("No storage configuration found, using noop implementation")
		return &storageNoop{}, nil
	}

	// Initialize GCP storage client
	client, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, fmt.Errorf("NewStorageClient: %v", err)
	}

	return &storageService{
		client: client,
		config: cfg,
	}, nil
}

// UploadFile uploads an object to GCP storage.
func (s *storageService) UploadFile(ctx context.Context, bucket, object, filePath string) (string, error) {
	defer tracing.TraceDefaultSegment(ctx, "StorageService.UploadFile")()

	log.Ctx(ctx).Info().Msgf("Uploading %v to %v/%v.\n", filePath, bucket, object)

	// Open local file
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("os.Open: %w", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := s.client.Bucket(bucket).Object(object)

	// Ensure that we don't overwrite an existing object.
	o = o.If(storage.Conditions{DoesNotExist: true})

	// Upload the file to the cloud
	wc := o.NewWriter(ctx)
	if _, err := io.Copy(wc, f); err != nil {
		return "", fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %w", err)
	}

	log.Ctx(ctx).Info().Msgf("Blob %v uploaded.\n", object)

	// Set the ACL to make the file publicly accessible
	if err := o.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("ACL().Set: %w", err)
	}

	// Construct the public URL
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, object)
	log.Ctx(ctx).Info().Msgf("Blob is publicly accessible at %v.\n", publicURL)
	return publicURL, nil
}

// StreamFileUpload uploads an object via a stream to GCP storage.
func (s *storageService) StreamFileUpload(w io.Writer, objectName, blob string) (string, string, error) {
	ctx := context.Background()

	b := []byte(blob)
	buf := bytes.NewBuffer(b)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload the object as a stream
	wc := s.client.Bucket(s.config.CloudStorageBucketName).Object(objectName).NewWriter(ctx)
	wc.ChunkSize = 0 // Note: retries are not supported for chunk size 0

	if _, err := io.Copy(wc, buf); err != nil {
		return "", "", fmt.Errorf("io.Copy: %w", err)
	}

	if err := wc.Close(); err != nil {
		return "", "", fmt.Errorf("Writer.Close: %w", err)
	}

	log.Ctx(ctx).Info().Msgf("%v uploaded to %v.\n", objectName, s.config.CloudStorageBucketName)

	return s.config.CloudStorageBucketName, objectName, nil
}

// GetFileUrl gets the public URL of a file from GCP storage.
func (s *storageService) GetFileUrl(ctx context.Context, bucketName, objectPath string) (string, error) {
	defer tracing.TraceDefaultSegment(ctx, "StorageService.GetFileUrl")()

	// Get the public URL of a file in a bucket
	attrs, err := s.client.Bucket(bucketName).Object(objectPath).Attrs(ctx)
	if err != nil {
		return "", fmt.Errorf("object.Attrs: %w", err)
	}

	publicURL := attrs.MediaLink
	log.Ctx(ctx).Info().Msgf("Public URL: %v", publicURL)
	return publicURL, nil
}

// GenerateSignedURL generates a signed URL for uploading to GCP storage.
func (s *storageService) GenerateSignedURL(bucketName, objectName string) (string, error) {
	expires := time.Now().Add(15 * time.Minute)
	url, err := s.client.Bucket(bucketName).SignedURL(objectName, &storage.SignedURLOptions{
		ContentType: "application/zip",
		Method:      "PUT",
		Expires:     expires,
	})
	if err != nil {
		return "", err
	}

	return url, nil
}

// storageNoop is a noop implementation of the StorageService interface.
type storageNoop struct{}

// Implement all StorageService methods for noop behavior.

func (s *storageNoop) UploadFile(ctx context.Context, bucket, object, filePath string) (string, error) {
	// No-op, return nil to avoid side-effects
	return "", nil
}

func (s *storageNoop) StreamFileUpload(w io.Writer, objectName, blob string) (string, string, error) {
	// No-op, return nil to avoid side-effects
	return "", "", nil
}

func (s *storageNoop) GetFileUrl(ctx context.Context, bucketName, objectPath string) (string, error) {
	// No-op, return empty string and nil error
	return "", nil
}

func (s *storageNoop) GenerateSignedURL(bucketName, objectName string) (string, error) {
	// No-op, return empty string and nil error
	return "", nil
}
