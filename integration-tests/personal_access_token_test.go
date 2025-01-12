package integration

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"registry-backend/config"
	"registry-backend/drip"
	authorization "registry-backend/server/middleware/authorization"
	"strings"
	"testing"
)

func TestRegistryPersonalAccessToken(t *testing.T) {
	// Setup the database and clean up after the test
	client, cleanup := setupDB(t, context.Background())
	defer cleanup()

	// Initialize server implementation and authorization middleware
	impl := NewStrictServerImplementationWithMocks(client, &config.Config{})
	authz := authorization.NewAuthorizationManager(client, impl.RegistryService).AuthorizationMiddleware()

	// Setup test user and publisher
	ctx, _ := setupTestUser(client)
	publisherId, err := setupPublisher(ctx, authz, impl)
	require.NoError(t, err, "Failed to create publisher")

	tokenName := "test-token-name"
	tokenDescription := "test-token-description"

	// Run the token creation test
	t.Run("Create Personal Access Token", func(t *testing.T) {
		// Create the personal access token for the publisher
		res, err := withMiddleware(authz, impl.CreatePersonalAccessToken)(
			ctx, drip.CreatePersonalAccessTokenRequestObject{
				PublisherId: publisherId,
				Body: &drip.PersonalAccessToken{
					Name:        &tokenName,
					Description: &tokenDescription,
				},
			})
		require.NoError(t, err, "Failed to create personal access token")
		require.NotNil(t, *res.(drip.CreatePersonalAccessToken201JSONResponse).Token, "Token should not be nil")
	})

	// List and validate personal access token
	t.Run("List Personal Access Token", func(t *testing.T) {
		// Fetch the personal access tokens for the publisher
		getPersonalAccessTokenResponse, err := withMiddleware(
			authz, impl.ListPersonalAccessTokens)(ctx, drip.ListPersonalAccessTokensRequestObject{
			PublisherId: publisherId,
		})
		require.NoError(t, err, "Failed to fetch personal access tokens")

		// Ensure the response contains the correct token
		tokenResponse := getPersonalAccessTokenResponse.(drip.ListPersonalAccessTokens200JSONResponse)
		require.Len(t, tokenResponse, 1, "Expected exactly one token")

		// Extract the token and check if the details match
		token := tokenResponse[0]
		assert.Equal(t, "test-token-name", *token.Name, "Token name should match")
		assert.Equal(t, "test-token-description", *token.Description, "Token description should match")

		// Verify the token is masked properly
		assert.True(t, isTokenMasked(*token.Token), "Token should be masked")
	})
}

func isTokenMasked(token string) bool {
	tokenLength := len(token)
	// Ensure that only the first 4 and last 4 characters are not asterisks.
	middle := token[4:tokenLength]
	return strings.Count(middle, "*") == len(middle)
}
