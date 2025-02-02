package integration

import (
	"context"
	"google.golang.org/protobuf/proto"
	"strings"
	"testing"
	"time"

	"registry-backend/config"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/node"
	drip_logging "registry-backend/logging"
	authorization "registry-backend/server/middleware/authorization"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistryNode(t *testing.T) {
	client, cleanup := setupDB(t, context.Background())
	defer cleanup()

	// Initialize server implementation and authorization middleware
	impl := NewStrictServerImplementationWithMocks(client, &config.Config{})
	authz := authorization.NewAuthorizationManager(client, impl.RegistryService, impl.NewRelicApp).AuthorizationMiddleware()

	// Setup user context and publisher
	ctx, _ := setupTestUser(client)
	publisherId, err := setupPublisher(ctx, authz, impl)
	require.NoError(t, err, "Failed to set up publisher")

	node := randomNode()

	// Test creating a node
	t.Run("Create Node", func(t *testing.T) {
		createResponse, err := withMiddleware(authz, impl.CreateNode)(ctx, drip.CreateNodeRequestObject{
			PublisherId: publisherId,
			Body:        node,
		})
		require.NoError(t, err, "Node creation failed")
		require.NotNil(t, createResponse, "Node creation returned nil response")

		createdNode := createResponse.(drip.CreateNode201JSONResponse)

		// Validate the node creation response
		assert.Equal(t, node.Id, createdNode.Id)
		assert.Equal(t, node.Description, createdNode.Description)
		assert.Equal(t, node.Author, createdNode.Author)
		assert.Equal(t, node.License, createdNode.License)
		assert.Equal(t, node.Name, createdNode.Name)
		assert.Equal(t, node.Tags, createdNode.Tags)
		assert.Equal(t, node.Icon, createdNode.Icon)
		assert.Equal(t, node.Repository, createdNode.Repository)
		assert.Equal(t, drip.NodeStatusActive, *createdNode.Status)
	})

	// Test fetching the created node
	t.Run("Get Node", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.GetNode)(ctx, drip.GetNodeRequestObject{NodeId: *node.Id})
		require.NoError(t, err, "Failed to fetch node")
		require.IsType(t, drip.GetNode200JSONResponse{}, res)

		fetchedNode := res.(drip.GetNode200JSONResponse)
		assert.Equal(t, node.Id, fetchedNode.Id)
		assert.Equal(t, node.Description, fetchedNode.Description)
		assert.Equal(t, node.Author, fetchedNode.Author)
		assert.Equal(t, node.License, fetchedNode.License)
		assert.Equal(t, node.Name, fetchedNode.Name)
		assert.Equal(t, node.Tags, fetchedNode.Tags)
		assert.Equal(t, node.Icon, fetchedNode.Icon)
		assert.Equal(t, node.Repository, fetchedNode.Repository)
	})

	// Test filtering nodes after a specific timestamp
	t.Run("Filter Nodes After Timestamp", func(t *testing.T) {
		futureTime := time.Now()
		res, err := withMiddleware(authz, impl.ListAllNodes)(ctx, drip.ListAllNodesRequestObject{
			Params: drip.ListAllNodesParams{
				Timestamp: &futureTime,
			},
		})
		require.NoError(t, err, "Failed to filter nodes by timestamp")
		require.IsType(t, drip.ListAllNodes200JSONResponse{}, res)

		nodesResponse := res.(drip.ListAllNodes200JSONResponse)
		require.Len(t, *nodesResponse.Nodes, 0, "Expected no nodes to be returned")

		pastTime := time.Now().Add(-time.Hour)
		res, err = withMiddleware(authz, impl.ListAllNodes)(ctx, drip.ListAllNodesRequestObject{
			Params: drip.ListAllNodesParams{
				Timestamp: &pastTime,
			},
		})
		require.NoError(t, err, "Failed to filter nodes by timestamp")
		require.IsType(t, drip.ListAllNodes200JSONResponse{}, res)

		nodesResponse = res.(drip.ListAllNodes200JSONResponse)
		require.Len(t, *nodesResponse.Nodes, 1, "Expected one node to be returned")
	})

	// Test listing nodes for the publisher
	t.Run("Get Publisher Nodes", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.ListNodesForPublisher)(ctx, drip.ListNodesForPublisherRequestObject{
			PublisherId: publisherId,
		})
		require.NoError(t, err, "Failed to fetch nodes for publisher")
		require.IsType(t, drip.ListNodesForPublisher200JSONResponse{}, res)

		nodes := res.(drip.ListNodesForPublisher200JSONResponse)
		require.Len(t, nodes, 1)

		// Compare each field individually
		expectedNode := *node

		// Compare each field individually
		assert.Equal(t, expectedNode.Id, nodes[0].Id)
		assert.Equal(t, expectedNode.Name, nodes[0].Name)
		assert.Equal(t, expectedNode.Description, nodes[0].Description)
		assert.Equal(t, expectedNode.Author, nodes[0].Author)
		assert.Equal(t, expectedNode.License, nodes[0].License)
		assert.Equal(t, expectedNode.Tags, nodes[0].Tags)
		assert.Equal(t, expectedNode.Icon, nodes[0].Icon)
		assert.Equal(t, expectedNode.Repository, nodes[0].Repository)
	})

	// Test updating the node
	t.Run("Update Node", func(t *testing.T) {
		updatedNode := randomNode()
		updatedNode.Id = node.Id // Retain the same ID for updating
		node = updatedNode       // Update the reference for further tests

		updateResponse, err := withMiddleware(authz, impl.UpdateNode)(ctx, drip.UpdateNodeRequestObject{
			PublisherId: publisherId,
			NodeId:      *node.Id,
			Body:        node,
		})
		require.NoError(t, err, "Node update failed")
		require.NotNil(t, updateResponse, "Node update returned nil response")

		updatedResponse := updateResponse.(drip.UpdateNode200JSONResponse)
		assert.Equal(t, node.Id, updatedResponse.Id)
		assert.Equal(t, node.Description, updatedResponse.Description)
		assert.Equal(t, node.Author, updatedResponse.Author)
		assert.Equal(t, node.License, updatedResponse.License)
		assert.Equal(t, node.Name, updatedResponse.Name)
		assert.Equal(t, node.Tags, updatedResponse.Tags)
		assert.Equal(t, node.Icon, updatedResponse.Icon)
		assert.Equal(t, node.Repository, updatedResponse.Repository)
	})

	// Test for duplicate node_id with case-insensitive comparison
	t.Run("Create Duplicate Node ID (Case-Insensitive)", func(t *testing.T) {
		// Create a node with an ID that differs only by case
		duplicateNode := randomNode()
		duplicateNode.Id = proto.String(strings.ToUpper(*node.Id)) // Use uppercase version of the original node ID

		// Attempt to create the node
		_, err := withMiddleware(authz, impl.CreateNode)(ctx, drip.CreateNodeRequestObject{
			PublisherId: publisherId,
			Body:        duplicateNode,
		})

		// Expect an error indicating duplicate ID
		require.Error(t, err, "Creating a node with a duplicate ID (case-insensitive) should fail")
		assert.Contains(t, err.Error(), "duplicate", "Error should indicate duplicate node_id")
	})

	// Test deleting the node
	t.Run("Delete Node", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.DeleteNode)(ctx, drip.DeleteNodeRequestObject{
			PublisherId: publisherId,
			NodeId:      *node.Id,
		})
		require.NoError(t, err, "Node deletion failed")
		assert.IsType(t, drip.DeleteNode204Response{}, res)
	})

	// Test deleting a nonexistent node
	t.Run("Delete Nonexistent Node", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.DeleteNode)(ctx, drip.DeleteNodeRequestObject{
			PublisherId: publisherId,
			NodeId:      "nonexistent-id",
		})
		require.NoError(t, err, "Deleting nonexistent node should not return an error")
		assert.IsType(t, drip.DeleteNode204Response{}, res)
	})
}

