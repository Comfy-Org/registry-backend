package integration

import (
	"bytes"
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"registry-backend/config"
	"registry-backend/ent/schema"
)

func TestAdminAccessRejectsBannedVersionModeration(t *testing.T) {
	ctx := context.Background()
	client, cleanup := setupDB(t, ctx)
	defer cleanup()
	client.User.Create().SetID("banned-admin").SetIsAdmin(true).SetStatus(schema.UserStatusTypeBanned).SaveX(ctx)
	client.User.Create().SetID("admin").SetIsAdmin(true).SaveX(ctx)
	client.User.Create().SetID("owner").SaveX(ctx)
	publisher := client.Publisher.Create().SetID("publisher").SetName("Publisher").SaveX(ctx)
	client.PublisherPermission.Create().SetPublisherID(publisher.ID).SetUserID("owner").SetPermission(schema.PublisherPermissionTypeOwner).SaveX(ctx)
	node := client.Node.Create().SetID("example").SetNormalizedID("example").SetName("Example").SetPublisherID(publisher.ID).SetLicense("MIT").SetRepositoryURL("https://example.invalid").SaveX(ctx)
	version := client.NodeVersion.Create().SetNodeID(node.ID).SetVersion("1.0.0").SetPipDependencies([]string{}).SetStatus(schema.NodeVersionStatusFlagged).SetStatusReason("scan evidence").SetTagsAdmin([]string{"any-code-execute"}).SaveX(ctx)
	impl := NewStrictServerImplementationWithMocks(client, &config.Config{})
	e := newRegistryHTTPTestServer(impl.DripStrictServerImplementation)
	updatePath := "/admin/nodes/example/versions/1.0.0"
	updateBody := `{"tags_admin":[],"status":"NodeVersionStatusActive","status_reason":"approved"}`
	call := func(t *testing.T, actor, method, path, body string, expected int) {
		t.Helper()
		req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Test-User", actor)
		res := httptest.NewRecorder()
		e.ServeHTTP(res, req)
		require.Equal(t, expected, res.Code, "%s %s as %q: %s", method, path, actor, res.Body.String())
	}
	for _, actor := range []struct {
		id   string
		code int
	}{{"", 401}, {"deleted-admin", 401}, {"owner", 403}, {"banned-admin", 403}} {
		t.Run("denied-"+actor.id, func(t *testing.T) {
			call(t, actor.id, "PUT", updatePath, updateBody, actor.code)
			call(t, actor.id, "POST", "/publishers/publisher/ban", "", actor.code)
			call(t, actor.id, "POST", "/publishers/publisher/nodes/example/ban", "", actor.code)
			stored := client.NodeVersion.GetX(ctx, version.ID)
			require.Equal(t, version.TagsAdmin, stored.TagsAdmin, "denied writes must preserve installation policy")
			require.Equal(t, version.Status, stored.Status)
			require.Equal(t, version.StatusReason, stored.StatusReason)
			require.Equal(t, publisher.Status, client.Publisher.GetX(ctx, publisher.ID).Status)
			require.Equal(t, node.Status, client.Node.GetX(ctx, node.ID).Status)
		})
	}
	t.Run("active-admin", func(t *testing.T) {
		call(t, "admin", "PUT", updatePath, updateBody, 200)
		stored := client.NodeVersion.GetX(ctx, version.ID)
		require.Empty(t, stored.TagsAdmin)
		require.Equal(t, schema.NodeVersionStatusActive, stored.Status)
	})
}
