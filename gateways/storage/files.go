package storage

import (
	"bytes"
	"context"
	"fmt"

	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/rs/zerolog/log"
)

const BucketName = "comfy-workflow-json"

type StorageService interface {
	UploadFile(ctx context.Context, bucket, object, filePath string) (string, error)
	StreamFileUpload(w io.Writer, objectName, blob string) (string, string, error)
	GetFileUrl(ctx context.Context, bucketName, objectPath string) (string, error)
	GenerateSignedURL(bucketName, objectName string) (string, error)
}

type GCPStorageService struct {
	client *storage.Client
}

func NewGCPStorageService(ctx context.Context) (*GCPStorageService, error) {
	StorageClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewStorageClient: %v", err)
	}

	return &GCPStorageService{
		client: StorageClient,
	}, nil
}

// uploadFile uploads an object.
func (s *GCPStorageService) UploadFile(ctx context.Context, bucket, object string, filePath string) (string, error) {
	log.Ctx(ctx).Info().Msgf("Uploading %v to %v/%v.\n", filePath, bucket, object)
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	// Open local file.
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("os.Open: %w", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := client.Bucket(bucket).Object(object)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	// For an object that does not yet exist, set the DoesNotExist precondition.
	o = o.If(storage.Conditions{DoesNotExist: true})
	// If the live object already exists in your bucket, set instead a
	// generation-match precondition using the live object's generation number.
	// attrs, err := o.Attrs(ctx)
	// if err != nil {
	//      return fmt.Errorf("object.Attrs: %w", err)
	// }
	// o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	// Upload an object with storage.Writer.
	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return "", fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %w", err)
	}
	log.Ctx(ctx).Info().Msgf("Blob %v uploaded.\n", object)
	// Make the file publicly accessible
	if err := o.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("ACL().Set: %w", err)
	}

	// Construct the public URL
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, object)
	log.Ctx(ctx).Info().Msgf("Blob is publicly accessible at %v.\n", publicURL)
	return publicURL, nil
}

// StreamFileUpload uploads an object via a stream.
func (s *GCPStorageService) StreamFileUpload(w io.Writer, objectName string, blob string) (string, string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", "", fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	b := []byte(blob)
	buf := bytes.NewBuffer(b)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(BucketName).Object(objectName).NewWriter(ctx)
	wc.ChunkSize = 0 // note retries are not supported for chunk size 0.

	if _, err = io.Copy(wc, buf); err != nil {
		return "", "", fmt.Errorf("io.Copy: %w", err)
	}
	// Data can continue to be added to the file until the writer is closed.
	if err := wc.Close(); err != nil {
		return "", "", fmt.Errorf("Writer.Close: %w", err)
	}
	log.Ctx(ctx).Info().Msgf("%v uploaded to %v.\n", objectName, BucketName)

	return BucketName, objectName, nil
}

func (s *GCPStorageService) GetFileUrl(ctx context.Context, bucketName, objectPath string) (string, error) {
	// Get public url of a file in a bucket
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	// Get Public URL
	attrs, err := client.Bucket(bucketName).Object(objectPath).Attrs(ctx)
	if err != nil {
		return "", fmt.Errorf("object.Attrs: %w", err)
	}
	publicURL := attrs.MediaLink
	log.Ctx(ctx).Info().Msgf("Public URL: %v", publicURL)
	return publicURL, nil
}

func (s *GCPStorageService) GenerateSignedURL(bucketName, objectName string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	expires := time.Now().Add(15 * time.Minute)
	url, err := client.Bucket(bucketName).SignedURL(objectName, &storage.SignedURLOptions{
		ContentType: "application/zip",
		Method:      "PUT",
		Expires:     expires,
	})
	if err != nil {
		return "", err
	}

	return url, nil
}
