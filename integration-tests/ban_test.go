package integration

import (
	"context"
	"net/http"
	"registry-backend/config"
	"registry-backend/drip"
	"registry-backend/ent/schema"
	"registry-backend/mock/gateways"
	"registry-backend/server/implementation"
	drip_authorization "registry-backend/server/middleware/authorization"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBan(t *testing.T) {
	clientCtx := context.Background()
	client, cleanup := setupDB(t, clientCtx)
	defer cleanup()

	// Initialize the Service
	mockStorageService := new(gateways.MockStorageService)
	mockSlackService := new(gateways.MockSlackService)
	mockDiscordService := new(gateways.MockDiscordService)
	mockSlackService.
		On("SendRegistryMessageToSlack", mock.Anything).
		Return(nil) // Do nothing for all slack messsage calls.
	mockAlgolia := new(gateways.MockAlgoliaService)
	mockAlgolia.
		On("IndexNodes", mock.Anything, mock.Anything).
		Return(nil)

	impl := implementation.NewStrictServerImplementation(
		client, &config.Config{}, mockStorageService, mockSlackService, mockDiscordService, mockAlgolia)

	authz := drip_authorization.NewAuthorizationManager(client, impl.RegistryService).AuthorizationMiddleware()

	t.Run("Publisher", func(t *testing.T) {
		t.Run("Ban", func(t *testing.T) {
			ctx, user := setUpTest(client)

			publisherId := "test-publisher"
			description := "test-description"
			source_code_repo := "test-source-code-repo"
			website := "test-website"
			support := "test-support"
			logo := "test-logo"
			name := "test-name"
			_, err := withMiddleware(authz, impl.CreatePublisher)(ctx, drip.CreatePublisherRequestObject{
				Body: &drip.Publisher{
					Id:             &publisherId,
					Description:    &description,
					SourceCodeRepo: &source_code_repo,
					Website:        &website,
					Support:        &support,
					Logo:           &logo,
					Name:           &name,
				},
			})
			require.NoError(t, err, "should return created publisher")

			nodeId := "test-node"
			nodeDescription := "test-node-description"
			nodeAuthor := "test-node-author"
			nodeLicense := "test-node-license"
			nodeName := "test-node-name"
			nodeTags := []string{"test-node-tag"}
			icon := "https://wwww.github.com/test-icon.svg"
			githubUrl := "https://www.github.com/test-github-url"
			_, err = withMiddleware(authz, impl.CreateNode)(ctx, drip.CreateNodeRequestObject{
				PublisherId: publisherId,
				Body: &drip.Node{
					Id:          &nodeId,
					Name:        &nodeName,
					Description: &nodeDescription,
					Author:      &nodeAuthor,
					License:     &nodeLicense,
					Tags:        &nodeTags,
					Icon:        &icon,
					Repository:  &githubUrl,
				},
			})
			require.NoError(t, err, "should return created node")

			t.Run("By Non Admin", func(t *testing.T) {
				ctx, _ := setUpTest(client)
				res, err := withMiddleware(authz, impl.BanPublisher)(ctx, drip.BanPublisherRequestObject{PublisherId: publisherId})
				require.NoError(t, err, "should not ban publisher")
				require.IsType(t, drip.BanPublisher403JSONResponse{}, res)
			})

			t.Run("By Admin", func(t *testing.T) {
				ctx, admin := setUpTest(client)
				err = admin.Update().SetIsAdmin(true).Exec(clientCtx)
				require.NoError(t, err)
				_, err = withMiddleware(authz, impl.BanPublisher)(ctx, drip.BanPublisherRequestObject{PublisherId: publisherId})
				require.NoError(t, err)

				pub, err := client.Publisher.Get(ctx, publisherId)
				require.NoError(t, err)
				assert.Equal(t, schema.PublisherStatusTypeBanned, pub.Status, "should ban publisher")
				user, err := client.User.Get(ctx, user.ID)
				require.NoError(t, err)
				assert.Equal(t, schema.UserStatusTypeBanned, user.Status, "should ban user")
				node, err := client.Node.Get(ctx, nodeId)
				require.NoError(t, err)
				assert.Equal(t, schema.NodeStatusBanned, node.Status, "should ban node")
			})
		})

		t.Run("Access", func(t *testing.T) {
			testtable := []struct {
				name   string
				invoke func(ctx context.Context) error
			}{
				{
					name: "CreatePublisher",
					invoke: func(ctx context.Context) error {
						_, err := withMiddleware(authz, impl.CreatePublisher)(ctx, drip.CreatePublisherRequestObject{Body: &drip.Publisher{}})
						return err
					},
				},
				{
					name: "DeleteNodeVersion",
					invoke: func(ctx context.Context) error {
						_, err := withMiddleware(authz, impl.DeleteNodeVersion)(ctx, drip.DeleteNodeVersionRequestObject{})
						return err
					},
				},
			}

			t.Run("Banned", func(t *testing.T) {
				ctxBanned, testUserBanned := setUpTest(client)
				err := testUserBanned.Update().SetStatus(schema.UserStatusTypeBanned).Exec(ctxBanned)
				require.NoError(t, err)
				for _, tt := range testtable {
					t.Run(tt.name, func(t *testing.T) {
						err = tt.invoke(ctxBanned)
						require.Error(t, err, "should return error")
						require.IsType(t, &echo.HTTPError{}, err, "should return echo http error")
						echoErr := err.(*echo.HTTPError)
						assert.Equal(t, http.StatusForbidden, echoErr.Code, "should return 403")
					})
				}
			})

			t.Run("Not Banned", func(t *testing.T) {
				ctx, _ := setUpTest(client)
				for _, tt := range testtable {
					t.Run(tt.name, func(t *testing.T) {
						err := tt.invoke(ctx)
						_, ok := err.(*echo.HTTPError)
						assert.False(t, ok, err, "should pass the authorization middleware")
					})
				}
			})
		})
	})

	t.Run("Node", func(t *testing.T) {
		ctx, _ := setUpTest(client)

		publisherId := "test-publisher-1"
		description := "test-description"
		source_code_repo := "test-source-code-repo"
		website := "test-website"
		support := "test-support"
		logo := "test-logo"
		name := "test-name"
		_, err := withMiddleware(authz, impl.CreatePublisher)(ctx, drip.CreatePublisherRequestObject{
			Body: &drip.Publisher{
				Id:             &publisherId,
				Description:    &description,
				SourceCodeRepo: &source_code_repo,
				Website:        &website,
				Support:        &support,
				Logo:           &logo,
				Name:           &name,
			},
		})
		require.NoError(t, err, "should return created publisher")

		nodeId := "test-node-1"
		nodeDescription := "test-node-description"
		nodeAuthor := "test-node-author"
		nodeLicense := "test-node-license"
		nodeName := "test-node-name"
		nodeTags := []string{"test-node-tag"}
		icon := "https://wwww.github.com/test-icon.svg"
		githubUrl := "https://www.github.com/test-github-url"
		_, err = withMiddleware(authz, impl.CreateNode)(ctx, drip.CreateNodeRequestObject{
			PublisherId: publisherId,
			Body: &drip.Node{
				Id:          &nodeId,
				Name:        &nodeName,
				Description: &nodeDescription,
				Author:      &nodeAuthor,
				License:     &nodeLicense,
				Tags:        &nodeTags,
				Icon:        &icon,
				Repository:  &githubUrl,
			},
		})
		require.NoError(t, err, "should return created node")

		tokenName := "name"
		tokenDescription := "name"
		res, err := withMiddleware(authz, impl.CreatePersonalAccessToken)(ctx, drip.CreatePersonalAccessTokenRequestObject{
			PublisherId: publisherId,
			Body: &drip.PersonalAccessToken{
				Name:        &tokenName,
				Description: &tokenDescription,
			},
		})
		require.NoError(t, err, "should return created token")
		require.IsType(t, drip.CreatePersonalAccessToken201JSONResponse{}, res)
		pat := res.(drip.CreatePersonalAccessToken201JSONResponse).Token

		t.Run("Ban", func(t *testing.T) {
			t.Run("By Non Admin", func(t *testing.T) {
				ctx, _ := setUpTest(client)
				res, err := withMiddleware(authz, impl.BanPublisherNode)(ctx, drip.BanPublisherNodeRequestObject{PublisherId: publisherId, NodeId: nodeId})
				require.NoError(t, err, "should not ban publisher node")
				require.IsType(t, drip.BanPublisherNode403JSONResponse{}, res)
			})

			t.Run("By Admin", func(t *testing.T) {
				ctx, admin := setUpTest(client)
				err = admin.Update().SetIsAdmin(true).Exec(clientCtx)
				require.NoError(t, err)
				_, err = withMiddleware(authz, impl.BanPublisherNode)(ctx, drip.BanPublisherNodeRequestObject{PublisherId: publisherId, NodeId: nodeId})
				require.NoError(t, err)

				node, err := client.Node.Get(ctx, nodeId)
				require.NoError(t, err)
				assert.Equal(t, schema.NodeStatusBanned, node.Status, "should ban node")
			})
		})

		t.Run("Operate", func(t *testing.T) {
			t.Run("Get", func(t *testing.T) {
				f := withMiddleware(authz, impl.GetNode)
				_, err := f(ctx, drip.GetNodeRequestObject{NodeId: nodeId})
				require.Error(t, err)
				require.IsType(t, &echo.HTTPError{}, err)
				assert.Equal(t, err.(*echo.HTTPError).Code, http.StatusForbidden)
			})
			t.Run("Update", func(t *testing.T) {
				f := withMiddleware(authz, impl.UpdateNode)
				_, err := f(ctx, drip.UpdateNodeRequestObject{PublisherId: publisherId, NodeId: nodeId})
				require.Error(t, err)
				require.IsType(t, &echo.HTTPError{}, err)
				assert.Equal(t, err.(*echo.HTTPError).Code, http.StatusForbidden)
			})
			t.Run("ListNodeVersion", func(t *testing.T) {
				f := withMiddleware(authz, impl.ListNodeVersions)
				_, err := f(ctx, drip.ListNodeVersionsRequestObject{NodeId: nodeId})
				require.Error(t, err)
				require.IsType(t, &echo.HTTPError{}, err)
				assert.Equal(t, err.(*echo.HTTPError).Code, http.StatusForbidden)
			})
			t.Run("PublishNodeVersion", func(t *testing.T) {
				f := withMiddleware(authz, impl.PublishNodeVersion)
				_, err := f(ctx, drip.PublishNodeVersionRequestObject{
					PublisherId: publisherId, NodeId: nodeId,
					Body: &drip.PublishNodeVersionJSONRequestBody{PersonalAccessToken: *pat},
				})
				require.Error(t, err)
				require.IsType(t, &echo.HTTPError{}, err)
				assert.Equal(t, err.(*echo.HTTPError).Code, http.StatusForbidden)
			})
			t.Run("InstallNode", func(t *testing.T) {
				f := withMiddleware(authz, impl.InstallNode)
				_, err := f(ctx, drip.InstallNodeRequestObject{NodeId: nodeId})
				require.Error(t, err)
				require.IsType(t, &echo.HTTPError{}, err)
				assert.Equal(t, err.(*echo.HTTPError).Code, http.StatusForbidden)
			})
			t.Run("SearchNodes", func(t *testing.T) {
				f := withMiddleware(authz, impl.SearchNodes)
				res, err := f(ctx, drip.SearchNodesRequestObject{
					Params: drip.SearchNodesParams{},
				})
				require.NoError(t, err)
				require.IsType(t, drip.SearchNodes200JSONResponse{}, res)
				require.Empty(t, res.(drip.SearchNodes200JSONResponse).Nodes)

				res, err = f(ctx, drip.SearchNodesRequestObject{
					Params: drip.SearchNodesParams{IncludeBanned: proto.Bool(true)},
				})
				require.NoError(t, err)
				require.IsType(t, drip.SearchNodes200JSONResponse{}, res)
				require.NotEmpty(t, res.(drip.SearchNodes200JSONResponse).Nodes)
			})
		})
	})

}
