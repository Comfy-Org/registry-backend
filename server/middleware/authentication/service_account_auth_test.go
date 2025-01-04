package authentication

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestServiceAccountAllowList(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middleware := ServiceAccountAuthMiddleware()

	tests := []struct {
		name    string
		path    string
		method  string
		allowed bool
	}{
		{"OpenAPI GET", "/openapi", "GET", true},
		{"Session DELETE", "/users/sessions", "DELETE", true},
		{"Health GET", "/health", "GET", true},
		{"VM ANY", "/vm", "POST", true},
		{"VM ANY GET", "/vm", "GET", true},
		{"Artifact POST", "/upload-artifact", "POST", true},
		{"Git Commit POST", "/gitcommit", "POST", true},
		{"Git Commit GET", "/gitcommit", "GET", true},
		{"Branch GET", "/branch", "GET", true},
		{"Node Version Path POST", "/publishers/pub123/nodes/node456/versions", "POST", true},
		{"Publisher POST", "/publishers", "POST", true},
		{"Unauthorized Path", "/nonexistent", "GET", true},
		{"Get All Nodes", "/nodes", "GET", true},
		{"Install Nodes", "/nodes/node-id/install", "GET", true},
		{"Get Comfy-Nodes", "/nodes/test/versions/1.0.0/comfy-nodes/test", "GET", true},

		{"Reindex Nodes", "/nodes/reindex", "POST", false},
		{"Reindex Nodes", "/security-scan", "GET", false},
		{"Create Comfy-Nodes", "/nodes/test/versions/1.0.0/comfy-nodes", "POST", false},
		{"Backfill Comfy-Nodes", "/comfy-nodes/backfill", "POST", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			c.SetRequest(req)
			handled := false
			next := echo.HandlerFunc(func(c echo.Context) error {
				handled = true
				return nil
			})
			err := middleware(next)(c)
			if tt.allowed {
				assert.True(t, handled, "Request should be allowed through")
				assert.Nil(t, err)
			} else {
				assert.False(t, handled, "Request should not be allowed through")
				assert.NotNil(t, err)
				httpError, ok := err.(*echo.HTTPError)
				assert.True(t, ok, "Error should be HTTPError")
				assert.Equal(t, http.StatusUnauthorized, httpError.Code)
			}
		})
	}
}
