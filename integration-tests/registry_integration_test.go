package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"registry-backend/config"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/nodeversion"
	"registry-backend/ent/schema"
	"registry-backend/mock/gateways"
	"registry-backend/server/implementation"
	drip_authorization "registry-backend/server/middleware/authorization"
	dripservices_registry "registry-backend/services/registry"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	strictecho "github.com/oapi-codegen/runtime/strictmiddleware/echo"

	"github.com/google/uuid"
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

func randomPublisher() *drip.Publisher {
	suffix := uuid.New().String()
	publisherId := "test-publisher-" + suffix
	description := "test-description" + suffix
	source_code_repo := "test-source-code-repo" + suffix
	website := "test-website" + suffix
	support := "test-support" + suffix
	logo := "test-logo" + suffix
	name := "test-name" + suffix

	return &drip.Publisher{
		Id:             &publisherId,
		Description:    &description,
		SourceCodeRepo: &source_code_repo,
		Website:        &website,
		Support:        &support,
		Logo:           &logo,
		Name:           &name,
	}
}

func randomNode() *drip.Node {
	suffix := uuid.New().String()
	nodeId := "test-node" + suffix
	nodeDescription := "test-node-description" + suffix
	nodeAuthor := "test-node-author" + suffix
	nodeLicense := "test-node-license" + suffix
	nodeName := "test-node-name" + suffix
	nodeTags := []string{"test-node-tag"}
	icon := "https://wwww.github.com/test-icon-" + suffix + ".svg"
	githubUrl := "https://www.github.com/test-github-url-" + suffix

	return &drip.Node{
		Id:          &nodeId,
		Name:        &nodeName,
		Description: &nodeDescription,
		Author:      &nodeAuthor,
		License:     &nodeLicense,
		Tags:        &nodeTags,
		Icon:        &icon,
		Repository:  &githubUrl,
	}
}

func randomNodeVersion(revision int) *drip.NodeVersion {
	suffix := uuid.New().String()

	version := fmt.Sprintf("1.0.%d", revision)
	changelog := "test-changelog-" + suffix
	dependencies := []string{"test-dependency" + suffix}
	return &drip.NodeVersion{
		Version:      &version,
		Changelog:    &changelog,
		Dependencies: &dependencies,
	}
}

type mockedImpl struct {
	*implementation.DripStrictServerImplementation

	mockStorageService *gateways.MockStorageService
	mockSlackService   *gateways.MockSlackService
	mockDiscordService *gateways.MockDiscordService
	mockAlgolia        *gateways.MockAlgoliaService
}

func newMockedImpl(client *ent.Client, cfg *config.Config) (impl mockedImpl, authz strictecho.StrictEchoMiddlewareFunc) {
	// Initialize the Service
	mockStorageService := new(gateways.MockStorageService)

	mockDiscordService := new(gateways.MockDiscordService)
	mockDiscordService.On("SendSecurityCouncilMessage", mock.Anything).
		Return(nil) // Do nothing for all discord messsage calls.

	mockSlackService := new(gateways.MockSlackService)
	mockSlackService.
		On("SendRegistryMessageToSlack", mock.Anything).
		Return(nil) // Do nothing for all slack messsage calls.

	mockAlgolia := new(gateways.MockAlgoliaService)
	mockAlgolia.
		On("IndexNodes", mock.Anything, mock.Anything).
		Return(nil).
		On("DeleteNode", mock.Anything, mock.Anything).
		Return(nil)

	impl = mockedImpl{
		DripStrictServerImplementation: implementation.NewStrictServerImplementation(
			client, cfg, mockStorageService, mockSlackService, mockDiscordService, mockAlgolia),
		mockStorageService: mockStorageService,
		mockSlackService:   mockSlackService,
		mockDiscordService: mockDiscordService,
		mockAlgolia:        mockAlgolia,
	}
	authz = drip_authorization.NewAuthorizationManager(client, impl.RegistryService).
		AuthorizationMiddleware()
	return
}

