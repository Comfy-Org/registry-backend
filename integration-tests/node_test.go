package integration

import (
	"context"
	"testing"
	"time"

	"registry-backend/config"
	"registry-backend/drip"
	authorization "registry-backend/server/middleware/authorization"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistryNode(t *testing.T) {
	client, cleanup := setupDB(t, context.Background())
	defer cleanup()

	// Initialize server implementation and authorization middleware
	impl := NewStrictServerImplementationWithMocks(client, &config.Config{})
	authz := authorization.NewAuthorizationManager(client, impl.RegistryService).AuthorizationMiddleware()

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

	// Test reindexing nodes
	t.Run("Reindex Nodes", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.ReindexNodes)(ctx, drip.ReindexNodesRequestObject{})
		require.NoError(t, err, "Node reindexing failed")
		assert.IsType(t, drip.ReindexNodes200Response{}, res)

		time.Sleep(1 * time.Second)
		nodes := impl.mockAlgolia.LastIndexedNodes
		require.Equal(t, 1, len(nodes))
		assert.Equal(t, *node.Id, nodes[0].ID)
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
