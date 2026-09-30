package integration

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"registry-backend/drip"
	"registry-backend/ent/schema"
	"registry-backend/server/middleware/authentication"
	registry "registry-backend/services/registry"
	"testing"
)

func TestFeedbackWritesStayBoundToDisplayedRecipient(t *testing.T) {
	ctx := context.Background()
	client, cleanup := setupDB(t, ctx)
	defer cleanup()
	client.User.Create().SetID("admin").SetIsAdmin(true).SaveX(ctx)
	for _, id := range []string{"original-publisher", "new-publisher"} {
		client.Publisher.Create().SetID(id).SetName(id).SaveX(ctx)
	}
	client.Node.Create().SetID("example").SetNormalizedID("example").SetName("Example").SetPublisherID("original-publisher").SetLicense("MIT").SetRepositoryURL("https://example.invalid").SaveX(ctx)
	v := client.NodeVersion.Create().SetNodeID("example").SetVersion("1.0.0").SetPipDependencies([]string{}).SetStatus(schema.NodeVersionStatusFlagged).SaveX(ctx)
	adminCtx := context.WithValue(ctx, authentication.UserContextKey, &authentication.UserDetails{ID: "admin"})
	service := registry.FeedbackService{Client: client}
	target := registry.FeedbackTarget{NodeID: "example", VersionID: v.ID, Admin: true}
	status := func(err error, code int) {
		t.Helper()
		var httpErr *echo.HTTPError
		require.ErrorAs(t, err, &httpErr)
		require.Equal(t, code, httpErr.Code)
	}
	initial, err := service.Get(adminCtx, target, nil)
	require.NoError(t, err)
	require.Equal(t, "original-publisher", initial.Target.PublisherId)
	require.Nil(t, initial.Target.ThreadId)
	_, err = service.Send(adminCtx, target, &drip.FeedbackMessageInput{Body: "No binding", ClientMessageId: uuid.New()})
	status(err, http.StatusUnprocessableEntity)
	input := &drip.FeedbackMessageInput{Target: initial.Target, Body: "Private message intended only for original-publisher", ClientMessageId: uuid.New()}
	first, err := service.Send(adminCtx, target, input)
	require.NoError(t, err)
	// A lost response is safe to retry against the same recipient, including null thread binding.
	retry, err := service.Send(adminCtx, target, input)
	require.NoError(t, err)
	require.Len(t, retry.Messages, 1)
	require.Equal(t, first.Thread.Id, *first.Target.ThreadId)
	_, err = service.Send(adminCtx, target, &drip.FeedbackMessageInput{Target: initial.Target, Body: "Another stale first message", ClientMessageId: uuid.New()})
	status(err, http.StatusConflict)
	tx, err := client.Tx(ctx)
	require.NoError(t, err)
	_, err = tx.Node.UpdateOneID("example").SetPublisherID("new-publisher").Save(ctx)
	require.NoError(t, err)
	require.NoError(t, tx.Commit())
	// Replaying a successful POST must not copy its body to the new Publisher.
	_, err = service.Send(adminCtx, target, input)
	status(err, http.StatusConflict)
	require.Equal(t, 1, client.FeedbackThread.Query().CountX(ctx))
	require.Equal(t, 1, client.FeedbackMessage.Query().CountX(ctx))
	fresh, err := service.Get(adminCtx, target, nil)
	require.NoError(t, err)
	second, err := service.Send(adminCtx, target, &drip.FeedbackMessageInput{Target: fresh.Target, Body: "New recipient message", ClientMessageId: uuid.New()})
	require.NoError(t, err)
	require.Equal(t, first.Thread.Revision, second.Thread.Revision, "revision alone cannot identify a conversation")
	for _, stale := range []drip.FeedbackWriteTarget{
		first.Target,
		{PublisherId: "new-publisher", ThreadId: first.Target.ThreadId},
		fresh.Target, // No thread is not a wildcard for an existing conversation.
	} {
		_, err = service.Send(adminCtx, target, &drip.FeedbackMessageInput{Target: stale, Body: "Stale draft", ClientMessageId: uuid.New()})
		status(err, http.StatusConflict)
		_, err = service.SetState(adminCtx, target, &drip.FeedbackStateInput{Target: stale, State: "resolved", ExpectedRevision: first.Thread.Revision})
		status(err, http.StatusConflict)
		_, err = service.MarkRead(adminCtx, target, &drip.FeedbackReadInput{Target: stale, LastReadMessageSeq: 1})
		status(err, http.StatusConflict)
	}
	current, err := service.Get(adminCtx, target, nil)
	require.NoError(t, err)
	require.Len(t, current.Messages, 1)
	require.Equal(t, "New recipient message", current.Messages[0].Body)
	require.Equal(t, drip.FeedbackState("awaiting_author"), current.Thread.State)
	require.Zero(t, client.FeedbackRead.Query().CountX(ctx))
	_, err = service.MarkRead(adminCtx, target, &drip.FeedbackReadInput{Target: current.Target, LastReadMessageSeq: 1})
	require.NoError(t, err)
	_, err = service.SetState(adminCtx, target, &drip.FeedbackStateInput{Target: current.Target, State: "resolved", ExpectedRevision: current.Thread.Revision})
	require.NoError(t, err)
}
