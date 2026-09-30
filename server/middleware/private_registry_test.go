package middleware

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPrivateRegistryErrorsAreNotCacheable(t *testing.T) {
	e := echo.New()
	e.Use(PrivateRegistryData())
	e.GET("/admin/nodeversions", func(c echo.Context) error { return echo.NewHTTPError(401) })
	res := httptest.NewRecorder()
	e.ServeHTTP(res, httptest.NewRequest("GET", "/admin/nodeversions", nil))
	require.Equal(t, 401, res.Code)
	require.Equal(t, "private, no-store", res.Header().Get("Cache-Control"))
}

func TestPublicFeedbackNamesRetainTheirRequestPolicy(t *testing.T) {
	for _, route := range []struct{ method, path string }{
		{"GET", "/nodes/feedback"},
		{"GET", "/nodes/feedback/install"},
		{"GET", "/publishers/feedback"},
		{"GET", "/publishers/feedback/nodes"},
		{"POST", "/publishers/feedback/nodes/example/versions"},
		{"POST", "/publishers/p/nodes/feedback/versions"},
	} {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			e := echo.New()
			e.Use(PrivateRegistryData())
			e.Add(route.method, route.path, func(c echo.Context) error { return c.NoContent(204) })
			request := httptest.NewRequest(route.method, route.path, strings.NewReader(strings.Repeat("x", 33000)))
			response := httptest.NewRecorder()
			e.ServeHTTP(response, request)
			require.Equal(t, 204, response.Code, "feedback write limits must not apply to version publishing")
			require.Empty(t, response.Header().Get("Cache-Control"), "public names must not acquire the private cache policy")
		})
	}
}
func TestPrivateRegistryDoesNotLogContentOrTokens(t *testing.T) {
	for _, path := range []string{
		"/admin/nodes/example/versions/version/feedback/messages",
		"/publishers/feedback/nodes/feedback/versions/version/feedback/messages",
	} {
		t.Run(path, func(t *testing.T) {
			var output bytes.Buffer
			logger := zerolog.New(&output)
			e := echo.New()
			e.Use(RequestLoggerMiddleware())
			e.Use(ResponseLoggerMiddleware())
			e.POST(path, func(c echo.Context) error { return c.String(200, "private-reply") })
			req := httptest.NewRequest("POST", path, strings.NewReader("private-message"))
			req.Header.Set("Authorization", "Bearer private-token")
			req = req.WithContext(logger.WithContext(req.Context()))
			e.ServeHTTP(httptest.NewRecorder(), req)
			require.NotContains(t, output.String(), "private-message")
			require.NotContains(t, output.String(), "private-reply")
			require.NotContains(t, output.String(), "private-token")
		})
	}
}

func TestVersionFeedbackBodyLimit(t *testing.T) {
	for _, base := range []string{
		"/admin/nodes/feedback/versions/version/feedback",
		"/publishers/feedback/nodes/feedback/versions/version/feedback",
	} {
		for _, action := range []struct{ method, suffix string }{
			{"POST", "/messages"}, {"POST", "/read"}, {"PATCH", ""},
		} {
			path := base + action.suffix
			t.Run(action.method+" "+path, func(t *testing.T) {
				e := echo.New()
				e.Use(PrivateRegistryData())
				e.Add(action.method, path, func(c echo.Context) error { return c.NoContent(204) })
				for _, size := range []int{100, 33000} {
					res := httptest.NewRecorder()
					e.ServeHTTP(res, httptest.NewRequest(action.method, path, strings.NewReader(strings.Repeat("x", size))))
					want := 204
					if size > 32768 {
						want = 413
					}
					require.Equal(t, want, res.Code)
					require.Equal(t, "private, no-store", res.Header().Get("Cache-Control"))
				}
			})
		}
	}
}
func TestPublicPrivateFiltersAreRejected(t *testing.T) {
	for _, query := range []string{"status_reason=secret", "feedback_status=needs_response"} {
		t.Run(query, func(t *testing.T) {
			e := echo.New()
			e.Use(PrivateRegistryData())
			e.GET("/versions", func(c echo.Context) error { return c.NoContent(200) })
			res := httptest.NewRecorder()
			e.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/versions?"+query, nil))
			require.Equal(t, 400, res.Code)
		})
	}
}
