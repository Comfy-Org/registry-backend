package integration

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"registry-backend/drip"
	"registry-backend/ent/feedbackmessage"
	"registry-backend/ent/feedbackthread"
	"registry-backend/ent/schema"
	"registry-backend/server/middleware/authentication"
	registry "registry-backend/services/registry"
	"testing"
)

func TestFeedbackConcurrentRetriesAndStateChanges(t *testing.T) {
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
	ownerCtx := context.WithValue(ctx, authentication.UserContextKey, &authentication.UserDetails{ID: "owner"})
	service := registry.FeedbackService{Client: client}
	target := registry.FeedbackTarget{NodeID: "example", VersionID: version.ID, Admin: true}
	ownerTarget := registry.FeedbackTarget{NodeID: "example", VersionID: version.ID, PublisherID: "publisher"}
	initial, err := service.Get(adminCtx, target, nil)
	require.NoError(t, err)
	retryID := uuid.New()
	results := make(chan error, 8)
	for i := 0; i < 8; i++ {
		go func() {
			_, err := service.Send(adminCtx, target, &drip.FeedbackMessageInput{Target: initial.Target, Body: "First message", ClientMessageId: retryID})
			results <- err
		}()
	}
	for i := 0; i < 8; i++ {
		require.NoError(t, <-results)
	}
	require.Equal(t, 1, client.FeedbackThread.Query().CountX(ctx))
	require.Equal(t, 1, client.FeedbackMessage.Query().CountX(ctx))
	current, err := service.Get(adminCtx, target, nil)
	require.NoError(t, err)
	for i := 0; i < 8; i++ {
		go func(i int) {
			_, err := service.Send(adminCtx, target, &drip.FeedbackMessageInput{Target: current.Target, Body: fmt.Sprintf("Follow-up %d", i), ClientMessageId: uuid.New()})
			results <- err
		}(i)
	}
	for i := 0; i < 8; i++ {
		require.NoError(t, <-results)
	}
	thread := client.FeedbackThread.Query().Where(feedbackthread.VersionIDEQ(version.ID)).OnlyX(ctx)
	require.Equal(t, 9, thread.LastMessageSeq)
	require.Equal(t, 9, client.FeedbackMessage.Query().Where(feedbackmessage.ThreadIDEQ(thread.ID)).CountX(ctx))
	// Exactly one competing operation succeeds; neither silently loses a reply.
	go func() {
		_, err := service.SetState(adminCtx, target, &drip.FeedbackStateInput{Target: current.Target, State: "resolved", ExpectedRevision: thread.Revision})
		results <- err
	}()
	go func() {
		_, err := service.Send(ownerCtx, ownerTarget, &drip.FeedbackMessageInput{Target: current.Target, Body: "Updated", ClientMessageId: uuid.New()})
		results <- err
	}()
	successes, conflicts := 0, 0
	for i := 0; i < 2; i++ {
		err := <-results
		if err == nil {
			successes++
		} else {
			httpErr, ok := err.(*echo.HTTPError)
			require.True(t, ok)
			require.Equal(t, 409, httpErr.Code)
			conflicts++
		}
	}
	require.Equal(t, 1, successes)
	require.Equal(t, 1, conflicts)
	// Read positions never regress. The caller chooses the consumed sequence.
	_, err = service.MarkRead(ownerCtx, ownerTarget, &drip.FeedbackReadInput{Target: current.Target, LastReadMessageSeq: 5})
	require.NoError(t, err)
	read, err := service.MarkRead(ownerCtx, ownerTarget, &drip.FeedbackReadInput{Target: current.Target, LastReadMessageSeq: 2})
	require.NoError(t, err)
	require.Equal(t, 5, read.LastReadMessageSeq)
	client.NodeVersion.UpdateOne(version).SetStatus(schema.NodeVersionStatusDeleted).ExecX(ctx)
	view, err := service.Get(ownerCtx, ownerTarget, nil)
	require.NoError(t, err)
	require.False(t, view.Permissions.CanReply)
	_, err = service.Send(adminCtx, target, &drip.FeedbackMessageInput{Target: current.Target, Body: "No write after deletion", ClientMessageId: uuid.New()})
	require.Error(t, err)
}
