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

// IndexNodes implements algolia.AlgoliaService.
func (m *MockAlgoliaService) IndexNodes(ctx context.Context, n ...*ent.Node) error {
	args := m.Called(ctx, n)
	return args.Error(0)
}

// DeleteNode implements algolia.AlgoliaService.
func (m *MockAlgoliaService) DeleteNode(ctx context.Context, n *ent.Node) error {
	args := m.Called(ctx, n)
	return args.Error(0)
}

// SearchNodes implements algolia.AlgoliaService.
func (m *MockAlgoliaService) SearchNodes(ctx context.Context, query string, opts ...interface{}) (nodes []*ent.Node, err error) {
	args := m.Called(ctx, query, opts)
	return args.Get(0).([]*ent.Node), args.Error(1)
}

// IndexNodeVersions implements algolia.AlgoliaService.
func (m *MockAlgoliaService) IndexNodeVersions(ctx context.Context, nodes ...*ent.NodeVersion) error {
	args := m.Called(ctx, nodes)
	return args.Error(0)
}

// DeleteNodeVersion implements algolia.AlgoliaService.
func (m *MockAlgoliaService) DeleteNodeVersions(ctx context.Context, node ...*ent.NodeVersion) error {
	args := m.Called(ctx, node)
	return args.Error(0)
}

// SearchNodeVersions implements algolia.AlgoliaService.
func (m *MockAlgoliaService) SearchNodeVersions(ctx context.Context, query string, opts ...interface{}) ([]*ent.NodeVersion, error) {
	args := m.Called(ctx, query, opts)
	return args.Get(0).([]*ent.NodeVersion), args.Error(1)
}