func TestRegistryPublisher(t *testing.T) {
	client, cleanup := setupDB(t, context.Background())
	defer cleanup()
	impl, authz := newMockedImpl(client, &config.Config{})

	ctx, testUser := setUpTest(client)
	pub := randomPublisher()

	createPublisherResponse, err := withMiddleware(authz, impl.CreatePublisher)(ctx, drip.CreatePublisherRequestObject{
		Body: pub,
	})
	require.NoError(t, err, "should return created publisher")
	require.NotNil(t, createPublisherResponse, "should return created publisher")
	assert.Equal(t, pub.Id, createPublisherResponse.(drip.CreatePublisher201JSONResponse).Id)
	assert.Equal(t, pub.Description, createPublisherResponse.(drip.CreatePublisher201JSONResponse).Description)
	assert.Equal(t, pub.SourceCodeRepo, createPublisherResponse.(drip.CreatePublisher201JSONResponse).SourceCodeRepo)
	assert.Equal(t, pub.Website, createPublisherResponse.(drip.CreatePublisher201JSONResponse).Website)
	assert.Equal(t, pub.Support, createPublisherResponse.(drip.CreatePublisher201JSONResponse).Support)
	assert.Equal(t, pub.Logo, createPublisherResponse.(drip.CreatePublisher201JSONResponse).Logo)

	t.Run("Reject New Publisher With The Same Name", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.CreatePublisher)(ctx, drip.CreatePublisherRequestObject{
			Body: pub,
		})
		require.NoError(t, err, "should return error")
		assert.IsType(t, drip.CreatePublisher400JSONResponse{}, res)
	})

	t.Run("Validate Publisher", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.ValidatePublisher)(ctx, drip.ValidatePublisherRequestObject{
			Params: drip.ValidatePublisherParams{Username: *pub.Name},
		})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.ValidatePublisher200JSONResponse{}, res, "should return 200")
		require.True(t, *res.(drip.ValidatePublisher200JSONResponse).IsAvailable, "should be available")
	})

	t.Run("Get Publisher", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.GetPublisher)(ctx, drip.GetPublisherRequestObject{
			PublisherId: *pub.Id})
		require.NoError(t, err, "should return created publisher")
		assert.Equal(t, pub.Id, res.(drip.GetPublisher200JSONResponse).Id)
		assert.Equal(t, pub.Description, res.(drip.GetPublisher200JSONResponse).Description)
		assert.Equal(t, pub.SourceCodeRepo, res.(drip.GetPublisher200JSONResponse).SourceCodeRepo)
		assert.Equal(t, pub.Website, res.(drip.GetPublisher200JSONResponse).Website)
		assert.Equal(t, pub.Support, res.(drip.GetPublisher200JSONResponse).Support)
		assert.Equal(t, pub.Logo, res.(drip.GetPublisher200JSONResponse).Logo)
		assert.Equal(t, pub.Name, res.(drip.GetPublisher200JSONResponse).Name)

		// Check the number of members returned
		expectedMembersCount := 1 // Adjust to your expected count
		assert.Equal(t, expectedMembersCount,
			len(*res.(drip.GetPublisher200JSONResponse).Members),
			"should return the correct number of members")

		// Check specific properties of each member, adjust indices accordingly
		for _, member := range *res.(drip.GetPublisher200JSONResponse).Members {
			expectedUserId := testUser.ID
			expectedUserName := testUser.Name
			expectedUserEmail := testUser.Email

			assert.Equal(t, expectedUserId, *member.User.Id, "User ID should match")
			assert.Equal(t, expectedUserName, *member.User.Name, "User name should match")
			assert.Equal(t, expectedUserEmail, *member.User.Email, "User email should match")
		}
	})

	t.Run("Get Non-Exist Publisher", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.GetPublisher)(ctx, drip.GetPublisherRequestObject{PublisherId: *pub.Id + "invalid"})
		require.NoError(t, err, "should not return error")
		assert.IsType(t, drip.GetPublisher404JSONResponse{}, res)
	})

	t.Run("List Publishers", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.ListPublishers)(ctx, drip.ListPublishersRequestObject{})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.ListPublishers200JSONResponse{}, res, "should return 200 status code")
		res200 := res.(drip.ListPublishers200JSONResponse)
		require.Len(t, res200, 1, "should return all stored publlishers")
		assert.Equal(t, drip.Publisher{
			Id:             pub.Id,
			Description:    pub.Description,
			SourceCodeRepo: pub.SourceCodeRepo,
			Website:        pub.Website,
			Support:        pub.Support,
			Logo:           pub.Logo,
			Name:           pub.Name,

			// generated thus ignored in comparison
			Members:   res200[0].Members,
			CreatedAt: res200[0].CreatedAt,
			Status:    res200[0].Status,
		}, res200[0], "should return correct publishers")
	})

	t.Run("Update Publisher", func(t *testing.T) {
		pubUpdated := randomPublisher()
		pubUpdated.Id, pubUpdated.Name = pub.Id, pub.Name
		pub = pubUpdated

		res, err := withMiddleware(authz, impl.UpdatePublisher)(ctx, drip.UpdatePublisherRequestObject{
			PublisherId: *pubUpdated.Id,
			Body:        pubUpdated,
		})
		require.NoError(t, err, "should return created publisher")
		assert.Equal(t, pubUpdated.Id, res.(drip.UpdatePublisher200JSONResponse).Id)
		assert.Equal(t, pubUpdated.Description, res.(drip.UpdatePublisher200JSONResponse).Description)
		assert.Equal(t, pubUpdated.SourceCodeRepo, res.(drip.UpdatePublisher200JSONResponse).SourceCodeRepo)
		assert.Equal(t, pubUpdated.Website, res.(drip.UpdatePublisher200JSONResponse).Website)
		assert.Equal(t, pubUpdated.Support, res.(drip.UpdatePublisher200JSONResponse).Support)
		assert.Equal(t, pubUpdated.Logo, res.(drip.UpdatePublisher200JSONResponse).Logo)

		_, err = withMiddleware(authz, impl.ListPublishersForUser)(ctx, drip.ListPublishersForUserRequestObject{})
		require.NoError(t, err, "should return created publisher")
	})

	t.Run("Delete Publisher", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.DeletePublisher)(ctx, drip.DeletePublisherRequestObject{PublisherId: *pub.Id})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.DeletePublisher204Response{}, res, "should return 204")
	})
}

