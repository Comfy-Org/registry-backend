package dripservices

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"registry-backend/common"
	"registry-backend/config"
	"registry-backend/db"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/comfynode"
	"registry-backend/ent/node"
	"registry-backend/ent/nodeversion"
	"registry-backend/ent/personalaccesstoken"
	"registry-backend/ent/predicate"
	"registry-backend/ent/publisher"
	"registry-backend/ent/publisherpermission"
	"registry-backend/ent/schema"
	"registry-backend/ent/user"
	"registry-backend/entity"
	"registry-backend/gateways/algolia"
	"registry-backend/gateways/discord"
	"registry-backend/gateways/pubsub"
	gateway "registry-backend/gateways/slack"
	"registry-backend/gateways/storage"
	"registry-backend/mapper"
	drip_metric "registry-backend/server/middleware/metric"
	"strings"
	"sync"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/Masterminds/semver/v3"
	"github.com/google/uuid"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

type RegistryService struct {
	storageService storage.StorageService
	pubsubService  pubsub.PubSubService
	slackService   gateway.SlackService
	algolia        algolia.AlgoliaService
	discordService discord.DiscordService
	config         *config.Config
	newRelicApp    *newrelic.Application
}

func NewRegistryService(storageSvc storage.StorageService, pubsubService pubsub.PubSubService, slackSvc gateway.SlackService, discordSvc discord.DiscordService, algoliaSvc algolia.AlgoliaService, config *config.Config, newRelicApp *newrelic.Application) *RegistryService {
	return &RegistryService{
		storageService: storageSvc,
		pubsubService:  pubsubService,
		slackService:   slackSvc,
		discordService: discordSvc,
		algolia:        algoliaSvc,
		config:         config,
		newRelicApp:    newRelicApp,
	}
}

// ListNodes retrieves a paginated list of nodes with optional filtering.
func (s *RegistryService) ListNodes(ctx context.Context, client *ent.Client, page, limit int, filter *entity.NodeFilter) (*entity.ListNodesResult, error) {
	// Start New Relic transaction segment
	var txn *newrelic.Transaction
	if txnCtx := newrelic.FromContext(ctx); txnCtx != nil {
		txn = txnCtx
		txn.Application().RecordCustomMetric(
			"Custom/ListNodes/Limit",
			float64(limit),
		)
		segment := txn.StartSegment("RegistryService.ListNodes")
		defer segment.End()
	}

	// Ensure valid pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Initialize the query with relationships
	query := client.Node.Query().WithPublisher()

	// Apply filters if provided
	if filter != nil {
		var predicates []predicate.Node

		// Filter by PublisherID
		if filter.PublisherID != "" {
			predicates = append(predicates, node.PublisherID(filter.PublisherID))
		}

		// Filter by search term across multiple fields
		if filter.Search != "" {
			predicates = append(predicates, node.Or(
				node.IDContainsFold(filter.Search),
				node.NameContainsFold(filter.Search),
				node.DescriptionContainsFold(filter.Search),
				node.AuthorContainsFold(filter.Search),
			))
		}

		// Exclude banned nodes if not requested
		if !filter.IncludeBanned {
			predicates = append(predicates, node.StatusNEQ(schema.NodeStatusBanned))
		}

		// Apply predicates to the query
		if len(predicates) > 1 {
			query.Where(node.And(predicates...))
		} else if len(predicates) == 1 {
			query.Where(predicates[0])
		}
	}

	// Calculate pagination offset
	offset := (page - 1) * limit

	// Count total nodes with New Relic datastore segment
	var total int
	var err error
	if txn != nil {
		segment := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres, // Change based on your DB
			Collection: node.Table,                 // Table name
			Operation:  "COUNT",
			StartTime:  txn.StartSegmentNow(),
		}
		total, err = query.Count(ctx)
		segment.End()
	} else {
		total, err = query.Count(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to count nodes: %w", err)
	}

	// Fetch nodes with pagination and New Relic datastore segment
	query = s.decorateNodeQueryWithLatestVersion(query).Offset(offset).Limit(limit)
	var nodes []*ent.Node
	if txn != nil {
		segment := newrelic.DatastoreSegment{
			Product:    newrelic.DatastorePostgres,
			Collection: node.Table,
			Operation:  "SELECT",
			StartTime:  txn.StartSegmentNow(),
		}
		nodes, err = query.All(ctx)
		segment.End()
	} else {
		nodes, err = query.All(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	// Calculate total pages
	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	// Return the result
	return &entity.ListNodesResult{
		Total:      total,
		Nodes:      nodes,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// ListPublishers queries the Publisher table with an optional user ID filter via PublisherPermission
func (s *RegistryService) ListPublishers(ctx context.Context, client *ent.Client, filter *entity.PublisherFilter) ([]*ent.Publisher, error) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.ListPublishers")
		defer segment.End()
	}
	log.Ctx(ctx).Info().Msg("Listing publishers")

	query := client.Publisher.Query()

	if filter != nil && filter.UserID != "" {
		query = query.Where(
			// Ensure that the publisher has the associated permission with the specific user ID
			publisher.HasPublisherPermissionsWith(publisherpermission.UserIDEQ(filter.UserID)),
		)
	}

	publishers, err := query.All(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to list publishers")
		return nil, err
	}

	log.Ctx(ctx).Info().Msgf("Found %d publishers", len(publishers))
	return publishers, nil
}

func (s *RegistryService) CreatePublisher(ctx context.Context, client *ent.Client, userId string, publisher *drip.Publisher) (*ent.Publisher, error) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.CreatePublisher")
		defer segment.End()
	}
	publisherValid := mapper.ValidatePublisher(publisher)
	if publisherValid != nil {
		return nil, fmt.Errorf("invalid publisher: %w", publisherValid)
	}
	return db.WithTxResult(ctx, client, func(tx *ent.Tx) (*ent.Publisher, error) {
		newPublisher, err := mapper.ApiCreatePublisherToDb(publisher, tx.Client())
		log.Ctx(ctx).Info().Msgf("creating publisher with fields: %v", newPublisher.Mutation().Fields())
		if err != nil {
			return nil, fmt.Errorf("failed to map publisher: %w", err)
		}
		publisher, err := newPublisher.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create publisher: %w", err)
		}

		_, err = tx.PublisherPermission.Create().SetPublisherID(publisher.ID).
			SetUserID(userId).
			SetPermission(schema.PublisherPermissionTypeOwner).
			Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create publisher permissions: %w", err)
		}
		return publisher, nil
	})
}

