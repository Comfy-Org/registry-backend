package dripservices

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"registry-backend/db"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/feedbackevent"
	"registry-backend/ent/feedbackmessage"
	"registry-backend/ent/feedbackread"
	"registry-backend/ent/feedbackthread"
	"registry-backend/ent/node"
	"registry-backend/ent/nodeversion"
	"registry-backend/ent/predicate"
	"registry-backend/ent/publisherpermission"
	"registry-backend/ent/schema"
)

// FeedbackService never indexes, logs, or copies scan reports into conversations.
type FeedbackService struct{ Client *ent.Client }
type FeedbackTarget struct {
	NodeID      string
	VersionID   uuid.UUID
	PublisherID string
	Admin       bool
}
type feedbackResource struct {
	user        *ent.User
	version     *ent.NodeVersion
	thread      *ent.FeedbackThread
	publisherID string
	admin       bool
}

// Private endpoints conceal inaccessible resources while sharing the same
// active-user/admin policy as existing Registry operations.
func FeedbackUser(ctx context.Context, client *ent.Client, admin bool) (*ent.User, error) {
	check := RequireActiveUser
	if admin {
		check = RequireAdmin
	}
	user, err := check(ctx, client)
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) && httpErr.Code == http.StatusForbidden {
		return nil, feedbackNotFound()
	}
	return user, err
}
func feedbackNotFound() error { return echo.NewHTTPError(http.StatusNotFound, "Resource not found") }
func feedbackConflict() error {
	return echo.NewHTTPError(http.StatusConflict, "Conversation changed or is not open. Refresh and try again.")
}
func feedbackInvalid() error {
	return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid feedback input")
}

func feedbackAccess(ctx context.Context, c *ent.Client, target FeedbackTarget, lock bool) (*feedbackResource, error) {
	u, err := FeedbackUser(ctx, c, target.Admin)
	if err != nil {
		return nil, err
	}
	nq := c.Node.Query().Where(node.IDEQ(target.NodeID))
	if lock {
		nq.ForUpdate()
	}
	n, err := nq.Only(ctx)
	if ent.IsNotFound(err) {
		return nil, feedbackNotFound()
	}
	if err != nil {
		return nil, err
	}
	if !target.Admin {
		if target.PublisherID != n.PublisherID {
			return nil, feedbackNotFound()
		}
		err := assertPublisherPermissions(ctx, c, n.PublisherID, u.ID, []schema.PublisherPermissionType{schema.PublisherPermissionTypeOwner})
		if IsPermissionError(err) || ent.IsNotFound(err) {
			return nil, feedbackNotFound()
		}
		if err != nil {
			return nil, err
		}
	}
	vq := c.NodeVersion.Query().Where(nodeversion.IDEQ(target.VersionID), nodeversion.NodeIDEQ(n.ID))
	if lock {
		vq.ForUpdate()
	}
	v, err := vq.Only(ctx)
	if ent.IsNotFound(err) {
		return nil, feedbackNotFound()
	}
	if err != nil {
		return nil, err
	}
	thread, err := c.FeedbackThread.Query().Where(feedbackthread.VersionIDEQ(v.ID), feedbackthread.PublisherIDEQ(n.PublisherID)).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, err
	}
	return &feedbackResource{user: u, version: v, thread: thread, publisherID: n.PublisherID, admin: target.Admin}, nil
}

