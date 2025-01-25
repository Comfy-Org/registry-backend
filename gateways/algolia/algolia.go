package algolia

import (
	"context"
	"fmt"
	"registry-backend/config" // assuming a config package exists to hold config values
	"registry-backend/ent"
	"registry-backend/entity"
	"registry-backend/mapper"
	"registry-backend/tracing"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/rs/zerolog/log"
)

// AlgoliaService defines the interface for interacting with Algolia search.
type AlgoliaService interface {
	IndexNodes(ctx context.Context, nodes ...*ent.Node) error
	DeleteNode(ctx context.Context, node *ent.Node) error
	SearchNodes(ctx context.Context, query string, opts ...interface{}) ([]*ent.Node, error)
	IndexNodeVersions(ctx context.Context, nodes ...*ent.NodeVersion) error
	DeleteNodeVersions(ctx context.Context, nodes ...*ent.NodeVersion) error
	SearchNodeVersions(ctx context.Context, query string, opts ...interface{}) ([]*ent.NodeVersion, error)
}

// Ensure the `algolia` struct implements the `AlgoliaService` interface.
var _ AlgoliaService = (*algolia)(nil)

// `algolia` struct holds the Algolia client and an executor for managing tasks.
type algolia struct {
	client   *search.Client
	executor *executor
}

// NewAlgoliaService creates a new Algolia service or returns a noop implementation if configuration is missing.
func NewAlgoliaService(cfg *config.Config) (AlgoliaService, error) {
	if cfg == nil || cfg.AlgoliaAppID == "" || cfg.AlgoliaAPIKey == "" {
		// If configuration is missing, use a noop implementation.
		log.Info().Msg("No Algolia configuration found, using noop implementation")
		return &algolianoop{}, nil
	}

	// Initialize the Algolia client.
	client := search.NewClient(cfg.AlgoliaAppID, cfg.AlgoliaAPIKey)

	// Create an executor for managing background tasks.
	executor := newExecutor(100)
	go executor.start(context.Background())

	return &algolia{client: client, executor: executor}, nil
}

// IndexNodes indexes the provided nodes in the "nodes_index" index on Algolia.
func (a *algolia) IndexNodes(ctx context.Context, nodes ...*ent.Node) error {
	defer tracing.TraceDefaultSegment(ctx, "AlgoliaService.IndexNodes")()

	// Initialize the index and map nodes to Algolia-friendly format.
	index := a.client.InitIndex("nodes_index")
	objects := make([]entity.AlgoliaNode, len(nodes))
	for i, n := range nodes {
		objects[i] = mapper.AlgoliaNodeFromEntNode(n)
	}

	// Schedule the indexing task.
	a.executor.schedule(ctx, "AlgoliaService.IndexNodes", func() error {
		res, err := index.SaveObjects(objects)
		if err != nil {
			return fmt.Errorf("failed to index nodes: %w", err)
		}
		return res.Wait()
	})

	return nil
}

// SearchNodes searches for nodes in the "nodes_index" index using the provided query.
func (a *algolia) SearchNodes(ctx context.Context, query string, opts ...interface{}) (nodes []*ent.Node, err error) {
	defer tracing.TraceDefaultSegment(ctx, "AlgoliaService.SearchNodes")()

	// Perform the search query.
	index := a.client.InitIndex("nodes_index")
	res, err := index.Search(query, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to search nodes: %w", err)
	}

	// Unmarshal the results into AlgoliaNode and convert them to ent.Node.
	var algoliaNodes []entity.AlgoliaNode
	if err := res.UnmarshalHits(&algoliaNodes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %w", err)
	}
	for _, n := range algoliaNodes {
		nodes = append(nodes, n.ToEntNode())
	}
	return
}

// DeleteNode removes a node from the "nodes_index" index by its ID.
func (a *algolia) DeleteNode(ctx context.Context, node *ent.Node) error {
	defer tracing.TraceDefaultSegment(ctx, "AlgoliaService.DeleteNode")()

	index := a.client.InitIndex("nodes_index")
	a.executor.schedule(ctx, "AlgoliaService.DeleteNode", func() error {
		res, err := index.DeleteObject(node.ID)
		if err != nil {
			return fmt.Errorf("failed to delete node: %w", err)
		}
		return res.Wait()
	})

	return nil
}

// IndexNodeVersions indexes the provided node versions in the "node_versions_index" index.
func (a *algolia) IndexNodeVersions(ctx context.Context, nodes ...*ent.NodeVersion) error {
	defer tracing.TraceDefaultSegment(ctx, "AlgoliaService.IndexNodeVersions")()

	// Initialize the index and prepare objects for indexing.
	index := a.client.InitIndex("node_versions_index")
	objects := make([]struct {
		ObjectID string `json:"objectID"`
		*ent.NodeVersion
	}, len(nodes))

	for i, n := range nodes {
		objects[i] = struct {
			ObjectID string `json:"objectID"`
			*ent.NodeVersion
		}{
			ObjectID:    n.ID.String(),
			NodeVersion: n,
		}
		objects[i].Status = "" // Exclude the status field from indexing.
	}

	// Schedule the indexing task.
	a.executor.schedule(ctx, "AlgoliaService.IndexNodeVersions", func() error {
		res, err := index.SaveObjects(objects)
		if err != nil {
			return fmt.Errorf("failed to index node versions: %w", err)
		}
		return res.Wait()
	})

	return nil
}

// DeleteNodeVersions removes node versions from the "node_versions_index" index by their IDs.
func (a *algolia) DeleteNodeVersions(ctx context.Context, nodes ...*ent.NodeVersion) error {
	defer tracing.TraceDefaultSegment(ctx, "AlgoliaService.DeleteNodeVersions")()

	// Extract IDs of the node versions to delete.
	index := a.client.InitIndex("node_versions_index")
	ids := make([]string, len(nodes))
	for i, node := range nodes {
		ids[i] = node.ID.String()
	}

	// Schedule the deletion task.
	a.executor.schedule(ctx, "AlgoliaService.DeleteNodeVersions", func() error {
		res, err := index.DeleteObjects(ids)
		if err != nil {
			return fmt.Errorf("failed to delete node versions: %w", err)
		}
		return res.Wait()
	})

	return nil
}

// SearchNodeVersions searches for node versions in the "node_versions_index" index using the provided query.
func (a *algolia) SearchNodeVersions(ctx context.Context, query string, opts ...interface{}) ([]*ent.NodeVersion, error) {
	defer tracing.TraceDefaultSegment(ctx, "AlgoliaService.SearchNodeVersions")()

	// Perform the search query.
	index := a.client.InitIndex("node_versions_index")
	res, err := index.Search(query, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to search node versions: %w", err)
	}

	// Unmarshal the results into ent.NodeVersion.
	var nodes []*ent.NodeVersion
	if err := res.UnmarshalHits(&nodes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %w", err)
	}
	return nodes, nil
}