func (s *RegistryService) UpdatePublisher(ctx context.Context, client *ent.Client, update *ent.PublisherUpdateOne) (*ent.Publisher, error) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.UpdatePublisher")
		defer segment.End()
	}
	log.Ctx(ctx).Info().Msgf("updating publisher fields: %v", update.Mutation().Fields())
	publisher, err := update.Save(ctx)
	log.Ctx(ctx).Info().Msgf("success: updated publisher: %v", publisher)
	if err != nil {
		return nil, fmt.Errorf("failed to create publisher: %w", err)
	}

	return publisher, nil
}

func (s *RegistryService) GetPublisher(ctx context.Context, client *ent.Client, publisherID string) (*ent.Publisher, error) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.GetPublisher")
		defer segment.End()
	}
	log.Ctx(ctx).Info().Msgf("getting publisher: %v", publisherID)
	publisher, err := client.Publisher.
		Query().
		Where(publisher.IDEQ(publisherID)).
		WithPublisherPermissions(func(ppq *ent.PublisherPermissionQuery) {
			ppq.WithUser()
		}).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get publisher: %w", err)
	}
	return publisher, nil
}

func (s *RegistryService) CreatePersonalAccessToken(ctx context.Context, client *ent.Client, publisherID, name, description string) (*ent.PersonalAccessToken, error) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.CreatePersonalAccessToken")
		defer segment.End()
	}
	log.Ctx(ctx).Info().Msgf("creating personal access token for publisher: %v", publisherID)
	token := uuid.New().String()
	pat, err := client.PersonalAccessToken.
		Create().
		SetPublisherID(publisherID).
		SetName(name).
		SetDescription(description).
		SetToken(token).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create personal access token: %w", err)
	}
	return pat, nil
}

func (s *RegistryService) ListPersonalAccessTokens(ctx context.Context, client *ent.Client, publisherID string) ([]*ent.PersonalAccessToken, error) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.ListPersonalAccessTokens")
		defer segment.End()
	}
	pats, err := client.PersonalAccessToken.Query().
		Where(personalaccesstoken.PublisherIDEQ(publisherID)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list personal access tokens: %w", err)
	}
	return pats, nil
}

func (s *RegistryService) DeletePersonalAccessToken(ctx context.Context, client *ent.Client, tokenID uuid.UUID) error {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.DeletePersonalAccessToken")
		defer segment.End()
	}
	log.Ctx(ctx).Info().Msgf("deleting personal access token: %v", tokenID)
	err := client.PersonalAccessToken.
		DeleteOneID(tokenID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete personal access token: %w", err)
	}
	return nil
}

func (s *RegistryService) CreateNode(ctx context.Context, client *ent.Client, publisherId string, node *drip.Node) (*ent.Node, error) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.CreateNode")
		defer segment.End()
	}
	validNode := mapper.ValidateNode(node)
	if validNode != nil {
		return nil, fmt.Errorf("invalid node: %w", validNode)
	}

	var createdNode *ent.Node
	err := db.WithTx(ctx, client, func(tx *ent.Tx) (err error) {
		createNode, err := mapper.ApiCreateNodeToDb(publisherId, node, tx.Client())
		log.Ctx(ctx).Info().Msgf("creating node with fields: %v", createNode.Mutation().Fields())
		if err != nil {
			return fmt.Errorf("failed to map node: %w", err)
		}

		createdNode, err = createNode.Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create node: %w", err)
		}

		err = s.algolia.IndexNodes(ctx, createdNode)
		if err != nil {
			return fmt.Errorf("failed to index node: %w", err)
		}

		return
	})

	return createdNode, err
}

func (s *RegistryService) UpdateNode(ctx context.Context, client *ent.Client, updateFunc func(client *ent.Client) *ent.NodeUpdateOne) (*ent.Node, error) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.UpdateNode")
		defer segment.End()
	}
	var n *ent.Node
	err := db.WithTx(ctx, client, func(tx *ent.Tx) (err error) {
		update := updateFunc(tx.Client())
		log.Ctx(ctx).Info().Msgf("updating node fields: %v", update.Mutation().Fields())

		n, err = update.Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to update node: %w", err)
		}

		_, err = s.indexNodeWithLatestVersion(ctx, tx.Client(), n.ID)
		if err != nil {
			return fmt.Errorf("failed to index node: %w", err)
		}

		return err
	})
	return n, err
}

