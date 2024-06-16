package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"registry-backend/config"
	"registry-backend/drip"
	"registry-backend/ent/schema"
	"registry-backend/mock/gateways"
	"registry-backend/server/implementation"
	drip_authorization "registry-backend/server/middleware/authorization"
	"testing"

	strictecho "github.com/oapi-codegen/runtime/strictmiddleware/echo"
	"google.golang.org/protobuf/proto"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBan(t *testing.T) {
	clientCtx := context.Background()
	client, postgresContainer := setupDB(t, clientCtx)
	// Cleanup
	defer func() {
		if err := postgresContainer.Terminate(clientCtx); err != nil {
			log.Ctx(clientCtx).Error().Msgf("failed to terminate container: %s", err)
		}
	}()

	// Initialize the Service
	mockStorageService := new(gateways.MockStorageService)
	mockSlackService := new(gateways.MockSlackService)
	mockSlackService.
		On("SendRegistryMessageToSlack", mock.Anything).
		Return(nil) // Do nothing for all slack messsage calls.

	impl := implementation.NewStrictServerImplementation(
		client, &config.Config{}, mockStorageService, mockSlackService)

	errNotBanned := errors.New("passed authorization middleware")
	notBanned := func(f strictecho.StrictEchoHandlerFunc, operationID string) strictecho.StrictEchoHandlerFunc {
		return func(ctx echo.Context, request interface{}) (response interface{}, err error) {
			return nil, errNotBanned
		}
	}
	authorizationManager := drip_authorization.NewAuthorizationManager(client, impl.RegistryService)
	authz := authorizationManager.AuthorizationMiddleware()
	wrapped := drip.NewStrictHandler(impl, []strictecho.StrictEchoMiddlewareFunc{
		notBanned,
		authz,
	})

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
			_, err := impl.CreatePublisher(ctx, drip.CreatePublisherRequestObject{
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
			_, err = impl.CreateNode(ctx, drip.CreateNodeRequestObject{
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
				res, err := impl.BanPublisher(ctx, drip.BanPublisherRequestObject{PublisherId: publisherId})
				require.NoError(t, err, "should not ban publisher")
				require.IsType(t, drip.BanPublisher403JSONResponse{}, res)
			})

			t.Run("By Admin", func(t *testing.T) {
				ctx, admin := setUpTest(client)
				err = admin.Update().SetIsAdmin(true).Exec(clientCtx)
				require.NoError(t, err)
				_, err = impl.BanPublisher(ctx, drip.BanPublisherRequestObject{PublisherId: publisherId})
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
				name string
				req  func(ctx context.Context) *http.Request
				fn   func(ctx echo.Context) error
			}{
				{
					name: "CreatePublisher",
					req: func(ctx context.Context) (req *http.Request) {
						payloadBuf := new(bytes.Buffer)
						json.NewEncoder(payloadBuf).Encode(drip.CreatePublisherJSONRequestBody{})

						req = httptest.NewRequest(http.MethodPost, "/", payloadBuf).WithContext(ctx)
						req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
						return
					},
					fn: wrapped.CreatePublisher,
				},
				{
					name: "DeleteNodeVersion",
					req: func(ctx context.Context) (req *http.Request) {
						payloadBuf := new(bytes.Buffer)
						json.NewEncoder(payloadBuf).Encode(drip.DeleteNodeVersionRequestObject{})

						req = httptest.NewRequest(http.MethodPost, "/", payloadBuf).WithContext(ctx)
						req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
						return
					},
					fn: func(ctx echo.Context) error {
						return wrapped.DeleteNodeVersion(ctx, "", "", "")
					},
				},
			}

			e := echo.New()

			t.Run("Banned", func(t *testing.T) {
				ctxBanned, testUserBanned := setUpTest(client)
				err := testUserBanned.Update().SetStatus(schema.UserStatusTypeBanned).Exec(ctxBanned)
				require.NoError(t, err)
				for _, tt := range testtable {
					t.Run(tt.name, func(t *testing.T) {
						c := e.NewContext(tt.req(ctxBanned), httptest.NewRecorder())
						err = tt.fn(c)
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
						c := e.NewContext(tt.req(ctx), httptest.NewRecorder())
						err := tt.fn(c)
						assert.Equal(t, errNotBanned, err, "should pass the authorization middleware")
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
		_, err := impl.CreatePublisher(ctx, drip.CreatePublisherRequestObject{
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
		_, err = withMiddleware(authz, "CreateNode", impl.CreateNode)(ctx, drip.CreateNodeRequestObject{
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
		res, err := impl.CreatePersonalAccessToken(ctx, drip.CreatePersonalAccessTokenRequestObject{
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
				res, err := impl.BanPublisherNode(ctx, drip.BanPublisherNodeRequestObject{PublisherId: publisherId, NodeId: nodeId})
				require.NoError(t, err, "should not ban publisher node")
				require.IsType(t, drip.BanPublisherNode403JSONResponse{}, res)
			})

			t.Run("By Admin", func(t *testing.T) {
				ctx, admin := setUpTest(client)
				err = admin.Update().SetIsAdmin(true).Exec(clientCtx)
				require.NoError(t, err)
				_, err = impl.BanPublisherNode(ctx, drip.BanPublisherNodeRequestObject{PublisherId: publisherId, NodeId: nodeId})
				require.NoError(t, err)

				node, err := client.Node.Get(ctx, nodeId)
				require.NoError(t, err)
				assert.Equal(t, schema.NodeStatusBanned, node.Status, "should ban node")
			})
		})

		t.Run("Operate", func(t *testing.T) {
			t.Run("Get", func(t *testing.T) {
				f := withMiddleware(authz, "GetNode", impl.GetNode)
				_, err := f(ctx, drip.GetNodeRequestObject{NodeId: nodeId})
				require.Error(t, err)
				require.IsType(t, &echo.HTTPError{}, err)
				assert.Equal(t, err.(*echo.HTTPError).Code, http.StatusForbidden)
			})
			t.Run("Update", func(t *testing.T) {
				f := withMiddleware(authz, "UpdateNode", impl.UpdateNode)
				_, err := f(ctx, drip.UpdateNodeRequestObject{PublisherId: publisherId, NodeId: nodeId})
				require.Error(t, err)
				require.IsType(t, &echo.HTTPError{}, err)
				assert.Equal(t, err.(*echo.HTTPError).Code, http.StatusForbidden)
			})
			t.Run("ListNodeVersion", func(t *testing.T) {
				f := withMiddleware(authz, "ListNodeVersions", impl.ListNodeVersions)
				_, err := f(ctx, drip.ListNodeVersionsRequestObject{NodeId: nodeId})
				require.Error(t, err)
				require.IsType(t, &echo.HTTPError{}, err)
				assert.Equal(t, err.(*echo.HTTPError).Code, http.StatusForbidden)
			})
			t.Run("PublishNodeVersion", func(t *testing.T) {
				f := withMiddleware(authz, "PublishNodeVersion", impl.PublishNodeVersion)
				_, err := f(ctx, drip.PublishNodeVersionRequestObject{
					PublisherId: publisherId, NodeId: nodeId,
					Body: &drip.PublishNodeVersionJSONRequestBody{PersonalAccessToken: *pat},
				})
				require.Error(t, err)
				require.IsType(t, &echo.HTTPError{}, err)
				assert.Equal(t, err.(*echo.HTTPError).Code, http.StatusForbidden)
			})
			t.Run("InstallNode", func(t *testing.T) {
				f := withMiddleware(authz, "InstallNode", impl.InstallNode)
				_, err := f(ctx, drip.InstallNodeRequestObject{NodeId: nodeId})
				require.Error(t, err)
				require.IsType(t, &echo.HTTPError{}, err)
				assert.Equal(t, err.(*echo.HTTPError).Code, http.StatusForbidden)
			})
			t.Run("SearchNodes", func(t *testing.T) {
				f := withMiddleware(authz, "SearchNodes", impl.SearchNodes)
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