func TestRegistryPersonalAccessToken(t *testing.T) {
	client, cleanup := setupDB(t, context.Background())
	defer cleanup()
	impl, authz := newMockedImpl(client, &config.Config{})

	ctx, _ := setUpTest(client)
	pub := randomPublisher()
	_, err := withMiddleware(authz, impl.CreatePublisher)(ctx, drip.CreatePublisherRequestObject{
		Body: pub,
	})
	require.NoError(t, err, "should return created publisher")

	tokenName := "test-token-name"
	tokenDescription := "test-token-description"
	res, err := withMiddleware(authz, impl.CreatePersonalAccessToken)(
		ctx, drip.CreatePersonalAccessTokenRequestObject{
			PublisherId: *pub.Id,
			Body: &drip.PersonalAccessToken{
				Name:        &tokenName,
				Description: &tokenDescription,
			},
		})
	require.NoError(t, err, "should return created token")
	require.NotNil(t,
		*res.(drip.CreatePersonalAccessToken201JSONResponse).Token,
		"Token should have a value.")

	t.Run("List Personal Access Token", func(t *testing.T) {
		getPersonalAccessTokenResponse, err := withMiddleware(authz, impl.ListPersonalAccessTokens)(ctx, drip.ListPersonalAccessTokensRequestObject{
			PublisherId: *pub.Id,
		})
		require.NoError(t, err, "should return created token")
		assert.Equal(t, tokenName,
			*getPersonalAccessTokenResponse.(drip.ListPersonalAccessTokens200JSONResponse)[0].Name)
		assert.Equal(t, tokenDescription,
			*getPersonalAccessTokenResponse.(drip.ListPersonalAccessTokens200JSONResponse)[0].Description)
		assert.True(t,
			isTokenMasked(*getPersonalAccessTokenResponse.(drip.ListPersonalAccessTokens200JSONResponse)[0].Token))
	})
}