func (s *RegistryService) GetNode(ctx context.Context, client *ent.Client, nodeID string) (*ent.Node, error) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.GetNode")
		defer segment.End()
	}

	log.Ctx(ctx).Info().Msgf("getting node: %v", nodeID)
	node, err := client.Node.Get(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}
	return node, nil
}

func (s *RegistryService) CreateNodeVersion(ctx context.Context, client *ent.Client, publisherID, nodeID string, nodeVersion *drip.NodeVersion) (*NodeVersionCreation, error) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.CreateNodeVersion")
		defer segment.End()
	}
	log.Ctx(ctx).Info().Msgf("creating node version: %v for nodeId %v", nodeVersion, nodeID)
	bucketName := "comfy-registry"
	return db.WithTxResult(ctx, client, func(tx *ent.Tx) (*NodeVersionCreation, error) {
		// If the node version is not provided, we will generate a new version
		if nodeVersion.Version != nil {
			defaultVersion, err := semver.NewVersion(*nodeVersion.Version)
			if err != nil {
				return nil, err
			}

			nodeVersion.Version = proto.String(defaultVersion.String())
		}

		// Create a new storage file for the node version
		objectPath := fmt.Sprintf("%s/%s/%s/%s", publisherID, nodeID, *nodeVersion.Version, "node.zip")
		storageFile := tx.StorageFile.Create().
			SetBucketName(bucketName).
			SetFilePath(objectPath).
			SetFileType("zip").
			// Sample URL: https://storage.googleapis.com/comfy-registry/james-test-publisher/comfyui-inspire-pack/1.0.0/node.zip
			SetFileURL(fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectPath)).
			SaveX(ctx)
		signedUrl, err := s.storageService.GenerateSignedURL(bucketName, objectPath)
		if err != nil {
			return nil, fmt.Errorf("failed to generate signed url: %w", err)
		}
		log.Ctx(ctx).Info().Msgf("generated signed url: %v", signedUrl)

		newNodeVersion := mapper.ApiCreateNodeVersionToDb(nodeID, nodeVersion, tx.Client())
		newNodeVersion.SetStorageFile(storageFile)
		createdNodeVersion, err := newNodeVersion.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create node version: %w", err)
		}

		err = s.algolia.IndexNodeVersions(ctx, createdNodeVersion)
		if err != nil {
			return nil, fmt.Errorf("failed to index node version: %w", err)
		}

		message := fmt.Sprintf("Version %s of node %s was published successfully. Publisher: %s. https://registry.comfy.org/nodes/%s", createdNodeVersion.Version, createdNodeVersion.NodeID, publisherID, nodeID)
		slackErr := s.slackService.SendRegistryMessageToSlack(message)
		// Send the message to the private channel
		s.discordService.SendSecurityCouncilMessage(message, true)
		if slackErr != nil {
			log.Ctx(ctx).Error().Msgf("Failed to send message to Slack w/ err: %v", slackErr)
			drip_metric.IncrementCustomCounterMetric(ctx, drip_metric.CustomCounterIncrement{
				Type:   "slack-send-error",
				Val:    1,
				Labels: map[string]string{},
			})
		}

		return &NodeVersionCreation{
			NodeVersion: createdNodeVersion,
			SignedUrl:   signedUrl,
		}, nil
	})
}

type NodeVersionCreation struct {
	NodeVersion *ent.NodeVersion
	SignedUrl   string
}

func (s *RegistryService) ListNodeVersions(ctx context.Context, client *ent.Client, filter *entity.NodeVersionFilter) (*entity.ListNodeVersionsResult, error) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.ListNodeVersions")
		defer segment.End()
	}
	query := client.NodeVersion.Query().
		WithStorageFile().
		Order(ent.Desc(nodeversion.FieldVersion))

	if filter.NodeId != "" {
		log.Ctx(ctx).Info().Msgf("listing node versions: %v", filter.NodeId)
		query.Where(nodeversion.NodeIDEQ(filter.NodeId))
	}

	if len(filter.Status) > 0 {
		log.Ctx(ctx).Info().Msgf("listing node versions with status: %v", filter.Status)
		query.Where(nodeversion.StatusIn(filter.Status...))
	}

	if filter.MinAge > 0 {
		log.Ctx(ctx).Info().Msgf("listing node versions with min age: %v", filter.MinAge)
		query.Where(nodeversion.CreateTimeLT(time.Now().Add(-filter.MinAge)))
	}

	// Note: custom SELECT statement cause errors in the ent framework when using the Count method.
	// We need to include the logic to exclude certain fields after the count query is executed.
	total, err := query.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count node versions: %w", err)
	}

	// Apply pagination
	// Note: the pagination logic needs to be applied after the count query is executed
	if filter.Page > 0 && filter.PageSize > 0 {
		log.Ctx(ctx).Info().Msgf(
			"listing node versions with pagination: page %v, limit %v", filter.Page, filter.PageSize)
		query.Offset((filter.Page - 1) * filter.PageSize).Limit(filter.PageSize)
	}

	// By default, we are selecting all fields. If the status reason is not required, we will exclude it
	if !filter.IncludeStatusReason {
		columns := make([]string, 0, len(nodeversion.Columns))
		for _, column := range nodeversion.Columns {
			if column != nodeversion.FieldStatusReason {
				columns = append(columns, column)
			}
		}

		query.Select(columns...)
	}

	versions, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list node versions: %w", err)
	}

	// Calculate total pages
	totalPages := 0
	if filter.PageSize > 0 {
		totalPages = (total + filter.PageSize - 1) / filter.PageSize // Use ceiling division for total pages
	}

	return &entity.ListNodeVersionsResult{
		Total:        total,
		NodeVersions: versions,
		Page:         filter.Page,
		Limit:        filter.PageSize,
		TotalPages:   totalPages,
	}, nil
}

