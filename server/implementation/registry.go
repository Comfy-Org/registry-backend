package implementation

import (
	"context"
	"errors"
	"fmt"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/publisher"
	"registry-backend/ent/schema"
	"registry-backend/entity"
	drip_logging "registry-backend/logging"
	"registry-backend/mapper"
	drip_services "registry-backend/services/registry"
	"registry-backend/tracing"
	"time"

	"github.com/google/uuid"
	"github.com/mixpanel/mixpanel-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

func (impl *DripStrictServerImplementation) ListPublishersForUser(
	ctx context.Context, request drip.ListPublishersForUserRequestObject) (drip.ListPublishersForUserResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.ListPublishersForUser")()

	// Extract user ID from context
	userId, err := mapper.GetUserIDFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to get user ID from context w/ err: %v", err)
		return drip.ListPublishersForUser400JSONResponse{Message: "Invalid user ID"}, err
	}

	// Call the service to list publishers
	log.Ctx(ctx).Info().Msgf("Fetching publishers for user %s", userId)
	publishers, err := impl.RegistryService.ListPublishers(ctx, impl.Client, &entity.PublisherFilter{
		UserID: userId,
	})
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to list publishers w/ err: %v", err)
		return drip.ListPublishersForUser500JSONResponse{
			Message: "Failed to fetch list of publishers", Error: err.Error()}, err
	}

	// Map the publishers to API format
	apiPublishers := make([]drip.Publisher, 0, len(publishers))
	log.Ctx(ctx).Info().Msgf(
		"Successfully fetched publishers for user %s, count %d", userId, len(apiPublishers))
	for _, dbPublisher := range publishers {
		apiPublishers = append(apiPublishers, *mapper.DbPublisherToApiPublisher(dbPublisher, true))
	}

	return drip.ListPublishersForUser200JSONResponse(apiPublishers), nil
}

func (s *DripStrictServerImplementation) ValidatePublisher(
	ctx context.Context, request drip.ValidatePublisherRequestObject) (drip.ValidatePublisherResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.ValidatePublisher")()

	// Check if the username is empty
	name := request.Params.Username
	if name == "" {
		log.Ctx(ctx).Warn().Msg("Username parameter is missing")
		return drip.ValidatePublisher400JSONResponse{Message: "Username parameter is required"}, nil
	}

	isValid := mapper.IsValidPublisherID(name)
	if !isValid {
		return drip.ValidatePublisher400JSONResponse{
			Message: "Must start with a lowercase letter and can only contain lowercase letters, digits, and hyphens.",
		}, nil
	}

	// Note: username = id field in publisher table, display = name field in publisher table
	count, err := s.Client.Publisher.Query().Where(publisher.ID(name)).Count(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to query username %s w/ err: %v", name, err)
		return drip.ValidatePublisher500JSONResponse{Message: "Failed to query username", Error: err.Error()}, err
	}

	// Log the result of the count query
	log.Ctx(ctx).Info().Msgf("Count for username %s: %d", name, count)
	if count > 0 {
		return drip.ValidatePublisher400JSONResponse{
			Message: "Publisher ID already exists.",
		}, nil
	}

	return drip.ValidatePublisher200JSONResponse{
		IsAvailable: proto.Bool(true),
	}, nil
}

func (s *DripStrictServerImplementation) CreatePublisher(
	ctx context.Context, request drip.CreatePublisherRequestObject) (drip.CreatePublisherResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.CreatePublisher")()

	// Extract user ID from context
	userId, err := mapper.GetUserIDFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to get user ID from context w/ err: %v", err)
		return drip.CreatePublisher400JSONResponse{Message: "Invalid user ID"}, err
	}

	log.Ctx(ctx).Info().Msgf("Checking if user ID %s has reached the maximum number of publishers", userId)
	userPublishers, err := s.RegistryService.ListPublishers(
		ctx, s.Client, &entity.PublisherFilter{UserID: userId})
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to list publishers for user ID %s w/ err: %v", userId, err)
		return drip.CreatePublisher500JSONResponse{Message: "Failed to list publishers", Error: err.Error()}, err
	}
	if len(userPublishers) >= 5 {
		log.Ctx(ctx).Info().Msgf("User ID %s has reached the maximum number of publishers", userId)
		return drip.CreatePublisher403JSONResponse{
			Message: "User has reached the maximum number of publishers.",
		}, nil
	}

	// Create a new publisher
	log.Ctx(ctx).Info().Msgf("Creating publisher for user ID %s", userId)
	publisher, err := s.RegistryService.CreatePublisher(ctx, s.Client, userId, request.Body)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to create publisher for user ID %s w/ err: %v", userId, err)
		if ent.IsConstraintError(err) {
			return drip.CreatePublisher400JSONResponse{Message: "Constraint error", Error: err.Error()}, nil
		}

		return drip.CreatePublisher500JSONResponse{Message: "Internal server error", Error: err.Error()}, err
	}

	// Log the successful creation
	log.Ctx(ctx).Info().Msgf("Publisher created successfully for user ID: %s", userId)
	return drip.CreatePublisher201JSONResponse(*mapper.DbPublisherToApiPublisher(publisher, true)), nil
}

func (s *DripStrictServerImplementation) ListPublishers(
	ctx context.Context, request drip.ListPublishersRequestObject) (drip.ListPublishersResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.ListPublishers")()

	pubs, err := s.RegistryService.ListPublishers(ctx, s.Client, nil)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to retrieve list of publishers w/ err: %v", err)
		return drip.ListPublishers500JSONResponse{Message: "Failed to get publisher", Error: err.Error()}, err
	}

	res := drip.ListPublishers200JSONResponse{}
	for _, pub := range pubs {
		res = append(res, *mapper.DbPublisherToApiPublisher(pub, false))
	}

	log.Ctx(ctx).Info().Msgf("List of Publishers retrieved successfully")
	return res, nil
}

func (s *DripStrictServerImplementation) DeletePublisher(
	ctx context.Context, request drip.DeletePublisherRequestObject) (drip.DeletePublisherResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.DeletePublisher")()

	err := s.RegistryService.DeletePublisher(ctx, s.Client, request.PublisherId)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to delete publisher with ID %s w/ err: %v", request.PublisherId, err)
		return drip.DeletePublisher500JSONResponse{}, nil
	}

	log.Ctx(ctx).Info().Msgf("Publisher with ID %s deleted successfully", request.PublisherId)
	return drip.DeletePublisher204Response{}, nil
}

