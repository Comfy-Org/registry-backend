package algolia

import (
	"context"
	"fmt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/rs/zerolog/log"
	"registry-backend/config" // assuming a config package exists to hold config values
	"registry-backend/ent"
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

// Ensure algolia struct implements AlgoliaService interface
var _ AlgoliaService = (*algolia)(nil)

// algolia struct holds the Algolia client.
type algolia struct {
	client *search.Client
}

// NewAlgoliaService creates a new Algolia service using the provided config or returns a noop implementation if the config is missing.
func NewAlgoliaService(cfg *config.Config) (AlgoliaService, error) {
	if cfg == nil || cfg.AlgoliaAppID == "" || cfg.AlgoliaAPIKey == "" {
		// Return a noop implementation if config is nil or missing keys
		log.Info().Msg("No Algolia configuration found, using noop implementation")
		return &algolianoop{}, nil
	}

	// Fetch the Algolia app ID and API key from the provided config
	appID := cfg.AlgoliaAppID
	apiKey := cfg.AlgoliaAPIKey

	// Initialize the Algolia client
	client := search.NewClient(appID, apiKey)
	return &algolia{client: client}, nil
}

// IndexNodes indexes the provided nodes in Algolia.
func (a *algolia) IndexNodes(ctx context.Context, nodes ...*ent.Node) error {
	index := a.client.InitIndex("nodes_index")
	objects := make([]map[string]interface{}, len(nodes))

	for i, n := range nodes {
		o := map[string]interface{}{
			"objectID":       n.ID,
			"name":           n.Name,
			"publisher_id":   n.PublisherID,
			"description":    n.Description,
			"id":             n.ID,
			"create_time":    n.CreateTime,
			"update_time":    n.UpdateTime,
			"license":        n.License,
			"repository_url": n.RepositoryURL,
			"total_install":  n.TotalInstall,
			"status":         n.Status,
			"author":         n.Author,
			"category":       n.Category,
			"total_star":     n.TotalStar,
			"total_review":   n.TotalReview,
		}
		objects[i] = o
	}

	res, err := index.SaveObjects(objects)
	if err != nil {
		return fmt.Errorf("failed to index nodes: %w", err)
	}

	return res.Wait()
}

// SearchNodes searches for nodes in Algolia matching the query.
func (a *algolia) SearchNodes(ctx context.Context, query string, opts ...interface{}) ([]*ent.Node, error) {
	index := a.client.InitIndex("nodes_index")
	res, err := index.Search(query, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to search nodes: %w", err)
	}

	var nodes []*ent.Node
	if err := res.UnmarshalHits(&nodes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %w", err)
	}
	return nodes, nil
}

// DeleteNode deletes the specified node from Algolia.
func (a *algolia) DeleteNode(ctx context.Context, node *ent.Node) error {
	index := a.client.InitIndex("nodes_index")
	res, err := index.DeleteObject(node.ID)
	if err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}
	return res.Wait()
}

// IndexNodeVersions implements AlgoliaService.
func (a *algolia) IndexNodeVersions(ctx context.Context, nodes ...*ent.NodeVersion) error {
	index := a.client.InitIndex("node_versions_index")
	objects := make([]struct {
		ObjectID string `json:"objectID"`
		*ent.NodeVersion
	}, len(nodes))

	for i, n := range nodes {
		o := struct {
			ObjectID string `json:"objectID"`
			*ent.NodeVersion
		}{
			ObjectID:    n.ID.String(),
			NodeVersion: n,
		}
		o.StatusReason = ""
		objects[i] = o
	}

	res, err := index.SaveObjects(objects)
	if err != nil {
		return fmt.Errorf("failed to index node versions: %w", err)
	}

	return res.Wait()
}

// DeleteNodeVersions implements AlgoliaService.
func (a *algolia) DeleteNodeVersions(ctx context.Context, nodes ...*ent.NodeVersion) error {
	index := a.client.InitIndex("node_versions_index")
	ids := []string{}
	for _, node := range nodes {
		ids = append(ids, node.ID.String())
	}
	res, err := index.DeleteObjects(ids)
	if err != nil {
		return fmt.Errorf("failed to delete node versions: %w", err)
	}
	return res.Wait()
}

// SearchNodeVersions implements AlgoliaService.
func (a *algolia) SearchNodeVersions(ctx context.Context, query string, opts ...interface{}) ([]*ent.NodeVersion, error) {
	index := a.client.InitIndex("node_versions_index")
	res, err := index.Search(query, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to search node versions: %w", err)
	}

	var nodes []*ent.NodeVersion
	if err := res.UnmarshalHits(&nodes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %w", err)
	}
	return nodes, nil
}
