package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"net/url"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"registry-backend/config"
	"registry-backend/drip"
	"registry-backend/ent/schema"
	"registry-backend/server/middleware/authentication"
	registry "registry-backend/services/registry"
)

func TestFeedbackHTTPPagination(t *testing.T) {
	ctx := context.Background()
	client, cleanup := setupDB(t, ctx)
	defer cleanup()
	client.User.Create().SetID("admin").SetIsAdmin(true).SaveX(ctx)
	client.User.Create().SetID("owner").SaveX(ctx)
	client.Publisher.Create().SetID("publisher").SetName("Publisher").SaveX(ctx)
	client.PublisherPermission.Create().SetPublisherID("publisher").SetUserID("owner").SetPermission(schema.PublisherPermissionTypeOwner).SaveX(ctx)
	client.Node.Create().SetID("example").SetNormalizedID("example").SetName("Example").SetPublisherID("publisher").SetLicense("MIT").SetRepositoryURL("https://example.invalid").SaveX(ctx)
	version := client.NodeVersion.Create().SetNodeID("example").SetVersion("1.0.0").SetPipDependencies([]string{}).SetStatus(schema.NodeVersionStatusFlagged).SaveX(ctx)
	adminCtx := context.WithValue(ctx, authentication.UserContextKey, &authentication.UserDetails{ID: "admin"})
	service := registry.FeedbackService{Client: client}
	target := registry.FeedbackTarget{NodeID: "example", VersionID: version.ID, Admin: true}
	binding := drip.FeedbackWriteTarget{PublisherId: "publisher"}
	for seq := 1; seq <= 65; seq++ {
		response, err := service.Send(adminCtx, target, &drip.FeedbackMessageInput{Target: binding, Body: fmt.Sprintf("Message %d", seq), ClientMessageId: uuid.New()})
		require.NoError(t, err)
		binding = response.Target
	}
	impl := NewStrictServerImplementationWithMocks(client, &config.Config{})
	server := newRegistryHTTPTestServer(impl.DripStrictServerImplementation)
	get := func(path string, status int, result any) {
		t.Helper()
		request := httptest.NewRequest("GET", path, nil)
		request.Header.Set("X-Test-User", "owner")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		require.Equal(t, status, response.Code, response.Body.String())
		if result != nil {
			require.NoError(t, json.Unmarshal(response.Body.Bytes(), result))
		}
	}
	path := fmt.Sprintf("/publishers/publisher/nodes/example/versions/%s/feedback", version.ID)
	// Exact boundaries are a public API contract, independent of the SQL query.
	pages := []struct {
		query             string
		first, last, next int
	}{
		{"", 36, 65, 36}, {"?before_seq=36", 6, 35, 6}, {"?before_seq=6", 1, 5, 0},
	}
	seen := map[uuid.UUID]bool{}
	for _, page := range pages {
		var response drip.FeedbackResponse
		get(path+page.query, 200, &response)
		require.Len(t, response.Messages, page.last-page.first+1)
		for i, message := range response.Messages {
			require.Equal(t, page.first+i, message.Seq)
			require.Equal(t, fmt.Sprintf("Message %d", page.first+i), message.Body)
			require.False(t, seen[message.Id], "pages must not repeat a message")
			seen[message.Id] = true
		}
		if page.next == 0 {
			require.Nil(t, response.NextBeforeSeq)
		} else {
			require.NotNil(t, response.NextBeforeSeq)
			require.Equal(t, page.next, *response.NextBeforeSeq)
		}
		require.Zero(t, response.LastReadMessageSeq, "fetching history must not mark it read")
		require.Equal(t, 65, response.UnreadCount)
	}
	require.Len(t, seen, 65)
	get(path+"?before_seq=0", 422, nil)

	// Exercise real inbox cursors across more than twenty conversations.
	when := time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)
	client.FeedbackThread.UpdateOneID(*binding.ThreadId).SetLastMessageAt(when).ExecX(ctx)
	expected := []uuid.UUID{*binding.ThreadId}
	for i := 1; i <= 22; i++ {
		v := client.NodeVersion.Create().SetNodeID("example").SetVersion(fmt.Sprintf("1.0.%d", i)).SetPipDependencies([]string{}).SetStatus(schema.NodeVersionStatusFlagged).SaveX(ctx)
		response, err := service.Send(adminCtx, registry.FeedbackTarget{NodeID: "example", VersionID: v.ID, Admin: true}, &drip.FeedbackMessageInput{Target: drip.FeedbackWriteTarget{PublisherId: "publisher"}, Body: "Review", ClientMessageId: uuid.New()})
		require.NoError(t, err)
		client.FeedbackThread.UpdateOneID(response.Thread.Id).SetLastMessageAt(when.Add(time.Duration(i) * time.Minute)).ExecX(ctx)
		expected = append([]uuid.UUID{response.Thread.Id}, expected...)
	}
	checkInboxPages := func() {
		t.Helper()
		var inbox drip.FeedbackInbox
		get("/users/me/node-version-feedback", 200, &inbox)
		require.Len(t, inbox.Threads, 20)
		require.NotNil(t, inbox.NextCursor)
		require.NotEmpty(t, *inbox.NextCursor)
		// Fresh response: an omitted nullable cursor must not retain page one.
		var older drip.FeedbackInbox
		get("/users/me/node-version-feedback?cursor="+url.QueryEscape(*inbox.NextCursor), 200, &older)
		require.Len(t, older.Threads, 3)
		require.Nil(t, older.NextCursor)
		actual := []uuid.UUID{}
		for _, entry := range append(inbox.Threads, older.Threads...) {
			actual = append(actual, entry.Thread.Id)
			unread := 1
			if entry.Thread.Id == *binding.ThreadId {
				unread = 65
			}
			require.Equal(t, unread, entry.UnreadCount)
		}
		require.Equal(t, expected, actual, "every conversation must appear once in newest-first order")
	}
	checkInboxPages()
	// Equal timestamps must not drop conversations at the page boundary. The
	// documented secondary order is descending thread ID, for both pages.
	client.FeedbackThread.Update().SetLastMessageAt(when).ExecX(ctx)
	sort.Slice(expected, func(i, j int) bool { return expected[i].String() > expected[j].String() })
	checkInboxPages()
	get("/users/me/node-version-feedback?cursor=invalid", 422, nil)
}