func (s *DripStrictServerImplementation) GetPublisher(
	ctx context.Context, request drip.GetPublisherRequestObject) (drip.GetPublisherResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.GetPublisher")()

	publisherId := request.PublisherId
	publisher, err := s.RegistryService.GetPublisher(ctx, s.Client, request.PublisherId)
	if ent.IsNotFound(err) {
		log.Ctx(ctx).Info().Msgf("Publisher with ID %s not found", publisherId)
		return drip.GetPublisher404JSONResponse{Message: "Publisher not found"}, nil
	}
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to retrieve publisher with ID %s w/ err: %v", publisherId, err)
		return drip.GetPublisher500JSONResponse{Message: "Failed to get publisher", Error: err.Error()}, err
	}

	log.Ctx(ctx).Info().Msgf("Publisher with ID %s retrieved successfully", publisherId)
	return drip.GetPublisher200JSONResponse(*mapper.DbPublisherToApiPublisher(publisher, false)), nil
}

func (s *DripStrictServerImplementation) UpdatePublisher(
	ctx context.Context, request drip.UpdatePublisherRequestObject) (drip.UpdatePublisherResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.UpdatePublisher")()

	updateOne := mapper.ApiUpdatePublisherToUpdateFields(request.PublisherId, request.Body, s.Client)
	updatedPublisher, err := s.RegistryService.UpdatePublisher(ctx, s.Client, updateOne)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to update publisher with ID %s w/ err: %v", request.PublisherId, err)
		return drip.UpdatePublisher500JSONResponse{Message: "Failed to update publisher", Error: err.Error()}, err
	}

	log.Ctx(ctx).Info().Msgf("Publisher with ID %s updated successfully", request.PublisherId)
	return drip.UpdatePublisher200JSONResponse(*mapper.DbPublisherToApiPublisher(updatedPublisher, true)), nil
}

func (s *DripStrictServerImplementation) CreateNode(
	ctx context.Context, request drip.CreateNodeRequestObject) (drip.CreateNodeResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.CreateNode")()

	node, err := s.RegistryService.CreateNode(ctx, s.Client, request.PublisherId, request.Body)
	if mapper.IsErrorBadRequest(err) || ent.IsConstraintError(err) {
		log.Ctx(ctx).Error().Msgf(
			"Failed to create node for publisher ID %s w/ err: %v", request.PublisherId, err)
		return drip.CreateNode400JSONResponse{Message: "The node already exists", Error: err.Error()}, err
	}
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to create node for publisher ID %s w/ err: %v", request.PublisherId, err)
		return drip.CreateNode500JSONResponse{Message: "Failed to create node", Error: err.Error()}, err
	}

	log.Ctx(ctx).Info().Msgf("Node created successfully for publisher ID: %s", request.PublisherId)
	return drip.CreateNode201JSONResponse(*mapper.DbNodeToApiNode(node)), nil
}

func (s *DripStrictServerImplementation) ListNodesForPublisher(
	ctx context.Context, request drip.ListNodesForPublisherRequestObject) (drip.ListNodesForPublisherResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.ListNodesForPublisher")()

	nodeResults, err := s.RegistryService.ListNodes(
		ctx, s.Client /*page=*/, 1 /*limit=*/, 10, &entity.NodeFilter{
			PublisherID: request.PublisherId,
		})
	if err != nil {
		log.Ctx(ctx).Error().Msgf(
			"Failed to list nodes for publisher ID %s w/ err: %v", request.PublisherId, err)
		return drip.ListNodesForPublisher500JSONResponse{Message: "Failed to list nodes", Error: err.Error()}, err
	}

	if len(nodeResults.Nodes) == 0 {
		log.Ctx(ctx).Info().Msgf("No nodes found for publisher ID: %s", request.PublisherId)
		return drip.ListNodesForPublisher200JSONResponse([]drip.Node{}), nil
	}

	apiNodes := make([]drip.Node, 0, len(nodeResults.Nodes))
	for _, dbNode := range nodeResults.Nodes {
		apiNodes = append(apiNodes, *mapper.DbNodeToApiNode(dbNode))
	}

	log.Ctx(ctx).Info().Msgf(
		"Found %d nodes for publisher ID: %s", len(apiNodes), request.PublisherId)
	return drip.ListNodesForPublisher200JSONResponse(apiNodes), nil
}

func (s *DripStrictServerImplementation) ListAllNodes(
	ctx context.Context, request drip.ListAllNodesRequestObject) (drip.ListAllNodesResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.ListAllNodes")()

	err := s.MixpanelService.Track(ctx, []*mixpanel.Event{
		s.MixpanelService.NewEvent("List All Nodes", "", map[string]any{
			"page":  request.Params.Page,
			"limit": request.Params.Limit,
		}),
	})
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to track event w/ err: %v", err)
	}

	// Set default values for pagination parameters
	page := 1
	if request.Params.Page != nil {
		page = *request.Params.Page
	}
	limit := 10
	if request.Params.Limit != nil {
		limit = *request.Params.Limit
	}

	// Initialize the node filter
	filter := &entity.NodeFilter{
		Timestamp:     request.Params.Timestamp,
		IncludeBanned: request.Params.IncludeBanned,
	}

	// List nodes from the registry service
	latest := false
	if request.Params.Latest != nil {
		latest = *request.Params.Latest
	}
	nodeResults, err := s.RegistryService.ListNodesWithCache(ctx, s.Client, page, limit, filter, latest)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to list nodes w/ err: %v", err)
		return drip.ListAllNodes500JSONResponse{Message: "Failed to list nodes", Error: err.Error()}, err
	}

	// Handle case when no nodes are found
	if len(nodeResults.Nodes) == 0 {
		log.Ctx(ctx).Info().Msg("No nodes found")
		return drip.ListAllNodes200JSONResponse{
			Nodes:      &[]drip.Node{},
			Total:      &nodeResults.Total,
			Page:       &nodeResults.Page,
			Limit:      &nodeResults.Limit,
			TotalPages: &nodeResults.TotalPages,
		}, nil
	}

	// Convert database nodes to API nodes
	apiNodes := make([]drip.Node, 0, len(nodeResults.Nodes))
	for _, dbNode := range nodeResults.Nodes {
		apiNode := mapper.DbNodeToApiNode(dbNode)

		// attach information of latest version if available
		if len(dbNode.Edges.Versions) > 0 {
			apiNode.LatestVersion = mapper.DbNodeVersionToApiNodeVersion(dbNode.Edges.Versions[0])
			apiNode.LatestVersion.StatusReason = nil
		}

		// Map publisher information
		apiNode.Publisher = mapper.DbPublisherToApiPublisher(dbNode.Edges.Publisher, false)
		apiNodes = append(apiNodes, *apiNode)
	}

	log.Ctx(ctx).Info().Msgf("Found %d nodes", len(apiNodes))
	return drip.ListAllNodes200JSONResponse{
		Nodes:      &apiNodes,
		Total:      &nodeResults.Total,
		Page:       &nodeResults.Page,
		Limit:      &nodeResults.Limit,
		TotalPages: &nodeResults.TotalPages,
	}, nil
}

