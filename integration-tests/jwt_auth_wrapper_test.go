package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	drip_middleware "registry-backend/server/middleware"
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
	newAuthMW := func() (echo.MiddlewareFunc, *bool) {
		invoked := false

		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				invoked = true
				return next(c)
			}
		}, &invoked
	}

	e := echo.New()

	t.Run("No JWT", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		c := e.NewContext(req, httptest.NewRecorder())

		next, nextivc := newHandler()
		mw, mwivc := newAuthMW()
		jwtmw := drip_middleware.JWTWrapperMiddleware(client, jwtsecret, mw)
		err := jwtmw(next)(c)
		require.NoError(t, err)
		assert.True(t, *mwivc, "should invoke the wrapped middleware")
		assert.True(t, *nextivc, "should invoke the handler")
	})

	t.Run("Invalid JWT", func(t *testing.T) {
		_, user := setUpTest(client)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID,
			"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		})
		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString([]byte("invallid"))
		require.NoError(t, err, "should not return error")

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		c := e.NewContext(req, httptest.NewRecorder())
		next, nextivc := newHandler()
		mw, mwivc := newAuthMW()
		jwtmw := drip_middleware.JWTWrapperMiddleware(client, jwtsecret, mw)
		err = jwtmw(next)(c)
		require.NoError(t, err, "should not return error")
		assert.True(t, *nextivc, "should invoke the wrapped middleware")
		assert.True(t, *mwivc, "should invoke the handler ")
	})

	t.Run("Valid JWT Invalid User", func(t *testing.T) {
		_, user := setUpTest(client)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID + "Invalids",
			"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		})
		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString([]byte(jwtsecret))
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		c := e.NewContext(req, httptest.NewRecorder())
		next, nextivc := newHandler()
		mw, mwivc := newAuthMW()
		jwtmw := drip_middleware.JWTWrapperMiddleware(client, jwtsecret, mw)
		err = jwtmw(next)(c)
		require.Error(t, err, "should return error")
		assert.False(t, *nextivc, "should not invoke the handler")
		assert.False(t, *mwivc, "should not invoke the wrapped middleware")
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

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		c := e.NewContext(req, httptest.NewRecorder())
		next, nextivc := newHandler()
		mw, mwivc := newAuthMW()
		jwtmw := drip_middleware.JWTWrapperMiddleware(client, jwtsecret, mw)
		err = jwtmw(next)(c)
		require.NoError(t, err, "should not return error")
		assert.True(t, *nextivc, "should invoke the handler")
		assert.False(t, *mwivc, "should not invoke the wrapped middleware")
	})

}
