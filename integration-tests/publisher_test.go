package integration

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"registry-backend/config"
	"registry-backend/drip"
	authorization "registry-backend/server/middleware/authorization"
	"testing"
)

func TestRegistryPublisher(t *testing.T) {
	// Setup the database and mocked implementation
	client, cleanup := setupDB(t, context.Background())
	defer cleanup()

	impl := NewStrictServerImplementationWithMocks(client, &config.Config{})
	authz := authorization.NewAuthorizationManager(client, impl.RegistryService).AuthorizationMiddleware()

	// Create a test user and a random publisher for testing
	ctx, testUser := setupTestUser(client)
	pub := randomPublisher()

	t.Run("Create Publisher", func(t *testing.T) {
		// Test creating a publisher
		createPublisherResponse, err := withMiddleware(
			authz, impl.CreatePublisher)(ctx, drip.CreatePublisherRequestObject{
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
	})

	t.Run("Reject New Publisher With The Same Name", func(t *testing.T) {
		// Test rejecting duplicate publisher creation
		res, err := withMiddleware(authz, impl.CreatePublisher)(ctx, drip.CreatePublisherRequestObject{
			Body: pub,
		})
		require.NoError(t, err, "should return error")
		assert.IsType(t, drip.CreatePublisher400JSONResponse{}, res)
	})

	t.Run("Validate Publisher", func(t *testing.T) {
		// Test validating a publisher's name
		res, err := withMiddleware(authz, impl.ValidatePublisher)(ctx, drip.ValidatePublisherRequestObject{
			Params: drip.ValidatePublisherParams{Username: *pub.Name},
		})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.ValidatePublisher200JSONResponse{}, res, "should return 200")
		require.True(t, *res.(drip.ValidatePublisher200JSONResponse).IsAvailable, "should be available")
	})

	t.Run("Get Publisher", func(t *testing.T) {
		// Test retrieving the created publisher
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

		// Verify the members of the publisher
		expectedMembersCount := 1 // Adjust as needed
		assert.Equal(t, expectedMembersCount, len(*res.(drip.GetPublisher200JSONResponse).Members),
			"should return the correct number of members")
		for _, member := range *res.(drip.GetPublisher200JSONResponse).Members {
			assert.Equal(t, testUser.ID, *member.User.Id, "User ID should match")
			assert.Equal(t, testUser.Name, *member.User.Name, "User name should match")
			assert.Equal(t, testUser.Email, *member.User.Email, "User email should match")
		}
	})

	t.Run("Get Non-Exist Publisher", func(t *testing.T) {
		// Test retrieving a non-existent publisher
		res, err := withMiddleware(authz, impl.GetPublisher)(ctx, drip.GetPublisherRequestObject{
			PublisherId: *pub.Id + "invalid"})
		require.NoError(t, err, "should not return error")
		assert.IsType(t, drip.GetPublisher404JSONResponse{}, res)
	})

	t.Run("List Publishers", func(t *testing.T) {
		// Test listing all publishers
		res, err := withMiddleware(authz, impl.ListPublishers)(ctx, drip.ListPublishersRequestObject{})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.ListPublishers200JSONResponse{}, res, "should return 200 status code")
		res200 := res.(drip.ListPublishers200JSONResponse)
		require.Len(t, res200, 1, "should return all stored publishers")
		assert.Equal(t, drip.Publisher{
			Id:             pub.Id,
			Description:    pub.Description,
			SourceCodeRepo: pub.SourceCodeRepo,
			Website:        pub.Website,
			Support:        pub.Support,
			Logo:           pub.Logo,
			Name:           pub.Name,
			Members:        res200[0].Members, // Ignore dynamically generated fields
			CreatedAt:      res200[0].CreatedAt,
			Status:         res200[0].Status,
		}, res200[0], "should return correct publishers")
	})

	t.Run("Update Publisher", func(t *testing.T) {
		// Test updating a publisher
		pubUpdated := randomPublisher()
		pubUpdated.Id, pubUpdated.Name = pub.Id, pub.Name
		pub = pubUpdated

		res, err := withMiddleware(authz, impl.UpdatePublisher)(ctx, drip.UpdatePublisherRequestObject{
			PublisherId: *pubUpdated.Id,
			Body:        pubUpdated,
		})
		require.NoError(t, err, "should return updated publisher")
		assert.Equal(t, pubUpdated.Id, res.(drip.UpdatePublisher200JSONResponse).Id)
		assert.Equal(t, pubUpdated.Description, res.(drip.UpdatePublisher200JSONResponse).Description)
		assert.Equal(t, pubUpdated.SourceCodeRepo, res.(drip.UpdatePublisher200JSONResponse).SourceCodeRepo)
		assert.Equal(t, pubUpdated.Website, res.(drip.UpdatePublisher200JSONResponse).Website)
		assert.Equal(t, pubUpdated.Support, res.(drip.UpdatePublisher200JSONResponse).Support)
		assert.Equal(t, pubUpdated.Logo, res.(drip.UpdatePublisher200JSONResponse).Logo)
	})

	t.Run("Delete Publisher", func(t *testing.T) {
		// Test deleting a publisher
		res, err := withMiddleware(authz, impl.DeletePublisher)(ctx, drip.DeletePublisherRequestObject{
			PublisherId: *pub.Id})
		require.NoError(t, err, "should not return error")
		require.IsType(t, drip.DeletePublisher204Response{}, res, "should return 204")
	})
}
