package algolia

import (
	"context"
	"registry-backend/ent"

	"github.com/rs/zerolog/log"
)

var _ AlgoliaService = (*algolianoop)(nil)

// No-op implementation for AlgoliaService, logging calls instead of performing operations
type algolianoop struct{}

func (a *algolianoop) DeleteNode(ctx context.Context, node *ent.Node) error {
	log.Ctx(ctx).Info().Msgf("algolia noop: delete node: %s", node.ID)
	return nil
}

func (a *algolianoop) IndexNodes(ctx context.Context, nodes ...*ent.Node) error {
	log.Ctx(ctx).Info().Msgf("algolia noop: index nodes: %d number of nodes", len(nodes))
	return nil
}

func (a *algolianoop) SearchNodes(ctx context.Context, query string, opts ...interface{}) ([]*ent.Node, error) {
	log.Ctx(ctx).Info().Msgf("algolia noop: search nodes: %s", query)
	return nil, nil
}

func (a *algolianoop) DeleteNodeVersions(ctx context.Context, nodes ...*ent.NodeVersion) error {
	log.Ctx(ctx).Info().Msgf("algolia noop: delete node version: %d number of node versions", len(nodes))
	return nil
}

func (a *algolianoop) IndexNodeVersions(ctx context.Context, nodes ...*ent.NodeVersion) error {
	log.Ctx(ctx).Info().Msgf("algolia noop: index node versions: %d number of node versions", len(nodes))
	return nil
}

func (a *algolianoop) SearchNodeVersions(ctx context.Context, query string, opts ...interface{}) ([]*ent.NodeVersion, error) {
	log.Ctx(ctx).Info().Msgf("algolia noop: search node versions: %s", query)
	return nil, nil
}
