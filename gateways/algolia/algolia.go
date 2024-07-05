package algolia

import (
	"context"
	"fmt"
	"os"
	"registry-backend/ent"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/rs/zerolog/log"
)

// AlgoliaService defines the interface for interacting with Algolia search.
type AlgoliaService interface {
	IndexNodes(ctx context.Context, nodes ...*ent.Node) error
	SearchNodes(ctx context.Context, query string, opts ...interface{}) ([]*ent.Node, error)
	DeleteNode(ctx context.Context, node *ent.Node) error
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

// IndexNodes indexes the provided nodes in Algolia.
func (a *algolia) IndexNodes(ctx context.Context, nodes ...*ent.Node) error {
	index := a.client.InitIndex("nodes_index")
	objects := make([]struct {
		ObjectID string `json:"objectID"`
		*ent.Node
	}, len(nodes))

	for i, n := range nodes {
		objects[i] = struct {
			ObjectID string `json:"objectID"`
			*ent.Node
		}{
			ObjectID: n.ID,
			Node:     n,
		}
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

var _ AlgoliaService = (*algolianoop)(nil)

type algolianoop struct{}

func NewFromEnvOrNoop() (AlgoliaService, error) {
	id := os.Getenv("ALGOLIA_APP_ID")
	key := os.Getenv("ALGOLIA_API_KEY")
	if id == "" && key == "" {
		return &algolianoop{}, nil
	}

	return NewFromEnv()
}

// DeleteNode implements AlgoliaService.
func (a *algolianoop) DeleteNode(ctx context.Context, node *ent.Node) error {
	log.Ctx(ctx).Info().Msgf("algolia noop: delete node: %s", node.ID)
	return nil
}

// IndexNodes implements AlgoliaService.
func (a *algolianoop) IndexNodes(ctx context.Context, nodes ...*ent.Node) error {
	log.Ctx(ctx).Info().Msgf("algolia noop: index nodes: %d number of nodes", len(nodes))
	return nil
}

// SearchNodes implements AlgoliaService.
func (a *algolianoop) SearchNodes(ctx context.Context, query string, opts ...interface{}) ([]*ent.Node, error) {
	log.Ctx(ctx).Info().Msgf("algolia noop: search nodes: %s", query)
	return nil, nil
}
