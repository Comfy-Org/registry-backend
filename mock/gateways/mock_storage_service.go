package gateways

import (
	"context"
	"io"

	"github.com/stretchr/testify/mock"
)

// MockStorageService is a mock of StorageService interface
type MockStorageService struct {
	mock.Mock
}

func (m *MockStorageService) UploadFile(ctx context.Context, bucket, object, filePath string) (string, error) {
	args := m.Called(ctx, bucket, object, filePath)
	return args.String(0), args.Error(1)
}

func (m *MockStorageService) StreamFileUpload(w io.Writer, objectName, blob string) (string, string, error) {
	args := m.Called(w, objectName, blob)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockStorageService) GetFileUrl(ctx context.Context, bucketName, objectPath string) (string, error) {
	args := m.Called(ctx, bucketName, objectPath)
	return args.String(0), args.Error(1)
}

func (m *MockStorageService) GenerateSignedURL(bucketName, objectName string) (string, error) {
	args := m.Called(bucketName, objectName)
	return args.String(0), args.Error(1)
}
