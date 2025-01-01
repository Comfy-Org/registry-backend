package gateways

import (
	"context"
	"registry-backend/gateways/pubsub"

	"github.com/stretchr/testify/mock"
)

var _ pubsub.PubSubService = &MockPubSubService{}

type MockPubSubService struct {
	mock.Mock
}

// PublishNodePack implements pubsub.PubSubService.
func (m *MockPubSubService) PublishNodePack(ctx context.Context, storageURL string) error {
	args := m.Called(ctx, storageURL)
	return args.Error(0)
}