// SearchNodes implements drip.StrictServerInterface.
func (s *DripStrictServerImplementation) SearchNodes(ctx context.Context, request drip.SearchNodesRequestObject) (drip.SearchNodesResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.SearchNodes")()

	// Set default values for pagination parameters
	page := 1
	if request.Params.Page != nil {
		page = *request.Params.Page
	}
	limit := 10
	if request.Params.Limit != nil {
		limit = *request.Params.Limit
	}

	f := &entity.NodeFilter{
		IncludeBanned: request.Params.IncludeBanned,
	}
	if request.Params.Search != nil {
		f.Search = *request.Params.Search
	}
	// List nodes from the registry service
	nodeResults, err := s.RegistryService.ListNodes(ctx, s.Client, page, limit, f)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to search nodes w/ err: %v", err)
		return drip.SearchNodes500JSONResponse{Message: "Failed to search nodes", Error: err.Error()}, err
	}

	if len(nodeResults.Nodes) == 0 {
		log.Ctx(ctx).Info().Msg("No nodes found")
		return drip.SearchNodes200JSONResponse{
			Nodes:      &[]drip.Node{},
			Total:      &nodeResults.Total,
			Page:       &nodeResults.Page,
			Limit:      &nodeResults.Limit,
			TotalPages: &nodeResults.TotalPages,
		}, nil
	}

	apiNodes := make([]drip.Node, 0, len(nodeResults.Nodes))
	for _, dbNode := range nodeResults.Nodes {
		apiNode := mapper.DbNodeToApiNode(dbNode)
		if len(dbNode.Edges.Versions) > 0 {
			latestVersion := dbNode.Edges.Versions[0]
			apiNode.LatestVersion = mapper.DbNodeVersionToApiNodeVersion(latestVersion)
		}
		apiNode.Publisher = mapper.DbPublisherToApiPublisher(dbNode.Edges.Publisher, false)
		apiNodes = append(apiNodes, *apiNode)
	}

	log.Ctx(ctx).Info().Msgf("Found %d nodes", len(apiNodes))
	return drip.SearchNodes200JSONResponse{
		Nodes:      &apiNodes,
		Total:      &nodeResults.Total,
		Page:       &nodeResults.Page,
		Limit:      &nodeResults.Limit,
		TotalPages: &nodeResults.TotalPages,
	}, nil
}

func (s *DripStrictServerImplementation) DeleteNode(
	ctx context.Context, request drip.DeleteNodeRequestObject) (drip.DeleteNodeResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.DeleteNode")()

	err := s.RegistryService.DeleteNode(ctx, s.Client, request.NodeId)
	if err != nil && !ent.IsNotFound(err) {
		log.Ctx(ctx).Error().Msgf("Failed to delete node %s w/ err: %v", request.NodeId, err)
		return drip.DeleteNode500JSONResponse{Message: "Internal server error"}, err
	}

	log.Ctx(ctx).Info().Msgf("Node %s deleted successfully", request.NodeId)
	return drip.DeleteNode204Response{}, nil
}

func (s *DripStrictServerImplementation) GetNode(
	ctx context.Context, request drip.GetNodeRequestObject) (drip.GetNodeResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.GetNode")()

	node, err := s.RegistryService.GetNode(ctx, s.Client, request.NodeId)
	if ent.IsNotFound(err) {
		log.Ctx(ctx).Error().Msgf("Failed to get node %s w/ err: %v", request.NodeId, err)
		return drip.GetNode404JSONResponse{Message: "Node not found"}, nil
	}
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to get node %s w/ err: %v", request.NodeId, err)
		return drip.GetNode500JSONResponse{Message: "Failed to get node"}, err
	}

	nodeVersion, err := s.RegistryService.GetLatestNodeVersion(ctx, s.Client, request.NodeId)
	if err != nil {
		log.Ctx(ctx).Error().Msgf(
			"Failed to get latest node version for node %s w/ err: %v", request.NodeId, err)
		return drip.GetNode500JSONResponse{Message: "Failed to get latest node version", Error: err.Error()}, err
	}

	apiNode := mapper.DbNodeToApiNode(node)
	apiNode.LatestVersion = mapper.DbNodeVersionToApiNodeVersion(nodeVersion)

	log.Ctx(ctx).Info().Msgf("Node %s retrieved successfully", request.NodeId)
	return drip.GetNode200JSONResponse(*apiNode), nil
}

func (s *DripStrictServerImplementation) UpdateNode(
	ctx context.Context, request drip.UpdateNodeRequestObject) (drip.UpdateNodeResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.UpdateNode")()

	updateOneFunc := func(client *ent.Client) *ent.NodeUpdateOne {
		return mapper.ApiUpdateNodeToUpdateFields(request.NodeId, request.Body, client)
	}
	updatedNode, err := s.RegistryService.UpdateNode(ctx, s.Client, updateOneFunc)
	if ent.IsNotFound(err) {
		log.Ctx(ctx).Error().Msgf("Node %s not found w/ err: %v", request.NodeId, err)
		return drip.UpdateNode404JSONResponse{Message: "Not Found"}, nil
	}
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to update node %s w/ err: %v", request.NodeId, err)
		return drip.UpdateNode500JSONResponse{Message: "Failed to update node", Error: err.Error()}, err
	}

	log.Ctx(ctx).Info().Msgf("Node %s updated successfully", request.NodeId)
	return drip.UpdateNode200JSONResponse(*mapper.DbNodeToApiNode(updatedNode)), nil
}