func TestRegistryNode(t *testing.T) {
	client, cleanup := setupDB(t, context.Background())
	defer cleanup()
	impl, authz := newMockedImpl(client, &config.Config{})

	ctx, _ := setUpTest(client)
	pub := randomPublisher()

	_, err := withMiddleware(authz, impl.CreatePublisher)(ctx, drip.CreatePublisherRequestObject{
		Body: pub,
	})
	require.NoError(t, err, "should return created publisher")

	node := randomNode()
	res, err := withMiddleware(authz, impl.CreateNode)(ctx, drip.CreateNodeRequestObject{
		PublisherId: *pub.Id,
		Body:        node,
	})
	require.NoError(t, err, "should return created node")
	require.NotNil(t, res, "should return created node")
	assert.Equal(t, node.Id, res.(drip.CreateNode201JSONResponse).Id)
	assert.Equal(t, node.Description, res.(drip.CreateNode201JSONResponse).Description)
	assert.Equal(t, node.Author, res.(drip.CreateNode201JSONResponse).Author)
	assert.Equal(t, node.License, res.(drip.CreateNode201JSONResponse).License)
	assert.Equal(t, node.Name, res.(drip.CreateNode201JSONResponse).Name)
	assert.Equal(t, node.Tags, res.(drip.CreateNode201JSONResponse).Tags)
	assert.Equal(t, node.Icon, res.(drip.CreateNode201JSONResponse).Icon)
	assert.Equal(t, node.Repository, res.(drip.CreateNode201JSONResponse).Repository)
	assert.Equal(t, drip.NodeStatusActive, *res.(drip.CreateNode201JSONResponse).Status)

	t.Run("Get Node", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.GetNode)(ctx, drip.GetNodeRequestObject{NodeId: *node.Id})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.GetNode200JSONResponse{}, res)
		res200 := res.(drip.GetNode200JSONResponse)
		expDl, expRate := 0, float32(0)
		nodeStatus := drip.NodeStatusActive
		assert.Equal(t, drip.GetNode200JSONResponse{
			Id:          node.Id,
			Name:        node.Name,
			Description: node.Description,
			Author:      node.Author,
			Tags:        node.Tags,
			License:     node.License,
			Icon:        node.Icon,
			Repository:  node.Repository,

			Downloads:    &expDl,
			Rating:       &expRate,
			Status:       &nodeStatus,
			StatusDetail: proto.String(""),
			Category:     proto.String(""),
		}, res200, "should return stored node data")
	})

	t.Run("Get Publisher Nodes", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.ListNodesForPublisher)(ctx, drip.ListNodesForPublisherRequestObject{
			PublisherId: *pub.Id,
		})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.ListNodesForPublisher200JSONResponse{}, res)
		res200 := res.(drip.ListNodesForPublisher200JSONResponse)
		require.Len(t, res200, 1)
		expDl, expRate := 0, float32(0)
		nodeStatus := drip.NodeStatusActive
		assert.Equal(t, drip.Node{
			Id:          node.Id,
			Name:        node.Name,
			Description: node.Description,
			Author:      node.Author,
			Tags:        node.Tags,
			License:     node.License,
			Icon:        node.Icon,
			Repository:  node.Repository,

			Downloads:    &expDl,
			Rating:       &expRate,
			Status:       &nodeStatus,
			StatusDetail: proto.String(""),
			Category:     proto.String(""),
		}, res200[0], "should return stored node data")
	})

	t.Run("Get Not Exist Node", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.GetNode)(ctx, drip.GetNodeRequestObject{NodeId: *node.Id + "fake"})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.GetNode404JSONResponse{}, res)
	})

	t.Run("Update Node", func(t *testing.T) {
		unode := randomNode()
		unode.Id = node.Id
		node = unode

		updateNodeResponse, err := withMiddleware(authz, impl.UpdateNode)(ctx, drip.UpdateNodeRequestObject{
			PublisherId: *pub.Id,
			NodeId:      *node.Id,
			Body:        node,
		})
		require.NoError(t, err, "should return created node")
		assert.Equal(t, node.Id, updateNodeResponse.(drip.UpdateNode200JSONResponse).Id)
		assert.Equal(t, node.Description, updateNodeResponse.(drip.UpdateNode200JSONResponse).Description)
		assert.Equal(t, node.Author, updateNodeResponse.(drip.UpdateNode200JSONResponse).Author)
		assert.Equal(t, node.License, updateNodeResponse.(drip.UpdateNode200JSONResponse).License)
		assert.Equal(t, node.Name, updateNodeResponse.(drip.UpdateNode200JSONResponse).Name)
		assert.Equal(t, node.Tags, updateNodeResponse.(drip.UpdateNode200JSONResponse).Tags)
		assert.Equal(t, node.Icon, updateNodeResponse.(drip.UpdateNode200JSONResponse).Icon)
		assert.Equal(t, node.Repository, updateNodeResponse.(drip.UpdateNode200JSONResponse).Repository)

		resUpdated, err := withMiddleware(authz, impl.GetNode)(ctx, drip.GetNodeRequestObject{NodeId: *node.Id})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.GetNode200JSONResponse{}, resUpdated)
		res200Updated := resUpdated.(drip.GetNode200JSONResponse)
		expDl, expRate := 0, float32(0)
		nodeStatus := drip.NodeStatusActive
		assert.Equal(t, drip.GetNode200JSONResponse{
			Id:          node.Id,
			Description: node.Description,
			Author:      node.Author,
			License:     node.License,
			Name:        node.Name,
			Tags:        node.Tags,
			Icon:        node.Icon,
			Repository:  node.Repository,

			Downloads:    &expDl,
			Rating:       &expRate,
			Status:       &nodeStatus,
			StatusDetail: proto.String(""),
			Category:     proto.String(""),
		}, res200Updated, "should return updated node data")
	})

	t.Run("Update Not Exist Node", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.UpdateNode)(ctx, drip.UpdateNodeRequestObject{PublisherId: *pub.Id, NodeId: *node.Id + "fake", Body: &drip.Node{}})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.UpdateNode404JSONResponse{}, res)
	})

	t.Run("Index Nodes", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.ReindexNodes)(ctx, drip.ReindexNodesRequestObject{})
		require.NoError(t, err, "should not return error")
		assert.IsType(t, drip.ReindexNodes200Response{}, res)
	})

	t.Run("Delete Node", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.DeleteNode)(ctx, drip.DeleteNodeRequestObject{PublisherId: *pub.Id, NodeId: *node.Id})
		require.NoError(t, err, "should not return error")
		assert.IsType(t, drip.DeleteNode204Response{}, res)
	})

	t.Run("Delete Not Exist Node", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.DeleteNode)(ctx, drip.DeleteNodeRequestObject{PublisherId: *pub.Id, NodeId: *node.Id + "fake"})
		require.NoError(t, err, "should not return error")
		assert.IsType(t, drip.DeleteNode204Response{}, res)
	})
}

