package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"registry-backend/config"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/mock/gateways"
	"registry-backend/server/implementation"
	drip_authorization "registry-backend/server/middleware/authorization"
	dripservices_registry "registry-backend/services/registry"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func setUpTest(client *ent.Client) (context.Context, *ent.User) {
	ctx := context.Background()
	// create a User and attach to context
	testUser := createTestUser(ctx, client)
	ctx = decorateUserInContext(ctx, testUser)
	return ctx, testUser
}

func setUpAdminTest(client *ent.Client) (context.Context, *ent.User) {
	ctx := context.Background()
	testUser := createAdminUser(ctx, client)
	ctx = decorateUserInContext(ctx, testUser)
	return ctx, testUser
}

func TestRegistry(t *testing.T) {
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
	mockDiscordService := new(gateways.MockDiscordService)
	mockSlackService.
		On("SendRegistryMessageToSlack", mock.Anything).
		Return(nil) // Do nothing for all slack messsage calls.
	mockAlgolia := new(gateways.MockAlgoliaService)
	mockAlgolia.
		On("IndexNodes", mock.Anything, mock.Anything).
		Return(nil).
		On("DeleteNode", mock.Anything, mock.Anything).
		Return(nil)
	impl := implementation.NewStrictServerImplementation(
		client, &config.Config{}, mockStorageService, mockSlackService, mockDiscordService, mockAlgolia)
	authz := drip_authorization.NewAuthorizationManager(client, impl.RegistryService).
		AuthorizationMiddleware()

	t.Run("Publisher", func(t *testing.T) {
		ctx, testUser := setUpTest(client)
		publisherId := "test-publisher"
		description := "test-description"
		source_code_repo := "test-source-code-repo"
		website := "test-website"
		support := "test-support"
		logo := "test-logo"
		name := "test-name"

		t.Run("Create Publisher", func(t *testing.T) {
			createPublisherResponse, err := withMiddleware(authz, impl.CreatePublisher)(ctx, drip.CreatePublisherRequestObject{
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
			require.NotNil(t, createPublisherResponse, "should return created publisher")
			assert.Equal(t, publisherId, *createPublisherResponse.(drip.CreatePublisher201JSONResponse).Id)
			assert.Equal(t, description, *createPublisherResponse.(drip.CreatePublisher201JSONResponse).Description)
			assert.Equal(t, source_code_repo, *createPublisherResponse.(drip.CreatePublisher201JSONResponse).SourceCodeRepo)
			assert.Equal(t, website, *createPublisherResponse.(drip.CreatePublisher201JSONResponse).Website)
			assert.Equal(t, support, *createPublisherResponse.(drip.CreatePublisher201JSONResponse).Support)
			assert.Equal(t, logo, *createPublisherResponse.(drip.CreatePublisher201JSONResponse).Logo)
		})

		t.Run("Validate Publisher", func(t *testing.T) {
			res, err := withMiddleware(authz, impl.ValidatePublisher)(ctx, drip.ValidatePublisherRequestObject{
				Params: drip.ValidatePublisherParams{Username: name},
			})
			require.NoError(t, err, "should not return error")
			require.IsType(t, drip.ValidatePublisher200JSONResponse{}, res, "should return 200")
			require.True(t, *res.(drip.ValidatePublisher200JSONResponse).IsAvailable, "should be available")
		})

		t.Run("Get Publisher", func(t *testing.T) {
			getPublisherResponse, err := withMiddleware(authz, impl.GetPublisher)(ctx, drip.GetPublisherRequestObject{
				PublisherId: publisherId})
			require.NoError(t, err, "should return created publisher")
			assert.Equal(t, publisherId, *getPublisherResponse.(drip.GetPublisher200JSONResponse).Id)
			assert.Equal(t, description, *getPublisherResponse.(drip.GetPublisher200JSONResponse).Description)
			assert.Equal(t, source_code_repo, *getPublisherResponse.(drip.GetPublisher200JSONResponse).SourceCodeRepo)
			assert.Equal(t, website, *getPublisherResponse.(drip.GetPublisher200JSONResponse).Website)
			assert.Equal(t, support, *getPublisherResponse.(drip.GetPublisher200JSONResponse).Support)
			assert.Equal(t, logo, *getPublisherResponse.(drip.GetPublisher200JSONResponse).Logo)
			assert.Equal(t, name, *getPublisherResponse.(drip.GetPublisher200JSONResponse).Name)

			// Check the number of members returned
			expectedMembersCount := 1 // Adjust to your expected count
			assert.Equal(t, expectedMembersCount,
				len(*getPublisherResponse.(drip.GetPublisher200JSONResponse).Members),
				"should return the correct number of members")

			// Check specific properties of each member, adjust indices accordingly
			for _, member := range *getPublisherResponse.(drip.GetPublisher200JSONResponse).Members {
				expectedUserId := testUser.ID
				expectedUserName := testUser.Name
				expectedUserEmail := testUser.Email

				assert.Equal(t, expectedUserId, *member.User.Id, "User ID should match")
				assert.Equal(t, expectedUserName, *member.User.Name, "User name should match")
				assert.Equal(t, expectedUserEmail, *member.User.Email, "User email should match")
			}
		})

		t.Run("List Publishers", func(t *testing.T) {
			res, err := withMiddleware(authz, impl.ListPublishers)(ctx, drip.ListPublishersRequestObject{})
			require.NoError(t, err, "should not return error")
			require.IsType(t, drip.ListPublishers200JSONResponse{}, res, "should return 200 status code")
			res200 := res.(drip.ListPublishers200JSONResponse)
			require.Len(t, res200, 1, "should return all stored publlishers")
			assert.Equal(t, drip.Publisher{
				Id:             &publisherId,
				Description:    &description,
				SourceCodeRepo: &source_code_repo,
				Website:        &website,
				Support:        &support,
				Logo:           &logo,
				Name:           &name,

				// generated thus ignored in comparison
				Members:   res200[0].Members,
				CreatedAt: res200[0].CreatedAt,
				Status:    res200[0].Status,
			}, res200[0], "should return correct publishers")
		})

		t.Run("Get Non-Exist Publisher", func(t *testing.T) {
			res, err := withMiddleware(authz, impl.GetPublisher)(ctx, drip.GetPublisherRequestObject{PublisherId: publisherId + "invalid"})
			require.NoError(t, err, "should not return error")
			assert.IsType(t, drip.GetPublisher404JSONResponse{}, res)
		})

		t.Run("Update Publisher", func(t *testing.T) {
			update_description := "update-test-description"
			update_source_code_repo := "update-test-source-code-repo"
			update_website := "update-test-website"
			update_support := "update-test-support"
			update_logo := "update-test-logo"
			update_name := "update-test-name"

			updatePublisherResponse, err := withMiddleware(authz, impl.UpdatePublisher)(ctx, drip.UpdatePublisherRequestObject{
				PublisherId: publisherId,
				Body: &drip.Publisher{
					Description:    &update_description,
					SourceCodeRepo: &update_source_code_repo,
					Website:        &update_website,
					Support:        &update_support,
					Logo:           &update_logo,
					Name:           &update_name,
				},
			})
			require.NoError(t, err, "should return created publisher")
			assert.Equal(t, publisherId, *updatePublisherResponse.(drip.UpdatePublisher200JSONResponse).Id)
			assert.Equal(t, update_description, *updatePublisherResponse.(drip.UpdatePublisher200JSONResponse).Description)
			assert.Equal(t, update_source_code_repo, *updatePublisherResponse.(drip.UpdatePublisher200JSONResponse).SourceCodeRepo)
			assert.Equal(t, update_website, *updatePublisherResponse.(drip.UpdatePublisher200JSONResponse).Website)
			assert.Equal(t, update_support, *updatePublisherResponse.(drip.UpdatePublisher200JSONResponse).Support)
			assert.Equal(t, update_logo, *updatePublisherResponse.(drip.UpdatePublisher200JSONResponse).Logo)

			_, err = withMiddleware(authz, impl.ListPublishersForUser)(ctx, drip.ListPublishersForUserRequestObject{})
			require.NoError(t, err, "should return created publisher")
		})

		t.Run("Reject New Publisher With The Same Name", func(t *testing.T) {
			duplicateCreatePublisherResponse, err := withMiddleware(authz, impl.CreatePublisher)(ctx, drip.CreatePublisherRequestObject{
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
			require.NoError(t, err, "should return error")
			assert.IsType(t, drip.CreatePublisher400JSONResponse{}, duplicateCreatePublisherResponse)
		})

		t.Run("Delete Publisher", func(t *testing.T) {
			res, err := withMiddleware(authz, impl.DeletePublisher)(ctx, drip.DeletePublisherRequestObject{PublisherId: publisherId})
			require.NoError(t, err, "should not return error")
			require.IsType(t, drip.DeletePublisher204Response{}, res, "should return 204")
		})
	})

	t.Run("Personal Access Token", func(t *testing.T) {
		ctx, _ := setUpTest(client)
		publisherId := "test-publisher-pat"
		description := "test-description"
		source_code_repo := "test-source-code-repo"
		website := "test-website"
		support := "test-support"
		logo := "test-logo"
		name := "test-name"
		tokenName := "test-token-name"
		tokenDescription := "test-token-description"

		t.Run("Create Publisher", func(t *testing.T) {
			createPublisherResponse, err := withMiddleware(authz, impl.CreatePublisher)(ctx, drip.CreatePublisherRequestObject{
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
			require.NotNil(t, createPublisherResponse, "should return created publisher")
		})

		t.Run("List Personal Access Token Before Create", func(t *testing.T) {
			none, err := withMiddleware(authz, impl.ListPersonalAccessTokens)(ctx, drip.ListPersonalAccessTokensRequestObject{
				PublisherId: publisherId,
			})
			require.NoError(t, err, "should not return error")
			require.IsType(t, drip.ListPersonalAccessTokens200JSONResponse{}, none, "should return 200")
			assert.Empty(t, none.(drip.ListPersonalAccessTokens200JSONResponse))
		})

		t.Run("Create Personal Acccess Token", func(t *testing.T) {
			createPersonalAccessTokenResponse, err := withMiddleware(authz, impl.CreatePersonalAccessToken)(
				ctx, drip.CreatePersonalAccessTokenRequestObject{
					PublisherId: publisherId,
					Body: &drip.PersonalAccessToken{
						Name:        &tokenName,
						Description: &tokenDescription,
					},
				})
			require.NoError(t, err, "should return created token")
			require.NotNil(t,
				*createPersonalAccessTokenResponse.(drip.CreatePersonalAccessToken201JSONResponse).Token,
				"Token should have a value.")
		})

		t.Run("List Personal Access Token", func(t *testing.T) {
			getPersonalAccessTokenResponse, err := withMiddleware(authz, impl.ListPersonalAccessTokens)(ctx, drip.ListPersonalAccessTokensRequestObject{
				PublisherId: publisherId,
			})
			require.NoError(t, err, "should return created token")
			assert.Equal(t, tokenName,
				*getPersonalAccessTokenResponse.(drip.ListPersonalAccessTokens200JSONResponse)[0].Name)
			assert.Equal(t, tokenDescription,
				*getPersonalAccessTokenResponse.(drip.ListPersonalAccessTokens200JSONResponse)[0].Description)
			assert.True(t,
				isTokenMasked(*getPersonalAccessTokenResponse.(drip.ListPersonalAccessTokens200JSONResponse)[0].Token))
		})
	})

	t.Run("Node", func(t *testing.T) {
		ctx, _ := setUpTest(client)
		publisherId := "test-publisher-node"
		description := "test-description"
		sourceCodeRepo := "test-source-code-repo"
		website := "test-website"
		support := "test-support"
		logo := "test-logo"
		name := "test-name"

		createPublisherResponse, err := withMiddleware(authz, impl.CreatePublisher)(ctx, drip.CreatePublisherRequestObject{
			Body: &drip.Publisher{
				Id:             &publisherId,
				Description:    &description,
				SourceCodeRepo: &sourceCodeRepo,
				Website:        &website,
				Support:        &support,
				Logo:           &logo,
				Name:           &name,
			},
		})
		require.NoError(t, err, "should return created publisher")
		require.NotNil(t, createPublisherResponse, "should return created publisher")

		nodeId := "test-node"
		nodeDescription := "test-node-description"
		nodeAuthor := "test-node-author"
		nodeLicense := "test-node-license"
		nodeName := "test-node-name"
		nodeTags := []string{"test-node-tag"}
		icon := "https://wwww.github.com/test-icon.svg"
		githubUrl := "https://www.github.com/test-github-url"

		var real_node_id *string
		t.Run("Create Node", func(t *testing.T) {
			createNodeResponse, err := withMiddleware(authz, impl.CreateNode)(ctx, drip.CreateNodeRequestObject{
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
			require.NotNil(t, createNodeResponse, "should return created node")
			assert.Equal(t, nodeId, *createNodeResponse.(drip.CreateNode201JSONResponse).Id)
			assert.Equal(t, nodeDescription, *createNodeResponse.(drip.CreateNode201JSONResponse).Description)
			assert.Equal(t, nodeAuthor, *createNodeResponse.(drip.CreateNode201JSONResponse).Author)
			assert.Equal(t, nodeLicense, *createNodeResponse.(drip.CreateNode201JSONResponse).License)
			assert.Equal(t, nodeName, *createNodeResponse.(drip.CreateNode201JSONResponse).Name)
			assert.Equal(t, nodeTags, *createNodeResponse.(drip.CreateNode201JSONResponse).Tags)
			assert.Equal(t, icon, *createNodeResponse.(drip.CreateNode201JSONResponse).Icon)
			assert.Equal(t, githubUrl, *createNodeResponse.(drip.CreateNode201JSONResponse).Repository)
			assert.Equal(t, drip.NodeStatusActive, *createNodeResponse.(drip.CreateNode201JSONResponse).Status)
			real_node_id = createNodeResponse.(drip.CreateNode201JSONResponse).Id

		})

		t.Run("Get Node", func(t *testing.T) {
			res, err := withMiddleware(authz, impl.GetNode)(ctx, drip.GetNodeRequestObject{NodeId: nodeId})
			require.NoError(t, err, "should not return error")
			require.IsType(t, drip.GetNode200JSONResponse{}, res)
			res200 := res.(drip.GetNode200JSONResponse)
			expDl, expRate := 0, float32(0)
			nodeStatus := drip.NodeStatusActive
			assert.Equal(t, drip.GetNode200JSONResponse{
				Id:          &nodeId,
				Name:        &nodeName,
				Description: &nodeDescription,
				Author:      &nodeAuthor,
				Tags:        &nodeTags,
				License:     &nodeLicense,
				Icon:        &icon,
				Repository:  &githubUrl,

				Downloads:    &expDl,
				Rating:       &expRate,
				Status:       &nodeStatus,
				StatusDetail: proto.String(""),
				Category:     proto.String(""),
			}, res200, "should return stored node data")
		})

		t.Run("Get Not Exist Node", func(t *testing.T) {
			res, err := withMiddleware(authz, impl.GetNode)(ctx, drip.GetNodeRequestObject{NodeId: nodeId + "fake"})
			require.NoError(t, err, "should not return error")
			require.IsType(t, drip.GetNode404JSONResponse{}, res)
		})

		t.Run("Get Publisher Nodes", func(t *testing.T) {
			res, err := withMiddleware(authz, impl.ListNodesForPublisher)(ctx, drip.ListNodesForPublisherRequestObject{
				PublisherId: publisherId,
			})
			require.NoError(t, err, "should not return error")
			require.IsType(t, drip.ListNodesForPublisher200JSONResponse{}, res)
			res200 := res.(drip.ListNodesForPublisher200JSONResponse)
			require.Len(t, res200, 1)
			expDl, expRate := 0, float32(0)
			nodeStatus := drip.NodeStatusActive
			assert.Equal(t, drip.Node{
				Id:          &nodeId,
				Name:        &nodeName,
				Description: &nodeDescription,
				Author:      &nodeAuthor,
				Tags:        &nodeTags,
				License:     &nodeLicense,
				Icon:        &icon,
				Repository:  &githubUrl,

				Downloads:    &expDl,
				Rating:       &expRate,
				Status:       &nodeStatus,
				StatusDetail: proto.String(""),
				Category:     proto.String(""),
			}, res200[0], "should return stored node data")
		})

		t.Run("Update Node", func(t *testing.T) {
			updateNodeDescription := "update_test-node-description"
			updateNodeAuthor := "update_test-node-author"
			updateNodeLicense := "update_test-node-license"
			updateNodeName := "update_test-node-name"
			updateNodeTags := []string{"update-test-node-tag"}
			updateIcon := "https://wwww.github.com/update-icon.svg"
			updateGithubUrl := "https://www.github.com/update-github-url"

			updateNodeResponse, err := withMiddleware(authz, impl.UpdateNode)(ctx, drip.UpdateNodeRequestObject{
				PublisherId: publisherId,
				NodeId:      *real_node_id,
				Body: &drip.Node{
					Id:          &nodeId,
					Description: &updateNodeDescription,
					Author:      &updateNodeAuthor,
					License:     &updateNodeLicense,
					Name:        &updateNodeName,
					Tags:        &updateNodeTags,
					Icon:        &updateIcon,
					Repository:  &updateGithubUrl,
				},
			})
			require.NoError(t, err, "should return created node")
			assert.Equal(t, nodeId, *updateNodeResponse.(drip.UpdateNode200JSONResponse).Id)
			assert.Equal(t, updateNodeDescription, *updateNodeResponse.(drip.UpdateNode200JSONResponse).Description)
			assert.Equal(t, updateNodeAuthor, *updateNodeResponse.(drip.UpdateNode200JSONResponse).Author)
			assert.Equal(t, updateNodeLicense, *updateNodeResponse.(drip.UpdateNode200JSONResponse).License)
			assert.Equal(t, updateNodeName, *updateNodeResponse.(drip.UpdateNode200JSONResponse).Name)
			assert.Equal(t, updateNodeTags, *updateNodeResponse.(drip.UpdateNode200JSONResponse).Tags)
			assert.Equal(t, updateIcon, *updateNodeResponse.(drip.UpdateNode200JSONResponse).Icon)
			assert.Equal(t, updateGithubUrl, *updateNodeResponse.(drip.UpdateNode200JSONResponse).Repository)

			resUpdated, err := withMiddleware(authz, impl.GetNode)(ctx, drip.GetNodeRequestObject{NodeId: nodeId})
			require.NoError(t, err, "should not return error")
			require.IsType(t, drip.GetNode200JSONResponse{}, resUpdated)
			res200Updated := resUpdated.(drip.GetNode200JSONResponse)
			expDl, expRate := 0, float32(0)
			nodeStatus := drip.NodeStatusActive
			assert.Equal(t, drip.GetNode200JSONResponse{
				Id:          &nodeId,
				Description: &updateNodeDescription,
				Author:      &updateNodeAuthor,
				License:     &updateNodeLicense,
				Name:        &updateNodeName,
				Tags:        &updateNodeTags,
				Icon:        &updateIcon,
				Repository:  &updateGithubUrl,

				Downloads:    &expDl,
				Rating:       &expRate,
				Status:       &nodeStatus,
				StatusDetail: proto.String(""),
				Category:     proto.String(""),
			}, res200Updated, "should return updated node data")
		})

		t.Run("Update Not Exist Node", func(t *testing.T) {
			res, err := withMiddleware(authz, impl.UpdateNode)(ctx, drip.UpdateNodeRequestObject{PublisherId: publisherId, NodeId: nodeId + "fake", Body: &drip.Node{}})
			require.NoError(t, err, "should not return error")
			require.IsType(t, drip.UpdateNode404JSONResponse{}, res)
		})

		t.Run("Delete Node", func(t *testing.T) {
			res, err := withMiddleware(authz, impl.DeleteNode)(ctx, drip.DeleteNodeRequestObject{PublisherId: publisherId, NodeId: nodeId})
			require.NoError(t, err, "should not return error")
			assert.IsType(t, drip.DeleteNode204Response{}, res)
		})
		t.Run("Delete Not Exist Node", func(t *testing.T) {
			res, err := withMiddleware(authz, impl.DeleteNode)(ctx, drip.DeleteNodeRequestObject{PublisherId: publisherId, NodeId: nodeId + "fake"})
			require.NoError(t, err, "should not return error")
			assert.IsType(t, drip.DeleteNode204Response{}, res)
		})
	})

	t.Run("Node Version", func(t *testing.T) {
		ctx, _ := setUpTest(client)
		publisherId := "test-publisher-node-version"
		description := "test-description"
		source_code_repo := "test-source-code-repo"
		website := "test-website"
		support := "test-support"
		logo := "test-logo"
		name := "test-name"

		createPublisherResponse, err := withMiddleware(authz, impl.CreatePublisher)(ctx, drip.CreatePublisherRequestObject{
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
		require.NotNil(t, createPublisherResponse, "should return created publisher")
		assert.Equal(t, publisherId, *createPublisherResponse.(drip.CreatePublisher201JSONResponse).Id)

		tokenName := "test-token-name"
		tokenDescription := "test-token-description"
		createPersonalAccessTokenResponse, err := withMiddleware(authz, impl.CreatePersonalAccessToken)(ctx, drip.CreatePersonalAccessTokenRequestObject{
			PublisherId: publisherId,
			Body: &drip.PersonalAccessToken{
				Name:        &tokenName,
				Description: &tokenDescription,
			},
		})
		require.NoError(t, err, "should return created token")
		require.NotNil(t, *createPersonalAccessTokenResponse.(drip.CreatePersonalAccessToken201JSONResponse).Token, "Token should have a value.")

		nodeId := "test-node1"
		nodeDescription := "test-node-description"
		nodeAuthor := "test-node-author"
		nodeLicense := "test-node-license"
		nodeName := "test-node-name"
		nodeTags := []string{"test-node-tag"}
		nodeVersionLiteral := "1.0.0"
		changelog := "test-changelog"
		dependencies := []string{"test-dependency"}
		downloadUrl := "https://storage.googleapis.com/comfy-registry/test-publisher-node-version/test-node1/1.0.0/node.tar.gz"

		createdPublisher := createPublisherResponse.(drip.CreatePublisher201JSONResponse)
		var createdNodeVersion drip.NodeVersion

		t.Run("List Node Version Before Create", func(t *testing.T) {
			resVersions, err := withMiddleware(authz, impl.ListNodeVersions)(ctx, drip.ListNodeVersionsRequestObject{NodeId: nodeId})
			require.NoError(t, err, "should return error since node version doesn't exists")
			require.IsType(t, drip.ListNodeVersions200JSONResponse{}, resVersions)
			assert.Empty(t, resVersions.(drip.ListNodeVersions200JSONResponse), "should not return any node versions")
		})

		t.Run("Create Node Version with Fake Token", func(t *testing.T) {
			_, err := withMiddleware(authz, impl.PublishNodeVersion)(ctx, drip.PublishNodeVersionRequestObject{
				PublisherId: publisherId,
				NodeId:      nodeId,
				Body: &drip.PublishNodeVersionJSONRequestBody{
					Node: drip.Node{
						Id:          &nodeId,
						Description: &nodeDescription,
						Author:      &nodeAuthor,
						License:     &nodeLicense,
						Name:        &nodeName,
						Tags:        &nodeTags,
						Repository:  &source_code_repo,
					},
					NodeVersion: drip.NodeVersion{
						Version:      &nodeVersionLiteral,
						Changelog:    &changelog,
						Dependencies: &dependencies,
					},
					PersonalAccessToken: "faketoken",
				},
			})
			require.Error(t, err)
			assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code, "should return 400 bad request")
		})

		t.Run("Create Node Version", func(t *testing.T) {
			mockStorageService.On("GenerateSignedURL", mock.Anything, mock.Anything).Return("test-url", nil)
			mockStorageService.On("GetFileUrl", mock.Anything, mock.Anything, mock.Anything).Return("test-url", nil)
			mockDiscordService.On("SendSecurityCouncilMessage", mock.Anything, mock.Anything).Return(nil)
			createNodeVersionResp, err := withMiddleware(authz, impl.PublishNodeVersion)(ctx, drip.PublishNodeVersionRequestObject{
				PublisherId: publisherId,
				NodeId:      nodeId,
				Body: &drip.PublishNodeVersionJSONRequestBody{
					Node: drip.Node{
						Id:          &nodeId,
						Description: &nodeDescription,
						Author:      &nodeAuthor,
						License:     &nodeLicense,
						Name:        &nodeName,
						Tags:        &nodeTags,
						Repository:  &source_code_repo,
					},
					NodeVersion: drip.NodeVersion{
						Version:      &nodeVersionLiteral,
						Changelog:    &changelog,
						Dependencies: &dependencies,
					},
					PersonalAccessToken: *createPersonalAccessTokenResponse.(drip.CreatePersonalAccessToken201JSONResponse).Token,
				},
			})
			require.NoError(t, err, "should return created node version")
			assert.Equal(t, nodeVersionLiteral, *createNodeVersionResp.(drip.PublishNodeVersion201JSONResponse).NodeVersion.Version)
			require.Equal(t, "test-url", *createNodeVersionResp.(drip.PublishNodeVersion201JSONResponse).SignedUrl, "should return signed url")
			require.Equal(t, dependencies, *createNodeVersionResp.(drip.PublishNodeVersion201JSONResponse).NodeVersion.Dependencies, "should return pip dependencies")
			require.Equal(t, changelog, *createNodeVersionResp.(drip.PublishNodeVersion201JSONResponse).NodeVersion.Changelog, "should return changelog")
			versionStatus := drip.NodeVersionStatusPending
			require.Equal(t, versionStatus, *createNodeVersionResp.(drip.PublishNodeVersion201JSONResponse).NodeVersion.Status, "should return pending status")
			createdNodeVersion = *createNodeVersionResp.(drip.PublishNodeVersion201JSONResponse).NodeVersion // Needed for downstream tests.

			adminCtx, _ := setUpAdminTest(client)
			activeStatus := drip.NodeVersionStatusActive
			adminUpdateNodeVersionResp, err := impl.AdminUpdateNodeVersion(adminCtx, drip.AdminUpdateNodeVersionRequestObject{
				NodeId:        nodeId,
				VersionNumber: *createdNodeVersion.Version,
				Body: &drip.AdminUpdateNodeVersionJSONRequestBody{
					Status: &activeStatus,
				},
			})
			require.NoError(t, err, "should return updated node version")
			assert.Equal(t, activeStatus, *adminUpdateNodeVersionResp.(drip.AdminUpdateNodeVersion200JSONResponse).Status)
		})

		t.Run("Get not exist Node Version ", func(t *testing.T) {
			res, err := withMiddleware(authz, impl.GetNodeVersion)(ctx, drip.GetNodeVersionRequestObject{NodeId: nodeId + "fake", VersionId: nodeVersionLiteral})
			require.NoError(t, err, "should not return error")
			assert.IsType(t, drip.GetNodeVersion404JSONResponse{}, res)
		})

		t.Run("Create Node Version of Not Exist Node", func(t *testing.T) {
			_, err := withMiddleware(authz, impl.PublishNodeVersion)(ctx, drip.PublishNodeVersionRequestObject{
				PublisherId: publisherId,
				NodeId:      nodeId + "fake",
				Body:        &drip.PublishNodeVersionJSONRequestBody{},
			})
			require.Error(t, err)
			assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code, "should return 400 bad request")
		})

		t.Run("List Node Versions", func(t *testing.T) {
			resVersions, err := withMiddleware(authz, impl.ListNodeVersions)(ctx, drip.ListNodeVersionsRequestObject{NodeId: nodeId})
			require.NoError(t, err, "should not return error")
			require.IsType(t, drip.ListNodeVersions200JSONResponse{}, resVersions, "should return 200")
			resVersions200 := resVersions.(drip.ListNodeVersions200JSONResponse)
			require.Len(t, resVersions200, 1, "should return only one version")
			nodeVersionStatus := drip.NodeVersionStatusActive
			println("Download URL: ", *resVersions200[0].DownloadUrl)
			println("Download URL: ", downloadUrl)
			println("Status: ", *resVersions200[0].Status)
			println("Status: ", nodeVersionStatus)
			assert.Equal(t, drip.NodeVersion{
				// generated attribute
				Id:        resVersions200[0].Id,
				CreatedAt: resVersions200[0].CreatedAt,

				Deprecated:   proto.Bool(false),
				Version:      &nodeVersionLiteral,
				Changelog:    &changelog,
				Dependencies: &dependencies,
				DownloadUrl:  &downloadUrl,
				Status:       &nodeVersionStatus,
				StatusReason: proto.String(""),
			}, resVersions200[0], "should be equal")
		})

		t.Run("Update Node Version", func(t *testing.T) {
			updatedChangelog := "test-changelog-2"
			resUNV, err := withMiddleware(authz, impl.UpdateNodeVersion)(ctx, drip.UpdateNodeVersionRequestObject{
				PublisherId: publisherId,
				NodeId:      nodeId,
				VersionId:   *createdNodeVersion.Id,
				Body: &drip.NodeVersionUpdateRequest{
					Changelog:  &updatedChangelog,
					Deprecated: proto.Bool(true),
				},
			})
			require.NoError(t, err, "should not return error")
			require.IsType(t, drip.UpdateNodeVersion200JSONResponse{}, resUNV, "should return 200")

			res, err := withMiddleware(authz, impl.ListNodeVersions)(ctx, drip.ListNodeVersionsRequestObject{NodeId: nodeId})
			require.NoError(t, err, "should not return error")
			require.IsType(t, drip.ListNodeVersions200JSONResponse{}, res, "should return 200")
			res200 := res.(drip.ListNodeVersions200JSONResponse)
			require.Len(t, res200, 1, "should return only one version")
			status := drip.NodeVersionStatusActive
			updatedNodeVersion := drip.NodeVersion{
				// generated attribute
				Id:        res200[0].Id,
				CreatedAt: res200[0].CreatedAt,

				Deprecated:   proto.Bool(true),
				Version:      &nodeVersionLiteral,
				Changelog:    &updatedChangelog,
				Dependencies: &dependencies,
				DownloadUrl:  &downloadUrl,
				Status:       &status,
				StatusReason: proto.String(""),
			}
			assert.Equal(t, updatedNodeVersion, res200[0], "should be equal")
			createdNodeVersion = res200[0]
		})

		t.Run("List Nodes", func(t *testing.T) {
			resNodes, err := withMiddleware(authz, impl.ListAllNodes)(ctx, drip.ListAllNodesRequestObject{})
			require.NoError(t, err, "should not return error")
			require.IsType(t, drip.ListAllNodes200JSONResponse{}, resNodes, "should return 200 server response")
			resNodes200 := resNodes.(drip.ListAllNodes200JSONResponse)
			assert.Len(t, *resNodes200.Nodes, 1, "should only contain 1 node")

			expDl, expRate := 0, float32(0)
			nodeStatus := drip.NodeStatusActive
			expectedNode := drip.Node{
				Id:            &nodeId,
				Name:          &nodeName,
				Repository:    &source_code_repo,
				Description:   &nodeDescription,
				Author:        &nodeAuthor,
				License:       &nodeLicense,
				Tags:          &nodeTags,
				LatestVersion: &createdNodeVersion,
				Icon:          proto.String(""),
				Publisher:     (*drip.Publisher)(&createdPublisher),
				Downloads:     &expDl,
				Rating:        &expRate,
				Status:        &nodeStatus,
				StatusDetail:  proto.String(""),
				Category:      proto.String(""),
			}
			expectedNode.LatestVersion.DownloadUrl = (*resNodes200.Nodes)[0].LatestVersion.DownloadUrl // generated
			expectedNode.LatestVersion.Deprecated = (*resNodes200.Nodes)[0].LatestVersion.Deprecated   // generated
			expectedNode.Publisher.CreatedAt = (*resNodes200.Nodes)[0].Publisher.CreatedAt
			assert.Equal(t, expectedNode, (*resNodes200.Nodes)[0])
		})

		t.Run("Install Node", func(t *testing.T) {
			resIns, err := withMiddleware(authz, impl.InstallNode)(ctx, drip.InstallNodeRequestObject{NodeId: nodeId})
			require.NoError(t, err, "should not return error")
			require.IsType(t, drip.InstallNode200JSONResponse{}, resIns, "should return 200")

			resIns, err = withMiddleware(authz, impl.InstallNode)(ctx, drip.InstallNodeRequestObject{
				NodeId: nodeId, Params: drip.InstallNodeParams{Version: &nodeVersionLiteral}})
			require.NoError(t, err, "should not return error")
			require.IsType(t, drip.InstallNode200JSONResponse{}, resIns, "should return 200")
		})

		t.Run("Install Node Version on not exist node or version", func(t *testing.T) {
			resIns, err := withMiddleware(authz, impl.InstallNode)(ctx, drip.InstallNodeRequestObject{NodeId: nodeId + "fake"})
			require.NoError(t, err, "should not return error")
			require.IsType(t, drip.InstallNode404JSONResponse{}, resIns, "should return 404")
			resIns, err = withMiddleware(authz, impl.InstallNode)(ctx, drip.InstallNodeRequestObject{
				NodeId: nodeId, Params: drip.InstallNodeParams{Version: proto.String(nodeVersionLiteral + "fake")}})
			require.NoError(t, err, "should not return error")
			require.IsType(t, drip.InstallNode404JSONResponse{}, resIns, "should return 404")
		})

		t.Run("Get Total Install", func(t *testing.T) {
			res, err := withMiddleware(authz, impl.GetNode)(ctx, drip.GetNodeRequestObject{
				NodeId: nodeId,
			})
			require.NoError(t, err, "should not return error")
			require.IsType(t, drip.GetNode200JSONResponse{}, res)
			assert.Equal(t, int(2), *res.(drip.GetNode200JSONResponse).Downloads)
		})

		t.Run("Add review", func(t *testing.T) {
			res, err := withMiddleware(authz, impl.PostNodeReview)(ctx, drip.PostNodeReviewRequestObject{
				NodeId: nodeId,
				Params: drip.PostNodeReviewParams{Star: 5},
			})
			require.NoError(t, err)
			require.IsType(t, drip.PostNodeReview200JSONResponse{}, res)
			res200 := res.(drip.PostNodeReview200JSONResponse)
			assert.Equal(t, float32(5), *res200.Rating)
		})

		t.Run("Scan Node", func(t *testing.T) {
			nodeId := nodeId + "-scan"
			mockStorageService.On("GenerateSignedURL", mock.Anything, mock.Anything).Return("test-url", nil)
			mockStorageService.On("GetFileUrl", mock.Anything, mock.Anything, mock.Anything).Return("test-url", nil)
			mockDiscordService.On("SendSecurityCouncilMessage", mock.Anything, mock.Anything).Return(nil)
			_, err := withMiddleware(authz, impl.PublishNodeVersion)(ctx, drip.PublishNodeVersionRequestObject{
				PublisherId: publisherId,
				NodeId:      nodeId,
				Body: &drip.PublishNodeVersionJSONRequestBody{
					Node: drip.Node{
						Id:          &nodeId,
						Description: &nodeDescription,
						Author:      &nodeAuthor,
						License:     &nodeLicense,
						Name:        &nodeName,
						Tags:        &nodeTags,
						Repository:  &source_code_repo,
					},
					NodeVersion: drip.NodeVersion{
						Version:      &nodeVersionLiteral,
						Changelog:    &changelog,
						Dependencies: &dependencies,
					},
					PersonalAccessToken: *createPersonalAccessTokenResponse.(drip.CreatePersonalAccessToken201JSONResponse).Token,
				},
			})
			require.NoError(t, err, "should return created node version")

			res, err := impl.GetNodeVersion(ctx, drip.GetNodeVersionRequestObject{NodeId: nodeId, VersionId: nodeVersionLiteral})
			require.NoError(t, err)
			require.IsType(t, drip.GetNodeVersion200JSONResponse{}, res)

			handled := true
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				req := dripservices_registry.ScanRequest{}
				require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
				assert.Equal(t, *res.(drip.GetNodeVersion200JSONResponse).DownloadUrl, req.URL)

				handled = true
				w.WriteHeader(http.StatusOK)
			}))
			t.Cleanup(s.Close)

			cfg := &config.Config{SecretScannerURL: s.URL}
			impl := implementation.NewStrictServerImplementation(
				client, cfg, mockStorageService, mockSlackService, mockDiscordService, mockAlgolia)
			dur := time.Duration(0)
			scanres, err := withMiddleware(authz, impl.SecurityScan)(ctx, drip.SecurityScanRequestObject{
				Params: drip.SecurityScanParams{
					MinAge: &dur,
				},
			})
			require.NoError(t, err)
			require.IsType(t, drip.SecurityScan200Response{}, scanres)
			assert.True(t, handled)
		})
	})
}

func isTokenMasked(token string) bool {
	tokenLength := len(token)
	// Ensure that only the first 4 and last 4 characters are not asterisks.
	middle := token[4:tokenLength]
	return strings.Count(middle, "*") == len(middle)
}
