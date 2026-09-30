package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"registry-backend/config"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/feedbackevent"
	"registry-backend/ent/feedbackthread"
	"registry-backend/ent/publisherpermission"
	"registry-backend/ent/schema"
	"registry-backend/server/middleware/authentication"
	registry "registry-backend/services/registry"
)

func TestOwnerSupersedesEarlierFeedback(t *testing.T) {
	ctx := context.Background()
	client, cleanup := setupDB(t, ctx)
	defer cleanup()
	for _, id := range []string{"admin", "owner", "co-owner", "member", "other", "banned"} {
		c := client.User.Create().SetID(id).SetIsAdmin(id == "admin")
		if id == "banned" {
			c.SetStatus(schema.UserStatusTypeBanned)
		}
		c.SaveX(ctx)
	}
	for _, id := range []string{"publisher", "other-publisher"} {
		client.Publisher.Create().SetID(id).SetName(id).SaveX(ctx)
	}
	for _, id := range []string{"owner", "co-owner", "member", "banned"} {
		permission := schema.PublisherPermissionTypeOwner
		if id == "member" {
			permission = schema.PublisherPermissionTypeMember
		}
		client.PublisherPermission.Create().SetPublisherID("publisher").SetUserID(id).SetPermission(permission).SaveX(ctx)
	}
	for _, id := range []string{"example", "another"} {
		client.Node.Create().SetID(id).SetNormalizedID(id).SetName(id).SetPublisherID("publisher").SetLicense("MIT").SetRepositoryURL("https://example.invalid").SaveX(ctx)
	}
	when := time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)
	makeVersion := func(nodeID, number string, offset int) *ent.NodeVersion {
		return client.NodeVersion.Create().SetNodeID(nodeID).SetVersion(number).SetPipDependencies([]string{}).SetStatus(schema.NodeVersionStatusFlagged).SetStatusReason("private-scan").SetTagsAdmin([]string{"any-code-execute"}).SetCreateTime(when.Add(time.Duration(offset) * time.Hour)).SaveX(ctx)
	}
	older := makeVersion("example", "1.0.0", 0)
	second := makeVersion("example", "1.0.1", 1)
	noFeedback := makeVersion("example", "1.0.2", 2)
	replacement := makeVersion("example", "1.0.3", 3)
	foreign := makeVersion("another", "1.0.0", 0)
	service := registry.FeedbackService{Client: client}
	adminCtx := context.WithValue(ctx, authentication.UserContextKey, &authentication.UserDetails{ID: "admin"})
	start := func(v *ent.NodeVersion) *drip.FeedbackResponse {
		t.Helper()
		result, err := service.Send(adminCtx, registry.FeedbackTarget{NodeID: v.NodeID, VersionID: v.ID, Admin: true}, &drip.FeedbackMessageInput{Target: drip.FeedbackWriteTarget{PublisherId: "publisher"}, Body: "Please revise this version.", ClientMessageId: uuid.New()})
		require.NoError(t, err)
		return result
	}
	firstFeedback, secondFeedback, foreignFeedback := start(older), start(second), start(foreign)
	item := func(v *ent.NodeVersion, f *drip.FeedbackResponse) map[string]any {
		return map[string]any{"version_id": v.ID, "target": f.Target, "expected_revision": f.Thread.Revision}
	}
	inputs := []map[string]any{item(older, firstFeedback), item(second, secondFeedback)}
	impl := NewStrictServerImplementationWithMocks(client, &config.Config{})
	server := newRegistryHTTPTestServer(impl.DripStrictServerImplementation)
	path := fmt.Sprintf("/publishers/publisher/nodes/example/versions/%s/feedback/supersede", replacement.ID)
	calls := 0
	call := func(actor, url string, versions []map[string]any, want int) {
		t.Helper()
		raw, err := json.Marshal(map[string]any{"versions": versions})
		require.NoError(t, err)
		req := httptest.NewRequest("POST", url, bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Test-User", actor)
		calls++
		req.RemoteAddr = fmt.Sprintf("192.0.2.%d:1234", calls)
		res := httptest.NewRecorder()
		server.ServeHTTP(res, req)
		require.Equal(t, want, res.Code, res.Body.String())
		require.Equal(t, "private, no-store", res.Header().Get("Cache-Control"))
	}
	untouched := func() {
		t.Helper()
		require.False(t, client.NodeVersion.GetX(ctx, older.ID).Deprecated)
		require.False(t, client.NodeVersion.GetX(ctx, second.ID).Deprecated)
		require.Equal(t, feedbackthread.StateAwaitingAuthor, client.FeedbackThread.GetX(ctx, firstFeedback.Thread.Id).State)
		require.Equal(t, 0, client.FeedbackEvent.Query().CountX(ctx))
	}
	// A stale member or a user from another Publisher cannot close/deprecate anything.
	for _, actor := range []string{"", "member", "other", "banned", "admin"} {
		want := 404
		if actor == "" {
			want = 401
		}
		call(actor, path, inputs, want)
		untouched()
	}
	// Validate the whole displayed batch before changing either resource.
	stale := item(second, secondFeedback)
	stale["expected_revision"] = 0
	call("owner", path, []map[string]any{inputs[0], stale}, 409)
	untouched()
	call("owner", path, []map[string]any{inputs[0], item(foreign, foreignFeedback)}, 404)
	untouched()
	call("owner", path, []map[string]any{inputs[0], {"version_id": noFeedback.ID, "target": firstFeedback.Target, "expected_revision": 1}}, 404)
	untouched()
	call("owner", path, []map[string]any{inputs[0], inputs[0]}, 422)
	untouched()
	call("owner", path, nil, 422)
	untouched()
	oldReplacementPath := fmt.Sprintf("/publishers/publisher/nodes/example/versions/%s/feedback/supersede", older.ID)
	call("owner", oldReplacementPath, []map[string]any{inputs[1]}, 409)
	untouched()
	client.NodeVersion.UpdateOne(replacement).SetDeprecated(true).ExecX(ctx)
	call("owner", path, inputs, 409)
	untouched()
	client.NodeVersion.UpdateOne(replacement).SetDeprecated(false).ExecX(ctx)
	client.FeedbackThread.UpdateOneID(firstFeedback.Thread.Id).SetArchivedAt(time.Now()).ExecX(ctx)
	call("owner", path, inputs, 409)
	untouched()
	client.FeedbackThread.UpdateOneID(firstFeedback.Thread.Id).ClearArchivedAt().ExecX(ctx)
	client.NodeVersion.UpdateOne(older).SetStatus(schema.NodeVersionStatusDeleted).ExecX(ctx)
	call("owner", path, inputs, 409)
	untouched()
	client.NodeVersion.UpdateOne(older).SetStatus(schema.NodeVersionStatusFlagged).ExecX(ctx)
	client.NodeVersion.UpdateOne(replacement).SetStatus(schema.NodeVersionStatusDeleted).ExecX(ctx)
	call("owner", path, inputs, 409)
	untouched()
	client.NodeVersion.UpdateOne(replacement).SetStatus(schema.NodeVersionStatusFlagged).ExecX(ctx)
	// Current ownership, rather than a remembered permission, authorizes the write.
	client.PublisherPermission.Delete().Where(publisherpermission.UserIDEQ("owner")).ExecX(ctx)
	call("owner", path, inputs, 404)
	untouched()

	// Failure while writing the audit record rolls back deprecation and closure.
	failAudit := true
	client.FeedbackEvent.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, mutation ent.Mutation) (ent.Value, error) {
			if failAudit {
				return nil, errors.New("injected audit failure")
			}
			return next.Mutate(ctx, mutation)
		})
	})
	call("co-owner", path, inputs, 500)
	untouched()
	impl.mockAlgolia.AssertNumberOfCalls(t, "IndexNodeVersions", 0)
	failAudit = false
	call("co-owner", path, inputs, 204)
	for _, pair := range []struct {
		version  *ent.NodeVersion
		feedback *drip.FeedbackResponse
	}{{older, firstFeedback}, {second, secondFeedback}} {
		v := client.NodeVersion.GetX(ctx, pair.version.ID)
		require.True(t, v.Deprecated)
		require.Equal(t, schema.NodeVersionStatusFlagged, v.Status)
		require.Equal(t, "private-scan", v.StatusReason)
		require.Equal(t, []string{"any-code-execute"}, v.TagsAdmin)
		thread := client.FeedbackThread.GetX(ctx, pair.feedback.Thread.Id)
		require.Equal(t, feedbackthread.StateResolved, thread.State)
		require.Equal(t, 2, thread.Revision)
		require.Equal(t, "co-owner", thread.ResolvedByUserID)
		event := client.FeedbackEvent.Query().Where(feedbackevent.ThreadIDEQ(thread.ID)).OnlyX(ctx)
		require.Equal(t, "superseded", string(event.EventType))
		require.Equal(t, "co-owner", event.ActorUserID)
		// Validate the public feedback response, including its replacement reference.
		view, err := service.Get(adminCtx, registry.FeedbackTarget{NodeID: "example", VersionID: v.ID, Admin: true}, nil)
		require.NoError(t, err)
		encoded, err := json.Marshal(view.Events[0])
		require.NoError(t, err)
		var wire map[string]any
		require.NoError(t, json.Unmarshal(encoded, &wire))
		require.Equal(t, replacement.ID.String(), wire["replacement_version_id"])
		require.Equal(t, replacement.Version, wire["replacement_version"])
	}
	require.False(t, client.NodeVersion.GetX(ctx, replacement.ID).Deprecated)
	require.False(t, client.NodeVersion.GetX(ctx, noFeedback.ID).Deprecated)
	require.False(t, client.NodeVersion.GetX(ctx, foreign.ID).Deprecated)
	require.Equal(t, 3, client.FeedbackMessage.Query().CountX(ctx), "closing must not post a message or change the last sender")
	require.Equal(t, 3, client.FeedbackThread.Query().CountX(ctx), "a reviewer starts feedback on the replacement separately")
	// A lost success response can be retried without duplicate events or revisions.
	call("co-owner", path, inputs, 204)
	require.Equal(t, 2, client.FeedbackEvent.Query().CountX(ctx))
	require.Equal(t, 2, client.FeedbackThread.GetX(ctx, firstFeedback.Thread.Id).Revision)
	impl.mockAlgolia.AssertNumberOfCalls(t, "IndexNodeVersions", 2)
	for _, indexCall := range impl.mockAlgolia.Calls {
		if indexCall.Method == "IndexNodeVersions" {
			indexed := indexCall.Arguments.Get(1).([]*ent.NodeVersion)
			require.Len(t, indexed, 2)
			for _, version := range indexed {
				require.True(t, version.Deprecated)
			}
		}
	}
	alreadyClosed := item(older, firstFeedback)
	alreadyClosed["expected_revision"] = 2
	call("co-owner", path, []map[string]any{alreadyClosed}, 409)
	require.Equal(t, 2, client.FeedbackThread.GetX(ctx, firstFeedback.Thread.Id).Revision)
	wrongRevision := item(older, firstFeedback)
	wrongRevision["expected_revision"] = 0
	call("co-owner", path, []map[string]any{wrongRevision}, 409)
	newer := makeVersion("example", "1.0.4", 4)
	call("co-owner", fmt.Sprintf("/publishers/publisher/nodes/example/versions/%s/feedback/supersede", newer.ID), inputs, 409)
	// Reopening invalidates a lost-response retry; it cannot silently close new work.
	_, err := service.SetState(adminCtx, registry.FeedbackTarget{NodeID: "example", VersionID: older.ID, Admin: true}, &drip.FeedbackStateInput{Target: firstFeedback.Target, State: "open", ExpectedRevision: 2})
	require.NoError(t, err)
	call("co-owner", path, inputs, 409)
	require.Equal(t, feedbackthread.StateAwaitingAuthor, client.FeedbackThread.GetX(ctx, firstFeedback.Thread.Id).State)
	// The old route cannot act after a Publisher transfer.
	client.Node.UpdateOneID("example").SetPublisherID("other-publisher").ExecX(ctx)
	call("co-owner", path, inputs, 404)
}