func (s *DripStrictServerImplementation) ListNodeVersions(
	ctx context.Context, request drip.ListNodeVersionsRequestObject) (drip.ListNodeVersionsResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.ListNodeVersions")()

	apiStatus := mapper.ApiNodeVersionStatusesToDbNodeVersionStatuses(request.Params.Statuses)

	nodeVersionsResult, err := s.RegistryService.ListNodeVersions(ctx, s.Client, &entity.NodeVersionFilter{
		NodeId:              request.NodeId,
		Status:              apiStatus,
		IncludeStatusReason: mapper.BoolPtrToBool(request.Params.IncludeStatusReason),
	})
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to list node versions for node %s w/ err: %v", request.NodeId, err)
		return drip.ListNodeVersions500JSONResponse{Message: "Failed to list node versions", Error: err.Error()}, err
	}
	nodeVersions := nodeVersionsResult.NodeVersions
	apiNodeVersions := make([]drip.NodeVersion, 0, len(nodeVersions))
	for _, dbNodeVersion := range nodeVersions {
		apiNodeVersions = append(apiNodeVersions, *mapper.DbNodeVersionToApiNodeVersion(dbNodeVersion))
	}

	log.Ctx(ctx).Info().Msgf("Found %d versions for node %s", len(apiNodeVersions), request.NodeId)
	return drip.ListNodeVersions200JSONResponse(apiNodeVersions), nil
}

func (s *DripStrictServerImplementation) PublishNodeVersion(
	ctx context.Context, request drip.PublishNodeVersionRequestObject) (drip.PublishNodeVersionResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.PublishNodeVersion")()

	// Check if node exists, create if not
	node, err := s.RegistryService.GetNode(ctx, s.Client, request.NodeId)
	if err != nil && !ent.IsNotFound(err) {
		// Case #1: Internal server error when getting node
		log.Ctx(ctx).Error().Msgf("Failed to get node w/ err: %v", err)
		return drip.PublishNodeVersion500JSONResponse{}, err
	} else if err != nil {
		// Case #2: Node not found, create a new node
		node, err = s.RegistryService.CreateNode(ctx, s.Client, request.PublisherId, &request.Body.Node)
		if mapper.IsErrorBadRequest(err) || ent.IsConstraintError(err) {
			log.Ctx(ctx).Error().Msgf("Node creation failed w/ err: %v", err)
			return drip.PublishNodeVersion400JSONResponse{Message: "Failed to create node", Error: err.Error()}, nil
		}
		if err != nil {
			log.Ctx(ctx).Error().Msgf("Node creation failed w/ err: %v", err)
			return drip.PublishNodeVersion500JSONResponse{Message: "Failed to create node", Error: err.Error()}, nil
		}

		log.Ctx(ctx).Info().Msgf("Node %s created successfully", node.ID)
	} else {
		// Case #3: Node already exist, update the node
		updateOneFunc := func(client *ent.Client) *ent.NodeUpdateOne {
			return mapper.ApiUpdateNodeToUpdateFields(node.ID, &request.Body.Node, s.Client)
		}
		_, err = s.RegistryService.UpdateNode(ctx, s.Client, updateOneFunc)
		if err != nil {
			errMessage := "Failed to update node: " + err.Error()
			log.Ctx(ctx).Error().Msgf("Node update failed w/ err: %v", err)
			return drip.PublishNodeVersion400JSONResponse{Message: errMessage}, err
		}
		log.Ctx(ctx).Info().Msgf("Node %s updated successfully", node.ID)
	}

	// Create node version
	nodeVersionCreation, err := s.RegistryService.CreateNodeVersion(
		ctx, s.Client, request.PublisherId, node.ID, &request.Body.NodeVersion)
	if err != nil {
		if ent.IsConstraintError(err) {
			return drip.PublishNodeVersion400JSONResponse{Message: "The node version already exists"}, nil
		}
		log.Ctx(ctx).Error().Msgf("Node version creation failed w/ err: %v", err)
		return drip.PublishNodeVersion400JSONResponse{
			Message: "Failed to create node version: " + err.Error(),
		}, err
	}

	apiNodeVersion := mapper.DbNodeVersionToApiNodeVersion(nodeVersionCreation.NodeVersion)
	log.Ctx(ctx).Info().Msgf("Node version %s published successfully", nodeVersionCreation.NodeVersion.ID)
	return drip.PublishNodeVersion201JSONResponse{
		NodeVersion: apiNodeVersion,
		SignedUrl:   &nodeVersionCreation.SignedUrl,
	}, nil
}

func (s *DripStrictServerImplementation) UpdateNodeVersion(
	ctx context.Context, request drip.UpdateNodeVersionRequestObject) (drip.UpdateNodeVersionResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.UpdateNodeVersion")()

	// Update node version
	updateOne := mapper.ApiUpdateNodeVersionToUpdateFields(request.VersionId, request.Body, s.Client)
	version, err := s.RegistryService.UpdateNodeVersion(ctx, s.Client, updateOne)
	if ent.IsNotFound(err) {
		log.Ctx(ctx).Error().Msgf("Node %s or it's version not found w/ err: %v", request.NodeId, err)
		return drip.UpdateNodeVersion404JSONResponse{Message: "Not Found"}, nil
	}
	if err != nil {
		errMessage := "Failed to update node version"
		log.Ctx(ctx).Error().Msgf("Node version update failed w/ err: %v", err)
		return drip.UpdateNodeVersion500JSONResponse{Message: errMessage, Error: err.Error()}, err
	}

	log.Ctx(ctx).Info().Msgf("Node version %s updated successfully", request.VersionId)
	return drip.UpdateNodeVersion200JSONResponse{
		Changelog:  &version.Changelog,
		Deprecated: &version.Deprecated,
	}, nil
}

// PostNodeVersionReview implements drip.StrictServerInterface.
func (s *DripStrictServerImplementation) PostNodeReview(ctx context.Context, request drip.PostNodeReviewRequestObject) (drip.PostNodeReviewResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.PostNodeReview")()

	if request.Params.Star < 1 || request.Params.Star > 5 {
		log.Ctx(ctx).Error().Msgf("Invalid star received: %d", request.Params.Star)
		return drip.PostNodeReview400Response{}, nil
	}

	userId, err := mapper.GetUserIDFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to get user ID from context w/ err: %v", err)
		return drip.PostNodeReview404JSONResponse{}, err
	}

	nv, err := s.RegistryService.AddNodeReview(ctx, s.Client, request.NodeId, userId, request.Params.Star)
	if ent.IsNotFound(err) {
		log.Ctx(ctx).Error().Msgf("Error retrieving node version w/ err: %v", err)
		return drip.PostNodeReview404JSONResponse{}, nil
	}

	node := mapper.DbNodeToApiNode(nv)
	log.Ctx(ctx).Info().Msgf("Node review for %s stored successfully", request.NodeId)
	return drip.PostNodeReview200JSONResponse(*node), nil

}

