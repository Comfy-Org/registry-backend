package algolia

import (
	"context"
	"os"
	"registry-backend/ent"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	appid, ok := os.LookupEnv("ALGOLIA_APP_ID")
	if !ok {
		t.Skip("Required env variables `ALGOLIA_APP_ID` is not set")
	}
	apikey, ok := os.LookupEnv("ALGOLIA_API_KEY")
	if !ok {
		t.Skip("Required env variables `ALGOLIA_API_KEY` is not set")
	}

	algolia, err := New(appid, apikey)
	require.NoError(t, err)

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

}
