package authentication

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPrivateRegistryRoutesRequireFirebaseIdentity(t *testing.T) {
	for _, path := range []string{"/admin/nodeversions", "/admin/node-version-feedback", "/admin/nodes/n/versions/v/feedback", "/publishers/p/nodes/n/versions/v/feedback", "/users/me/node-version-feedback"} {
		t.Run(path, func(t *testing.T) {
			e := echo.New()
			e.Use(FirebaseAuthMiddleware(nil))
			e.GET(path, func(c echo.Context) error { return c.NoContent(200) })
			res := httptest.NewRecorder()
			e.ServeHTTP(res, httptest.NewRequest("GET", path, nil))
			require.Equal(t, 401, res.Code)
		})
	}
}

func TestPrivateRegistryWritesRejectMissingIdentityAndPATBody(t *testing.T) {
	routes := []struct{ method, path string }{
		{"POST", "/admin/nodes/n/versions/v/feedback/messages"},
		{"POST", "/admin/nodes/n/versions/v/feedback/read"},
		{"PATCH", "/admin/nodes/n/versions/v/feedback"},
		{"POST", "/publishers/p/nodes/n/versions/v/feedback/messages"},
		{"POST", "/publishers/p/nodes/n/versions/v/feedback/read"},
	}
	for _, route := range routes {
		for _, header := range []string{"", "Token personal-access-token"} {
			t.Run(route.method+" "+route.path+" "+header, func(t *testing.T) {
				e := echo.New()
				e.Use(FirebaseAuthMiddleware(nil))
				reached := false
				e.Add(route.method, route.path, func(c echo.Context) error { reached = true; return c.NoContent(200) })
				request := httptest.NewRequest(route.method, route.path, strings.NewReader(`{"personal_access_token":"publisher-pat","sender_user_id":"admin","is_admin":true}`))
				request.Header.Set("Content-Type", "application/json")
				request.Header.Set("Authorization", header)
				response := httptest.NewRecorder()
				e.ServeHTTP(response, request)
				require.Equal(t, 401, response.Code)
				require.False(t, reached, "PAT/role fields must not substitute for verified Firebase identity")
			})
		}
	}
}
