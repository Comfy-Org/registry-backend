package implementation

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/labstack/echo/v4"
	"net/http"
	"registry-backend/drip"
	"registry-backend/ent/feedbackmessage"
	"registry-backend/ent/feedbackthread"
	"registry-backend/ent/node"
	"registry-backend/ent/nodeversion"
	"registry-backend/ent/predicate"
	"registry-backend/entity"
	"registry-backend/mapper"
	registry "registry-backend/services/registry"
)

func (s *DripStrictServerImplementation) feedbackService() *registry.FeedbackService {
	return &registry.FeedbackService{Client: s.Client}
}
func (s *DripStrictServerImplementation) AdminGetVersionFeedback(ctx context.Context, r drip.AdminGetVersionFeedbackRequestObject) (drip.AdminGetVersionFeedbackResponseObject, error) {
	result, err := s.feedbackService().Get(ctx, registry.FeedbackTarget{NodeID: r.NodeId, VersionID: r.VersionId, Admin: true}, r.Params.BeforeSeq)
	if err != nil {
		return nil, err
	}
	return drip.AdminGetVersionFeedback200JSONResponse(*result), nil
}
func (s *DripStrictServerImplementation) AdminSendVersionFeedback(ctx context.Context, r drip.AdminSendVersionFeedbackRequestObject) (drip.AdminSendVersionFeedbackResponseObject, error) {
	result, err := s.feedbackService().Send(ctx, registry.FeedbackTarget{NodeID: r.NodeId, VersionID: r.VersionId, Admin: true}, r.Body)
	if err != nil {
		return nil, err
	}
	return drip.AdminSendVersionFeedback200JSONResponse(*result), nil
}
func (s *DripStrictServerImplementation) AdminReadVersionFeedback(ctx context.Context, r drip.AdminReadVersionFeedbackRequestObject) (drip.AdminReadVersionFeedbackResponseObject, error) {
	result, err := s.feedbackService().MarkRead(ctx, registry.FeedbackTarget{NodeID: r.NodeId, VersionID: r.VersionId, Admin: true}, r.Body)
	if err != nil {
		return nil, err
	}
	return drip.AdminReadVersionFeedback200JSONResponse(*result), nil
}
func (s *DripStrictServerImplementation) AdminSetVersionFeedbackState(ctx context.Context, r drip.AdminSetVersionFeedbackStateRequestObject) (drip.AdminSetVersionFeedbackStateResponseObject, error) {
	result, err := s.feedbackService().SetState(ctx, registry.FeedbackTarget{NodeID: r.NodeId, VersionID: r.VersionId, Admin: true}, r.Body)
	if err != nil {
		return nil, err
	}
	return drip.AdminSetVersionFeedbackState200JSONResponse(*result), nil
}
func (s *DripStrictServerImplementation) AdminListVersionFeedback(ctx context.Context, r drip.AdminListVersionFeedbackRequestObject) (drip.AdminListVersionFeedbackResponseObject, error) {
	result, err := s.feedbackService().Inbox(ctx, registry.FeedbackInboxFilter{NodeID: r.Params.NodeId, PublisherID: r.Params.PublisherId, State: r.Params.State, Cursor: r.Params.Cursor, Admin: true})
	if err != nil {
		return nil, err
	}
	return drip.AdminListVersionFeedback200JSONResponse(*result), nil
}
func (s *DripStrictServerImplementation) AuthorGetVersionFeedback(ctx context.Context, r drip.AuthorGetVersionFeedbackRequestObject) (drip.AuthorGetVersionFeedbackResponseObject, error) {
	result, err := s.feedbackService().Get(ctx, registry.FeedbackTarget{NodeID: r.NodeId, VersionID: r.VersionId, PublisherID: r.PublisherId}, r.Params.BeforeSeq)
	if err != nil {
		return nil, err
	}
	return drip.AuthorGetVersionFeedback200JSONResponse(*result), nil
}
func (s *DripStrictServerImplementation) AuthorSendVersionFeedback(ctx context.Context, r drip.AuthorSendVersionFeedbackRequestObject) (drip.AuthorSendVersionFeedbackResponseObject, error) {
	result, err := s.feedbackService().Send(ctx, registry.FeedbackTarget{NodeID: r.NodeId, VersionID: r.VersionId, PublisherID: r.PublisherId}, r.Body)
	if err != nil {
		return nil, err
	}
	return drip.AuthorSendVersionFeedback200JSONResponse(*result), nil
}
func (s *DripStrictServerImplementation) AuthorReadVersionFeedback(ctx context.Context, r drip.AuthorReadVersionFeedbackRequestObject) (drip.AuthorReadVersionFeedbackResponseObject, error) {
	result, err := s.feedbackService().MarkRead(ctx, registry.FeedbackTarget{NodeID: r.NodeId, VersionID: r.VersionId, PublisherID: r.PublisherId}, r.Body)
	if err != nil {
		return nil, err
	}
	return drip.AuthorReadVersionFeedback200JSONResponse(*result), nil
}
func (s *DripStrictServerImplementation) AuthorListVersionFeedback(ctx context.Context, r drip.AuthorListVersionFeedbackRequestObject) (drip.AuthorListVersionFeedbackResponseObject, error) {
	result, err := s.feedbackService().Inbox(ctx, registry.FeedbackInboxFilter{NodeID: r.Params.NodeId, PublisherID: r.Params.PublisherId, State: r.Params.State, Cursor: r.Params.Cursor})
	if err != nil {
		return nil, err
	}
	return drip.AuthorListVersionFeedback200JSONResponse(*result), nil
}
func (s *DripStrictServerImplementation) AdminListNodeVersions(ctx context.Context, r drip.AdminListNodeVersionsRequestObject) (drip.AdminListNodeVersionsResponseObject, error) {
	if _, err := registry.FeedbackUser(ctx, s.Client, true); err != nil {
		return nil, err
	}
	page, size := 1, 10
	if r.Params.Page != nil {
		page = *r.Params.Page
	}
	if r.Params.PageSize != nil {
		size = *r.Params.PageSize
	}
	if page < 1 || page > 100000 || size < 1 || size > 100 {
		return nil, echo.NewHTTPError(422, "Invalid pagination")
	}
	filter := &entity.NodeVersionFilter{
		NodeId: r.Params.NodeId,
		Status: mapper.ApiNodeVersionStatusesToDbNodeVersionStatuses(r.Params.Statuses),
		Page:   page, PageSize: size, IncludeStatusReason: true, NewestFirst: true,
	}
	if r.Params.Version != nil {
		filter.Predicates = append(filter.Predicates, nodeversion.VersionEQ(*r.Params.Version))
	}
	if r.Params.StatusReason != nil {
		filter.Predicates = append(filter.Predicates, nodeversion.StatusReasonContains(*r.Params.StatusReason))
	}
	if r.Params.FeedbackStatus != nil {
		switch *r.Params.FeedbackStatus {
		case drip.Unprocessed, drip.Processed, drip.NeedsResponse:
			filter.Predicates = append(filter.Predicates, versionFeedbackStatusPredicate(*r.Params.FeedbackStatus))
		default:
			return nil, echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid feedback status")
		}
	}
	versions, err := s.RegistryService.ListNodeVersions(ctx, s.Client, filter)
	if err != nil {
		return nil, err
	}
	result := drip.AdminListNodeVersions200JSONResponse{Versions: []drip.AdminNodeVersion{}, Total: versions.Total, Page: versions.Page, PageSize: versions.Limit, TotalPages: versions.TotalPages}
	for _, v := range versions.NodeVersions {
		result.Versions = append(result.Versions, *mapper.DbNodeVersionToAdminNodeVersion(v))
	}
	return result, nil
}

