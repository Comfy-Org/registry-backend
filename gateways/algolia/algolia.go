package algolia

import (
	"context"
	"fmt"
	"os"
	"registry-backend/ent"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
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

// New creates a new Algolia service with the provided app ID and API key.
func New(appid, apikey string) (AlgoliaService, error) {
	return &algolia{
		client: search.NewClient(appid, apikey),
	}, nil
}

// NewFromEnv creates a new Algolia service using environment variables for configuration.
func NewFromEnv() (AlgoliaService, error) {
	appid, ok := os.LookupEnv("ALGOLIA_APP_ID")
	if !ok {
		return nil, fmt.Errorf("required env variable ALGOLIA_APP_ID is not set")
	}
	apikey, ok := os.LookupEnv("ALGOLIA_API_KEY")
	if !ok {
		return nil, fmt.Errorf("required env variable ALGOLIA_API_KEY is not set")
	}
	return New(appid, apikey)
}

// NewFromEnvOrNoop creates a new Algolia service using environment variables or noop implementation if no environment found
func NewFromEnvOrNoop() (AlgoliaService, error) {
	id := os.Getenv("ALGOLIA_APP_ID")
	key := os.Getenv("ALGOLIA_API_KEY")
	if id == "" && key == "" {
		return &algolianoop{}, nil
	}

	return NewFromEnv()
}

// IndexNodes indexes the provided nodes in Algolia.
func (a *algolia) IndexNodes(ctx context.Context, nodes ...*ent.Node) error {
	index := a.client.InitIndex("nodes_index")
	objects := make([]struct {
		ObjectID string `json:"objectID"`
		*ent.Node
	}, len(nodes))

	for i, n := range nodes {
		o := struct {
			ObjectID string `json:"objectID"`
			*ent.Node
		}{
			ObjectID: n.ID,
			Node:     n,
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
		return fmt.Errorf("failed to index nodes: %w", err)
	}

	return res.Wait()
}

// DeleteNodeVersion implements AlgoliaService.
func (a *algolia) DeleteNodeVersions(ctx context.Context, nodes ...*ent.NodeVersion) error {
	index := a.client.InitIndex("node_versions_index")
	ids := []string{}
	for _, node := range nodes {
		ids = append(ids, node.ID.String())
	}
	res, err := index.DeleteObjects(ids)
	if err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}
	return res.Wait()
}

// SearchNodeVersions implements AlgoliaService.
func (a *algolia) SearchNodeVersions(ctx context.Context, query string, opts ...interface{}) ([]*ent.NodeVersion, error) {
	index := a.client.InitIndex("node_versions_index")
	res, err := index.Search(query, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to search nodes: %w", err)
	}

	var nodes []*ent.NodeVersion
	if err := res.UnmarshalHits(&nodes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %w", err)
	}
	return nodes, nil
}
