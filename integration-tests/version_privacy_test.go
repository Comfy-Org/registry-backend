package integration

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"registry-backend/config"
	"registry-backend/drip"
	"registry-backend/ent/schema"
	"registry-backend/server/middleware"
)

func TestVersionPrivacyPublicContracts(t *testing.T) {
	ctx := context.Background()
	client, cleanup := setupDB(t, ctx)
	defer cleanup()
	impl := NewStrictServerImplementationWithMocks(client, &config.Config{})
	adminCtx, _ := setupAdminUser(client)
	client.Publisher.Create().SetID("publisher").SetName("Publisher").SaveX(ctx)
	client.Node.Create().SetID("example").SetNormalizedID("example").SetName("Example").SetPublisherID("publisher").SetLicense("MIT").SetRepositoryURL("https://example.invalid/repo").SaveX(ctx)
	version := client.NodeVersion.Create().SetNodeID("example").SetVersion("1.0.0").SetPipDependencies([]string{}).SetStatus(schema.NodeVersionStatusActive).SetStatusReason("private-scan-evidence").SetTagsAdmin([]string{"any-code-execute"}).SaveX(ctx)
	e := echo.New()
	e.Use(middleware.PrivateRegistryData())
	drip.RegisterHandlers(e, drip.NewStrictHandler(impl.DripStrictServerImplementation, nil))
	for _, path := range []string{"/nodes/example", "/nodes/example/versions?include_status_reason=true", "/nodes/example/versions/1.0.0", "/versions?include_status_reason=true", "/nodes/example/install", "/nodes/example/install?version=1.0.0"} {
		t.Run(path, func(t *testing.T) {
			res := httptest.NewRecorder()
			e.ServeHTTP(res, httptest.NewRequest("GET", path, nil))
			require.Equal(t, 200, res.Code, res.Body.String())
			require.NotContains(t, res.Body.String(), "status_reason")
			require.NotContains(t, res.Body.String(), "private-scan-evidence")
			var payload any
			require.NoError(t, json.Unmarshal(res.Body.Bytes(), &payload))
			if path == "/nodes/example" {
				// The standalone node page needs this public identity to check
				// ownership and discover the owner's private feedback.
				publisher, ok := payload.(map[string]any)["publisher"].(map[string]any)
				require.True(t, ok, "node details must identify their Publisher")
				require.Equal(t, "publisher", publisher["id"])
			}
			found := 0
			var check func(any)
			check = func(value any) {
				switch v := value.(type) {
				case map[string]any:
					if v["version"] == "1.0.0" {
						require.Equal(t, []any{"any-code-execute"}, v["tags_admin"])
						found++
					}
					for _, child := range v {
						check(child)
					}
				case []any:
					for _, child := range v {
						check(child)
					}
				}
			}
			check(payload)
			require.Positive(t, found, "version and its installation tags must remain public")
		})
	}
	t.Run("tag-only admin update returns public tags without erasing scanner evidence", func(t *testing.T) {
		tags := []string{"any-code-execute", "requires-review"}
		response, err := impl.AdminUpdateNodeVersion(adminCtx, drip.AdminUpdateNodeVersionRequestObject{NodeId: "example", VersionNumber: "1.0.0", Body: &drip.AdminUpdateNodeVersionJSONRequestBody{TagsAdmin: &tags}})
		require.NoError(t, err)
		require.IsType(t, drip.AdminUpdateNodeVersion200JSONResponse{}, response)
		result := response.(drip.AdminUpdateNodeVersion200JSONResponse)
		require.NotNil(t, result.TagsAdmin)
		require.Equal(t, tags, *result.TagsAdmin)
		stored := client.NodeVersion.GetX(ctx, version.ID)
		require.Equal(t, "private-scan-evidence", stored.StatusReason)
		require.Equal(t, schema.NodeVersionStatusActive, stored.Status)
	})
}