func (s *DripStrictServerImplementation) DeleteNodeVersion(
	ctx context.Context, request drip.DeleteNodeVersionRequestObject) (drip.DeleteNodeVersionResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.DeleteNodeVersion")()

	nodeVersion, err := s.RegistryService.GetNodeVersion(ctx, s.Client, request.VersionId)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to get node version w/ err: %v", err)
		return drip.DeleteNodeVersion404JSONResponse{Message: proto.String("Node version not found")}, nil
	}

	err = s.RegistryService.DeleteNodeVersion(ctx, s.Client, nodeVersion.ID.String())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to delete node version w/ err: %v", err)
		return drip.DeleteNodeVersion500JSONResponse{Message: "Failed to delete node version", Error: err.Error()}, err
	}

	log.Ctx(ctx).Info().Msgf("Node version %s deleted successfully", request.VersionId)
	return drip.DeleteNodeVersion204Response{}, nil
}

func (s *DripStrictServerImplementation) GetNodeVersion(
	ctx context.Context, request drip.GetNodeVersionRequestObject) (drip.GetNodeVersionResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.GetNodeVersion")()

	nodeVersion, err := s.RegistryService.GetNodeVersionByVersion(ctx, s.Client, request.NodeId, request.VersionId)
	if ent.IsNotFound(err) {
		log.Ctx(ctx).Error().Msgf("Error retrieving node version w/ err: %v", err)
		return drip.GetNodeVersion404JSONResponse{}, nil
	}
	if err != nil {
		errMessage := "Failed to get node version"
		log.Ctx(ctx).Error().Msgf("Error retrieving node version w/ err: %v", err)
		return drip.GetNodeVersion500JSONResponse{
			Message: errMessage,
			Error:   err.Error(),
		}, err
	}

	apiNodeVersion := mapper.DbNodeVersionToApiNodeVersion(nodeVersion)
	log.Ctx(ctx).Info().Msgf("Node version %s retrieved successfully", request.VersionId)
	return drip.GetNodeVersion200JSONResponse(*apiNodeVersion), nil
}

func (s *DripStrictServerImplementation) ListPersonalAccessTokens(
	ctx context.Context, request drip.ListPersonalAccessTokensRequestObject) (drip.ListPersonalAccessTokensResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.ListPersonalAccessTokens")()

	// List personal access tokens
	personalAccessTokens, err := s.RegistryService.ListPersonalAccessTokens(ctx, s.Client, request.PublisherId)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to list personal access tokens w/ err: %v", err)
		errMessage := "Failed to list personal access tokens."
		return drip.ListPersonalAccessTokens500JSONResponse{Message: errMessage, Error: err.Error()}, err
	}

	apiTokens := make([]drip.PersonalAccessToken, 0, len(personalAccessTokens))
	for _, dbToken := range personalAccessTokens {
		apiTokens = append(apiTokens, *mapper.DbToApiPersonalAccessToken(dbToken))
	}

	log.Ctx(ctx).Info().Msgf("Listed %d personal access tokens for "+
		"publisher ID: %s", len(apiTokens), request.PublisherId)
	return drip.ListPersonalAccessTokens200JSONResponse(apiTokens), nil
}

func (s *DripStrictServerImplementation) CreatePersonalAccessToken(
	ctx context.Context, request drip.CreatePersonalAccessTokenRequestObject) (drip.CreatePersonalAccessTokenResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.CreatePersonalAccessToken")()

	// Create personal access token
	description := ""
	if request.Body.Description != nil {
		description = *request.Body.Description
	}

	personalAccessToken, err := s.RegistryService.CreatePersonalAccessToken(
		ctx, s.Client, request.PublisherId, *request.Body.Name, description)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to create personal access token w/ err: %v", err)
		errMessage := "Failed to create personal access token: " + err.Error()
		return drip.CreatePersonalAccessToken500JSONResponse{Message: errMessage}, err
	}

	log.Ctx(ctx).Info().Msgf("Personal access token created "+
		"successfully for publisher ID: %s", request.PublisherId)
	return drip.CreatePersonalAccessToken201JSONResponse{
		Token: &personalAccessToken.Token,
	}, nil
}

func (s *DripStrictServerImplementation) DeletePersonalAccessToken(
	ctx context.Context, request drip.DeletePersonalAccessTokenRequestObject) (drip.DeletePersonalAccessTokenResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.DeletePersonalAccessToken")()

	// Retrieve user ID from context
	userId, err := mapper.GetUserIDFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to get user ID from context w/ err: %v", err)
		return drip.DeletePersonalAccessToken404JSONResponse{Message: "Invalid user ID"}, err
	}

	// Assert access token belongs to publisher
	err = s.RegistryService.AssertAccessTokenBelongsToPublisher(ctx, s.Client, request.PublisherId, uuid.MustParse(request.TokenId))
	switch {
	case ent.IsNotFound(err):
		log.Ctx(ctx).Warn().Msgf("Publisher with ID %s not found", request.PublisherId)
		return drip.DeletePersonalAccessToken404JSONResponse{Message: "Publisher not found"}, nil

	case drip_services.IsPermissionError(err):
		log.Ctx(ctx).Warn().Msgf(
			"Permission denied for user ID %s on publisher ID %s w/ err: %v", userId, request.PublisherId, err)
		return drip.DeletePersonalAccessToken403JSONResponse{}, err

	case err != nil:
		log.Ctx(ctx).Warn().Msgf("Failed to assert publisher permission %s w/ err: %v", request.PublisherId, err)
		return drip.DeletePersonalAccessToken500JSONResponse{Message: "Failed to assert publisher permission", Error: err.Error()}, err
	}

	// Delete personal access token
	err = s.RegistryService.DeletePersonalAccessToken(ctx, s.Client, uuid.MustParse(request.TokenId))
	if err != nil {
		errMessage := "Failed to delete personal access token: " + err.Error()
		log.Ctx(ctx).Error().Msgf("Token deletion failed w/ err: %v", err)
		return drip.DeletePersonalAccessToken500JSONResponse{Message: errMessage}, err
	}

	log.Ctx(ctx).Info().Msgf("Personal access token %s deleted successfully", request.TokenId)
	return drip.DeletePersonalAccessToken204Response{}, nil
}