func TestRegistryNodeVersion(t *testing.T) {
	client, cleanup := setupDB(t, context.Background())
	defer cleanup()
	impl, authz := newMockedImpl(client, &config.Config{})

	ctx, _ := setUpTest(client)
	pub := randomPublisher()

	respub, err := withMiddleware(authz, impl.CreatePublisher)(ctx, drip.CreatePublisherRequestObject{
		Body: pub,
	})
	require.NoError(t, err, "should return created publisher")
	createdPublisher := (respub.(drip.CreatePublisher201JSONResponse))

	tokenName := "test-token-name"
	tokenDescription := "test-token-description"
	respat, err := withMiddleware(authz, impl.CreatePersonalAccessToken)(ctx, drip.CreatePersonalAccessTokenRequestObject{
		PublisherId: *pub.Id,
		Body: &drip.PersonalAccessToken{
			Name:        &tokenName,
			Description: &tokenDescription,
		},
	})
	require.NoError(t, err, "should return created token")
	token := *respat.(drip.CreatePersonalAccessToken201JSONResponse).Token

	node := randomNode()
	nodeVersion := randomNodeVersion(0)
	signedUrl := "test-url"
	downloadUrl := fmt.Sprintf("https://storage.googleapis.com/comfy-registry/%s/%s/%s/node.tar.gz", *pub.Id, *node.Id, *nodeVersion.Version)
	var createdNodeVersion drip.NodeVersion

	impl.mockStorageService.On("GenerateSignedURL", mock.Anything, mock.Anything).Return(signedUrl, nil)
	impl.mockStorageService.On("GetFileUrl", mock.Anything, mock.Anything, mock.Anything).Return(signedUrl, nil)
	impl.mockDiscordService.On("SendSecurityCouncilMessage", mock.Anything, mock.Anything).Return(nil)
	res, err := withMiddleware(authz, impl.PublishNodeVersion)(ctx, drip.PublishNodeVersionRequestObject{
		PublisherId: *pub.Id,
		NodeId:      *node.Id,
		Body: &drip.PublishNodeVersionJSONRequestBody{
			PersonalAccessToken: token,
			Node:                *node,
			NodeVersion:         *nodeVersion,
		},
	})
	require.NoError(t, err, "should return created node version")
	assert.Equal(t, nodeVersion.Version, res.(drip.PublishNodeVersion201JSONResponse).NodeVersion.Version)
	require.Equal(t, nodeVersion.Dependencies, res.(drip.PublishNodeVersion201JSONResponse).NodeVersion.Dependencies, "should return pip dependencies")
	require.Equal(t, nodeVersion.Changelog, res.(drip.PublishNodeVersion201JSONResponse).NodeVersion.Changelog, "should return changelog")
	require.Equal(t, signedUrl, *res.(drip.PublishNodeVersion201JSONResponse).SignedUrl, "should return signed url")
	versionStatus := drip.NodeVersionStatusPending
	require.Equal(t, versionStatus, *res.(drip.PublishNodeVersion201JSONResponse).NodeVersion.Status, "should return pending status")
	createdNodeVersion = *res.(drip.PublishNodeVersion201JSONResponse).NodeVersion // Needed for downstream tests.

	t.Run("Admin Update", func(t *testing.T) {
		adminCtx, _ := setUpAdminTest(client)
		activeStatus := drip.NodeVersionStatusActive
		adminUpdateNodeVersionResp, err := impl.AdminUpdateNodeVersion(adminCtx, drip.AdminUpdateNodeVersionRequestObject{
			NodeId:        *node.Id,
			VersionNumber: *createdNodeVersion.Version,
			Body: &drip.AdminUpdateNodeVersionJSONRequestBody{
				Status: &activeStatus,
			},
		})
		require.NoError(t, err, "should return updated node version")
		assert.Equal(t, activeStatus, *adminUpdateNodeVersionResp.(drip.AdminUpdateNodeVersion200JSONResponse).Status)
	})

	t.Run("List Node Version Before Create", func(t *testing.T) {
		node := randomNode()
		resVersions, err := withMiddleware(authz, impl.ListNodeVersions)(ctx, drip.ListNodeVersionsRequestObject{NodeId: *node.Id})
		require.NoError(t, err, "should return error since node version doesn't exists")
		require.IsType(t, drip.ListNodeVersions200JSONResponse{}, resVersions)
		assert.Empty(t, resVersions.(drip.ListNodeVersions200JSONResponse), "should not return any node versions")
	})

	t.Run("Create Node Version with Fake Token", func(t *testing.T) {
		_, err := withMiddleware(authz, impl.PublishNodeVersion)(ctx, drip.PublishNodeVersionRequestObject{
			PublisherId: *pub.Id,
			NodeId:      *node.Id,
			Body: &drip.PublishNodeVersionJSONRequestBody{
				Node:                *node,
				NodeVersion:         *nodeVersion,
				PersonalAccessToken: "faketoken",
			},
		})
		require.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code, "should return 400 bad request")
	})

	t.Run("Create Node Version with invalid node id", func(t *testing.T) {
		for _, suffix := range []string{"LOWERCASEONLY", "invalidCharacter&"} {
			node := randomNode()
			*node.Id = *node.Id + suffix
			res, err := withMiddleware(authz, impl.PublishNodeVersion)(ctx, drip.PublishNodeVersionRequestObject{
				PublisherId: *pub.Id,
				NodeId:      *node.Id,
				Body: &drip.PublishNodeVersionJSONRequestBody{
					Node:                *node,
					NodeVersion:         *randomNodeVersion(0),
					PersonalAccessToken: token,
				},
			})
			require.NoError(t, err)
			require.IsType(t, drip.PublishNodeVersion400JSONResponse{}, res)
		}
	})

	t.Run("Get not exist Node Version ", func(t *testing.T) {
		res, err := withMiddleware(authz, impl.GetNodeVersion)(ctx, drip.GetNodeVersionRequestObject{NodeId: *node.Id + "fake", VersionId: *nodeVersion.Version})
		require.NoError(t, err, "should not return error")
		assert.IsType(t, drip.GetNodeVersion404JSONResponse{}, res)
	})

	t.Run("Create Node Version of Not Exist Node", func(t *testing.T) {
		_, err := withMiddleware(authz, impl.PublishNodeVersion)(ctx, drip.PublishNodeVersionRequestObject{
			PublisherId: *pub.Id,
			NodeId:      *node.Id + "fake",
			Body:        &drip.PublishNodeVersionJSONRequestBody{},
		})
		require.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code, "should return 400 bad request")
	})

	t.Run("List Node Versions", func(t *testing.T) {
		resVersions, err := withMiddleware(authz, impl.ListNodeVersions)(ctx, drip.ListNodeVersionsRequestObject{NodeId: *node.Id})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.ListNodeVersions200JSONResponse{}, resVersions, "should return 200")
		resVersions200 := resVersions.(drip.ListNodeVersions200JSONResponse)
		require.Len(t, resVersions200, 1, "should return only one version")
		nodeVersionStatus := drip.NodeVersionStatusActive
		t.Log("Download URL: ", *resVersions200[0].DownloadUrl)
		t.Log("Download URL: ", downloadUrl)
		t.Log("Status: ", *resVersions200[0].Status)
		t.Log("Status: ", nodeVersionStatus)
		assert.Equal(t, drip.NodeVersion{
			// generated attribute
			Id:        resVersions200[0].Id,
			CreatedAt: resVersions200[0].CreatedAt,

			Deprecated:   proto.Bool(false),
			Version:      nodeVersion.Version,
			Changelog:    nodeVersion.Changelog,
			Dependencies: nodeVersion.Dependencies,
			DownloadUrl:  &downloadUrl,
			Status:       &nodeVersionStatus,
			StatusReason: proto.String(""),
		}, resVersions200[0], "should be equal")
	})

	t.Run("Update Node Version", func(t *testing.T) {
		updatedChangelog := "test-changelog-2"
		resUNV, err := withMiddleware(authz, impl.UpdateNodeVersion)(ctx, drip.UpdateNodeVersionRequestObject{
			PublisherId: *pub.Id,
			NodeId:      *node.Id,
			VersionId:   *createdNodeVersion.Id,
			Body: &drip.NodeVersionUpdateRequest{
				Changelog:  &updatedChangelog,
				Deprecated: proto.Bool(true),
			},
		})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.UpdateNodeVersion200JSONResponse{}, resUNV, "should return 200")

		res, err := withMiddleware(authz, impl.ListNodeVersions)(ctx, drip.ListNodeVersionsRequestObject{NodeId: *node.Id})
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
			Version:      nodeVersion.Version,
			Dependencies: nodeVersion.Dependencies,
			Changelog:    &updatedChangelog,
			DownloadUrl:  &downloadUrl,
			Status:       &status,
			StatusReason: proto.String(""),
		}
		assert.Equal(t, updatedNodeVersion, res200[0], "should be equal")
		createdNodeVersion = res200[0]
	})

	t.Run("List Nodes", func(t *testing.T) {
		nodeIDs := map[string]*drip.NodeVersion{
			*node.Id:        &createdNodeVersion,
			*node.Id + "-1": nil,
			*node.Id + "-2": nil,
		}

		for nodeId := range nodeIDs {
			for i := 0; i < 2; i++ {
				version := fmt.Sprintf("2.0.%d", i)
				res, err := withMiddleware(authz, impl.PublishNodeVersion)(ctx, drip.PublishNodeVersionRequestObject{
					PublisherId: *pub.Id,
					NodeId:      nodeId,
					Body: &drip.PublishNodeVersionJSONRequestBody{
						Node: drip.Node{
							Id:          &nodeId,
							Description: node.Description,
							Author:      node.Author,
							License:     node.License,
							Name:        node.Name,
							Tags:        node.Tags,
							Repository:  node.Repository,
						},
						NodeVersion: drip.NodeVersion{
							Version:      &version,
							Changelog:    createdNodeVersion.Changelog,
							Dependencies: createdNodeVersion.Dependencies,
						},
						PersonalAccessToken: *respat.(drip.CreatePersonalAccessToken201JSONResponse).Token,
					},
				})
				require.NoError(t, err, "should return created node version")
				require.IsType(t, drip.PublishNodeVersion201JSONResponse{}, res)
				res200 := res.(drip.PublishNodeVersion201JSONResponse)
				nodeIDs[nodeId] = res200.NodeVersion
			}
		}

		resNodes, err := withMiddleware(authz, impl.ListAllNodes)(ctx, drip.ListAllNodesRequestObject{})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.ListAllNodes200JSONResponse{}, resNodes, "should return 200 server response")
		resNodes200 := resNodes.(drip.ListAllNodes200JSONResponse)
		assert.Len(t, *resNodes200.Nodes, len(nodeIDs), "should only contain 1 node")

		for _, node := range *resNodes200.Nodes {
			expDl, expRate := 0, float32(0)
			nodeStatus := drip.NodeStatusActive
			expectedNode := drip.Node{
				Id:            node.Id,
				Name:          node.Name,
				Repository:    node.Repository,
				Description:   node.Description,
				Author:        node.Author,
				License:       node.License,
				Tags:          node.Tags,
				LatestVersion: nodeIDs[*node.Id],
				Icon:          node.Icon,
				Publisher:     (*drip.Publisher)(&createdPublisher),
				Downloads:     &expDl,
				Rating:        &expRate,
				Status:        &nodeStatus,
				StatusDetail:  proto.String(""),
				Category:      proto.String(""),
			}
			expectedNode.LatestVersion.DownloadUrl = node.LatestVersion.DownloadUrl // generated
			expectedNode.LatestVersion.Deprecated = node.LatestVersion.Deprecated   // generated
			expectedNode.LatestVersion.CreatedAt = node.LatestVersion.CreatedAt     // generated
			expectedNode.Publisher.CreatedAt = node.Publisher.CreatedAt
			assert.Equal(t, expectedNode, node)
		}
	})

	t.Run("Node Installation", func(t *testing.T) {
		resIns, err := withMiddleware(authz, impl.InstallNode)(ctx, drip.InstallNodeRequestObject{NodeId: *node.Id})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.InstallNode200JSONResponse{}, resIns, "should return 200")

		resIns, err = withMiddleware(authz, impl.InstallNode)(ctx, drip.InstallNodeRequestObject{
			NodeId: *node.Id, Params: drip.InstallNodeParams{Version: createdNodeVersion.Version}})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.InstallNode200JSONResponse{}, resIns, "should return 200")

		t.Run("Get Total Install", func(t *testing.T) {
			res, err := withMiddleware(authz, impl.GetNode)(ctx, drip.GetNodeRequestObject{
				NodeId: *node.Id,
			})
			require.NoError(t, err, "should not return error")
			require.IsType(t, drip.GetNode200JSONResponse{}, res)
			assert.Equal(t, int(2), *res.(drip.GetNode200JSONResponse).Downloads)
		})

		t.Run("Add review", func(t *testing.T) {
			res, err := withMiddleware(authz, impl.PostNodeReview)(ctx, drip.PostNodeReviewRequestObject{
				NodeId: *node.Id,
				Params: drip.PostNodeReviewParams{Star: 5},
			})
			require.NoError(t, err)
			require.IsType(t, drip.PostNodeReview200JSONResponse{}, res)
			res200 := res.(drip.PostNodeReview200JSONResponse)
			assert.Equal(t, float32(5), *res200.Rating)
		})
	})

	t.Run("Node installation on not exist node or version", func(t *testing.T) {
		resIns, err := withMiddleware(authz, impl.InstallNode)(ctx, drip.InstallNodeRequestObject{NodeId: *node.Id + "fake"})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.InstallNode404JSONResponse{}, resIns, "should return 404")

		resIns, err = withMiddleware(authz, impl.InstallNode)(ctx, drip.InstallNodeRequestObject{
			NodeId: *node.Id, Params: drip.InstallNodeParams{Version: proto.String(*createdNodeVersion.Version + "fake")}})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.InstallNode404JSONResponse{}, resIns, "should return 404")
	})

	t.Run("Scan Node", func(t *testing.T) {
		node := randomNode()
		nodeVersion := randomNodeVersion(0)
		downloadUrl := fmt.Sprintf("https://storage.googleapis.com/comfy-registry/%s/%s/%s/node.tar.gz", *pub.Id, *node.Id, *nodeVersion.Version)

		impl.mockStorageService.On("GenerateSignedURL", mock.Anything, mock.Anything).Return("test-url", nil)
		impl.mockStorageService.On("GetFileUrl", mock.Anything, mock.Anything, mock.Anything).Return("test-url", nil)
		impl.mockDiscordService.On("SendSecurityCouncilMessage", mock.Anything, mock.Anything).Return(nil)
		_, err := withMiddleware(authz, impl.PublishNodeVersion)(ctx, drip.PublishNodeVersionRequestObject{
			PublisherId: *pub.Id,
			NodeId:      *node.Id,
			Body: &drip.PublishNodeVersionJSONRequestBody{
				Node:                *node,
				NodeVersion:         *nodeVersion,
				PersonalAccessToken: *respat.(drip.CreatePersonalAccessToken201JSONResponse).Token,
			},
		})
		require.NoError(t, err, "should return created node version")

		nodesToScans, err := client.NodeVersion.Query().Where(nodeversion.StatusEQ(schema.NodeVersionStatusPending)).Count(ctx)
		require.NoError(t, err)

		newNodeScanned := false
		nodesScanned := 0
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			req := dripservices_registry.ScanRequest{}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
			if downloadUrl == req.URL {
				newNodeScanned = true
			}
			nodesScanned++
		}))
		t.Cleanup(s.Close)

		impl, authz := newMockedImpl(client, &config.Config{SecretScannerURL: s.URL})
		dur := time.Duration(0)
		scanres, err := withMiddleware(authz, impl.SecurityScan)(ctx, drip.SecurityScanRequestObject{
			Params: drip.SecurityScanParams{
				MinAge: &dur,
			},
		})
		require.NoError(t, err)
		require.IsType(t, drip.SecurityScan200Response{}, scanres)
		assert.True(t, newNodeScanned)
		assert.Equal(t, nodesToScans, nodesScanned)
	})

}

func isTokenMasked(token string) bool {
	tokenLength := len(token)
	// Ensure that only the first 4 and last 4 characters are not asterisks.
	middle := token[4:tokenLength]
	return strings.Count(middle, "*") == len(middle)
}
