package integration

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"registry-backend/config"
	"registry-backend/db"
	"registry-backend/drip"
	"registry-backend/mock/gateways"
	"registry-backend/server/implementation"
	"registry-backend/server/middleware"
	auth "registry-backend/server/middleware/authentication"
	authorization "registry-backend/server/middleware/authorization"
	"runtime"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"registry-backend/ent"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mixpanel/mixpanel-go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/lib/pq"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type MockedServerImplementation struct {
	*implementation.DripStrictServerImplementation

	mockStorageService *gateways.MockStorageService
	mockSlackService   *gateways.MockSlackService
	mockDiscordService *gateways.MockDiscordService
	mockAlgolia        *gateways.MockAlgoliaService
	mockPubsubService  *gateways.MockPubSubService
}

// NewStrictServerImplementationWithMocks initializes and returns the implementation with mock services.
func NewStrictServerImplementationWithMocks(
	client *ent.Client, config *config.Config) *MockedServerImplementation {
	// Create mock services for dependencies.
	mockStorageService := new(gateways.MockStorageService)
	mockPubsubService := new(gateways.MockPubSubService)
	mockSlackService := new(gateways.MockSlackService)
	mockDiscordService := new(gateways.MockDiscordService)
	mockAlgolia := new(gateways.MockAlgoliaService)
	newRelicApp := new(newrelic.Application)

	// Set up mock service expectations.
	mockDiscordService.On("SendSecurityCouncilMessage", mock.Anything, mock.Anything).
		Return(nil) // Accept both string and bool parameters.
	mockSlackService.On("SendRegistryMessageToSlack", mock.Anything).
		Return(nil) // Do nothing for all Slack message calls.
	mockAlgolia.On("IndexNodes", mock.Anything, mock.Anything).
		Return(nil).
		On("DeleteNode", mock.Anything, mock.Anything).
		Return(nil).
		On("IndexNodeVersions", mock.Anything, mock.Anything).
		Return(nil).
		On("DeleteNodeVersions", mock.Anything, mock.Anything).
		Return(nil)

	// Initialize the mocked implementation with mocked services.
	result := &MockedServerImplementation{
		DripStrictServerImplementation: implementation.NewStrictServerImplementation(
			client, config, mockStorageService, mockPubsubService, mockSlackService, mockDiscordService, mockAlgolia, newRelicApp),
		mockStorageService: mockStorageService,
		mockSlackService:   mockSlackService,
		mockDiscordService: mockDiscordService,
		mockAlgolia:        mockAlgolia,
		mockPubsubService:  mockPubsubService,
	}
	// Integration tests must not submit installation or search telemetry.
	result.MixpanelService = mixpanel.NewApiClient("test", mixpanel.HttpClient(&http.Client{Transport: testAnalyticsTransport{}}))
	return result
}

type testAnalyticsTransport struct{}

// newRegistryHTTPTestServer injects the principal where verified Firebase
// authentication normally does. Routing, private policies and authorization are
// real; this helper deliberately does not validate Firebase credentials.
func newRegistryHTTPTestServer(impl *implementation.DripStrictServerImplementation) *echo.Echo {
	e := echo.New()
	e.Use(middleware.PrivateRegistryData())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if id := c.Request().Header.Get("X-Test-User"); id != "" {
				c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), auth.UserContextKey, &auth.UserDetails{ID: id})))
			}
			return next(c)
		}
	})
	authz := authorization.NewAuthorizationManager(impl.Client, impl.RegistryService, nil)
	drip.RegisterHandlers(e, drip.NewStrictHandler(impl, []drip.StrictMiddlewareFunc{authz.AuthorizationMiddleware()}))
	return e
}

func (testAnalyticsTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("1")), Request: request}, nil
}

func setupTestUser(client *ent.Client) (context.Context, *ent.User) {
	// Create a new context and a test user
	ctx := context.Background()
	testUser := createTestUser(ctx, client)

	// Attach the test user to the context
	ctx = decorateUserInContext(ctx, testUser)

	// Return the context and the created test user
	return ctx, testUser
}

func setupAdminUser(client *ent.Client) (context.Context, *ent.User) {
	// Create a new context to isolate the test setup
	ctx := context.Background()

	// Attempt to create the admin user
	testUser := createAdminUser(ctx, client)

	// Decorate the user in the context
	ctx = decorateUserInContext(ctx, testUser)

	// Return the decorated context and the created user
	return ctx, testUser
}