func (s *DripStrictServerImplementation) InstallNode(
	ctx context.Context, request drip.InstallNodeRequestObject) (drip.InstallNodeResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.InstallNode")()

	// TODO(robinhuang): Refactor to separate class
	// Get node
	node, err := s.RegistryService.GetNode(ctx, s.Client, request.NodeId)
	if ent.IsNotFound(err) {
		log.Ctx(ctx).Error().Msgf("Error retrieving node w/ err: %v", err)
		return drip.InstallNode404JSONResponse{Message: "Node not found"}, nil
	}
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Error retrieving node w/ err: %v", err)
		return drip.InstallNode500JSONResponse{Message: "Failed to get node"}, err
	}

	// Install node version
	if request.Params.Version == nil {
		s.MixpanelService.Track(ctx, []*mixpanel.Event{
			s.MixpanelService.NewEvent("Install Node", "", map[string]any{
				"Node ID": request.NodeId,
				"Version": "latest",
			}),
		})
		nodeVersion, err := s.RegistryService.GetLatestNodeVersion(ctx, s.Client, request.NodeId)
		if err == nil && nodeVersion == nil {
			log.Ctx(ctx).Error().Msgf("Latest node version not found")
			return drip.InstallNode404JSONResponse{Message: "Not found"}, nil
		}
		if err != nil {
			errMessage := "Failed to get latest node version: " + err.Error()
			log.Ctx(ctx).Error().Msgf("Error retrieving latest node version w/ err: %v", err)
			return drip.InstallNode500JSONResponse{Message: errMessage}, err
		}

		_, err = s.RegistryService.RecordNodeInstallation(ctx, s.Client, node)
		if err != nil {
			errMessage := "Failed to get increment number of node version install: " + err.Error()
			log.Ctx(ctx).Error().Msgf("Error incrementing number of latest node version install w/ err: %v", err)
			return drip.InstallNode500JSONResponse{Message: errMessage}, err
		}

		return drip.InstallNode200JSONResponse(
			*mapper.DbNodeVersionToApiNodeVersion(nodeVersion),
		), nil
	} else {

		nodeVersion, err := s.RegistryService.GetNodeVersionByVersion(ctx, s.Client, request.NodeId, *request.Params.Version)
		if ent.IsNotFound(err) {
			log.Ctx(ctx).Error().Msgf("Error retrieving node version w/ err: %v", err)
			return drip.InstallNode404JSONResponse{Message: "Not found"}, nil
		}
		if err != nil {
			errMessage := "Failed to get specified node version: " + err.Error()
			log.Ctx(ctx).Error().Msgf("Error retrieving node version w/ err: %v", err)
			return drip.InstallNode500JSONResponse{Message: errMessage}, err
		}
		s.MixpanelService.Track(ctx, []*mixpanel.Event{
			s.MixpanelService.NewEvent("Install Node", "", map[string]any{
				"Node ID": request.NodeId,
				"Version": request.Params.Version,
			}),
		})
		_, err = s.RegistryService.RecordNodeInstallation(ctx, s.Client, node)
		if err != nil {
			errMessage := "Failed to get increment number of node version install: " + err.Error()
			log.Ctx(ctx).Error().Msgf("Error incrementing number of latest node version install w/ err: %v", err)
			return drip.InstallNode500JSONResponse{Message: errMessage}, err
		}
		return drip.InstallNode200JSONResponse(
			*mapper.DbNodeVersionToApiNodeVersion(nodeVersion),
		), nil
	}
}

func (s *DripStrictServerImplementation) GetPermissionOnPublisherNodes(
	ctx context.Context, request drip.GetPermissionOnPublisherNodesRequestObject) (drip.GetPermissionOnPublisherNodesResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.GetPermissionOnPublisherNodes")()

	err := s.RegistryService.AssertNodeBelongsToPublisher(ctx, s.Client, request.PublisherId, request.NodeId)
	if err != nil {
		return drip.GetPermissionOnPublisherNodes200JSONResponse{CanEdit: proto.Bool(false)}, nil
	}

	return drip.GetPermissionOnPublisherNodes200JSONResponse{CanEdit: proto.Bool(true)}, nil
}

func (s *DripStrictServerImplementation) GetPermissionOnPublisher(
	ctx context.Context, request drip.GetPermissionOnPublisherRequestObject) (drip.GetPermissionOnPublisherResponseObject, error) {

	return drip.GetPermissionOnPublisher200JSONResponse{CanEdit: proto.Bool(true)}, nil
}

// BanPublisher implements drip.StrictServerInterface.
func (s *DripStrictServerImplementation) BanPublisher(ctx context.Context, request drip.BanPublisherRequestObject) (drip.BanPublisherResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.BanPublisher")()

	userId, err := mapper.GetUserIDFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to get user ID from context w/ err: %v", err)
		return drip.BanPublisher401Response{}, nil
	}
	user, err := s.Client.User.Get(ctx, userId)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to get user ID from context w/ err: %v", err)
		return drip.BanPublisher401Response{}, nil
	}
	if !user.IsAdmin {
		log.Ctx(ctx).Error().Msgf("User is not admin w/ err")
		return drip.BanPublisher403JSONResponse{
			Message: "User is not admin",
		}, nil
	}

	err = s.RegistryService.BanPublisher(ctx, s.Client, request.PublisherId)
	if ent.IsNotFound(err) {
		log.Ctx(ctx).Error().Msgf("Publisher '%s' not found  w/ err: %v", request.PublisherId, err)
		return drip.BanPublisher404JSONResponse{
			Message: "Publisher not found",
		}, nil
	}
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Error banning publisher w/ err: %v", err)
		return drip.BanPublisher500JSONResponse{
			Message: "Error banning publisher",
			Error:   err.Error(),
		}, nil
	}
	return drip.BanPublisher204Response{}, nil
}

