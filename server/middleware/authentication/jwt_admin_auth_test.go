package drip_authentication

import (
	"net/http"
	"net/http/httptest"
	"registry-backend/ent"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestJWTAdminAllowlist(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock ent.Client
	mockEntClient := &ent.Client{}

	middleware := JWTAdminAuthMiddleware(mockEntClient, "secret")

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

		{"Ban Publisher", "/publishers/publisher-id/ban", "POST", false},
		{"Ban Node", "/publishers/publisher-id/nodes/node-id/ban", "POST", false},
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