// Helper function to set up a personal access token
func setupPersonalAccessToken(
	ctx context.Context,
	authz drip.StrictMiddlewareFunc,
	impl *MockedServerImplementation,
	publisherId string) (*string, error) {

	tokenName := "test-token"
	tokenDescription := "test-description"
	res, err := withMiddleware(authz, impl.CreatePersonalAccessToken)(ctx, drip.CreatePersonalAccessTokenRequestObject{
		PublisherId: publisherId,
		Body: &drip.PersonalAccessToken{
			Name:        &tokenName,
			Description: &tokenDescription,
		},
	})
	if err != nil {
		return nil, err
	}

	// Extract the created token from the response
	pat := res.(drip.CreatePersonalAccessToken201JSONResponse).Token
	return pat, nil
}

// Helper function to generate a valid publisher ID based on the pattern "^[a-z][a-z0-9-]*$"
func generatePublisherId() string {
	// Generate a random UUID and use a portion of it for the publisher ID
	rawId := uuid.New().String()
	// Strip hyphens and convert to lowercase to fit the pattern
	id := strings.ToLower(strings.ReplaceAll(rawId, "-", ""))
	// Ensure the ID starts with a letter and follows the pattern
	if match, _ := regexp.MatchString("^[a-z][a-z0-9-]*$", id); match {
		return id
	}
	// If it doesn't match, regenerate a valid ID
	return generatePublisherId()
}

// Helper function to generate a valid node ID based on the pattern "^[a-z][a-z0-9-_]+(\\.[a-z0-9-_]+)*$"
func generateNodeId() string {
	// Generate a random UUID and use a portion of it for the node ID
	rawId := uuid.New().String()
	// Strip hyphens and convert to lowercase to fit the pattern
	id := strings.ToLower(strings.ReplaceAll(rawId, "-", ""))
	// Ensure the ID starts with a letter and follows the pattern
	if match, _ := regexp.MatchString("^[a-z][a-z0-9-_]+(\\.[a-z0-9-_]+)*$", id); match {
		return id
	}
	// If it doesn't match, regenerate a valid ID
	return generateNodeId()
}

// Helper function to generate a random Publisher for testing
func randomPublisher() *drip.Publisher {
	suffix := uuid.New().String()
	publisherId := generatePublisherId()

	description := "test-description-" + suffix
	sourceCodeRepo := "https://github.com/test-repo-" + suffix
	website := "https://test-website-" + suffix + ".com"
	support := "test-support-" + suffix
	logo := "https://test-logo-" + suffix + ".png"
	name := "test-name-" + suffix

	return &drip.Publisher{
		Id:             &publisherId,
		Name:           &name,
		Description:    &description,
		SourceCodeRepo: &sourceCodeRepo,
		Website:        &website,
		Support:        &support,
		Logo:           &logo,
	}
}

// Helper function to generate a random Node for testing
func randomNode() *drip.Node {
	suffix := uuid.New().String()
	nodeId := generateNodeId()

	description := "test-node-description-" + suffix
	author := "test-node-author-" + suffix
	license := "test-node-license-" + suffix
	name := "test-node-name-" + suffix
	tags := []string{"test-node-tag"}
	icon := "https://www.github.com/test-icon-" + suffix + ".svg"
	repository := "https://www.github.com/test-repo-" + suffix

	return &drip.Node{
		Id:          &nodeId,
		Name:        &name,
		Description: &description,
		Author:      &author,
		License:     &license,
		Tags:        &tags,
		Icon:        &icon,
		Repository:  &repository,
	}
}

// Helper function to generate a random NodeVersion for testing
func randomNodeVersion(revision int) *drip.NodeVersion {
	suffix := uuid.New().String()

	version := fmt.Sprintf("1.0.%d", revision)
	changelog := "test-changelog-" + suffix
	dependencies := []string{"test-dependency-" + suffix}

	return &drip.NodeVersion{
		Version:      &version,
		Changelog:    &changelog,
		Dependencies: &dependencies,
	}
}

// Helper function to generate a random comfy node with a random name
func randomComfyNode() drip.ComfyNode {
	return drip.ComfyNode{
		ComfyNodeId:  proto.String(uuid.New().String()),
		InputTypes:   proto.String(`{"required":{"param":"string"}}`),
		Function:     proto.String("test Func"),
		Category:     proto.String("test Category"),
		Description:  proto.String("test Description"),
		Deprecated:   proto.Bool(false),
		Experimental: proto.Bool(false),
		ReturnNames:  proto.String(`["result1", "result2"]`),
		ReturnTypes:  proto.String(`["string", "string"]`),
		OutputIsList: &[]bool{false, false},
	}
}

// Helper function to set up a publisher with a random ID
func setupPublisher(
	ctx context.Context,
	authz drip.StrictMiddlewareFunc,
	impl *MockedServerImplementation) (string, error) {

	publisher := randomPublisher()

	_, err := withMiddleware(authz, impl.CreatePublisher)(ctx, drip.CreatePublisherRequestObject{
		Body: publisher,
	})
	return *publisher.Id, err
}

