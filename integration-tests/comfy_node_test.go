package integration

import (
	"context"
	"registry-backend/config"
	"registry-backend/drip"
	"registry-backend/ent/nodeversion"
	"registry-backend/mapper"
	authorization "registry-backend/server/middleware/authorization"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestRegistryComfyNode(t *testing.T) {
	// Setup DB and clean up
	client, cleanup := setupDB(t, context.Background())
	defer cleanup()

	// Initialize server implementation and authorization middleware
	impl := NewStrictServerImplementationWithMocks(client, &config.Config{})
	authz := authorization.NewAuthorizationManager(client, impl.RegistryService, impl.NewRelicApp).AuthorizationMiddleware()

	// Setup test user context and publisher
	ctx, _ := setupTestUser(client)
	publisherId, err := setupPublisher(ctx, authz, impl)
	require.NoError(t, err, "Failed to set up publisher")

	// Setup personal access token
	pat, err := setupPersonalAccessToken(ctx, authz, impl, publisherId)
	require.NoError(t, err, "Failed to create personal access token")

	// Setup a node
	node, err := setupNode(ctx, authz, impl, publisherId)
	require.NoError(t, err, "Failed to create node")

	// Mock external service responses for storage and Discord
	signedUrl := "test-url"
	impl.mockStorageService.On(
		"GenerateSignedURL", mock.Anything, mock.Anything).Return(signedUrl, nil)
	impl.mockStorageService.On(
		"GetFileUrl", mock.Anything, mock.Anything, mock.Anything).Return(signedUrl, nil)
	impl.mockDiscordService.On(
		"SendSecurityCouncilMessage", mock.Anything, mock.Anything).Return(nil)

	// Create test node versions
	nodeVersions := []*drip.NodeVersion{
		randomNodeVersion(0),
		randomNodeVersion(1),
		randomNodeVersion(2),
		randomNodeVersion(3),
		randomNodeVersion(4),
		randomNodeVersion(5),
	}

	// Identify specific node versions
	nodeVersionExtractionSucceeded1 := nodeVersions[len(nodeVersions)-1]
	nodeVersionExtractionSucceeded2 := nodeVersions[len(nodeVersions)-2]
	nodeVersionExtractionFailed := nodeVersions[len(nodeVersions)-3]
	backfilledNodeVersions := nodeVersions[:len(nodeVersions)-3]

	// Cloud Build Info
	cloudBuildInfo := &drip.ComfyNodeCloudBuildInfo{
		BuildId:       proto.String("test-build-id"),
		Location:      proto.String("test-location"),
		ProjectId:     proto.String("test-project-id"),
		ProjectNumber: proto.String("12345"),
	}

	// Create comfy nodes
	comfyNode1 := randomComfyNode()
	comfyNode2 := randomComfyNode()
	comfyNodes := drip.CreateComfyNodesJSONRequestBody{
		CloudBuildInfo: cloudBuildInfo,
		Nodes: &map[string]drip.ComfyNode{
			*comfyNode1.ComfyNodeId: comfyNode1,
			*comfyNode2.ComfyNodeId: comfyNode1,
		}}

	// Test case: Create Node Versions
	t.Run("CreateNodeVersions", func(t *testing.T) {
		for _, nv := range nodeVersions {
			// Publish each node version and ensure no error occurs
			_, err = withMiddleware(authz, impl.PublishNodeVersion)(ctx, drip.PublishNodeVersionRequestObject{
				PublisherId: publisherId,
				NodeId:      *node.Id,
				Body: &drip.PublishNodeVersionJSONRequestBody{
					PersonalAccessToken: *pat,
					Node:                *node,
					NodeVersion:         *nv,
				},
			})
			require.NoError(t, err, "should not return error")
		}
	})

	// Test case: Ensure no comfy nodes exist initially
	t.Run("NoComfyNode", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.GetNodeVersion)(ctx, drip.GetNodeVersionRequestObject{
			NodeId:    *node.Id,
			VersionId: *nodeVersionExtractionSucceeded1.Version,
		})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.GetNodeVersion200JSONResponse{}, res)
		assert.Empty(t, res.(drip.GetNodeVersion200JSONResponse).ComfyNodes)
	})

	t.Run("CreateComfyNodes", func(t *testing.T) {
		// Test case: Create comfy nodes for existing versions
		for _, nv := range []*drip.NodeVersion{nodeVersionExtractionSucceeded1, nodeVersionExtractionSucceeded2} {
			body := comfyNodes
			res, err := withMiddleware(authz, impl.CreateComfyNodes)(ctx, drip.CreateComfyNodesRequestObject{
				NodeId:  *node.Id,
				Version: *nv.Version,
				Body:    &body,
			})
			require.NoError(t, err)
			require.IsType(t, drip.CreateComfyNodes204Response{}, res)

			nv, err := client.NodeVersion.Query().Where(nodeversion.Version(*nv.Version)).Only(ctx)
			require.NoError(t, err)
			assert.Equal(t, mapper.ApiComfyNodeCloudBuildToDbComfyNodeCloudBuild(cloudBuildInfo), &nv.ComfyNodeCloudBuildInfo)
		}
	})

	// Test case: Get comfy nodes
	t.Run("GetComfyNodes", func(t *testing.T) {
		for k, v := range *comfyNodes.Nodes {
			v.ComfyNodeId = proto.String(k)
			t.Run(k, func(t *testing.T) {
				res, err := withMiddleware(authz, impl.GetComfyNode)(ctx, drip.GetComfyNodeRequestObject{
					NodeId:      *node.Id,
					Version:     *nodeVersionExtractionSucceeded1.Version,
					ComfyNodeId: k,
				})
				require.NoError(t, err, "should return created node version")
				require.IsType(t, drip.GetComfyNode200JSONResponse{}, res)
				assert.Equal(t, drip.GetComfyNode200JSONResponse(v), res.(drip.GetComfyNode200JSONResponse))
			})
		}
	})

	// Test case: Mark comfy node extraction as failed
	t.Run("FailedComfyNodesExtraction", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.CreateComfyNodes)(ctx, drip.CreateComfyNodesRequestObject{
			NodeId:  *node.Id,
			Version: *nodeVersionExtractionFailed.Version,
			Body:    &drip.CreateComfyNodesJSONRequestBody{Success: proto.Bool(false), CloudBuildInfo: cloudBuildInfo},
		})
		require.NoError(t, err)
		require.IsType(t, drip.CreateComfyNodes204Response{}, res)

		nv, err := client.NodeVersion.Query().Where(nodeversion.Version(*nodeVersionExtractionFailed.Version)).Only(ctx)
		require.NoError(t, err)
		assert.Equal(t, mapper.ApiComfyNodeCloudBuildToDbComfyNodeCloudBuild(cloudBuildInfo), &nv.ComfyNodeCloudBuildInfo)
	})

	// Test case: Conflict in creating comfy nodes
	t.Run("Conflict creating comfy nodes", func(t *testing.T) {
		body := comfyNodes
		res, err := withMiddleware(authz, impl.CreateComfyNodes)(ctx, drip.CreateComfyNodesRequestObject{
			NodeId:  *node.Id,
			Version: *nodeVersionExtractionSucceeded1.Version,
			Body:    &body,
		})
		require.NoError(t, err)
		require.IsType(t, drip.CreateComfyNodes409JSONResponse{}, res)
	})

	// Test case: Retrieve node version
	t.Run("GetNodeVersion contain comfy nodes", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.GetNodeVersion)(ctx, drip.GetNodeVersionRequestObject{
			NodeId:    *node.Id,
			VersionId: *nodeVersionExtractionSucceeded1.Version,
		})
		require.NoError(t, err, "should return created node version")
		require.IsType(t, drip.GetNodeVersion200JSONResponse{}, res)
		for k, v := range *res.(drip.GetNodeVersion200JSONResponse).ComfyNodes {
			ev := (*comfyNodes.Nodes)[k]
			ev.ComfyNodeId = proto.String(k)
			assert.Equal(t, ev, v)
		}
	})

	// Test case: List node versions
	t.Run("ListNodeVersion contains comfy nodes", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.ListNodeVersions)(ctx, drip.ListNodeVersionsRequestObject{
			NodeId: *node.Id,
		})
		require.NoError(t, err, "should return created node version")
		require.IsType(t, drip.ListNodeVersions200JSONResponse{}, res)
		found := false
		for _, nv := range res.(drip.ListNodeVersions200JSONResponse) {
			if *nv.Version == *nodeVersionExtractionSucceeded1.Version ||
				*nv.Version == *nodeVersionExtractionSucceeded2.Version {
				for k, v := range *nv.ComfyNodes {
					found = true
					ev := (*comfyNodes.Nodes)[k]
					ev.ComfyNodeId = proto.String(k)
					assert.Equal(t, ev, v)
				}
			} else {
				assert.Empty(t, nv.ComfyNodes)
			}
		}

		assert.True(t, found)
	})

	// Test case: Assert Algolia indexing
	t.Run("AssertAlgolia contains comfy nodes", func(t *testing.T) {
		indexed := impl.mockAlgolia.LastIndexedNodes
		require.Len(t, indexed, 1)

		// Verify the indexed node matches the created one
		assert.Equal(t, *node.Id, indexed[0].ID)

		// Compare comfy nodes in Algolia index with expected values
		for _, node := range indexed[0].Edges.Versions[0].Edges.ComfyNodes {
			// Map the indexed node to the API node
			cn := *(mapper.DBComfyNodeToApiComfyNode(node))

			// Compare with the corresponding expected comfy node
			expectedNode := (*comfyNodes.Nodes)[node.Name]

			// Compare all fields except ComfyNodeId if it's irrelevant
			assert.Equal(t, expectedNode.Category, cn.Category)
			assert.Equal(t, expectedNode.Function, cn.Function)
			assert.Equal(t, expectedNode.Description, cn.Description)
			assert.Equal(t, expectedNode.Deprecated, cn.Deprecated)
			assert.Equal(t, expectedNode.Experimental, cn.Experimental)
			assert.Equal(t, expectedNode.ReturnNames, cn.ReturnNames)
			assert.Equal(t, expectedNode.ReturnTypes, cn.ReturnTypes)
			assert.Equal(t, expectedNode.OutputIsList, cn.OutputIsList)
		}
	})

	// Test case: Trigger backfill for comfy nodes
	t.Run("TriggerBackfill", func(t *testing.T) {
		mockCalled := 0

		// Test case: Unlimited backfill
		t.Run("Unlimited", func(t *testing.T) {
			impl.mockPubsubService.
				On("PublishNodePack", mock.Anything, mock.Anything).
				Return(nil)
			res, err := withMiddleware(authz, impl.ComfyNodesBackfill)(ctx, drip.ComfyNodesBackfillRequestObject{})
			require.NoError(t, err, "should return created node version")
			require.IsType(t, drip.ComfyNodesBackfill204Response{}, res)
			impl.mockPubsubService.AssertNumberOfCalls(
				t, "PublishNodePack", len(backfilledNodeVersions)+mockCalled)
			mockCalled += len(backfilledNodeVersions)
		})

		// Test case: Limited backfill
		t.Run("Limited", func(t *testing.T) {
			limit := 2
			impl.mockPubsubService.
				On("PublishNodePack", mock.Anything, mock.Anything).
				Return(nil)
			res, err := withMiddleware(authz, impl.ComfyNodesBackfill)(
				ctx, drip.ComfyNodesBackfillRequestObject{Params: drip.ComfyNodesBackfillParams{MaxNode: &limit}})
			require.NoError(t, err, "should return created node version")
			require.IsType(t, drip.ComfyNodesBackfill204Response{}, res)
			impl.mockPubsubService.AssertNumberOfCalls(t, "PublishNodePack", limit+mockCalled)
			mockCalled += limit
		})
	})
}
