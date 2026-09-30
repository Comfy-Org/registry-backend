package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"registry-backend/config"
	"registry-backend/drip"
	"registry-backend/ent/publisherpermission"
	"registry-backend/ent/schema"
	"registry-backend/server/middleware/authentication"
	registry "registry-backend/services/registry"
)

// The principal is injected at the boundary normally populated by verified
// Firebase authentication. Routes, strict middleware and service checks are real.
func TestFeedbackAuthorizationMatrix(t *testing.T) {
	ctx := context.Background()
	client, cleanup := setupDB(t, ctx)
	defer cleanup()
	for _, id := range []string{"admin", "banned-admin", "former-admin", "owner", "co-owner", "member", "other-owner", "banned-owner"} {
		create := client.User.Create().SetID(id).SetIsAdmin(id == "admin" || id == "banned-admin")
		if strings.HasPrefix(id, "banned-") {
			create.SetStatus(schema.UserStatusTypeBanned)
		}
		create.SaveX(ctx)
	}
	for _, id := range []string{"publisher-a", "publisher-b"} {
		client.Publisher.Create().SetID(id).SetName(id).SaveX(ctx)
	}
	for _, id := range []string{"owner", "co-owner", "member", "banned-owner"} {
		permission := schema.PublisherPermissionTypeOwner
		if id == "member" {
			permission = schema.PublisherPermissionTypeMember
		}
		client.PublisherPermission.Create().SetPublisherID("publisher-a").SetUserID(id).SetPermission(permission).SaveX(ctx)
	}
	for _, id := range []string{"other-owner", "co-owner"} {
		client.PublisherPermission.Create().SetPublisherID("publisher-b").SetUserID(id).SetPermission(schema.PublisherPermissionTypeOwner).SaveX(ctx)
	}
	var snapshots []*drip.FeedbackResponse
	service := registry.FeedbackService{Client: client}
	adminCtx := context.WithValue(ctx, authentication.UserContextKey, &authentication.UserDetails{ID: "admin"})
	for _, suffix := range []string{"a", "b"} {
		nodeID := "node-" + suffix
		publisherID := "publisher-" + suffix
		client.Node.Create().SetID(nodeID).SetNormalizedID(nodeID).SetName(nodeID).SetPublisherID(publisherID).SetLicense("MIT").SetRepositoryURL("https://example.invalid/repo").SaveX(ctx)
		version := client.NodeVersion.Create().SetNodeID(nodeID).SetVersion("1.0.0").SetPipDependencies([]string{}).SetStatus(schema.NodeVersionStatusFlagged).SetStatusReason("private-scan-" + suffix).SetTagsAdmin([]string{"any-code-execute"}).SaveX(ctx)
		response, err := service.Send(adminCtx, registry.FeedbackTarget{NodeID: nodeID, VersionID: version.ID, Admin: true}, &drip.FeedbackMessageInput{Target: drip.FeedbackWriteTarget{PublisherId: publisherID}, Body: "private-message-" + suffix, ClientMessageId: uuid.New()})
		require.NoError(t, err)
		snapshots = append(snapshots, response)
	}
	impl := NewStrictServerImplementationWithMocks(client, &config.Config{})
	e := newRegistryHTTPTestServer(impl.DripStrictServerImplementation)
	requests := 0
	call := func(actor, method, path string, body any, expected int) map[string]any {
		t.Helper()
		raw, err := json.Marshal(body)
		require.NoError(t, err)
		request := httptest.NewRequest(method, path, bytes.NewReader(raw))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("X-Test-User", actor)
		requests++
		request.RemoteAddr = fmt.Sprintf("192.0.%d.%d:1234", requests/250, requests%250+1)
		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		require.Equal(t, expected, recorder.Code, "%s %s as %s: %s", method, path, actor, recorder.Body.String())
		require.Equal(t, "private, no-store", recorder.Header().Get("Cache-Control"))
		require.Contains(t, recorder.Header().Values("Vary"), "Authorization")
		if expected != 200 {
			require.NotContains(t, recorder.Body.String(), "private-message-")
			require.NotContains(t, recorder.Body.String(), "private-scan-")
		}
		result := map[string]any{}
		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &result))
		return result
	}
	versionID := snapshots[0].Thread.VersionId
	adminPath := fmt.Sprintf("/admin/nodes/node-a/versions/%s/feedback", versionID)
	authorPath := fmt.Sprintf("/publishers/publisher-a/nodes/node-a/versions/%s/feedback", versionID)
	message := map[string]any{"target": snapshots[0].Target, "body": "Request body", "client_message_id": uuid.NewString()}
	read := map[string]any{"target": snapshots[0].Target, "last_read_message_seq": 1}
	state := map[string]any{"target": snapshots[0].Target, "state": "resolved", "expected_revision": snapshots[0].Thread.Revision}
	type endpoint struct {
		method, path string
		body         any
		inbox        bool
	}
	adminEndpoints := []endpoint{
		{"GET", adminPath, nil, false}, {"POST", adminPath + "/messages", message, false},
		{"POST", adminPath + "/read", read, false}, {"PATCH", adminPath, state, false},
		{"GET", "/admin/node-version-feedback?nodeId=node-a", nil, false},
		{"GET", "/admin/nodeversions?feedback_status=processed", nil, false},
	}
	authorEndpoints := []endpoint{
		{"GET", authorPath, nil, false}, {"POST", authorPath + "/messages", message, false},
		{"POST", authorPath + "/read", read, false},
		{"GET", "/users/me/node-version-feedback?publisherId=publisher-a", nil, true},
	}
	t.Run("all private routes reject missing deleted and banned principals", func(t *testing.T) {
		for _, actor := range []string{"", "missing-user", "banned-admin", "banned-owner"} {
			code := 404
			if actor == "" || actor == "missing-user" {
				code = 401
			}
			for _, route := range append(adminEndpoints, authorEndpoints...) {
				call(actor, route.method, route.path, route.body, code)
			}
		}
	})
	t.Run("all admin routes reject owners members strangers and former admins", func(t *testing.T) {
		for _, actor := range []string{"owner", "co-owner", "member", "other-owner", "former-admin"} {
			for _, route := range adminEndpoints {
				call(actor, route.method, route.path, route.body, 404)
			}
		}
	})
	t.Run("members and unrelated owners cannot read reply acknowledge or discover another thread", func(t *testing.T) {
		for _, actor := range []string{"member", "other-owner", "former-admin"} {
			for _, route := range authorEndpoints {
				if route.inbox {
					require.Empty(t, call(actor, route.method, route.path, route.body, 200)["threads"])
				} else {
					call(actor, route.method, route.path, route.body, 404)
				}
			}
		}
		require.Equal(t, 2, client.FeedbackMessage.Query().CountX(ctx), "denied writes must have no side effects")
		require.Zero(t, client.FeedbackRead.Query().CountX(ctx))
	})
	t.Run("resource identifiers cannot be mixed even when both publishers are owned", func(t *testing.T) {
		wrongPaths := []string{
			fmt.Sprintf("/publishers/publisher-b/nodes/node-a/versions/%s/feedback", versionID),
			fmt.Sprintf("/publishers/publisher-a/nodes/node-a/versions/%s/feedback", snapshots[1].Thread.VersionId),
		}
		for _, path := range wrongPaths {
			call("co-owner", "GET", path, nil, 404)
			call("co-owner", "POST", path+"/messages", message, 404)
			call("co-owner", "POST", path+"/read", read, 404)
		}
		call("admin", "GET", fmt.Sprintf("/admin/nodes/node-b/versions/%s/feedback", versionID), nil, 404)
	})
	t.Run("owners can read and reply but cannot impersonate staff or another reader", func(t *testing.T) {
		for _, actor := range []string{"owner", "co-owner"} {
			view := call(actor, "GET", authorPath, nil, 200)
			require.Len(t, view["messages"], 1)
			require.NotContains(t, fmt.Sprint(view), "private-scan-a")
			require.Len(t, call(actor, "GET", "/users/me/node-version-feedback?publisherId=publisher-a", nil, 200)["threads"], 1)
		}
		message["sender_role"], message["sender_user_id"], message["is_admin"] = "admin", "admin", true
		response := call("owner", "POST", authorPath+"/messages", message, 200)
		messages := response["messages"].([]any)
		last := messages[len(messages)-1].(map[string]any)
		require.Equal(t, "author", last["sender_role"])
		require.Equal(t, "owner", last["sender_user_id"])
		read["user_id"] = "co-owner"
		call("owner", "POST", authorPath+"/read", read, 200)
		require.Equal(t, float64(0), call("co-owner", "GET", authorPath, nil, 200)["last_read_message_seq"])
		admin := call("admin", "GET", adminPath, nil, 200)
		versions := call("admin", "GET", "/admin/nodeversions?nodeId=node-a&feedback_status=needs_response", nil, 200)["versions"].([]any)
		require.Len(t, versions, 1)
		require.Equal(t, versionID.String(), versions[0].(map[string]any)["id"])
		require.Equal(t, "private-scan-a", versions[0].(map[string]any)["status_reason"])
		state["expected_revision"] = admin["thread"].(map[string]any)["revision"]
		call("admin", "PATCH", adminPath, state, 200)
	})
	t.Run("revoked ownership and demoted admins cannot reuse an authenticated principal", func(t *testing.T) {
		client.PublisherPermission.Delete().Where(publisherpermission.UserIDEQ("owner")).ExecX(ctx)
		for _, route := range authorEndpoints {
			if route.inbox {
				require.Empty(t, call("owner", route.method, route.path, route.body, 200)["threads"])
			} else {
				call("owner", route.method, route.path, route.body, 404)
			}
		}
		client.User.UpdateOneID("admin").SetIsAdmin(false).ExecX(ctx)
		for _, route := range adminEndpoints {
			call("admin", route.method, route.path, route.body, 404)
		}
	})
}
