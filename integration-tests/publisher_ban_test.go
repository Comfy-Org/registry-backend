package integration

import (
	"context"
	"net/http"
	"registry-backend/config"
	"registry-backend/drip"
	"registry-backend/ent/schema"
	authorization "registry-backend/server/middleware/authorization"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublisherBan(t *testing.T) {
	clientCtx := context.Background()
	client, cleanup := setupDB(t, clientCtx)
	defer cleanup()

	// Setup the mock services and server
	impl := NewStrictServerImplementationWithMocks(client, &config.Config{})
	authz := authorization.NewAuthorizationManager(client, impl.RegistryService).AuthorizationMiddleware()

	t.Run("Publisher Ban Tests", func(t *testing.T) {
		userCtx, _ := setupTestUser(client)
		adminCtx, _ := setupAdminUser(client)

		// Setup a test publisher
		publisherId, err := setupPublisher(userCtx, authz, impl)
		require.NoError(t, err, "should set up publisher")

		// endpoints to test the authorization middleware
		testEndpoints := []struct {
			name   string
			invoke func(ctx context.Context) error
		}{
			{"CreatePublisher", func(ctx context.Context) error {
				_, err := setupPublisher(ctx, authz, impl)
				return err
			}},
			{"DeleteNodeVersion", func(ctx context.Context) error {
				_, err := withMiddleware(authz, impl.DeleteNodeVersion)(
					ctx, drip.DeleteNodeVersionRequestObject{
						PublisherId: publisherId,
					})
				return err
			}},
		}

		t.Run("Ban publisher by non-admin", func(t *testing.T) {
			// Use the same publisher and node created earlier
			res, err := withMiddleware(authz, impl.BanPublisher)(
				userCtx, drip.BanPublisherRequestObject{PublisherId: publisherId})
			require.NoError(t, err, "should not ban publisher by non-admin")
			require.IsType(t, drip.BanPublisher403JSONResponse{}, res)
		})

		t.Run("Ban publisher by admin", func(t *testing.T) {
			// Use the same publisher and node created earlier
			_, err := withMiddleware(authz, impl.BanPublisher)(
				adminCtx, drip.BanPublisherRequestObject{PublisherId: publisherId})
			require.NoError(t, err)

			pub, err := client.Publisher.Get(userCtx, publisherId)
			require.NoError(t, err)
			assert.Equal(t, schema.PublisherStatusTypeBanned, pub.Status, "should ban publisher")
		})

		t.Run("Calling endpoints with a banned user", func(t *testing.T) {
			for _, tc := range testEndpoints {
				t.Run(tc.name, func(t *testing.T) {
					err := tc.invoke(userCtx)
					require.Error(t, err, "should return error")
					assertHTTPError(t, err, http.StatusForbidden)
				})
			}
		})
	})
}