// Helper function to set up a node with a random ID
func setupNode(
	ctx context.Context,
	authz drip.StrictMiddlewareFunc,
	impl *MockedServerImplementation, publisherId string) (*drip.Node, error) {

	node := randomNode()
	node.Id = proto.String(generateNodeId())
	node.Publisher = &drip.Publisher{
		Id: proto.String(publisherId),
	}

	_, err := withMiddleware(authz, impl.CreateNode)(ctx, drip.CreateNodeRequestObject{
		PublisherId: publisherId,
		Body:        node,
	})
	return node, err
}

func createTestUser(ctx context.Context, client *ent.Client) *ent.User {
	return client.User.Create().
		SetID(uuid.New().String()).
		SetIsApproved(true).
		SetName("integration-test").
		SetEmail("integration-test@gmail.com").
		SaveX(ctx)
}

func createAdminUser(ctx context.Context, client *ent.Client) *ent.User {
	return client.User.Create().
		SetID(uuid.New().String()).
		SetIsApproved(true).
		SetIsAdmin(true).
		SetName("admin").
		SetEmail("admin@gmail.com").
		SaveX(ctx)
}

func decorateUserInContext(ctx context.Context, user *ent.User) context.Context {
	return context.WithValue(ctx, auth.UserContextKey, &auth.UserDetails{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	})
}

// setupDatabaseURL gives schema-based and migration-based tests the same isolated
// PostgreSQL fixture. CI uses Testcontainers when no local override is supplied.
func setupDatabaseURL(t *testing.T, ctx context.Context) (string, func()) {
	t.Helper()
	dsn := os.Getenv("REGISTRY_TEST_DATABASE_URL")
	stop := func() {}
	if dsn == "" {
		container, err := postgres.RunContainer(ctx,
			testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
			postgres.WithDatabase("postgres"), postgres.WithUsername("postgres"), postgres.WithPassword("password"),
			testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(time.Minute)),
		)
		require.NoError(t, err)
		dsn, err = container.ConnectionString(ctx, "sslmode=disable")
		require.NoError(t, err)
		stop = func() { require.NoError(t, container.Terminate(ctx)) }
	}
	database, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	schemaName := "registry_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	_, err = database.ExecContext(ctx, `CREATE SCHEMA "`+schemaName+`"`)
	require.NoError(t, err)
	u, err := url.Parse(dsn)
	require.NoError(t, err)
	q := u.Query()
	q.Set("search_path", schemaName)
	u.RawQuery = q.Encode()
	return u.String(), func() {
		_, err := database.ExecContext(context.Background(), `DROP SCHEMA "`+schemaName+`" CASCADE`)
		require.NoError(t, err)
		require.NoError(t, database.Close())
		stop()
	}
}

func setupDB(t *testing.T, ctx context.Context) (*ent.Client, func()) {
	t.Helper()
	dsn, cleanup := setupDatabaseURL(t, ctx)
	client, err := ent.Open("postgres", dsn)
	require.NoError(t, err)
	require.NoError(t, client.Schema.Create(ctx))
	db.RegisterFeedbackHooks(client)
	return client, func() {
		require.NoError(t, client.Close())
		cleanup()
	}
}

func withMiddleware[R any, S any](mw drip.StrictMiddlewareFunc, h func(ctx context.Context, req R) (res S, err error)) func(ctx context.Context, req R) (res S, err error) {
	// Adapt the provided handler `h` to the signature expected by the middleware.
	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		// Convert the `echo.Context` to a standard `context.Context` and cast the request to the expected type.
		return h(ctx.Request().Context(), request.(R))
	}

	// Use reflection to extract the operation name (function name) of the handler for logging or debugging purposes.
	nameParts := strings.Split(runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name(), ".")
	nameParts = strings.Split(nameParts[len(nameParts)-1], "-")
	opname := nameParts[0] // Isolate the operation name.

	return func(ctx context.Context, req R) (res S, err error) {
		// Create a simulated echo.Context with fake HTTP request and response.
		fakeReq := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
		fakeRes := httptest.NewRecorder()
		fakeCtx := echo.New().NewContext(fakeReq, fakeRes)

		// Wrap the adapted handler with the middleware.
		wrappedHandler := mw(handler, opname)

		// Invoke the middleware-wrapped handler and cast the result to the expected response type.
		result, err := wrappedHandler(fakeCtx, req)
		if result == nil {
			// Return a zero-value of type S if the result is nil.
			return *new(S), err
		}

		// Type assert the result to the expected response type S and return.
		return result.(S), err
	}
}

// Helper function for checking error type and code
func assertHTTPError(t *testing.T, err error, expectedCode int) {
	require.IsType(t, &echo.HTTPError{}, err)
	echoErr := err.(*echo.HTTPError)
	assert.Equal(t, expectedCode, echoErr.Code, "should return correct HTTP error code")
}
