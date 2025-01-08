package integration

import (
	"context"
	"google.golang.org/protobuf/proto"
	"net/http"
	"registry-backend/config"
	"registry-backend/drip"
	"registry-backend/ent/schema"
	drip_authorization "registry-backend/server/middleware/authorization"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeBan(t *testing.T) {
	clientCtx := context.Background()
	client, cleanup := setupDB(t, clientCtx)
	defer cleanup()

	// Setup the mock services and server
	impl := NewStrictServerImplementationWithMocks(client, &config.Config{})
	authz := drip_authorization.NewAuthorizationManager(client, impl.RegistryService).AuthorizationMiddleware()

	t.Run("Node Ban Tests", func(t *testing.T) {
		userCtx, _ := setupTestUser(client)
		adminCtx, _ := setupAdminUser(client)

		// Setup a test publisher
		publisherId, err := setupPublisher(userCtx, authz, impl)
		require.NoError(t, err, "should set up publisher")

		// Setup a test node
		nodeId, err := setupNode(userCtx, authz, impl, publisherId)
		require.NoError(t, err, "should set up node")

		// Setup a personal access token
		pat, err := setupPersonalAccessToken(userCtx, authz, impl, publisherId)
		require.NoError(t, err, "should set up personal access token")

		t.Run("Ban node by non-admin", func(t *testing.T) {
			// Attempt to ban the node as a non-admin user
			res, err := withMiddleware(authz, impl.BanPublisherNode)(
				userCtx, drip.BanPublisherNodeRequestObject{PublisherId: publisherId, NodeId: nodeId})
			require.NoError(t, err, "should not ban node by non-admin")
			require.IsType(t, drip.BanPublisherNode403JSONResponse{}, res)
		})

		t.Run("Ban node by admin", func(t *testing.T) {
			// Attempt to ban the node as an admin user
			_, err := withMiddleware(authz, impl.BanPublisherNode)(
				adminCtx, drip.BanPublisherNodeRequestObject{PublisherId: publisherId, NodeId: nodeId})
			require.NoError(t, err)

			// Verify that the node is banned
			node, err := client.Node.Get(adminCtx, nodeId)
			require.NoError(t, err)
			assert.Equal(t, schema.NodeStatusBanned, node.Status, "should ban node")
		})

		t.Run("Calling endpoints with a banned node", func(t *testing.T) {
			// endpoints to test the authorization middleware
			testEndpoints := []struct {
				name   string
				invoke func(ctx context.Context) error
			}{
				{"GetNode", func(ctx context.Context) error {
					f := withMiddleware(authz, impl.GetNode)
					_, err := f(ctx, drip.GetNodeRequestObject{NodeId: nodeId})
					return err
				}},
				{"UpdateNode", func(ctx context.Context) error {
					f := withMiddleware(authz, impl.UpdateNode)
					_, err := f(ctx, drip.UpdateNodeRequestObject{PublisherId: publisherId, NodeId: nodeId})
					return err
				}},
				{"ListNodeVersions", func(ctx context.Context) error {
					f := withMiddleware(authz, impl.ListNodeVersions)
					_, err := f(ctx, drip.ListNodeVersionsRequestObject{NodeId: nodeId})
					return err
				}},
				{"PublishNodeVersion", func(ctx context.Context) error {
					f := withMiddleware(authz, impl.PublishNodeVersion)
					_, err := f(ctx, drip.PublishNodeVersionRequestObject{
						PublisherId: publisherId, NodeId: nodeId,
						Body: &drip.PublishNodeVersionJSONRequestBody{PersonalAccessToken: *pat},
					})
					return err
				}},
				{"InstallNode", func(ctx context.Context) error {
					f := withMiddleware(authz, impl.InstallNode)
					_, err := f(ctx, drip.InstallNodeRequestObject{NodeId: nodeId})
					return err
				}},
			}

			for _, tc := range testEndpoints {
				t.Run(tc.name, func(t *testing.T) {
					err := tc.invoke(userCtx)
					require.Error(t, err, "should return error")
					assertHTTPError(t, err, http.StatusForbidden)
				})
			}
		})

		t.Run("SearchNodes with banned node", func(t *testing.T) {
			// Step 1: Perform a search without including banned nodes (default behavior).
			f := withMiddleware(authz, impl.SearchNodes)
			res, err := f(userCtx, drip.SearchNodesRequestObject{
				Params: drip.SearchNodesParams{IncludeBanned: proto.Bool(false)}, // Explicitly do not include banned nodes.
			})
			require.NoError(t, err)
			require.IsType(t, drip.SearchNodes200JSONResponse{}, res)

			// Assert that no nodes are returned when IncludeBanned is false (since the node should be banned).
			searchResponse := res.(drip.SearchNodes200JSONResponse)
			require.Empty(t, searchResponse.Nodes, "Search should not include banned nodes")

			// Step 2: Perform a search including banned nodes.
			res, err = f(userCtx, drip.SearchNodesRequestObject{
				Params: drip.SearchNodesParams{IncludeBanned: proto.Bool(true)}, // Explicitly include banned nodes.
			})
			require.NoError(t, err)
			require.IsType(t, drip.SearchNodes200JSONResponse{}, res)

			// Assert that nodes are returned when IncludeBanned is true.
			searchResponse = res.(drip.SearchNodes200JSONResponse)
			require.NotEmpty(t, searchResponse.Nodes, "Search should include banned nodes")

			// Step 3: Assert that the banned node is included in the search result.
			foundBannedNode := false
			for _, node := range *searchResponse.Nodes {
				if *node.Id == nodeId {
					foundBannedNode = true
					break
				}
			}
			require.True(t, foundBannedNode, "The banned node should be present in the search results")
		})
	})
}
