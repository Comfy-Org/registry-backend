package algolia

import (
	"context"
	"fmt"
	"os"
	"registry-backend/ent"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

type AlgoliaService interface {
	IndexNodes(ctx context.Context, n ...*ent.Node) error
	SearchNodes(ctx context.Context, query string, opts ...interface{}) (nodes []*ent.Node, err error)
	DeleteNode(ctx context.Context, n *ent.Node) error
}

var _ AlgoliaService = algolia{}

type algolia struct {
	client *search.Client
}

func New(appid, apikey string) (AlgoliaService, error) {
	return algolia{
		client: search.NewClient(appid, apikey),
	}, nil
}
func NewFromEnv() (AlgoliaService, error) {
	appid, ok := os.LookupEnv("ALGOLIA_APP_ID")
	if !ok {
		return nil, fmt.Errorf("Required env variables ALGOLIA_APP_ID is not set.")
	}
	apikey, ok := os.LookupEnv("ALGOLIA_API_KEY")
	if !ok {
		return nil, fmt.Errorf("Required env variables `ALGOLIA_API_KEY` is not set")
	}
	return algolia{
		client: search.NewClient(appid, apikey),
	}, nil
}

// IndexNodes implements AlgoliaService.
func (a algolia) IndexNodes(ctx context.Context, nodes ...*ent.Node) (err error) {
	index := a.client.InitIndex("comfy_registry_backend_node")
	objects := []struct {
		ObjectID string `json:"objectID"`
		*ent.Node
	}{}
	for _, n := range nodes {
		objects = append(objects, struct {
			ObjectID string `json:"objectID"`
			*ent.Node
		}{
			ObjectID: n.ID,
			Node:     n,
		})
	}
	res, err := index.SaveObjects(objects)
	if err != nil {
		return fmt.Errorf("failed to index node :%w", err)
	}

	return res.Wait()
}

func (a algolia) SearchNodes(ctx context.Context, query string, opts ...interface{}) (nodes []*ent.Node, err error) {
	index := a.client.InitIndex("comfy_registry_backend_node")
	res, err := index.Search(query, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to search nodes: %w", err)
	}
	nodes = make([]*ent.Node, 0)
	err = res.UnmarshalHits(&nodes)
	return
}

func (a algolia) DeleteNode(ctx context.Context, n *ent.Node) error {
	index := a.client.InitIndex("comfy_registry_backend_node")
	res, err := index.DeleteObject(n.ID)
	if err != nil {
		return fmt.Errorf("fail to delete node")
	}
	return res.Wait()
}