func (s *RegistryService) AddNodeReview(ctx context.Context, client *ent.Client, nodeId, userID string, star int) (n *ent.Node, err error) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.AddNodeReview")
		defer segment.End()
	}
	log.Ctx(ctx).Info().Msgf("add review to node: %v ", nodeId)

	err = db.WithTx(ctx, client, func(tx *ent.Tx) error {
		v, err := s.GetNode(ctx, tx.Client(), nodeId)
		if err != nil {
			return fmt.Errorf("fail to fetch node version")
		}

		err = tx.NodeReview.Create().
			SetNode(v).
			SetUserID(userID).
			SetStar(star).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("fail to add review to node ")
		}

		err = v.Update().AddTotalReview(1).AddTotalStar(int64(star)).Exec(ctx)
		if err != nil {
			return fmt.Errorf("fail to add review: %w", err)
		}

		n, err = s.indexNodeWithLatestVersion(ctx, tx.Client(), nodeId)
		if err != nil {
			return fmt.Errorf("failed to index node: %w", err)
		}
		n.Edges.Versions = nil
		return nil
	})

	return
}

func (s *RegistryService) GetNodeVersionByVersion(ctx context.Context, client *ent.Client, nodeId, nodeVersion string) (*ent.NodeVersion, error) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.GetNodeVersionByVersion")
		defer segment.End()
	}
	log.Ctx(ctx).Info().Msgf("getting node version %v@%v", nodeId, nodeVersion)
	return client.NodeVersion.
		Query().
		Where(nodeversion.VersionEQ(nodeVersion)).
		Where(nodeversion.NodeIDEQ(nodeId)).
		WithStorageFile().
		Only(ctx)
}

func (s *RegistryService) GetNodeVersion(ctx context.Context, client *ent.Client, nodeVersionId string) (*ent.NodeVersion, error) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.GetNodeVersion")
		defer segment.End()
	}
	log.Ctx(ctx).Info().Msgf("getting node version %v", nodeVersionId)
	return client.NodeVersion.
		Get(ctx, uuid.MustParse(nodeVersionId))
}

func (s *RegistryService) UpdateNodeVersion(ctx context.Context, client *ent.Client, update *ent.NodeVersionUpdateOne) (*ent.NodeVersion, error) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.UpdateNodeVersion")
		defer segment.End()
	}
	log.Ctx(ctx).Info().Msgf("updating node version fields: %v", update.Mutation().Fields())
	return db.WithTxResult(ctx, client, func(tx *ent.Tx) (*ent.NodeVersion, error) {
		node, err := update.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to update node version: %w", err)
		}

		err = s.algolia.IndexNodeVersions(ctx, node)
		if err != nil {
			return nil, fmt.Errorf("failed to index node version: %w", err)
		}

		return node, nil
	})
}

func (s *RegistryService) RecordNodeInstallation(ctx context.Context, client *ent.Client, node *ent.Node) (*ent.Node, error) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.RecordNodeInstallation")
		defer segment.End()
	}
	var n *ent.Node
	err := db.WithTx(ctx, client, func(tx *ent.Tx) (err error) {
		n, err = tx.Node.UpdateOne(node).AddTotalInstall(1).Save(ctx)
		if err != nil {
			return err
		}

		// _, err = s.indexNodeWithLatestVersion(ctx, tx.Client(), n.ID)
		// if err != nil {
		// 	return fmt.Errorf("failed to index node: %w", err)
		// }
		return
	})
	return n, err
}

func (s *RegistryService) GetLatestNodeVersion(ctx context.Context, client *ent.Client, nodeId string) (*ent.NodeVersion, error) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.GetLatestNodeVersion")
		defer segment.End()
	}
	log.Ctx(ctx).Info().Msgf("Getting latest version of node: %v", nodeId)
	nodeVersion, err := client.NodeVersion.
		Query().
		Where(nodeversion.NodeIDEQ(nodeId)).
		Where(nodeversion.StatusIn(
			schema.NodeVersionStatusActive,
			schema.NodeVersionStatusFlagged,
			schema.NodeVersionStatusPending,
		)).
		Order(ent.Desc(nodeversion.FieldVersion)).
		WithStorageFile().
		First(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			log.Ctx(ctx).Info().Msgf("No versions found for node %v", nodeId)
			return nil, nil
		}

		log.Ctx(ctx).Error().Msgf("Error fetching latest version for node %v: %v", nodeId, err)
		return nil, err
	}

	log.Ctx(ctx).Info().Msgf("Found latest version for node %v: %v", nodeId, nodeVersion)
	return nodeVersion, nil
}

var ErrComfyNodesAlreadyExist = errors.New("comfy nodes already exist")

