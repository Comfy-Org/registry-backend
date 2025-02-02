package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"registry-backend/config"
	"registry-backend/drip"
	"registry-backend/ent/nodeversion"
	"registry-backend/ent/schema"
	logging "registry-backend/logging"
	authorization "registry-backend/server/middleware/authorization"
	registry "registry-backend/services/registry"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestRegistryNodeVersion(t *testing.T) {
	client, cleanup := setupDB(t, context.Background())
	defer cleanup()

	// Initialize server implementation and authorization middleware
	impl := NewStrictServerImplementationWithMocks(client, &config.Config{
		CloudStorageBucketName: "test-bucket",
	})
	authz := authorization.NewAuthorizationManager(client, impl.RegistryService, impl.NewRelicApp).AuthorizationMiddleware()

	// Setup user context and publisher
	ctx, _ := setupTestUser(client)
	publisherId, err := setupPublisher(ctx, authz, impl)
	require.NoError(t, err, "Failed to set up publisher")

	// Setup personal access token
	pat, err := setupPersonalAccessToken(ctx, authz, impl, publisherId)
	require.NoError(t, err, "Failed to create personal access token")

	// Setup a node
	node, err := setupNode(ctx, authz, impl, publisherId)
	require.NoError(t, err, "Failed to create node")

	nodeVersion := randomNodeVersion(0)
	signedUrl := "test-url"
	downloadUrl := fmt.Sprintf(
		"https://storage.googleapis.com/test-bucket/%s/%s/%s/node.zip", publisherId, *node.Id, *nodeVersion.Version)

	// Mock external service responses for storage and Discord
	impl.mockStorageService.
		On("GenerateSignedURL", mock.Anything, mock.Anything).
		Return(signedUrl, nil)
	impl.mockStorageService.
		On("GetFileUrl", mock.Anything, mock.Anything, mock.Anything).
		Return(signedUrl, nil)
	impl.mockDiscordService.
		On("SendSecurityCouncilMessage", mock.Anything, mock.Anything).
		Return(nil)

	var createdNodeVersion drip.NodeVersion

	// Test case for listing node versions before creating any
	t.Run("List Node Version Before Create", func(t *testing.T) {
		resVersions, err := withMiddleware(authz, impl.ListNodeVersions)(
			ctx, drip.ListNodeVersionsRequestObject{NodeId: *node.Id})
		require.NoError(t, err, "should return error since node version doesn't exists")
		require.IsType(t, drip.ListNodeVersions200JSONResponse{}, resVersions)
		assert.Empty(t, resVersions.(drip.ListNodeVersions200JSONResponse),
			"should not return any node versions")
	})

	// Test case for creating a node version
	t.Run("Create Node Version", func(t *testing.T) {
		// Call the function with middleware
		res, err := withMiddleware(authz, impl.PublishNodeVersion)(ctx, drip.PublishNodeVersionRequestObject{
			PublisherId: publisherId,
			NodeId:      *node.Id,
			Body: &drip.PublishNodeVersionJSONRequestBody{
				PersonalAccessToken: *pat,
				Node:                *node,
				NodeVersion:         *nodeVersion,
			},
		})

		// Assert no error occurred
		require.NoError(t, err, "should return created node version")

		// Type assertion to get response object
		nodeVersionResp := res.(drip.PublishNodeVersion201JSONResponse)

		// Validate response fields
		assert.Equal(t, nodeVersion.Version,
			nodeVersionResp.NodeVersion.Version, "should return the correct node version")
		assert.Equal(t, nodeVersion.Dependencies,
			nodeVersionResp.NodeVersion.Dependencies, "should return the correct pip dependencies")
		assert.Equal(t, nodeVersion.Changelog,
			nodeVersionResp.NodeVersion.Changelog, "should return the correct changelog")
		assert.Equal(t, signedUrl, *nodeVersionResp.SignedUrl, "should return the correct signed URL")

		// Ensure the status is 'pending' when first created
		expectedStatus := drip.NodeVersionStatusPending
		assert.Equal(t, expectedStatus,
			*nodeVersionResp.NodeVersion.Status, "should return pending status")

		// Store created node version for later tests
		createdNodeVersion = *nodeVersionResp.NodeVersion
	})

	// Test case for admin updating node version status
	t.Run("Admin Update", func(t *testing.T) {
		// Setup admin context
		adminCtx, _ := setupAdminUser(client)
		activeStatus := drip.NodeVersionStatusActive

		// Request to update the node version status
		adminUpdateReq := drip.AdminUpdateNodeVersionRequestObject{
			NodeId:        *node.Id,
			VersionNumber: *createdNodeVersion.Version,
			Body: &drip.AdminUpdateNodeVersionJSONRequestBody{
				Status:       &activeStatus,
				StatusReason: proto.String("test reason"), // Provide reason for status change
			},
		}

		// Call the update function
		adminUpdateNodeVersionResp, err := impl.AdminUpdateNodeVersion(adminCtx, adminUpdateReq)

		// Assert no error occurred during update
		require.NoError(t, err, "should return updated node version")

		// Validate that the status was updated correctly
		assert.Equal(t, activeStatus, *adminUpdateNodeVersionResp.(drip.AdminUpdateNodeVersion200JSONResponse).Status, "should return active status")
	})

	t.Run("Create Node Version with Fake Token", func(t *testing.T) {
		// Attempt to create a node version with an invalid token ("faketoken")
		reqObj := drip.PublishNodeVersionRequestObject{
			PublisherId: publisherId,
			NodeId:      *node.Id,
			Body: &drip.PublishNodeVersionJSONRequestBody{
				Node:                *node,
				NodeVersion:         *nodeVersion,
				PersonalAccessToken: "faketoken", // Fake token used here
			},
		}

		// Perform the request with middleware
		_, err := withMiddleware(authz, impl.PublishNodeVersion)(ctx, reqObj)

		// Assert that an error occurs (invalid token should cause failure)
		require.Error(t, err, "should return error for fake token")

		// Check that the error is a bad request (HTTP 400)
		assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code, "should return 400 bad request")
	})

	t.Run("Get not exist Node Version", func(t *testing.T) {
		// Attempt to get a non-existent node version (using fake node ID)
		reqObj := drip.GetNodeVersionRequestObject{
			NodeId:    *node.Id + "fake", // Fake node ID
			VersionId: *nodeVersion.Version,
		}

		// Perform the request to fetch node version
		res, err := withMiddleware(authz, impl.GetNodeVersion)(ctx, reqObj)

		// Assert no error occurred during the request
		require.NoError(t, err, "should not return error for non-existent node version")

		// Check that the response is of type 'not found' (404)
		assert.IsType(t, drip.GetNodeVersion404JSONResponse{}, res, "should return 404 response")
	})

	t.Run("Create Node Version of Not Exist Node", func(t *testing.T) {
		// Attempt to create a node version for a non-existent node (using fake node ID)
		reqObj := drip.PublishNodeVersionRequestObject{
			PublisherId: publisherId,
			NodeId:      *node.Id + "fake",                         // Fake node ID
			Body:        &drip.PublishNodeVersionJSONRequestBody{}, // Empty body for simplicity
		}

		// Perform the request to create the node version
		_, err := withMiddleware(authz, impl.PublishNodeVersion)(ctx, reqObj)

		// Assert that an error occurs (should fail due to non-existent node)
		require.Error(t, err, "should return error for non-existent node")

		// Check that the error is a bad request (HTTP 400)
		assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code, "should return 400 bad request")
	})

	t.Run("List Node Versions", func(t *testing.T) {
		// Request to list node versions
		reqObj := drip.ListNodeVersionsRequestObject{
			NodeId: *node.Id,
		}

		// Perform the request with middleware
		resVersions, err := withMiddleware(authz, impl.ListNodeVersions)(ctx, reqObj)
		require.NoError(t, err)
		require.IsType(t, drip.ListNodeVersions200JSONResponse{}, resVersions)

		resVersions200 := resVersions.(drip.ListNodeVersions200JSONResponse)
		require.Len(t, resVersions200, 1)

		nodeVersionStatus := drip.NodeVersionStatusActive
		assert.Equal(t, drip.NodeVersion{
			Id:           resVersions200[0].Id,
			CreatedAt:    resVersions200[0].CreatedAt,
			Deprecated:   proto.Bool(false),
			Version:      nodeVersion.Version,
			Changelog:    nodeVersion.Changelog,
			Dependencies: nodeVersion.Dependencies,
			DownloadUrl:  &downloadUrl,
			Status:       &nodeVersionStatus,
			StatusReason: proto.String(""),
			NodeId:       node.Id,
		}, resVersions200[0])

		// Request node versions with status reason
		reqObjWithReason := drip.ListNodeVersionsRequestObject{
			NodeId: *node.Id,
			Params: drip.ListNodeVersionsParams{
				IncludeStatusReason: proto.Bool(true),
			},
		}

		resVersions, err = withMiddleware(authz, impl.ListNodeVersions)(ctx, reqObjWithReason)
		require.NoError(t, err)
		require.IsType(t, drip.ListNodeVersions200JSONResponse{}, resVersions)

		resVersions200 = resVersions.(drip.ListNodeVersions200JSONResponse)
		assert.Equal(t, "test reason", *resVersions200[0].StatusReason)
	})

	t.Run("Update Node Version", func(t *testing.T) {
		updatedChangelog := "test-changelog-2"

		// Request to update the node version
		resUNV, err := withMiddleware(authz, impl.UpdateNodeVersion)(ctx, drip.UpdateNodeVersionRequestObject{
			PublisherId: publisherId,
			NodeId:      *node.Id,
			VersionId:   *createdNodeVersion.Id,
			Body: &drip.NodeVersionUpdateRequest{
				Changelog:  &updatedChangelog,
				Deprecated: proto.Bool(true),
			},
		})
		require.NoError(t, err)
		require.IsType(t, drip.UpdateNodeVersion200JSONResponse{}, resUNV)

		// List node versions to verify the update
		res, err := withMiddleware(authz, impl.ListNodeVersions)(
			ctx, drip.ListNodeVersionsRequestObject{NodeId: *node.Id})
		require.NoError(t, err)
		require.IsType(t, drip.ListNodeVersions200JSONResponse{}, res)

		res200 := res.(drip.ListNodeVersions200JSONResponse)
		require.Len(t, res200, 1)

		status := drip.NodeVersionStatusActive
		updatedNodeVersion := drip.NodeVersion{
			Id:           res200[0].Id,
			CreatedAt:    res200[0].CreatedAt,
			Deprecated:   proto.Bool(true),
			Version:      nodeVersion.Version,
			Dependencies: nodeVersion.Dependencies,
			Changelog:    &updatedChangelog,
			DownloadUrl:  &downloadUrl,
			Status:       &status,
			StatusReason: proto.String(""),
			NodeId:       node.Id,
		}

		assert.Equal(t, updatedNodeVersion, res200[0])
		createdNodeVersion = res200[0]
	})

	t.Run("List Nodes", func(t *testing.T) {
		// Upload an additional node version
		nodeVersion2 := randomNodeVersion(1)
		res, err := withMiddleware(authz, impl.PublishNodeVersion)(ctx, drip.PublishNodeVersionRequestObject{
			PublisherId: publisherId,
			NodeId:      *node.Id,
			Body: &drip.PublishNodeVersionJSONRequestBody{
				PersonalAccessToken: *pat,
				Node:                *node,
				NodeVersion:         *nodeVersion2,
			},
		})
		nodeVersionResp := res.(drip.PublishNodeVersion201JSONResponse)
		require.NoError(t, err, "should return created node version")

		// Retrieve and verify the list of nodes
		resNodes, err := withMiddleware(authz, impl.ListAllNodes)(ctx, drip.ListAllNodesRequestObject{})
		require.NoError(t, err)
		require.IsType(t, drip.ListAllNodes200JSONResponse{}, resNodes)
		resNodes200 := resNodes.(drip.ListAllNodes200JSONResponse)

		// Iterate over each node and assert individual fields
		for _, node := range *resNodes200.Nodes {
			// Assertions for basic node attributes
			assert.Equal(t, *node.Id, *node.Id)
			assert.Equal(t, *node.Name, *node.Name)
			assert.Equal(t, *node.Repository, *node.Repository)
			assert.Equal(t, *node.Description, *node.Description)
			assert.Equal(t, *node.Author, *node.Author)
			assert.Equal(t, *node.License, *node.License)
			assert.Equal(t, *node.Tags, *node.Tags)

			// Assert the latest version
			assert.Equal(t, *node.LatestVersion.Version, *nodeVersionResp.NodeVersion.Version)
			assert.Equal(t, *node.LatestVersion.Changelog, *nodeVersionResp.NodeVersion.Changelog)
			assert.Equal(t, *node.LatestVersion.Dependencies, *nodeVersionResp.NodeVersion.Dependencies)
			assert.Equal(t, *node.LatestVersion.DownloadUrl, *nodeVersionResp.NodeVersion.DownloadUrl)
			assert.Equal(t, *node.LatestVersion.Status, *nodeVersionResp.NodeVersion.Status)

			// Status checks
			nodeStatus := drip.NodeStatusActive
			assert.Equal(t, *node.Status, nodeStatus)
		}
	})

	t.Run("Index Nodes", func(t *testing.T) {
		// Setting up logger context for the test
		ctx := logging.SetupLogger().WithContext(ctx)

		// Triggering reindexing of nodes and verifying the response
		res, err := withMiddleware(authz, impl.ReindexNodes)(ctx, drip.ReindexNodesRequestObject{})
		require.NoError(t, err, "should not return error")    // Ensure no error occurred
		assert.IsType(t, drip.ReindexNodes200Response{}, res) // Ensure response is of expected type
	})

	t.Run("Node Installation", func(t *testing.T) {
		// Installing the node without specifying a version
		resIns, err := withMiddleware(authz, impl.InstallNode)(ctx, drip.InstallNodeRequestObject{NodeId: *node.Id})
		require.NoError(t, err, "should not return error") // Ensure no error occurred
		require.IsType(t, drip.InstallNode200JSONResponse{}, resIns, "should return 200 response type")

		// Installing the node with a specific version
		resIns, err = withMiddleware(authz, impl.InstallNode)(ctx, drip.InstallNodeRequestObject{
			NodeId: *node.Id, Params: drip.InstallNodeParams{Version: createdNodeVersion.Version}})
		require.NoError(t, err, "should not return error") // Ensure no error occurred
		require.IsType(t, drip.InstallNode200JSONResponse{}, resIns, "should return 200 response type")

		// Verifying the total number of installs for the node
		t.Run("Get Total Install", func(t *testing.T) {
			res, err := withMiddleware(authz, impl.GetNode)(ctx, drip.GetNodeRequestObject{
				NodeId: *node.Id,
			})
			require.NoError(t, err, "should not return error")    // Ensure no error occurred
			require.IsType(t, drip.GetNode200JSONResponse{}, res) // Ensure response is of expected type
			// Verify the number of downloads (installs) for the node is 2
			assert.Equal(t, int(2), *res.(drip.GetNode200JSONResponse).Downloads)
		})

		// Adding a review for the node with a rating of 5 stars
		t.Run("Add review", func(t *testing.T) {
			res, err := withMiddleware(authz, impl.PostNodeReview)(ctx, drip.PostNodeReviewRequestObject{
				NodeId: *node.Id,
				Params: drip.PostNodeReviewParams{Star: 5},
			})
			require.NoError(t, err)                                      // Ensure no error occurred
			require.IsType(t, drip.PostNodeReview200JSONResponse{}, res) // Ensure response is of expected type
			res200 := res.(drip.PostNodeReview200JSONResponse)
			// Verify the rating returned is 5
			assert.Equal(t, float32(5), *res200.Rating)
		})
	})

	t.Run("Node installation on non-existent node or version", func(t *testing.T) {
		// Attempting to install a node with a non-existent ID
		resIns, err := withMiddleware(authz, impl.InstallNode)(ctx, drip.InstallNodeRequestObject{NodeId: *node.Id + "fake"})
		require.NoError(t, err, "should not return error") // Ensure no error occurred
		require.IsType(t, drip.InstallNode404JSONResponse{}, resIns, "should return 404 error")

		// Attempting to install a node with a non-existent version
		resIns, err = withMiddleware(authz, impl.InstallNode)(ctx, drip.InstallNodeRequestObject{
			NodeId: *node.Id, Params: drip.InstallNodeParams{Version: proto.String(*createdNodeVersion.Version + "fake")}})
		require.NoError(t, err, "should not return error") // Ensure no error occurred
		require.IsType(t, drip.InstallNode404JSONResponse{}, resIns, "should return 404 error")
	})

	t.Run("Scan Node", func(t *testing.T) {
		// Creating a new random node and version for scanning
		node := randomNode()
		nodeVersion := randomNodeVersion(0)
		downloadUrl := fmt.Sprintf("https://storage.googleapis.com/test-bucket/%s/%s/%s/node.zip", publisherId, *node.Id, *nodeVersion.Version)

		// Mocking the behavior of services for URL generation and message sending
		impl.mockStorageService.On("GenerateSignedURL", mock.Anything, mock.Anything).Return("test-url", nil)
		impl.mockStorageService.On("GetFileUrl", mock.Anything, mock.Anything, mock.Anything).Return("test-url", nil)
		impl.mockDiscordService.On("SendSecurityCouncilMessage", mock.Anything, mock.Anything).Return(nil)

		// Publishing the node version
		_, err := withMiddleware(authz, impl.PublishNodeVersion)(ctx, drip.PublishNodeVersionRequestObject{
			PublisherId: publisherId,
			NodeId:      *node.Id,
			Body: &drip.PublishNodeVersionJSONRequestBody{
				Node:                *node,
				NodeVersion:         *nodeVersion,
				PersonalAccessToken: *pat,
			},
		})
		require.NoError(t, err, "should return created node version")

		// Checking the number of nodes to be scanned
		nodesToScans, err := client.NodeVersion.Query().Where(nodeversion.StatusEQ(schema.NodeVersionStatusPending)).Count(ctx)
		require.NoError(t, err)

		// Mocking a server that handles the scan request and verifying the response
		newNodeScanned := false
		nodesScanned := 0
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			// Verifying the scan request
			req := registry.ScanRequest{}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
			if downloadUrl == req.URL {
				newNodeScanned = true
			}
			nodesScanned++
		}))
		t.Cleanup(s.Close)

		// Setting up the mocked implementation with the server URL
		impl := NewStrictServerImplementationWithMocks(client, &config.Config{SecretScannerURL: s.URL})
		dur := time.Duration(0)
		// Triggering the security scan
		scanres, err := withMiddleware(authz, impl.SecurityScan)(ctx, drip.SecurityScanRequestObject{
			Params: drip.SecurityScanParams{
				MinAge: &dur,
			},
		})
		require.NoError(t, err)
		require.IsType(t, drip.SecurityScan200Response{}, scanres)

		// Verifying that a new node was scanned and the count of scanned nodes is correct
		assert.True(t, newNodeScanned)
		assert.Equal(t, nodesToScans, nodesScanned)
	})
}