func TestRegistryNodeReindex(t *testing.T) {
	client, cleanup := setupDB(t, context.Background())
	defer cleanup()

	// Initialize server implementation and authorization middleware
	impl := NewStrictServerImplementationWithMocks(client, &config.Config{})
	authz := authorization.NewAuthorizationManager(
		client, impl.RegistryService, impl.NewRelicApp).AuthorizationMiddleware()

	// Setup user context and publisher
	ctx, _ := setupTestUser(client)
	ctx = drip_logging.SetupLogger().WithContext(ctx)
	publisherId, err := setupPublisher(ctx, authz, impl)
	require.NoError(t, err, "Failed to set up publisher")

	storeRandomNodes := func(t *testing.T, n int) []drip.Node {
		nodes := make([]drip.Node, 0, n)
		for i := 0; i < cap(nodes); i++ {
			node := randomNode()
			createResponse, err := withMiddleware(authz, impl.CreateNode)(ctx, drip.CreateNodeRequestObject{
				PublisherId: publisherId,
				Body:        node,
			})
			require.NoError(t, err, "Node creation failed")
			require.NotNil(t, createResponse, "Node creation returned nil response")

			createdNode := createResponse.(drip.CreateNode201JSONResponse)
			nodes = append(nodes, drip.Node(createdNode))
		}
		return nodes
	}

	fetchIndexed := func(t *testing.T, ctx context.Context, indexedAfter time.Time, expectedLen int) []*ent.Node {
		var indexed []*ent.Node
		for {
			require.NoError(t, ctx.Err())

			indexed, err = client.Node.Query().
				Where(node.LastAlgoliaIndexTimeGT(indexedAfter)).
				Where(node.LastAlgoliaIndexTimeLT(time.Now())).
				All(ctx)
			require.NoError(t, err)

			if len(indexed) < expectedLen {
				time.Sleep(time.Second)
				continue
			}

			return indexed
		}
	}

	now := time.Now()
	nodes := storeRandomNodes(t, 100)

	t.Run("AfterCreate", func(t *testing.T) {
		indexed := fetchIndexed(t, ctx, now, len(nodes))
		assert.Equal(t, len(nodes), len(indexed))
	})

	t.Run("AfterReindex", func(t *testing.T) {
		now, batch := time.Now(), 9
		res, err := withMiddleware(authz, impl.ReindexNodes)(ctx, drip.ReindexNodesRequestObject{
			Params: drip.ReindexNodesParams{
				MaxBatch: &batch,
			},
		})
		require.NoError(t, err, "Node reindexing failed")
		assert.IsType(t, drip.ReindexNodes200Response{}, res)

		// check last_algolia_index_time
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		indexed := fetchIndexed(t, ctx, now, len(nodes))

		assert.Equal(t, len(nodes), len(indexed), "should reindex all nodes")
		assert.Len(t, impl.mockAlgolia.LastIndexedNodes, len(nodes)%batch, "should index to algolia partially")

	})

	t.Run("AfterReindexWithMinAge", func(t *testing.T) {
		batch, age := 8, 3*time.Second
		time.Sleep(age)

		// add more nodes that will not be reindexed since it is too new
		{
			now := time.Now()
			newNodes := storeRandomNodes(t, 20)
			indexed := fetchIndexed(t, ctx, now, len(newNodes))
			assert.Equal(t, len(newNodes), len(indexed))
		}

		now = time.Now()
		res, err := withMiddleware(authz, impl.ReindexNodes)(ctx, drip.ReindexNodesRequestObject{
			Params: drip.ReindexNodesParams{
				MaxBatch: &batch,
				MinAge:   &age,
			},
		})
		require.NoError(t, err, "Node reindexing failed")
		assert.IsType(t, drip.ReindexNodes200Response{}, res)

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		indexed := fetchIndexed(t, ctx, now, len(nodes))
		assert.Equal(t, len(nodes), len(indexed), "should reindex some nodes")
		assert.Len(t, impl.mockAlgolia.LastIndexedNodes, len(nodes)%batch, "should index to algolia partially")
	})
}