func (s *RegistryService) MarkComfyNodeExtractionFailed(
	ctx context.Context,
	client *ent.Client,
	nodeID string,
	nodeVersion string,
	info *schema.ComfyNodeCloudBuildInfo,
) error {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.MarkComfyNodeExtractionFailed")
		defer segment.End()
	}
	u := client.NodeVersion.
		Update().
		Where(
			nodeversion.NodeIDEQ(nodeID),
			nodeversion.VersionEQ(nodeVersion),
		).
		SetComfyNodeExtractStatus(schema.ComfyNodeExtractStatusFailed)
	if info != nil {
		u = u.SetComfyNodeCloudBuildInfo(*info)
	}
	return u.Exec(ctx)
}

func (s *RegistryService) CreateComfyNodes(
	ctx context.Context,
	client *ent.Client,
	nodeID, nodeVersion string,
	comfyNodes map[string]drip.ComfyNode,
	info *schema.ComfyNodeCloudBuildInfo,
) error {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.CreateComfyNodes")
		defer segment.End()
	}
	return db.WithTx(ctx, client, func(tx *ent.Tx) error {
		// Query the NodeVersion with the given nodeID and nodeVersion, lock it for updates
		nv, err := tx.NodeVersion.Query().
			Where(nodeversion.VersionEQ(nodeVersion)).
			Where(nodeversion.NodeIDEQ(nodeID)).
			ForUpdate().
			Only(ctx)
		if err != nil {
			return err
		}

		// Return an error if comfy nodes already exist for this NodeVersion
		if len(nv.Edges.ComfyNodes) > 0 {
			return ErrComfyNodesAlreadyExist
		}

		// Prepare a slice for bulk creation of comfy nodes
		comfyNodesCreates := make([]*ent.ComfyNodeCreate, 0, len(comfyNodes))
		for name, node := range comfyNodes {
			// Initialize a new ComfyNodeCreate instance for each comfy node
			comfyNodeCreate := tx.ComfyNode.Create().
				SetName(name).
				SetNodeVersionID(nv.ID)

			// Set optional fields if they are provided
			if node.Category != nil {
				comfyNodeCreate.SetCategory(*node.Category)
			}
			if node.Description != nil {
				comfyNodeCreate.SetDescription(*node.Description)
			}
			if node.InputTypes != nil {
				comfyNodeCreate.SetInputTypes(*node.InputTypes)
			}
			if node.Deprecated != nil {
				comfyNodeCreate.SetDeprecated(*node.Deprecated)
			}
			if node.Experimental != nil {
				comfyNodeCreate.SetExperimental(*node.Experimental)
			}
			if node.OutputIsList != nil {
				comfyNodeCreate.SetOutputIsList(*node.OutputIsList)
			}
			if node.ReturnNames != nil {
				comfyNodeCreate.SetReturnNames(*node.ReturnNames)
			}
			if node.ReturnTypes != nil {
				comfyNodeCreate.SetReturnTypes(*node.ReturnTypes)
			}
			if node.Function != nil {
				comfyNodeCreate.SetFunction(*node.Function)
			}

			// Append the ComfyNodeCreate to the slice
			comfyNodesCreates = append(comfyNodesCreates, comfyNodeCreate)
		}

		// Execute the bulk creation of comfy nodes
		if err := tx.ComfyNode.CreateBulk(comfyNodesCreates...).Exec(ctx); err != nil {
			return fmt.Errorf("failed to create comfy nodes: %w", err)
		}

		// Update the comfy node extraction status to success
		u := nv.Update().
			SetComfyNodeExtractStatus(schema.ComfyNodeExtractStatusSuccess)
		if info != nil {
			u = u.SetComfyNodeCloudBuildInfo(*info)
		}
		if err := u.Exec(ctx); err != nil {
			return fmt.Errorf("failed to update comfy node extraction status: %w", err)
		}

		// Re-index the node with its latest version
		if _, err := s.indexNodeWithLatestVersion(ctx, tx.Client(), nodeID); err != nil {
			return fmt.Errorf("failed to update node index: %w", err)
		}

		return nil
	})
}

func (s *RegistryService) GetComfyNode(
	ctx context.Context,
	client *ent.Client,
	nodeID, nodeVersion, comfyNodeName string,
) (*ent.ComfyNode, error) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.GetComfyNode")
		defer segment.End()
	}
	// Query the NodeVersion with the given nodeID and nodeVersion, ensuring extraction status is success
	nv, err := client.NodeVersion.Query().
		Where(nodeversion.VersionEQ(nodeVersion)).
		Where(nodeversion.NodeIDEQ(nodeID)).
		Where(nodeversion.ComfyNodeExtractStatusEQ(schema.ComfyNodeExtractStatusSuccess)).
		WithComfyNodes(func(cnq *ent.ComfyNodeQuery) {
			// Filter to find the specific comfy node by name
			cnq.Where(comfynode.NameEQ(comfyNodeName))
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Ensure at least one comfy node is found
	if len(nv.Edges.ComfyNodes) == 0 {
		return nil, fmt.Errorf("comfy node not found")
	}

	// Return the first comfy node (should be unique per query)
	return nv.Edges.ComfyNodes[0], nil
}

func (s *RegistryService) TriggerComfyNodesBackfill(
	ctx context.Context, client *ent.Client, max *int) error {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.TriggerComfyNodesBackfill")
		defer segment.End()
	}
	// Query all NodeVersions with pending comfy node extraction status
	q := client.NodeVersion.
		Query().
		WithStorageFile().
		Where(nodeversion.ComfyNodeExtractStatusEQ(schema.ComfyNodeExtractStatusPending))
	if max != nil {
		// Apply a limit if specified
		q.Limit(*max)
	}

	// Fetch the filtered NodeVersions
	nvs, err := q.All(ctx)
	if err != nil {
		return fmt.Errorf("failed to query node versions: %w", err)
	}

	// Iterate through each NodeVersion and trigger backfill
	for i, nv := range nvs {
		// Skip if the associated storage file does not have a valid URL
		if nv.Edges.StorageFile.FileURL == "" {
			continue
		}

		// Log the backfilling process
		log.Ctx(ctx).Info().Msgf("backfilling comfy node: %s", nv.Edges.StorageFile.FileURL)

		// Publish the node pack for backfill
		if err := s.pubsubService.PublishNodePack(ctx, nv.Edges.StorageFile.FileURL); err != nil {
			return fmt.Errorf(
				"fail to trigger node pack backfill for node %s-%s at index %d", nv.NodeID, nv.Version, i)
		}
	}

	return nil
}

