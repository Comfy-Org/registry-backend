package algolia

import (
	"context"
	"os"
	"registry-backend/config"
	"registry-backend/ent"
	"registry-backend/ent/schema"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	_, ok := os.LookupEnv("ALGOLIA_APP_ID")
	if !ok {
		t.Skip("Required env variables `ALGOLIA_APP_ID` is not set")
	}
	_, ok = os.LookupEnv("ALGOLIA_API_KEY")
	if !ok {
		t.Skip("Required env variables `ALGOLIA_API_KEY` is not set")
	}

	algolia, err := NewAlgoliaService(&config.Config{
		AlgoliaAppID:  os.Getenv("ALGOLIA_APP_ID"),
		AlgoliaAPIKey: os.Getenv("ALGOLIA_API_KEY"),
	})
	require.NoError(t, err)

	t.Run("node", func(t *testing.T) {
		ctx := context.Background()
		node := &ent.Node{
			ID:          uuid.NewString(),
			Name:        t.Name() + "-" + uuid.NewString(),
			TotalStar:   98,
			TotalReview: 20,
		}
		for i := 0; i < 10; i++ {
			err = algolia.IndexNodes(ctx, node)
			require.NoError(t, err)
		}

		<-time.After(time.Second * 10)
		nodes, err := algolia.SearchNodes(ctx, node.Name)
		require.NoError(t, err)
		require.Len(t, nodes, 1)
		assert.Equal(t, node, nodes[0])
	})

	t.Run("nodeVersion", func(t *testing.T) {
		ctx := context.Background()
		now := time.Now()
		nv := &ent.NodeVersion{
			ID:              uuid.New(),
			NodeID:          uuid.NewString(),
			Version:         "v1.0.0-" + uuid.NewString(),
			Changelog:       "test",
			Status:          schema.NodeVersionStatusActive,
			StatusReason:    "test",
			PipDependencies: []string{"test"},
			CreateTime:      time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC),
			UpdateTime:      time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC),
		}
		for i := 0; i < 10; i++ {
			err = algolia.IndexNodeVersions(ctx, nv)
			require.NoError(t, err)
		}

		<-time.After(time.Second * 10)
		nodes, err := algolia.SearchNodeVersions(ctx, nv.Version)
		require.NoError(t, err)
		require.Len(t, nodes, 1)
		assert.Equal(t, nv, nodes[0])
	})
}

func TestNoop(t *testing.T) {
	a, err := NewAlgoliaService(&config.Config{})
	require.NoError(t, err)
	require.NoError(t, a.IndexNodes(context.Background(), &ent.Node{}))
	require.NoError(t, a.DeleteNode(context.Background(), &ent.Node{}))
}
