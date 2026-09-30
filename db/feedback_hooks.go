package db

import (
	"context"
	"registry-backend/ent"
	"registry-backend/ent/feedbackevent"
	"registry-backend/ent/feedbackthread"
	"registry-backend/ent/node"
	"registry-backend/ent/nodeversion"
	"time"
)

// RegisterFeedbackHooks prevents a transfer back from reopening historical
// conversations. Transfer services should wrap the node mutation in a transaction.
// Archiving happens before the ownership change, so failures remain fail-closed.
func RegisterFeedbackHooks(client *ent.Client) {
	client.Node.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, mutation ent.Mutation) (ent.Value, error) {
			m, ok := mutation.(*ent.NodeMutation)
			if !ok || !m.Op().Is(ent.OpUpdate|ent.OpUpdateOne) {
				return next.Mutate(ctx, mutation)
			}
			publisher, changed := m.PublisherID()
			if !changed {
				return next.Mutate(ctx, mutation)
			}
			c := m.Client()
			ids, err := m.IDs(ctx)
			if err != nil {
				return nil, err
			}
			// Share the node lock used by feedback writes. In a transfer transaction,
			// a concurrent first message cannot appear after the archive query.
			moved, err := c.Node.Query().Where(node.IDIn(ids...), node.PublisherIDNEQ(publisher)).ForUpdate().IDs(ctx)
			if err != nil {
				return nil, err
			}
			if len(moved) > 0 {
				versions, err := c.NodeVersion.Query().Where(nodeversion.NodeIDIn(moved...)).IDs(ctx)
				if err != nil {
					return nil, err
				}
				threads, err := c.FeedbackThread.Query().Where(feedbackthread.VersionIDIn(versions...), feedbackthread.ArchivedAtIsNil()).All(ctx)
				if err != nil {
					return nil, err
				}
				for _, thread := range threads {
					if err = thread.Update().SetArchivedAt(time.Now()).AddRevision(1).Exec(ctx); err != nil {
						return nil, err
					}
					if err = c.FeedbackEvent.Create().SetThreadID(thread.ID).SetActorUserID("system").SetEventType(feedbackevent.EventTypeArchived).Exec(ctx); err != nil {
						return nil, err
					}
				}
			}
			return next.Mutate(ctx, mutation)
		})
	})
}
