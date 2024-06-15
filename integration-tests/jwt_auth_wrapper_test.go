package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"registry-backend/server/middleware/authentication"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJwtAuthWrapper(t *testing.T) {
	clientCtx := context.Background()
	client, postgresContainer := setupDB(t, clientCtx)
	// Cleanup
	defer func() {
		if err := postgresContainer.Terminate(clientCtx); err != nil {
			log.Ctx(clientCtx).Error().Msgf("failed to terminate container: %s", err)
		}
	}()

	jwtsecret := "test"

	newHandler := func() (echo.HandlerFunc, *bool) {
		invoked := false

		return func(c echo.Context) error {
			invoked = true
			return nil
		}, &invoked
	}

	e := echo.New()
	jwtmw := drip_authentication.JWTAdminAuthMiddleware(client, jwtsecret)

	t.Run("No JWT", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/publishers/test-publisher/ban", nil)
		c := e.NewContext(req, httptest.NewRecorder())

		next, nextivc := newHandler()
		err := jwtmw(next)(c)
		require.NoError(t, err)
		assert.True(t, *nextivc, "should invoke the handler")
	})

	t.Run("Invalid JWT", func(t *testing.T) {
		_, user := setUpTest(client)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID,
			"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		})
		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString([]byte("invalid"))
		require.NoError(t, err, "should not return error")

		req := httptest.NewRequest(http.MethodGet, "/publishers/test-publisher/ban", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		c := e.NewContext(req, httptest.NewRecorder())
		next, nextivc := newHandler()
		err = jwtmw(next)(c)
		require.NoError(t, err, "should not return error")
		assert.True(t, *nextivc, "should invoke the wrapped middleware")
	})

	t.Run("Valid JWT Invalid User", func(t *testing.T) {
		_, user := setUpTest(client)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID + "Invalid",
			"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		})
		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString([]byte(jwtsecret))
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/publishers/test-publisher/ban", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		c := e.NewContext(req, httptest.NewRecorder())
		next, nextivc := newHandler()
		err = jwtmw(next)(c)
		require.Error(t, err, "should return error")
		assert.False(t, *nextivc, "should not invoke the handler")
	})

	t.Run("Valid JWT", func(t *testing.T) {
		_, user := setUpTest(client)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID,
			"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		})
		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString([]byte(jwtsecret))
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/publishers/test-publisher/ban", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		c := e.NewContext(req, httptest.NewRecorder())
		next, nextivc := newHandler()
		err = jwtmw(next)(c)
		require.NoError(t, err, "should not return error")
		assert.True(t, *nextivc, "should invoke the handler")
	})

	t.Run("Non-Protected Endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/non-protected-endpoint", nil)
		c := e.NewContext(req, httptest.NewRecorder())

		next, nextivc := newHandler()
		err := jwtmw(next)(c)
		require.NoError(t, err)
		assert.True(t, *nextivc, "should invoke the handler for non-protected endpoint")
	})
}