// Validate the displayed recipient under the same node/version locks as the write.
// A null thread is handled separately by Send to allow an idempotent first-message retry.
func feedbackCheckTarget(r *feedbackResource, target drip.FeedbackWriteTarget) error {
	if target.PublisherId == "" {
		return feedbackInvalid()
	}
	if target.PublisherId != r.publisherID || target.ThreadId != nil && (r.thread == nil || *target.ThreadId != r.thread.ID) {
		return feedbackConflict()
	}
	return nil
}
func feedbackWriteTarget(r *feedbackResource) drip.FeedbackWriteTarget {
	target := drip.FeedbackWriteTarget{PublisherId: r.publisherID}
	if r.thread != nil {
		id := r.thread.ID
		target.ThreadId = &id
	}
	return target
}
func feedbackThreadDTO(t *ent.FeedbackThread, v *ent.NodeVersion) drip.FeedbackThread {
	return drip.FeedbackThread{Id: t.ID, NodeId: v.NodeID, VersionId: v.ID, Version: v.Version, PublisherId: t.PublisherID, State: drip.FeedbackState(t.State), Revision: t.Revision, LastMessageSeq: t.LastMessageSeq, LastMessageAt: t.LastMessageAt, Archived: t.ArchivedAt != nil || v.Status == schema.NodeVersionStatusDeleted}
}
func feedbackPermissions(r *feedbackResource) drip.FeedbackPermissions {
	p := drip.FeedbackPermissions{CanStart: r.admin && r.thread == nil && r.version.Status == schema.NodeVersionStatusFlagged}
	if r.thread == nil || r.thread.ArchivedAt != nil || r.version.Status == schema.NodeVersionStatusDeleted {
		return p
	}
	p.CanReply = r.thread.State != feedbackthread.StateResolved
	p.CanResolve = r.admin && p.CanReply
	p.CanReopen = r.admin && !p.CanReply
	return p
}
func feedbackResponse(ctx context.Context, c *ent.Client, r *feedbackResource, before *int) (*drip.FeedbackResponse, error) {
	result := &drip.FeedbackResponse{Target: feedbackWriteTarget(r), Messages: []drip.FeedbackMessage{}, Events: []drip.FeedbackEvent{}, Permissions: feedbackPermissions(r)}
	if before != nil && *before < 1 {
		return nil, feedbackInvalid()
	}
	if r.thread == nil {
		return result, nil
	}
	dto := feedbackThreadDTO(r.thread, r.version)
	result.Thread = &dto
	read, err := c.FeedbackRead.Query().Where(feedbackread.ThreadIDEQ(r.thread.ID), feedbackread.UserIDEQ(r.user.ID)).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return nil, err
	}
	if read != nil {
		result.LastReadMessageSeq = read.LastReadMessageSeq
	}
	result.UnreadCount, err = c.FeedbackMessage.Query().Where(feedbackmessage.ThreadIDEQ(r.thread.ID), feedbackmessage.SenderUserIDNEQ(r.user.ID), feedbackmessage.SeqGT(result.LastReadMessageSeq)).Count(ctx)
	if err != nil {
		return nil, err
	}
	mq := c.FeedbackMessage.Query().Where(feedbackmessage.ThreadIDEQ(r.thread.ID)).Order(ent.Desc(feedbackmessage.FieldSeq)).Limit(31)
	if before != nil {
		mq.Where(feedbackmessage.SeqLT(*before))
	}
	messages, err := mq.All(ctx)
	if err != nil {
		return nil, err
	}
	if len(messages) > 30 {
		messages = messages[:30]
		seq := messages[29].Seq
		result.NextBeforeSeq = &seq
	}
	slices.Reverse(messages)
	for _, m := range messages {
		result.Messages = append(result.Messages, drip.FeedbackMessage{Id: m.ID, Seq: m.Seq, SenderUserId: m.SenderUserID, SenderName: m.SenderName, SenderRole: drip.FeedbackMessageSenderRole(m.SenderRole), Body: m.Body, CreatedAt: m.CreatedAt})
	}
	events, err := c.FeedbackEvent.Query().Where(feedbackevent.ThreadIDEQ(r.thread.ID)).Order(ent.Desc(feedbackevent.FieldCreatedAt)).Limit(100).All(ctx)
	if err != nil {
		return nil, err
	}
	slices.Reverse(events)
	for _, e := range events {
		result.Events = append(result.Events, drip.FeedbackEvent{Id: e.ID, EventType: drip.FeedbackEventEventType(e.EventType), CreatedAt: e.CreatedAt, ReplacementVersionId: e.ReplacementVersionID, ReplacementVersion: e.ReplacementVersion})
	}
	return result, nil
}
func (s *FeedbackService) Get(ctx context.Context, target FeedbackTarget, before *int) (*drip.FeedbackResponse, error) {
	r, err := feedbackAccess(ctx, s.Client, target, false)
	if err != nil {
		return nil, err
	}
	return feedbackResponse(ctx, s.Client, r, before)
}
func (s *FeedbackService) Send(ctx context.Context, target FeedbackTarget, input *drip.FeedbackMessageInput) (*drip.FeedbackResponse, error) {
	return db.WithTxResult(ctx, s.Client, func(tx *ent.Tx) (*drip.FeedbackResponse, error) {
		c := tx.Client()
		r, err := feedbackAccess(ctx, c, target, true)
		if err != nil {
			return nil, err
		}
		if input == nil || input.ClientMessageId == uuid.Nil {
			return nil, feedbackInvalid()
		}
		if err := feedbackCheckTarget(r, input.Target); err != nil {
			return nil, err
		}
		body := strings.TrimSpace(input.Body)
		length := utf8.RuneCountInString(body)
		if !utf8.ValidString(body) || length < 1 || length > 5000 {
			return nil, feedbackInvalid()
		}
		if r.thread != nil {
			existing, err := c.FeedbackMessage.Query().Where(feedbackmessage.ThreadIDEQ(r.thread.ID), feedbackmessage.SenderUserIDEQ(r.user.ID), feedbackmessage.ClientMessageIDEQ(input.ClientMessageId)).Only(ctx)
			if err != nil && !ent.IsNotFound(err) {
				return nil, err
			}
			if existing != nil {
				if existing.Body != body {
					return nil, feedbackConflict()
				}
				return feedbackResponse(ctx, c, r, nil)
			}
		}
		if r.thread != nil && input.Target.ThreadId == nil {
			return nil, feedbackConflict()
		}
		p := feedbackPermissions(r)
		if r.thread == nil {
			if !p.CanStart {
				return nil, feedbackConflict()
			}
			r.thread, err = c.FeedbackThread.Create().SetVersionID(r.version.ID).SetPublisherID(r.publisherID).SetCreatedByUserID(r.user.ID).Save(ctx)
			if err != nil {
				return nil, err
			}
		} else if !p.CanReply {
			return nil, feedbackConflict()
		}
		role := feedbackmessage.SenderRoleAuthor
		state := feedbackthread.StateAwaitingAdmin
		name := r.user.Name
		if target.Admin {
			role = feedbackmessage.SenderRoleAdmin
			state = feedbackthread.StateAwaitingAuthor
			name = "Registry team"
		}
		if name == "" {
			name = "Publisher owner"
		}
		seq := r.thread.LastMessageSeq + 1
		_, err = c.FeedbackMessage.Create().SetThreadID(r.thread.ID).SetSeq(seq).SetSenderUserID(r.user.ID).SetSenderName(name).SetSenderRole(role).SetBody(body).SetClientMessageID(input.ClientMessageId).Save(ctx)
		if err != nil {
			return nil, err
		}
		r.thread, err = r.thread.Update().SetLastMessageSeq(seq).SetLastMessageAt(time.Now()).SetState(state).AddRevision(1).Save(ctx)
		if err != nil {
			return nil, err
		}
		return feedbackResponse(ctx, c, r, nil)
	})
}
func (s *FeedbackService) SetState(ctx context.Context, target FeedbackTarget, input *drip.FeedbackStateInput) (*drip.FeedbackResponse, error) {
	target.Admin = true
	return db.WithTxResult(ctx, s.Client, func(tx *ent.Tx) (*drip.FeedbackResponse, error) {
		c := tx.Client()
		r, err := feedbackAccess(ctx, c, target, true)
		if err != nil {
			return nil, err
		}
		if input == nil || input.ExpectedRevision < 0 || input.State != "open" && input.State != "resolved" {
			return nil, feedbackInvalid()
		}
		if err := feedbackCheckTarget(r, input.Target); err != nil {
			return nil, err
		}
		if r.thread == nil || input.Target.ThreadId == nil || r.thread.Revision != input.ExpectedRevision {
			return nil, feedbackConflict()
		}
		p := feedbackPermissions(r)
		state, event := feedbackthread.StateResolved, feedbackevent.EventTypeResolved
		if input.State == "resolved" {
			if !p.CanResolve {
				return nil, feedbackConflict()
			}
		} else {
			if !p.CanReopen {
				return nil, feedbackConflict()
			}
			last, err := c.FeedbackMessage.Query().Where(feedbackmessage.ThreadIDEQ(r.thread.ID)).Order(ent.Desc(feedbackmessage.FieldSeq)).First(ctx)
			if err != nil {
				return nil, err
			}
			state, event = feedbackthread.StateAwaitingAdmin, feedbackevent.EventTypeReopened
			if last.SenderRole == feedbackmessage.SenderRoleAdmin {
				state = feedbackthread.StateAwaitingAuthor
			}
		}
		if err := feedbackTransition(ctx, c, r, state, event, nil); err != nil {
			return nil, err
		}
		return feedbackResponse(ctx, c, r, nil)
	})
}

