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
	drip_middleware "registry-backend/server/middleware"
	"testing"

	strictecho "github.com/oapi-codegen/runtime/strictmiddleware/echo"

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
	wrapped := drip.NewStrictHandler(impl, []strictecho.StrictEchoMiddlewareFunc{
		notBanned,
		drip_middleware.AuthorizationMiddleware(client),
	})

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
				name: "UpdatePublisher",
				req: func(ctx context.Context) (req *http.Request) {
					payloadBuf := new(bytes.Buffer)
					json.NewEncoder(payloadBuf).Encode(drip.UpdatePublisherJSONRequestBody{})

					req = httptest.NewRequest(http.MethodPost, "/", payloadBuf).WithContext(ctx)
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					return
				},
				fn: func(ctx echo.Context) error {
					return wrapped.UpdatePublisher(ctx, "")
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
}
