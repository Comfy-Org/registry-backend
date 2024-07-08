package implementation

import (
	"context"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/publisher"
	"registry-backend/ent/schema"
	"registry-backend/mapper"
	drip_services "registry-backend/services/registry"
	"time"

	"github.com/google/uuid"
	"github.com/mixpanel/mixpanel-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

func (impl *DripStrictServerImplementation) ListPublishersForUser(
	ctx context.Context, request drip.ListPublishersForUserRequestObject) (drip.ListPublishersForUserResponseObject, error) {
	log.Ctx(ctx).Debug().Msg("ListPublishersForUser called.")

	// Extract user ID from context
	userId, err := mapper.GetUserIDFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to get user ID from context w/ err: %v", err)
		return drip.ListPublishersForUser400JSONResponse{Message: "Invalid user ID"}, err
	}

	// Call the service to list publishers
	log.Ctx(ctx).Info().Msgf("Fetching publishers for user %s", userId)
	publishers, err := impl.RegistryService.ListPublishers(ctx, impl.Client, &drip_services.PublisherFilter{
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
	// Log the incoming request for validation
	log.Ctx(ctx).Info().Msgf("ValidatePublisher request with username: %s", request.Params.Username)

	// Check if the username is empty
	name := request.Params.Username
	if name == "" {
		log.Ctx(ctx).Error().Msg("Username parameter is missing")
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
	// Log the incoming request
	log.Ctx(ctx).Info().Msgf("CreatePublisher request called")

	// Extract user ID from context
	userId, err := mapper.GetUserIDFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to get user ID from context w/ err: %v", err)
		return drip.CreatePublisher400JSONResponse{Message: "Invalid user ID"}, err
	}

	log.Ctx(ctx).Info().Msgf("Checking if user ID %s has reached the maximum number of publishers", userId)
	userPublishers, err := s.RegistryService.ListPublishers(
		ctx, s.Client, &drip_services.PublisherFilter{UserID: userId})
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
	publisherId := request.PublisherId
	log.Ctx(ctx).Info().Msgf("GetPublisher request received for publisher ID: %s", publisherId)

	publisher, err := s.RegistryService.GetPublisher(ctx, s.Client, publisherId)
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
	log.Ctx(ctx).Info().Msgf("UpdatePublisher called with publisher ID: %s", request.PublisherId)

	log.Ctx(ctx).Info().Msgf("Updating publisher with ID %s", request.PublisherId)
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
	log.Ctx(ctx).Info().Msgf("CreateNode called with publisher ID: %s", request.PublisherId)

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
	log.Ctx(ctx).Info().Msgf("ListNodesForPublisher request received for publisher ID: %s", request.PublisherId)

	nodeResults, err := s.RegistryService.ListNodes(
		ctx, s.Client /*page=*/, 1 /*limit=*/, 10, &drip_services.NodeFilter{
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

	log.Ctx(ctx).Info().Msg("ListAllNodes request received")

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
	filter := &drip_services.NodeFilter{}
	if request.Params.IncludeBanned != nil {
		filter.IncludeBanned = *request.Params.IncludeBanned
	}

	// List nodes from the registry service
	nodeResults, err := s.RegistryService.ListNodes(ctx, s.Client, page, limit, filter)
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
	log.Ctx(ctx).Info().Msg("SearchNodes request received")

	// Set default values for pagination parameters
	page := 1
	if request.Params.Page != nil {
		page = *request.Params.Page
	}
	limit := 10
	if request.Params.Limit != nil {
		limit = *request.Params.Limit
	}

	f := &drip_services.NodeFilter{}
	if request.Params.Search != nil {
		f.Search = *request.Params.Search
	}
	if request.Params.IncludeBanned != nil {
		f.IncludeBanned = *request.Params.IncludeBanned
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
		if dbNode.Edges.Versions != nil && len(dbNode.Edges.Versions) > 0 {
			latestVersion, err := s.RegistryService.GetLatestNodeVersion(ctx, s.Client, dbNode.ID)
			if err == nil {
				apiNode.LatestVersion = mapper.DbNodeVersionToApiNodeVersion(latestVersion)
			} else {
				log.Ctx(ctx).Error().Msgf("Failed to get latest version for node %s w/ err: %v", dbNode.ID, err)
			}
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

	log.Ctx(ctx).Info().Msgf("DeleteNode request received for node ID: %s", request.NodeId)

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
	log.Ctx(ctx).Info().Msgf("GetNode request received for node ID: %s", request.NodeId)

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

	log.Ctx(ctx).Info().Msgf("UpdateNode request received for node ID: %s", request.NodeId)

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
	log.Ctx(ctx).Info().Msgf("ListNodeVersions request received for node ID: %s", request.NodeId)

	apiStatus := mapper.ApiNodeVersionStatusesToDbNodeVersionStatuses(request.Params.Statuses)

	nodeVersionsResult, err := s.RegistryService.ListNodeVersions(ctx, s.Client, &drip_services.NodeVersionFilter{
		NodeId: request.NodeId,
		Status: apiStatus,
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
	log.Ctx(ctx).Info().Msgf("PublishNodeVersion request received for node ID: %s", request.NodeId)

	// Check if node exists, create if not
	node, err := s.RegistryService.GetNode(ctx, s.Client, request.NodeId)
	if err != nil && !ent.IsNotFound(err) {
		log.Ctx(ctx).Error().Msgf("Failed to get node w/ err: %v", err)
		// TODO(James): create a new error code for this.
		return drip.PublishNodeVersion500JSONResponse{}, err
	} else if err != nil {
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
		// TODO(james): distinguish between not found vs. nodes that belong to other publishers
		// If node already exists, validate ownership
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
	nodeVersionCreation, err := s.RegistryService.CreateNodeVersion(ctx, s.Client, request.PublisherId, node.ID, &request.Body.NodeVersion)
	if err != nil {
		if ent.IsConstraintError(err) {
			return drip.PublishNodeVersion400JSONResponse{Message: "The node version already exists"}, nil
		}
		errMessage := "Failed to create node version: " + err.Error()
		log.Ctx(ctx).Error().Msgf("Node version creation failed w/ err: %v", err)
		return drip.PublishNodeVersion400JSONResponse{Message: errMessage}, err
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

	log.Ctx(ctx).Info().Msgf("UpdateNodeVersion request received for node ID: "+
		"%s, version ID: %s", request.NodeId, request.VersionId)

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
	log.Ctx(ctx).Info().Msgf("PostNodeReview request received for "+
		"node ID: %s", request.NodeId)

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
	log.Ctx(ctx).Info().Msgf("DeleteNodeVersion request received for node ID: "+
		"%s, version ID: %s", request.NodeId, request.VersionId)

	// Directly return the message that node versions cannot be deleted
	errMessage := "Cannot delete node versions. Please deprecate it instead."
	log.Ctx(ctx).Warn().Msg(errMessage)
	return drip.DeleteNodeVersion404JSONResponse{
		Message: proto.String(errMessage),
	}, nil
}

func (s *DripStrictServerImplementation) GetNodeVersion(
	ctx context.Context, request drip.GetNodeVersionRequestObject) (drip.GetNodeVersionResponseObject, error) {
	log.Ctx(ctx).Info().Msgf("GetNodeVersion request received for "+
		"node ID: %s, version ID: %s", request.NodeId, request.VersionId)

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
	log.Ctx(ctx).Info().Msgf("ListPersonalAccessTokens request received for publisher ID: %s", request.PublisherId)

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

	log.Ctx(ctx).Info().Msgf("CreatePersonalAccessToken request received "+
		"for publisher ID: %s", request.PublisherId)

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

	log.Ctx(ctx).Info().Msgf("DeletePersonalAccessToken request received for token ID: %s", request.TokenId)

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
		log.Ctx(ctx).Info().Msgf("Publisher with ID %s not found", request.PublisherId)
		return drip.DeletePersonalAccessToken404JSONResponse{Message: "Publisher not found"}, nil

	case drip_services.IsPermissionError(err):
		log.Ctx(ctx).Error().Msgf(
			"Permission denied for user ID %s on publisher ID %s w/ err: %v", userId, request.PublisherId, err)
		return drip.DeletePersonalAccessToken403JSONResponse{}, err

	case err != nil:
		log.Ctx(ctx).Error().Msgf("Failed to assert publisher permission %s w/ err: %v", request.PublisherId, err)
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
	// TODO(robinhuang): Refactor to separate class
	mp := mixpanel.NewApiClient("f919d1b9da9a57482453c72ef7b16d88")
	log.Ctx(ctx).Info().Msgf("InstallNode request received for node ID: %s", request.NodeId)

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
		_, err = s.RegistryService.RecordNodeInstalation(ctx, s.Client, node)
		if err != nil {
			errMessage := "Failed to get increment number of node version install: " + err.Error()
			log.Ctx(ctx).Error().Msgf("Error incrementing number of latest node version install w/ err: %v", err)
			return drip.InstallNode500JSONResponse{Message: errMessage}, err
		}
		mp.Track(ctx, []*mixpanel.Event{
			mp.NewEvent("Install Node Latest", "", map[string]any{
				"Node ID": request.NodeId,
				"Version": nodeVersion.Version,
			}),
		})
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
		_, err = s.RegistryService.RecordNodeInstalation(ctx, s.Client, node)
		if err != nil {
			errMessage := "Failed to get increment number of node version install: " + err.Error()
			log.Ctx(ctx).Error().Msgf("Error incrementing number of latest node version install w/ err: %v", err)
			return drip.InstallNode500JSONResponse{Message: errMessage}, err
		}
		mp.Track(ctx, []*mixpanel.Event{
			mp.NewEvent("Install Node", "", map[string]any{
				"Node ID": request.NodeId,
				"Version": nodeVersion.Version,
			}),
		})
		return drip.InstallNode200JSONResponse(
			*mapper.DbNodeVersionToApiNodeVersion(nodeVersion),
		), nil
	}
}

func (s *DripStrictServerImplementation) GetPermissionOnPublisherNodes(
	ctx context.Context, request drip.GetPermissionOnPublisherNodesRequestObject) (drip.GetPermissionOnPublisherNodesResponseObject, error) {

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
	minAge := 30 * time.Minute
	if request.Params.MinAge != nil {
		minAge = *request.Params.MinAge
	}
	maxNodes := 50
	if request.Params.MaxNodes != nil {
		maxNodes = *request.Params.MaxNodes
	}
	nodeVersionsResult, err := s.RegistryService.ListNodeVersions(ctx, s.Client, &drip_services.NodeVersionFilter{
		Status:   []schema.NodeVersionStatus{schema.NodeVersionStatusPending},
		MinAge:   minAge,
		PageSize: maxNodes,
		Page:     1,
	})
	nodeVersions := nodeVersionsResult.NodeVersions

	log.Ctx(ctx).Info().Msgf("Found %d node versions to scan", len(nodeVersions))

	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to list node versions w/ err: %v", err)
		return drip.SecurityScan500JSONResponse{}, err
	}

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
	log.Ctx(ctx).Info().Msg("ListAllNodeVersions request received")

	page := 1
	if request.Params.Page != nil {
		page = *request.Params.Page
	}
	pageSize := 10
	if request.Params.PageSize != nil && *request.Params.PageSize < 100 {
		pageSize = *request.Params.PageSize
	}
	f := &drip_services.NodeVersionFilter{
		Page:     page,
		PageSize: pageSize,
	}

	if request.Params.Statuses != nil {
		f.Status = mapper.ApiNodeVersionStatusesToDbNodeVersionStatuses(request.Params.Statuses)
	}

	// List nodes from the registry service
	nodeVersionResults, err := s.RegistryService.ListNodeVersions(ctx, s.Client, f)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to list node versions w/ err: %v", err)
		return drip.ListAllNodeVersions500JSONResponse{Message: "Failed to list node versions", Error: err.Error()}, nil
	}

	if len(nodeVersionResults.NodeVersions) == 0 {
		log.Ctx(ctx).Info().Msg("No node versions found")
		return drip.ListAllNodeVersions200JSONResponse{
			Versions:   &[]drip.NodeVersion{},
			Total:      &nodeVersionResults.Total,
			Page:       &nodeVersionResults.Page,
			PageSize:   &nodeVersionResults.Limit,
			TotalPages: &nodeVersionResults.TotalPages,
		}, nil
	}

	apiNodeVersions := make([]drip.NodeVersion, 0, len(nodeVersionResults.NodeVersions))
	for _, dbNodeVersion := range nodeVersionResults.NodeVersions {
		apiNodeVersions = append(apiNodeVersions, *mapper.DbNodeVersionToApiNodeVersion(dbNodeVersion))
	}

	log.Ctx(ctx).Info().Msgf("Found %d node versions", nodeVersionResults.Total)
	return drip.ListAllNodeVersions200JSONResponse{
		Versions:   &apiNodeVersions,
		Total:      &nodeVersionResults.Total,
		Page:       &nodeVersionResults.Page,
		PageSize:   &nodeVersionResults.Limit,
		TotalPages: &nodeVersionResults.TotalPages,
	}, nil
}

func (s *DripStrictServerImplementation) ReindexNodes(ctx context.Context, request drip.ReindexNodesRequestObject) (res drip.ReindexNodesResponseObject, err error) {
	log.Ctx(ctx).Info().Msg("ReindexNodes request received")
	err = s.RegistryService.ReindexAllNodes(ctx, s.Client)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("Failed to list node versions w/ err: %v", err)
		return drip.ReindexNodes500JSONResponse{Message: "Failed to reindex nodes", Error: err.Error()}, nil
	}

	log.Ctx(ctx).Info().Msgf("Reindex nodes successful")
	return drip.ReindexNodes200Response{}, nil
}