func (s *RegistryService) AssertPublisherPermissions(ctx context.Context,
	client *ent.Client,
	publisherID string,
	userID string,
	permissions []schema.PublisherPermissionType,
) (err error) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.AssertPublisherPermissions")
		defer segment.End()
	}
	w, err := client.Publisher.Get(ctx, publisherID)
	if err != nil {
		return fmt.Errorf("fail to query publisher by id: %s %w", publisherID, err)
	}
	wp, err := w.QueryPublisherPermissions().
		Where(
			publisherpermission.PermissionIn(permissions...),
			publisherpermission.UserIDEQ(userID),
		).
		Count(ctx)
	if err != nil {
		return fmt.Errorf("fail to query publisher permission :%w", err)
	}
	if wp < 1 {
		return newErrorPermission("user '%s' doesn't have required permission on publisher '%s' ", userID, publisherID)
	}
	return
}

func (s *RegistryService) IsPersonalAccessTokenValidForPublisher(ctx context.Context,
	client *ent.Client,
	publisherID string,
	accessToken string,
) (bool, error) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.IsPersonalAccessTokenValidForPublisher")
		defer segment.End()
	}
	w, err := client.Publisher.Get(ctx, publisherID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("fail to find publisher by id: %s", publisherID)
		return false, fmt.Errorf("fail to find publisher by id: %s", publisherID)
	}
	exists, err := w.QueryPersonalAccessTokens().
		Where(
			personalaccesstoken.And(
				personalaccesstoken.PublisherIDEQ(publisherID),
				personalaccesstoken.TokenEQ(accessToken),
			),
		).
		Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("fail to query publisher permission :%w", err)
	}
	return exists, nil
}

func (s *RegistryService) AssertNodeBelongsToPublisher(ctx context.Context, client *ent.Client, publisherID string, nodeID string) error {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.AssertNodeBelongsToPublisher")
		defer segment.End()
	}
	node, err := client.Node.Get(ctx, nodeID)
	if err != nil {
		return fmt.Errorf("failed to get node: %w", err)
	}
	if node.PublisherID != publisherID {
		return newErrorPermission("node %s does not belong to publisher %s", nodeID, publisherID)
	}
	return nil
}

func (s *RegistryService) AssertAccessTokenBelongsToPublisher(ctx context.Context, client *ent.Client, publisherID string, tokenId uuid.UUID) error {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.AssertAccessTokenBelongsToPublisher")
		defer segment.End()
	}
	pat, err := client.PersonalAccessToken.Query().Where(
		personalaccesstoken.IDEQ(tokenId),
		personalaccesstoken.PublisherIDEQ(publisherID),
	).Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to get personal access token: %w", err)
	}
	if pat.PublisherID != publisherID {
		return newErrorPermission("personal access token %s does not belong to publisher %s", tokenId, publisherID)
	}
	return nil
}

func (s *RegistryService) DeletePublisher(ctx context.Context, client *ent.Client, publisherID string) error {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.DeletePublisher")
		defer segment.End()
	}
	log.Ctx(ctx).Info().Msgf("deleting publisher: %v", publisherID)
	return db.WithTx(ctx, client, func(tx *ent.Tx) error {
		client = tx.Client()

		_, err := client.PublisherPermission.
			Delete().
			Where(publisherpermission.PublisherIDEQ(publisherID)).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete publisher permissions: %w", err)
		}

		_, err = client.PersonalAccessToken.Delete().
			Where(personalaccesstoken.
				PublisherIDEQ(publisherID)).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete publisher access token: %w", err)
		}

		_, err = client.Publisher.
			Delete().
			Where(publisher.IDEQ(publisherID)).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete publisher: %w", err)
		}
		return nil
	})

}

func (s *RegistryService) DeleteNode(ctx context.Context, client *ent.Client, nodeID string) error {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.DeleteNode")
		defer segment.End()
	}
	log.Ctx(ctx).Info().Msgf("deleting node: %v", nodeID)
	db.WithTx(ctx, client, func(tx *ent.Tx) error {
		nv, err := tx.Client().NodeVersion.Query().Where(nodeversion.NodeID(nodeID)).All(ctx)
		if err != nil {
			return fmt.Errorf("fail to fetch node version for algolia deletion: %w", err)
		}

		err = tx.Client().Node.DeleteOneID(nodeID).Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete node: %w", err)
		}

		if err = s.algolia.DeleteNode(ctx, &ent.Node{ID: nodeID}); err != nil {
			return fmt.Errorf("fail to delete node from algolia: %w", err)
		}

		if err = s.algolia.DeleteNodeVersions(ctx, nv...); err != nil {
			return fmt.Errorf("fail to delete node version from algolia: %w", err)
		}

		return nil
	})
	return nil
}