// The state, revision and audit event are always written in the caller's transaction.
func feedbackTransition(ctx context.Context, c *ent.Client, r *feedbackResource, state feedbackthread.State, event feedbackevent.EventType, replacement *ent.NodeVersion) error {
	update := r.thread.Update().SetState(state).AddRevision(1)
	if state == feedbackthread.StateResolved {
		update.SetResolvedAt(time.Now()).SetResolvedByUserID(r.user.ID)
	} else {
		update.ClearResolvedAt().ClearResolvedByUserID()
	}
	thread, err := update.Save(ctx)
	if err != nil {
		return err
	}
	r.thread = thread
	audit := c.FeedbackEvent.Create().SetThreadID(thread.ID).SetActorUserID(r.user.ID).SetEventType(event)
	if replacement != nil {
		audit.SetReplacementVersionID(replacement.ID).SetReplacementVersion(replacement.Version)
	}
	return audit.Exec(ctx)
}

// SupersedeVersionFeedback closes follow-up on selected older versions. It does
// not alter moderation, installation policy tags, or the replacement's feedback.
func (s *RegistryService) SupersedeVersionFeedback(ctx context.Context, client *ent.Client, target FeedbackTarget, input *drip.FeedbackSupersedeInput) error {
	target.Admin = false
	versions, err := db.WithTxResult(ctx, client, func(tx *ent.Tx) ([]*ent.NodeVersion, error) {
		c := tx.Client()
		replacement, err := feedbackAccess(ctx, c, target, true)
		if err != nil {
			return nil, err
		}
		if input == nil || len(input.Versions) == 0 || len(input.Versions) > 100 {
			return nil, feedbackInvalid()
		}
		if replacement.version.Status == schema.NodeVersionStatusDeleted || replacement.version.Deprecated {
			return nil, feedbackConflict()
		}
		ids := make([]uuid.UUID, 0, len(input.Versions))
		selected := make(map[uuid.UUID]drip.FeedbackSupersedeVersion, len(input.Versions))
		for _, item := range input.Versions {
			if _, duplicate := selected[item.VersionId]; duplicate || item.VersionId == uuid.Nil || item.ExpectedRevision < 0 {
				return nil, feedbackInvalid()
			}
			selected[item.VersionId] = item
			ids = append(ids, item.VersionId)
		}
		// All feedback writes lock the node first, including replies and transfers.
		versions, err := c.NodeVersion.Query().Where(nodeversion.IDIn(ids...), nodeversion.NodeIDEQ(target.NodeID)).Order(ent.Asc(nodeversion.FieldID)).ForUpdate().All(ctx)
		if err != nil {
			return nil, err
		}
		if len(versions) != len(ids) {
			return nil, feedbackNotFound()
		}
		threads, err := c.FeedbackThread.Query().Where(feedbackthread.VersionIDIn(ids...), feedbackthread.PublisherIDEQ(replacement.publisherID)).All(ctx)
		if err != nil {
			return nil, err
		}
		byVersion := make(map[uuid.UUID]*ent.FeedbackThread, len(threads))
		for _, thread := range threads {
			byVersion[thread.VersionID] = thread
		}
		for i, version := range versions {
			thread := byVersion[version.ID]
			if thread == nil || thread.LastMessageSeq == 0 {
				return nil, feedbackNotFound()
			}
			if !version.CreateTime.Before(replacement.version.CreateTime) || version.Status == schema.NodeVersionStatusDeleted || thread.ArchivedAt != nil {
				return nil, feedbackConflict()
			}
			item := selected[version.ID]
			resource := &feedbackResource{user: replacement.user, version: version, thread: thread, publisherID: replacement.publisherID}
			if err := feedbackCheckTarget(resource, item.Target); err != nil {
				return nil, err
			}
			if item.Target.ThreadId == nil {
				return nil, feedbackConflict()
			}
			// Only the exact completed handoff can bypass the old revision on retry.
			if version.Deprecated && thread.State == feedbackthread.StateResolved && thread.Revision == item.ExpectedRevision+1 {
				last, err := c.FeedbackEvent.Query().Where(feedbackevent.ThreadIDEQ(thread.ID)).Order(ent.Desc(feedbackevent.FieldCreatedAt)).First(ctx)
				if err != nil && !ent.IsNotFound(err) {
					return nil, err
				}
				if last != nil && last.EventType == feedbackevent.EventTypeSuperseded && last.ReplacementVersionID != nil && *last.ReplacementVersionID == replacement.version.ID {
					continue
				}
			}
			if version.Deprecated && thread.State == feedbackthread.StateResolved {
				return nil, feedbackConflict()
			}
			if thread.Revision != item.ExpectedRevision {
				return nil, feedbackConflict()
			}
			versions[i], err = version.Update().SetDeprecated(true).Save(ctx)
			if err != nil {
				return nil, err
			}
			if err = feedbackTransition(ctx, c, resource, feedbackthread.StateResolved, feedbackevent.EventTypeSuperseded, replacement.version); err != nil {
				return nil, err
			}
		}
		return versions, nil
	})
	if err != nil {
		return err
	}
	// The database handoff is atomic. Retrying after an index failure also retries indexing.
	return s.algolia.IndexNodeVersions(ctx, versions...)
}

