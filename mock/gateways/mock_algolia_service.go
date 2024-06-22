package gateways

import (
	"context"
	"registry-backend/ent"
	"registry-backend/gateways/algolia"

	"github.com/stretchr/testify/mock"
)

var _ algolia.AlgoliaService = &MockAlgoliaService{}

type MockAlgoliaService struct {
	mock.Mock
}

// DeleteNode implements algolia.AlgoliaService.
func (m *MockAlgoliaService) DeleteNode(ctx context.Context, n *ent.Node) error {
	args := m.Called(ctx, n)
	return args.Error(0)
}

// IndexNodes implements algolia.AlgoliaService.
func (m *MockAlgoliaService) IndexNodes(ctx context.Context, n ...*ent.Node) error {
	args := m.Called(ctx, n)
	return args.Error(0)
}

// SearchNodes implements algolia.AlgoliaService.
func (m *MockAlgoliaService) SearchNodes(ctx context.Context, query string, opts ...interface{}) (nodes []*ent.Node, err error) {
	args := m.Called(ctx, query, opts)
	return args.Get(0).([]*ent.Node), args.Error(1)
}