func (s *RegistryService) DeleteNodeVersion(ctx context.Context, client *ent.Client, nodeIDVersion string) error {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.DeleteNodeVersion")
		defer segment.End()
	}
	log.Ctx(ctx).Info().Msgf("deleting node version: %v", nodeIDVersion)
	db.WithTx(ctx, client, func(tx *ent.Tx) error {
		nv, err := tx.Client().NodeVersion.Get(ctx, uuid.MustParse(nodeIDVersion))
		if err != nil {
			return fmt.Errorf("fail to fetch node version while deleting node version: %w", err)
		}

		err = tx.Client().NodeVersion.UpdateOneID(nv.ID).
			SetStatus(schema.NodeVersionStatusDeleted).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to update node version status: %w", err)
		}

		if err = s.algolia.DeleteNodeVersions(ctx, nv); err != nil {
			return fmt.Errorf("fail to delete node version from algolia: %w", err)
		}

		return nil
	})
	return nil
}

type errorPermission string

// Error implements error.
func (p errorPermission) Error() string {
	return string(p)
}

func newErrorPermission(tmpl string, args ...interface{}) errorPermission {
	return errorPermission(fmt.Sprintf(tmpl, args...))
}

var _ error = errorPermission("")

func IsPermissionError(err error) bool {
	if err == nil {
		return false
	}
	var e errorPermission
	return errors.As(err, &e)
}

func (s *RegistryService) BanPublisher(ctx context.Context, client *ent.Client, id string) error {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.BanPublisher")
		defer segment.End()
	}
	log.Ctx(ctx).Info().Msgf("banning publisher: %v", id)
	pub, err := client.Publisher.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("fail to find publisher: %w", err)
	}

	err = db.WithTx(ctx, client, func(tx *ent.Tx) error {
		err = pub.Update().SetStatus(schema.PublisherStatusTypeBanned).Exec(ctx)
		if err != nil {
			return fmt.Errorf("fail to update publisher: %w", err)
		}

		err = tx.User.Update().
			Where(user.HasPublisherPermissionsWith(publisherpermission.HasPublisherWith(publisher.IDEQ(pub.ID)))).
			SetStatus(schema.UserStatusTypeBanned).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("fail to update users: %w", err)
		}

		err = tx.Node.Update().
			Where(node.PublisherIDEQ(pub.ID)).
			SetStatus(schema.NodeStatusBanned).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("fail to update users: %w", err)
		}

		nodes, err := s.decorateNodeQueryWithLatestVersion(
			tx.Node.Query().
				Where(node.PublisherID(id)),
		).All(ctx)
		if len(nodes) == 0 || ent.IsNotFound(err) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("fail to update nodes: %w", err)
		}

		err = s.algolia.IndexNodes(ctx, nodes...)
		if err != nil {
			return fmt.Errorf("failed to index node: %w", err)
		}

		return nil
	})

	return err
}

func (s *RegistryService) BanNode(ctx context.Context, client *ent.Client, publisherid, id string) error {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.BanNode")
		defer segment.End()
	}
	log.Ctx(ctx).Info().Msgf("banning publisher node: %v %v", publisherid, id)

	return db.WithTx(ctx, client, func(tx *ent.Tx) error {
		n, err := s.decorateNodeQueryWithLatestVersion(
			tx.Node.Query().
				Where(node.And(
					node.IDEQ(id),
					node.PublisherIDEQ(publisherid),
				))).
			Only(ctx)
		if ent.IsNotFound(err) {
			return nil
		}

		n, err = n.Update().
			SetStatus(schema.NodeStatusBanned).
			Save(ctx)
		if ent.IsNotFound(err) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("fail to ban node: %w", err)
		}

		err = s.algolia.IndexNodes(ctx, n)
		if err != nil {
			return fmt.Errorf("failed to index node: %w", err)
		}

		return err
	})

}

func (s *RegistryService) AssertNodeBanned(ctx context.Context, client *ent.Client, nodeID string) error {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.AssertNodeBanned")
		defer segment.End()
	}
	node, err := client.Node.Get(ctx, nodeID)
	if ent.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get node: %w", err)
	}
	if node.Status == schema.NodeStatusBanned {
		return newErrorPermission("node '%s' is currently banned", nodeID)
	}
	return nil
}

func (s *RegistryService) AssertPublisherBanned(ctx context.Context, client *ent.Client, publisherID string) error {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.AssertPublisherBanned")
		defer segment.End()
	}
	publisher, err := client.Publisher.Get(ctx, publisherID)
	if ent.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get node: %w", err)
	}
	if publisher.Status == schema.PublisherStatusTypeBanned {
		return newErrorPermission("node '%s' is currently banned", publisherID)
	}
	return nil
}

func (s *RegistryService) ReindexAllNodes(ctx context.Context, client *ent.Client) error {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.ReindexAllNodes")
		defer segment.End()
	}
	log.Ctx(ctx).Info().Msgf("reindexing nodes")
	nodes, err := s.decorateNodeQueryWithLatestVersion(client.Node.Query()).All(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch all nodes: %w", err)
	}

	log.Ctx(ctx).Info().Msgf("reindexing %d number of nodes", len(nodes))
	err = s.algolia.IndexNodes(ctx, nodes...)
	if err != nil {
		return fmt.Errorf("failed to reindex all nodes: %w", err)
	}

	var nvs []*ent.NodeVersion
	for _, n := range nodes {
		nvs = append(nvs, n.Edges.Versions...)
	}

	log.Ctx(ctx).Info().Msgf("reindexing %d number of node versions", len(nvs))
	err = s.algolia.IndexNodeVersions(ctx, nvs...)
	if err != nil {
		return fmt.Errorf("failed to reindex all node versions: %w", err)
	}

	return nil
}