func (s *FeedbackService) MarkRead(ctx context.Context, target FeedbackTarget, input *drip.FeedbackReadInput) (*drip.FeedbackReadResult, error) {
	return db.WithTxResult(ctx, s.Client, func(tx *ent.Tx) (*drip.FeedbackReadResult, error) {
		c := tx.Client()
		r, err := feedbackAccess(ctx, c, target, true)
		if err != nil {
			return nil, err
		}
		if input == nil {
			return nil, feedbackInvalid()
		}
		if err := feedbackCheckTarget(r, input.Target); err != nil {
			return nil, err
		}
		if r.thread == nil || input.Target.ThreadId == nil {
			return nil, feedbackConflict()
		}
		if input.LastReadMessageSeq < 0 || input.LastReadMessageSeq > r.thread.LastMessageSeq {
			return nil, feedbackInvalid()
		}
		read, err := c.FeedbackRead.Query().Where(feedbackread.ThreadIDEQ(r.thread.ID), feedbackread.UserIDEQ(r.user.ID)).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return nil, err
		}
		seq := input.LastReadMessageSeq
		if read == nil {
			err = c.FeedbackRead.Create().SetThreadID(r.thread.ID).SetUserID(r.user.ID).SetLastReadMessageSeq(seq).Exec(ctx)
		} else {
			seq = max(seq, read.LastReadMessageSeq)
			err = read.Update().SetLastReadMessageSeq(seq).SetReadAt(time.Now()).Exec(ctx)
		}
		return &drip.FeedbackReadResult{LastReadMessageSeq: seq}, err
	})
}