// BanPublisherNode implements drip.StrictServerInterface.
func (s *DripStrictServerImplementation) BanPublisherNode(ctx context.Context, request drip.BanPublisherNodeRequestObject) (drip.BanPublisherNodeResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.BanPublisherNode")()

	userId, err := mapper.GetUserIDFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to get user ID from context w/ err: %v", err)
		return drip.BanPublisherNode401Response{}, nil
	}
	user, err := s.Client.User.Get(ctx, userId)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to get user ID from context w/ err: %v", err)
		return drip.BanPublisherNode401Response{}, nil
	}
	if !user.IsAdmin {
		log.Ctx(ctx).Error().Msgf("User is not admin w/ err")
		return drip.BanPublisherNode403JSONResponse{
			Message: "User is not admin",
		}, nil
	}

	err = s.RegistryService.BanNode(ctx, s.Client, request.PublisherId, request.NodeId)
	if ent.IsNotFound(err) {
		log.Ctx(ctx).Error().Msgf("Publisher '%s' or node '%s' not found  w/ err: %v", request.PublisherId, request.NodeId, err)
		return drip.BanPublisherNode404JSONResponse{
			Message: "Publisher or Node not found",
		}, nil
	}
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Error banning node w/ err: %v", err)
		return drip.BanPublisherNode500JSONResponse{
			Message: "Error banning node",
			Error:   err.Error(),
		}, nil
	}
	return drip.BanPublisherNode204Response{}, nil

}

func (s *DripStrictServerImplementation) AdminUpdateNodeVersion(
	ctx context.Context, request drip.AdminUpdateNodeVersionRequestObject) (drip.AdminUpdateNodeVersionResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.AdminUpdateNodeVersion")()

	userId, err := mapper.GetUserIDFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to get user ID from context w/ err: %v", err)
		return drip.AdminUpdateNodeVersion401Response{}, nil
	}
	user, err := s.Client.User.Get(ctx, userId)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to get user ID from context w/ err: %v", err)
		return drip.AdminUpdateNodeVersion401Response{}, nil
	}
	if !user.IsAdmin {
		log.Ctx(ctx).Error().Msgf("User is not admin w/ err")
		return drip.AdminUpdateNodeVersion403JSONResponse{
			Message: "User is not admin",
		}, nil
	}

	nodeVersion, err := s.RegistryService.GetNodeVersionByVersion(ctx, s.Client, request.NodeId, request.VersionNumber)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Error retrieving node version w/ err: %v", err)
		if ent.IsNotFound(err) {
			return drip.AdminUpdateNodeVersion404JSONResponse{}, nil
		}
		return drip.AdminUpdateNodeVersion500JSONResponse{}, err
	}

	dbNodeVersion := mapper.ApiNodeVersionStatusToDbNodeVersionStatus(*request.Body.Status)
	statusReason := ""
	if request.Body.StatusReason != nil {
		statusReason = *request.Body.StatusReason
	}
	err = nodeVersion.Update().SetStatus(dbNodeVersion).SetStatusReason(statusReason).Exec(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to update node version w/ err: %v", err)
		return drip.AdminUpdateNodeVersion500JSONResponse{}, err
	}

	log.Ctx(ctx).Info().Msgf("Node version %s updated successfully", request.VersionNumber)
	return drip.AdminUpdateNodeVersion200JSONResponse{
		Status: request.Body.Status,
	}, nil
}

func (s *DripStrictServerImplementation) SecurityScan(
	ctx context.Context, request drip.SecurityScanRequestObject) (drip.SecurityScanResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.SecurityScan")()

	minAge := 30 * time.Minute
	if request.Params.MinAge != nil {
		minAge = *request.Params.MinAge
	}
	maxNodes := 50
	if request.Params.MaxNodes != nil {
		maxNodes = *request.Params.MaxNodes
	}

	nodeVersionsResult, err := s.RegistryService.ListNodeVersions(ctx, s.Client, &entity.NodeVersionFilter{
		Status:   []schema.NodeVersionStatus{schema.NodeVersionStatusPending},
		MinAge:   minAge,
		PageSize: maxNodes,
		Page:     1,
	})
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to list node versions w/ err: %v", err)
		return drip.SecurityScan500JSONResponse{}, err
	}

	nodeVersions := nodeVersionsResult.NodeVersions
	log.Ctx(ctx).Info().Msgf("Found %d node versions to scan", len(nodeVersions))
	for _, nodeVersion := range nodeVersions {
		err := s.RegistryService.PerformSecurityCheck(ctx, s.Client, nodeVersion)
		if err != nil {
			log.Ctx(ctx).Error().Msgf("Failed to perform security scan w/ err: %v", err)
		}
	}
	return drip.SecurityScan200Response{}, nil
}

func (s *DripStrictServerImplementation) ListAllNodeVersions(
	ctx context.Context, request drip.ListAllNodeVersionsRequestObject) (drip.ListAllNodeVersionsResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.ListAllNodeVersions")()

	// Default values for pagination
	page := 1
	if request.Params.Page != nil {
		page = *request.Params.Page
	}

	pageSize := 10
	maxPageSize := 100
	if request.Params.PageSize != nil {
		// Validate pageSize
		if *request.Params.PageSize > maxPageSize {
			log.Ctx(ctx).Error().Msgf(
				"Requested pageSize %d exceeds maximum allowed %d", *request.Params.PageSize, maxPageSize)
			return drip.ListAllNodeVersions400JSONResponse{
				Message: fmt.Sprintf(
					"Invalid pageSize: %d. The maximum allowed is %d.", *request.Params.PageSize, maxPageSize),
			}, nil
		}
		pageSize = *request.Params.PageSize
	}

	f := &entity.NodeVersionFilter{
		Page:                page,
		PageSize:            pageSize,
		IncludeStatusReason: mapper.BoolPtrToBool(request.Params.IncludeStatusReason),
	}

	if request.Params.Statuses != nil {
		f.Status = mapper.ApiNodeVersionStatusesToDbNodeVersionStatuses(request.Params.Statuses)
	}

	// Log the constructed filter
	log.Ctx(ctx).Info().Msgf("Constructed Filter: %+v", f)

	// List nodes from the registry service
	nodeVersionResults, err := s.RegistryService.ListNodeVersions(ctx, s.Client, f)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to list node versions w/ err: %v", err)
		return drip.ListAllNodeVersions500JSONResponse{
			Message: "Failed to list node versions",
			Error:   err.Error(),
		}, nil
	}

	// Handle no results
	if len(nodeVersionResults.NodeVersions) == 0 {
		log.Ctx(ctx).Info().Msgf("No node versions found. Total: %d", nodeVersionResults.Total)
		return drip.ListAllNodeVersions200JSONResponse{
			Versions:   &[]drip.NodeVersion{},
			Total:      &nodeVersionResults.Total,
			Page:       &nodeVersionResults.Page,
			PageSize:   &nodeVersionResults.Limit,
			TotalPages: &nodeVersionResults.TotalPages,
		}, nil
	}

	// Transform DB results into API results
	apiNodeVersions := make([]drip.NodeVersion, 0, len(nodeVersionResults.NodeVersions))
	for _, dbNodeVersion := range nodeVersionResults.NodeVersions {
		apiNodeVersions = append(apiNodeVersions, *mapper.DbNodeVersionToApiNodeVersion(dbNodeVersion))
	}

	// Log success
	log.Ctx(ctx).Info().Msgf(
		"Found %d node versions. Total: %d", len(nodeVersionResults.NodeVersions), nodeVersionResults.Total)
	return drip.ListAllNodeVersions200JSONResponse{
		Versions:   &apiNodeVersions,
		Total:      &nodeVersionResults.Total,
		Page:       &nodeVersionResults.Page,
		PageSize:   &nodeVersionResults.Limit,
		TotalPages: &nodeVersionResults.TotalPages,
	}, nil
}

