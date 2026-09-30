package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"registry-backend/config"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/schema"
	"registry-backend/server/middleware/authentication"
	registry "registry-backend/services/registry"
)

func TestAdminVersionFeedbackFilter(t *testing.T) {
	ctx := context.Background()
	client, cleanup := setupDB(t, ctx)
	defer cleanup()
	client.User.Create().SetID("admin").SetIsAdmin(true).SaveX(ctx)
	client.User.Create().SetID("owner").SaveX(ctx)
	for _, id := range []string{"publisher", "new-publisher"} {
		client.Publisher.Create().SetID(id).SetName(id).SaveX(ctx)
	}
	client.PublisherPermission.Create().SetPublisherID("publisher").SetUserID("owner").SetPermission(schema.PublisherPermissionTypeOwner).SaveX(ctx)
	for _, id := range []string{"example", "transferred"} {
		client.Node.Create().SetID(id).SetNormalizedID(id).SetName(id).SetPublisherID("publisher").SetLicense("MIT").SetRepositoryURL("https://example.invalid/repo").SaveX(ctx)
	}
	adminCtx := context.WithValue(ctx, authentication.UserContextKey, &authentication.UserDetails{ID: "admin"})
	ownerCtx := context.WithValue(ctx, authentication.UserContextKey, &authentication.UserDetails{ID: "owner"})
	service := registry.FeedbackService{Client: client}
	create := func(node string, number int, resolved bool, roles ...string) *ent.NodeVersion {
		t.Helper()
		v := client.NodeVersion.Create().SetNodeID(node).SetVersion(fmt.Sprintf("1.%d.0", number)).SetPipDependencies([]string{}).SetStatus(schema.NodeVersionStatusFlagged).SetStatusReason("scan-marker").SetCreateTime(time.Date(2026, 9, 30, 0, number, 0, 0, time.UTC)).SaveX(ctx)
		var response *drip.FeedbackResponse
		binding := drip.FeedbackWriteTarget{PublisherId: "publisher"}
		for _, role := range roles {
			actor := ownerCtx
			if role == "admin" {
				actor = adminCtx
			}
			var err error
			response, err = service.Send(actor, registry.FeedbackTarget{NodeID: node, VersionID: v.ID, PublisherID: "publisher", Admin: role == "admin"}, &drip.FeedbackMessageInput{Target: binding, Body: role + " message", ClientMessageId: uuid.New()})
			require.NoError(t, err)
			binding = response.Target
		}
		if resolved {
			_, err := service.SetState(adminCtx, registry.FeedbackTarget{NodeID: node, VersionID: v.ID, Admin: true}, &drip.FeedbackStateInput{Target: binding, State: "resolved", ExpectedRevision: response.Thread.Revision})
			require.NoError(t, err)
		}
		return v
	}
	untouched := create("example", 0, false)
	untouched = client.NodeVersion.UpdateOne(untouched).SetVersion("99.0.0").SaveX(ctx)
	adminLast := create("example", 1, false, "admin")
	authorLast := create("example", 2, false, "admin", "author")
	adminAgain := create("example", 3, false, "admin", "author", "admin")
	client.NodeVersion.UpdateOne(adminAgain).SetStatus(schema.NodeVersionStatusActive).SaveX(ctx)
	resolvedAuthor := create("example", 4, true, "admin", "author")
	resolvedAdmin := create("example", 5, true, "admin")
	emptyThread := create("example", 6, false)
	client.FeedbackThread.Create().SetVersionID(emptyThread.ID).SetPublisherID("publisher").SetCreatedByUserID("admin").SaveX(ctx)
	transferred := create("transferred", 7, false, "admin", "author")
	client.Node.UpdateOneID("transferred").SetPublisherID("new-publisher").SaveX(ctx)

	impl := NewStrictServerImplementationWithMocks(client, &config.Config{})
	e := newRegistryHTTPTestServer(impl.DripStrictServerImplementation)
	list := func(t *testing.T, user, query string, expectedStatus int) drip.AdminVersionList {
		t.Helper()
		request := httptest.NewRequest("GET", "/admin/nodeversions?"+query, nil)
		request.Header.Set("X-Test-User", user)
		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, request)
		require.Equal(t, expectedStatus, recorder.Code, recorder.Body.String())
		require.Equal(t, "private, no-store", recorder.Header().Get("Cache-Control"))
		var result drip.AdminVersionList
		if expectedStatus == 200 {
			require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &result))
		}
		return result
	}
	ids := func(result drip.AdminVersionList) []string {
		out := []string{}
		for _, v := range result.Versions {
			out = append(out, *v.Id)
		}
		return out
	}
	t.Run("three exclusive states use the latest sender including resolved conversations", func(t *testing.T) {
		for _, tc := range []struct {
			filter   string
			versions []*ent.NodeVersion
		}{
			{"unprocessed", []*ent.NodeVersion{untouched, emptyThread, transferred}},
			{"processed", []*ent.NodeVersion{adminLast, adminAgain, resolvedAdmin}},
			{"needs_response", []*ent.NodeVersion{authorLast, resolvedAuthor}},
		} {
			t.Run(tc.filter, func(t *testing.T) {
				result := list(t, "admin", "feedback_status="+tc.filter, 200)
				expected := []string{}
				for _, v := range tc.versions {
					expected = append(expected, v.ID.String())
				}
				require.ElementsMatch(t, expected, ids(result))
				require.Equal(t, len(expected), result.Total)
			})
		}
		require.Equal(t, 8, list(t, "admin", "", 200).Total)
	})
	t.Run("filtering precedes pagination and composes with moderation filters", func(t *testing.T) {
		first := list(t, "admin", "feedback_status=processed&pageSize=1&page=1", 200)
		second := list(t, "admin", "feedback_status=processed&pageSize=1&page=2", 200)
		require.Equal(t, 3, first.Total)
		require.Equal(t, 3, first.TotalPages)
		require.Equal(t, []string{resolvedAdmin.ID.String()}, ids(first))
		require.Equal(t, []string{adminAgain.ID.String()}, ids(second))
		filtered := list(t, "admin", "feedback_status=processed&nodeId=example&statuses=NodeVersionStatusActive&status_reason=scan-marker", 200)
		require.Equal(t, []string{adminAgain.ID.String()}, ids(filtered))
		require.Equal(t, 1, filtered.Total)
		require.Empty(t, list(t, "admin", "feedback_status=needs_response&nodeId=missing", 200).Versions)
		require.Empty(t, list(t, "admin", "nodeId=", 200).Versions, "an explicitly empty node filter must not broaden to all nodes")
	})
	t.Run("replacement version links filter before pagination", func(t *testing.T) {
		result := list(t, "admin", "nodeId=example&version="+untouched.Version+"&pageSize=1", 200)
		require.Equal(t, []string{untouched.ID.String()}, ids(result))
		require.Equal(t, 1, result.Total)
		require.Empty(t, list(t, "admin", "nodeId=example&version=", 200).Versions)
	})
	t.Run("admin chronology and public version ordering remain distinct", func(t *testing.T) {
		latest := list(t, "admin", "nodeId=example&pageSize=1", 200)
		require.Equal(t, []string{emptyThread.ID.String()}, ids(latest))
		require.Equal(t, "scan-marker", latest.Versions[0].StatusReason)
		response := httptest.NewRecorder()
		e.ServeHTTP(response, httptest.NewRequest("GET", "/nodes/example/versions", nil))
		require.Equal(t, 200, response.Code)
		var versions []drip.NodeVersion
		require.NoError(t, json.Unmarshal(response.Body.Bytes(), &versions))
		require.Equal(t, untouched.ID.String(), *versions[0].Id)
		require.NotContains(t, response.Body.String(), "scan-marker")
	})
	t.Run("old publisher history cannot classify the current publisher conversation", func(t *testing.T) {
		_, err := service.Send(adminCtx, registry.FeedbackTarget{NodeID: "transferred", VersionID: transferred.ID, Admin: true}, &drip.FeedbackMessageInput{Target: drip.FeedbackWriteTarget{PublisherId: "new-publisher"}, Body: "New owner feedback", ClientMessageId: uuid.New()})
		require.NoError(t, err)
		require.Equal(t, []string{transferred.ID.String()}, ids(list(t, "admin", "nodeId=transferred&feedback_status=processed", 200)))
		require.Empty(t, list(t, "admin", "nodeId=transferred&feedback_status=needs_response", 200).Versions)
		require.Empty(t, list(t, "admin", "nodeId=transferred&feedback_status=unprocessed", 200).Versions)
	})
	t.Run("invalid values and unauthorized callers are rejected", func(t *testing.T) {
		list(t, "admin", "feedback_status=unknown", 422)
		list(t, "", "feedback_status=needs_response", 401)
		list(t, "owner", "feedback_status=needs_response", 404)
	})
}