type FeedbackInboxFilter struct {
	NodeID, PublisherID, Cursor *string
	State                       *drip.FeedbackState
	Admin                       bool
}
type feedbackCursor struct {
	Time time.Time `json:"time"`
	ID   uuid.UUID `json:"id"`
}

func (s *FeedbackService) Inbox(ctx context.Context, f FeedbackInboxFilter) (*drip.FeedbackInbox, error) {
	u, err := FeedbackUser(ctx, s.Client, f.Admin)
	if err != nil {
		return nil, err
	}
	q := s.Client.FeedbackThread.Query().Where(func(sel *sql.Selector) {
		v := sql.Table(nodeversion.Table)
		n := sql.Table(node.Table)
		sel.Join(v).On(sel.C(feedbackthread.FieldVersionID), v.C(nodeversion.FieldID)).Join(n).On(v.C(nodeversion.FieldNodeID), n.C(node.FieldID))
		sel.Where(sql.ColumnsEQ(sel.C(feedbackthread.FieldPublisherID), n.C(node.FieldPublisherID)))
		if f.NodeID != nil {
			sel.Where(sql.EQ(n.C(node.FieldID), *f.NodeID))
		}
		if !f.Admin {
			p := sql.Table(publisherpermission.Table)
			owned := sql.Select(p.C(publisherpermission.FieldPublisherID)).From(p).Where(sql.And(sql.EQ(p.C(publisherpermission.FieldUserID), u.ID), sql.EQ(p.C(publisherpermission.FieldPermission), string(schema.PublisherPermissionTypeOwner))))
			sel.Where(sql.In(n.C(node.FieldPublisherID), owned))
		}
	})
	if f.PublisherID != nil {
		q.Where(feedbackthread.PublisherIDEQ(*f.PublisherID))
	}
	if f.State != nil {
		state := feedbackthread.State(*f.State)
		if feedbackthread.StateValidator(state) != nil {
			return nil, feedbackInvalid()
		}
		q.Where(feedbackthread.StateEQ(state))
	}
	if f.Cursor != nil && *f.Cursor != "" {
		raw, err := base64.RawURLEncoding.DecodeString(*f.Cursor)
		var cursor feedbackCursor
		if err != nil || len(raw) > 256 || json.Unmarshal(raw, &cursor) != nil || cursor.ID == uuid.Nil || cursor.Time.IsZero() {
			return nil, feedbackInvalid()
		}
		q.Where(feedbackthread.Or(feedbackthread.LastMessageAtLT(cursor.Time), feedbackthread.And(feedbackthread.LastMessageAtEQ(cursor.Time), feedbackthread.IDLT(cursor.ID))))
	}
	threads, err := q.WithReads(func(q *ent.FeedbackReadQuery) { q.Where(feedbackread.UserIDEQ(u.ID)) }).Order(ent.Desc(feedbackthread.FieldLastMessageAt), ent.Desc(feedbackthread.FieldID)).Limit(21).All(ctx)
	if err != nil {
		return nil, err
	}
	result := &drip.FeedbackInbox{Threads: []drip.FeedbackSummary{}}
	if len(threads) > 20 {
		threads = threads[:20]
		last := threads[19]
		raw, _ := json.Marshal(feedbackCursor{last.LastMessageAt, last.ID})
		cursor := base64.RawURLEncoding.EncodeToString(raw)
		result.NextCursor = &cursor
	}
	if len(threads) == 0 {
		return result, nil
	}
	ids := make([]uuid.UUID, 0, len(threads))
	unreadPredicates := make([]predicate.FeedbackMessage, 0, len(threads))
	for _, t := range threads {
		ids = append(ids, t.VersionID)
		seq := 0
		if len(t.Edges.Reads) > 0 {
			seq = t.Edges.Reads[0].LastReadMessageSeq
		}
		unreadPredicates = append(unreadPredicates, feedbackmessage.And(feedbackmessage.ThreadIDEQ(t.ID), feedbackmessage.SeqGT(seq)))
	}
	versions, err := s.Client.NodeVersion.Query().Where(nodeversion.IDIn(ids...)).All(ctx)
	if err != nil {
		return nil, err
	}
	versionMap := map[uuid.UUID]*ent.NodeVersion{}
	for _, v := range versions {
		versionMap[v.ID] = v
	}
	var counts []struct {
		ThreadID uuid.UUID `json:"thread_id"`
		Count    int       `json:"count"`
	}
	err = s.Client.FeedbackMessage.Query().Where(feedbackmessage.SenderUserIDNEQ(u.ID), feedbackmessage.Or(unreadPredicates...)).GroupBy(feedbackmessage.FieldThreadID).Aggregate(ent.Count()).Scan(ctx, &counts)
	if err != nil {
		return nil, err
	}
	countMap := map[uuid.UUID]int{}
	for _, c := range counts {
		countMap[c.ThreadID] = c.Count
	}
	for _, t := range threads {
		if v := versionMap[t.VersionID]; v != nil {
			result.Threads = append(result.Threads, drip.FeedbackSummary{Thread: feedbackThreadDTO(t, v), UnreadCount: countMap[t.ID]})
		}
	}
	return result, nil
}