func (s *DripStrictServerImplementation) ReindexNodes(ctx context.Context, request drip.ReindexNodesRequestObject) (res drip.ReindexNodesResponseObject, err error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.ReindexNodes")()

	// create new context with logger from original Context
	reindexCtx := drip_logging.ReuseContextLogger(ctx, context.Background())

	err = s.RegistryService.ReindexAllNodesBackground(
		reindexCtx, s.Client, request.Params.MaxBatch, request.Params.MinAge)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to trigger reindex all nodes w/ err: %v", err)
		return drip.ReindexNodes500JSONResponse{Message: "Failed to trigger reindex nodes", Error: err.Error()}, nil
	}

	log.Ctx(ctx).Info().Msgf("Triggering Reindex nodes successful")
	return drip.ReindexNodes200Response{}, nil
}

// CreateComfyNodes bulk-stores comfy-nodes extraction result for a node version
func (impl *DripStrictServerImplementation) CreateComfyNodes(
	ctx context.Context, request drip.CreateComfyNodesRequestObject) (res drip.CreateComfyNodesResponseObject, err error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.CreateComfyNodes")()

	cb := mapper.ApiComfyNodeCloudBuildToDbComfyNodeCloudBuild(request.Body.CloudBuildInfo)
	// Check if extraction was marked as unsuccessful
	if request.Body.Success != nil && !*request.Body.Success {
		reason := "unknown"
		if request.Body.Reason != nil {
			reason = *request.Body.Reason
		}
		log.Ctx(ctx).Warn().Msgf(
			"Comfy nodes extraction failed for %s %s: %s", request.NodeId, request.Version, reason)

		err = impl.RegistryService.MarkComfyNodeExtractionFailed(ctx, impl.Client, request.NodeId, request.Version, cb)

	} else {
		// Attempt to create comfy nodes in the registry
		err = impl.RegistryService.CreateComfyNodes(
			ctx, impl.Client, request.NodeId, request.Version, *request.Body.Nodes, cb)
	}

	// Handle specific error scenarios
	if ent.IsNotFound(err) {
		log.Ctx(ctx).Error().Stringer("cloud_build_info", cb).Msgf("Node or node version not found w/ err: %v", err)
		return drip.CreateComfyNodes404JSONResponse{
			Message: "Node or node version not found", Error: err.Error()}, nil
	}

	if errors.Is(err, drip_services.ErrComfyNodesAlreadyExist) {
		log.Ctx(ctx).Error().Stringer("cloud_build_info", cb).Msgf(
			"Comfy nodes extraction result for %s %s already set", request.NodeId, request.Version)
		return drip.CreateComfyNodes409JSONResponse{
			Message: "Comfy nodes extraction result already set", Error: err.Error()}, nil
	}

	if err != nil {
		log.Ctx(ctx).Error().Stringer("cloud_build_info", cb).Msgf("Failed to store comfy nodes extraction w/ err: %v", err)
		return drip.CreateComfyNodes500JSONResponse{
			Message: "Failed to store comfy nodes extraction", Error: err.Error()}, nil
	}

	log.Ctx(ctx).Info().Stringer("cloud_build_info", cb).Msg("CreateComfyNodes successful")
	return drip.CreateComfyNodes204Response{}, nil
}

// GetComfyNode returns a specific comfy-node of a certain node version
func (impl *DripStrictServerImplementation) GetComfyNode(
	ctx context.Context, request drip.GetComfyNodeRequestObject) (res drip.GetComfyNodeResponseObject, err error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.GetComfyNode")()

	// Retrieve the comfy-node from the registry
	n, err := impl.RegistryService.GetComfyNode(ctx, impl.Client, request.NodeId, request.Version, request.ComfyNodeId)

	// Handle node or version not found
	if ent.IsNotFound(err) {
		log.Ctx(ctx).Error().Msgf("Node or node version or comfy node not found w/ err: %v", err)
		return drip.GetComfyNode404JSONResponse{
			Message: "Node or node version or comfy node not found", Error: err.Error()}, nil
	}

	// Map database comfy-node to API representation
	cn := mapper.DBComfyNodeToApiComfyNode(n)
	if cn == nil {
		log.Ctx(ctx).Error().Msg("Comfy Node not found")
		return drip.GetComfyNode404JSONResponse{Message: "Comfy Node not found"}, nil
	}

	log.Ctx(ctx).Info().Msg("GetComfyNode successful")
	res = drip.GetComfyNode200JSONResponse(*cn)
	return
}

// ComfyNodesBackfill triggers a backfill process for comfy-nodes
func (impl *DripStrictServerImplementation) ComfyNodesBackfill(
	ctx context.Context, request drip.ComfyNodesBackfillRequestObject) (drip.ComfyNodesBackfillResponseObject, error) {
	defer tracing.TraceDefaultSegment(ctx, "DripStrictServerImplementation.ComfyNodesBackfill")()

	// Trigger the backfill process with a specified maximum node
	err := impl.RegistryService.TriggerComfyNodesBackfill(ctx, impl.Client, request.Params.MaxNode)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to trigger comfy nodes backfill w/ err: %v", err)
		return drip.ComfyNodesBackfill500JSONResponse{
			Message: "Failed to trigger comfy nodes backfill", Error: err.Error()}, nil
	}

	log.Ctx(ctx).Info().Msg("ComfyNodesBackfill successful")
	return drip.ComfyNodesBackfill204Response{}, nil
}