// Feedback handling status depends on the last message, not on resolution events.
// EXISTS keeps both the count and paginated rows free from duplicate versions.
func versionFeedbackStatusPredicate(status drip.VersionFeedbackStatus) predicate.NodeVersion {
	return func(versions *sql.Selector) {
		thread := sql.Table(feedbackthread.Table).As("feedback_thread")
		message := sql.Table(feedbackmessage.Table).As("feedback_last_message")
		currentNode := sql.Table(node.Table).As("feedback_node")
		latest := sql.Select(thread.C(feedbackthread.FieldID)).From(thread).
			Join(currentNode).On(thread.C(feedbackthread.FieldPublisherID), currentNode.C(node.FieldPublisherID)).
			Join(message).On(thread.C(feedbackthread.FieldID), message.C(feedbackmessage.FieldThreadID)).
			Where(sql.And(
				sql.ColumnsEQ(currentNode.C(node.FieldID), versions.C(nodeversion.FieldNodeID)),
				sql.ColumnsEQ(thread.C(feedbackthread.FieldVersionID), versions.C(nodeversion.FieldID)),
				sql.ColumnsEQ(message.C(feedbackmessage.FieldSeq), thread.C(feedbackthread.FieldLastMessageSeq)),
			))
		if status == drip.Unprocessed {
			versions.Where(sql.Not(sql.Exists(latest)))
			return
		}
		role := feedbackmessage.SenderRoleAdmin
		if status == drip.NeedsResponse {
			role = feedbackmessage.SenderRoleAuthor
		}
		latest.Where(sql.EQ(message.C(feedbackmessage.FieldSenderRole), role))
		versions.Where(sql.Exists(latest))
	}
}

func (s *DripStrictServerImplementation) AuthorSupersedeVersionFeedback(ctx context.Context, r drip.AuthorSupersedeVersionFeedbackRequestObject) (drip.AuthorSupersedeVersionFeedbackResponseObject, error) {
	err := s.RegistryService.SupersedeVersionFeedback(ctx, s.Client, registry.FeedbackTarget{NodeID: r.NodeId, VersionID: r.VersionId, PublisherID: r.PublisherId}, r.Body)
	if err != nil {
		return nil, err
	}
	return drip.AuthorSupersedeVersionFeedback204Response{}, nil
}
