package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"registry-backend/config"
	"registry-backend/drip"
	"registry-backend/ent/schema"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPrivateVersionFeedback(t *testing.T) {
	ctx := context.Background()
	client, cleanup := setupDB(t, ctx)
	defer cleanup()
	for _, id := range []string{"admin", "owner", "member", "stranger"} {
		client.User.Create().SetID(id).SetName(id).SetIsAdmin(id == "admin").SaveX(ctx)
	}
	client.Publisher.Create().SetID("publisher").SetName("Publisher").SaveX(ctx)
	for _, id := range []string{"owner", "member"} {
		permission := schema.PublisherPermissionTypeOwner
		if id == "member" {
			permission = schema.PublisherPermissionTypeMember
		}
		client.PublisherPermission.Create().SetPublisherID("publisher").SetUserID(id).SetPermission(permission).SaveX(ctx)
	}
	client.Node.Create().SetID("example").SetNormalizedID("example").SetLicense("MIT").SetRepositoryURL("https://example.invalid/repo").SetPublisherID("publisher").SetName("Example").SaveX(ctx)
	version := client.NodeVersion.Create().SetNodeID("example").SetVersion("1.0.0").SetPipDependencies([]string{}).SetStatus(schema.NodeVersionStatusFlagged).SetStatusReason("private-scan-evidence").SetTagsAdmin([]string{"any-code-execute"}).SaveX(ctx)
	impl := NewStrictServerImplementationWithMocks(client, &config.Config{})
	e := newRegistryHTTPTestServer(impl.DripStrictServerImplementation)
	spec, err := drip.GetSwagger()
	require.NoError(t, err)
	errorSchema := spec.Paths.Find("/admin/node-version-feedback").Get.Responses.Default().Value.Content.Get("application/json").Schema.Value
	call := func(user, method, path string, body any, status int) map[string]any {
		t.Helper()
		raw, err := json.Marshal(body)
		require.NoError(t, err)
		req := httptest.NewRequest(method, path, bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Test-User", user)
		res := httptest.NewRecorder()
		e.ServeHTTP(res, req)
		require.Equal(t, status, res.Code, res.Body.String())
		require.Equal(t, "private, no-store", res.Header().Get("Cache-Control"))
		var result map[string]any
		if res.Body.Len() > 0 {
			require.NoError(t, json.Unmarshal(res.Body.Bytes(), &result))
		}
		if status >= 400 {
			require.NoError(t, errorSchema.VisitJSON(result), "feedback errors must match the advertised API schema")
		}
		return result
	}
	adminPath := fmt.Sprintf("/admin/nodes/example/versions/%s/feedback", version.ID)
	ownerPath := fmt.Sprintf("/publishers/publisher/nodes/example/versions/%s/feedback", version.ID)
	assertAdminInboxState := func(expected string) {
		t.Helper()
		for _, state := range []string{"awaiting_author", "awaiting_admin", "resolved"} {
			result := call("admin", "GET", "/admin/node-version-feedback?state="+state, nil, 200)
			threads := result["threads"].([]any)
			if state == expected {
				require.Len(t, threads, 1)
				thread := threads[0].(map[string]any)["thread"].(map[string]any)
				require.Equal(t, version.ID.String(), thread["version_id"])
				require.Equal(t, expected, thread["state"])
			} else {
				require.Empty(t, threads, "a different conversation state must not match %s", state)
			}
		}
	}
	for _, user := range []string{"owner", "member", "stranger"} {
		call(user, "POST", adminPath+"/messages", map[string]any{"body": "not permitted", "client_message_id": uuid.NewString()}, 404)
		call(user, "GET", "/admin/node-version-feedback?state=awaiting_author", nil, 404)
	}
	call("", "GET", "/admin/node-version-feedback?state=awaiting_author", nil, 401)
	call("admin", "GET", "/admin/node-version-feedback?state=unknown", nil, 422)
	call("", "GET", ownerPath, nil, 401)
	call("owner", "POST", ownerPath+"/messages", map[string]any{"target": drip.FeedbackWriteTarget{PublisherId: "publisher"}, "body": "cannot start", "client_message_id": uuid.NewString()}, 409)
	initial := call("admin", "GET", adminPath, nil, 200)
	binding := initial["target"]
	messageID := uuid.NewString()
	first := call("admin", "POST", adminPath+"/messages", map[string]any{"target": binding, "body": "Please revise installation.", "client_message_id": messageID}, 200)
	require.Equal(t, "awaiting_author", first["thread"].(map[string]any)["state"])
	assertAdminInboxState("awaiting_author")
	for _, inbox := range []struct{ user, path string }{
		{"admin", "/admin/node-version-feedback"},
		{"owner", "/users/me/node-version-feedback"},
	} {
		matching := call(inbox.user, "GET", inbox.path+"?state=awaiting_author&nodeId=example&publisherId=publisher", nil, 200)
		require.Len(t, matching["threads"], 1)
		for _, filter := range []string{"nodeId", "publisherId"} {
			empty := call(inbox.user, "GET", inbox.path+"?state=awaiting_author&"+filter+"=", nil, 200)
			require.Empty(t, empty["threads"], "an explicit empty %s must not broaden the %s queue", filter, inbox.user)
		}
		call(inbox.user, "GET", inbox.path+"?state=", nil, 422)
	}
	require.NotContains(t, fmt.Sprint(first), "private-scan-evidence")
	duplicate := call("admin", "POST", adminPath+"/messages", map[string]any{"target": binding, "body": "Please revise installation.", "client_message_id": messageID}, 200)
	require.Len(t, duplicate["messages"], 1)
	call("admin", "POST", adminPath+"/messages", map[string]any{"target": binding, "body": "Changed payload", "client_message_id": messageID}, 409)
	binding = first["target"]
	for _, user := range []string{"member", "stranger"} {
		call(user, "GET", ownerPath, nil, 404)
	}
	call("owner", "GET", fmt.Sprintf("/publishers/publisher/nodes/wrong/versions/%s/feedback", version.ID), nil, 404)
	owner := call("owner", "GET", ownerPath, nil, 200)
	require.Equal(t, float64(1), owner["unread_count"])
	call("owner", "POST", ownerPath+"/read", map[string]any{"target": binding, "last_read_message_seq": 1}, 200)
	owner = call("owner", "GET", ownerPath, nil, 200)
	require.Equal(t, float64(0), owner["unread_count"])
	call("owner", "POST", ownerPath+"/messages", map[string]any{"target": binding, "body": "  ", "client_message_id": uuid.NewString()}, 422)
	reply := call("owner", "POST", ownerPath+"/messages", map[string]any{"target": binding, "body": "Fixed <script>alert(1)</script>", "client_message_id": uuid.NewString()}, 200)
	thread := reply["thread"].(map[string]any)
	require.Equal(t, "awaiting_admin", thread["state"])
	assertAdminInboxState("awaiting_admin")
	admin := call("admin", "GET", adminPath, nil, 200)
	require.Equal(t, float64(1), admin["unread_count"])
	call("admin", "PATCH", adminPath, map[string]any{"target": binding, "state": "resolved", "expected_revision": 0}, 409)
	resolved := call("admin", "PATCH", adminPath, map[string]any{"target": binding, "state": "resolved", "expected_revision": thread["revision"]}, 200)
	assertAdminInboxState("resolved")
	require.Equal(t, schema.NodeVersionStatusFlagged, client.NodeVersion.GetX(ctx, version.ID).Status)
	call("owner", "POST", ownerPath+"/messages", map[string]any{"target": binding, "body": "Late reply", "client_message_id": uuid.NewString()}, 409)
	call("admin", "PATCH", adminPath, map[string]any{"target": binding, "state": "open", "expected_revision": resolved["thread"].(map[string]any)["revision"]}, 200)
	assertAdminInboxState("awaiting_admin")
	scanPath := "/admin/nodeversions?nodeId=example"
	for _, user := range []string{"owner", "member", "stranger"} {
		call(user, "GET", scanPath, nil, 404)
	}
	versions := call("admin", "GET", scanPath, nil, 200)["versions"].([]any)
	require.Len(t, versions, 1)
	scan := versions[0].(map[string]any)
	require.Equal(t, version.ID.String(), scan["id"])
	require.Equal(t, "private-scan-evidence", scan["status_reason"])
	require.Equal(t, []any{"any-code-execute"}, scan["tags_admin"])
	inbox := call("owner", "GET", "/users/me/node-version-feedback", nil, 200)
	require.Len(t, inbox["threads"], 1)
	empty := call("stranger", "GET", "/users/me/node-version-feedback", nil, 200)
	require.Empty(t, empty["threads"])
	call("owner", "POST", ownerPath+"/messages", map[string]any{"target": binding, "body": strings.Repeat("x", 5001), "client_message_id": uuid.NewString()}, 422)
	for _, seq := range []int{-1, 9999} {
		call("owner", "POST", ownerPath+"/read", map[string]any{"target": binding, "last_read_message_seq": seq}, 422)
		require.Equal(t, float64(1), call("owner", "GET", ownerPath, nil, 200)["last_read_message_seq"], "invalid acknowledgements must not change the saved read position")
	}
	// Exactly 5,000 code points are valid after trimming, including characters
	// whose UTF-8 representation occupies more than one byte.
	for _, body := range []string{strings.Repeat("x", 5000), strings.Repeat("😀", 5000)} {
		call("owner", "POST", ownerPath+"/messages", map[string]any{"target": binding, "body": " \n" + body + "\t ", "client_message_id": uuid.NewString()}, 200)
		stored := call("owner", "GET", ownerPath, nil, 200)["messages"].([]any)
		require.Equal(t, body, stored[len(stored)-1].(map[string]any)["body"], "the complete trimmed message must survive a subsequent read")
	}
	client.Publisher.Create().SetID("new-publisher").SetName("New Publisher").SaveX(ctx)
	client.Node.UpdateOneID("example").SetPublisherID("new-publisher").ExecX(ctx)
	call("owner", "GET", ownerPath, nil, 404)
	client.Node.UpdateOneID("example").SetPublisherID("publisher").ExecX(ctx)
	archived := call("owner", "GET", ownerPath, nil, 200)
	require.Equal(t, true, archived["thread"].(map[string]any)["archived"])
	require.Equal(t, false, archived["permissions"].(map[string]any)["can_reply"])
	// Revoking ownership must revoke access even if the thread was previously read.
	client.PublisherPermission.Delete().ExecX(ctx)
	call("owner", "GET", ownerPath, nil, 404)
}