var reindexLock = sync.Mutex{}

func (s *RegistryService) ReindexAllNodesBackground(ctx context.Context, client *ent.Client) (err error) {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.ReindexAllNodesBackground")
		defer segment.End()
	}
	if !reindexLock.TryLock() {
		return fmt.Errorf("another reindex is in progress")
	}
	defer reindexLock.Unlock()

	go func() {
		err = s.ReindexAllNodes(ctx, client)
		if err != nil {
			log.Ctx(ctx).Err(err).Msg("failed to reindex all nodes in background")
		}
		log.Ctx(ctx).Info().Msg("reindexing nodes in background succesful")
	}()

	return nil
}

// indexNodeWithLatestVersion re-indexes a single node and its latest version
func (s *RegistryService) indexNodeWithLatestVersion(
	ctx context.Context,
	client *ent.Client,
	nodeID string) (*ent.Node, error) {
	n, err := s.decorateNodeQueryWithLatestVersion(
		client.Node.Query().
			Where(node.IDEQ(nodeID)),
	).Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query node: %w", err)
	}
	if err := s.algolia.IndexNodes(ctx, n); err != nil {
		return nil, fmt.Errorf("failed to update node: %w", err)
	}
	return n, nil
}

func (s *RegistryService) decorateNodeQueryWithLatestVersion(q *ent.NodeQuery) *ent.NodeQuery {
	return q.WithVersions(func(q *ent.NodeVersionQuery) {
		q.Modify(func(s *sql.Selector) {
			s.Where(sql.ExprP(
				`(node_id, create_time) IN (
					SELECT node_id, MAX(create_time)
					FROM node_versions
					GROUP BY node_id
				)`,
			))
		})
	})
}

func (s *RegistryService) PerformSecurityCheck(
	ctx context.Context, client *ent.Client, nodeVersion *ent.NodeVersion) error {
	if txn := newrelic.FromContext(ctx); txn != nil {
		segment := txn.StartSegment("RegistryService.PerformSecurityCheck")
		defer segment.End()
	}
	log.Ctx(ctx).Info().Msgf("Scanning node %s@%s w/ version ID: %s",
		nodeVersion.NodeID, nodeVersion.Version, nodeVersion.ID)

	if (nodeVersion.Edges.StorageFile == nil) || (nodeVersion.Edges.StorageFile.FileURL == "") {
		return fmt.Errorf("node version %s@%s does not have a storage file", nodeVersion.NodeID, nodeVersion.Version)
	}

	issues, err := sendScanRequest(s.config.SecretScannerURL, nodeVersion.Edges.StorageFile.FileURL)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Ctx(ctx).Info().Msgf(
				"Node zip file doesn't exist %s@%s. Updating to deleted.", nodeVersion.NodeID, nodeVersion.Version)
			err := nodeVersion.Update().
				SetStatus(schema.NodeVersionStatusDeleted).
				SetStatusReason("Node zip file doesn't exist").Exec(ctx)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msgf("failed to update node version status to active")
			}
		}
		return err
	}

	if issues == "" {
		log.Ctx(ctx).Info().Msgf(
			"No security issues found in node %s@%s. Updating to active.", nodeVersion.NodeID, nodeVersion.Version)
		err := nodeVersion.Update().
			SetStatus(schema.NodeVersionStatusActive).
			SetStatusReason("Passed automated checks").Exec(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msgf("failed to update node version status to active")
		}
		// Send the message to the private channel
		err = s.discordService.SendSecurityCouncilMessage(
			fmt.Sprintf("Node %s@%s has passed automated scans. Changing status to active.",
				nodeVersion.NodeID, nodeVersion.Version), true)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msgf("failed to send message to discord")
		}
	} else {
		log.Ctx(ctx).Info().Msgf(
			"Security issues found in node %s@%s. Updating to flagged.", nodeVersion.NodeID, nodeVersion.Version)
		log.Ctx(ctx).Info().Msgf(
			"List of security issues %s.", issues) // 500 character max.
		prettyIssues, err := common.PrettifyJSON(issues)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to prettify JSON issues")
			prettyIssues = issues // fallback to unprettified issues
		}
		err = s.discordService.SendSecurityCouncilMessage(
			fmt.Sprintf("Security issues were found in node %s@%s. Status is flagged. "+
				"Please check it here: https://registry.comfy.org/nodes/%s/versions/%s. \n "+
				"Issues are: \n%s", nodeVersion.NodeID, nodeVersion.Version, nodeVersion.NodeID, nodeVersion.Version,
				prettyIssues), false)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msgf("failed to send message to discord")
		}
	}

	return nil
}

type ScanRequest struct {
	URL string `json:"url"`
}

func sendScanRequest(apiURL, fileURL string) (string, error) {
	requestBody, err := json.Marshal(ScanRequest{URL: fileURL})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Println("Response Status:", resp.Status)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to scan file: %s", responseBody)
	}

	return string(responseBody), nil
}
