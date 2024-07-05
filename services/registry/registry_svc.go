package dripservices

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"registry-backend/config"
	"registry-backend/db"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/node"
	"registry-backend/ent/nodeversion"
	"registry-backend/ent/personalaccesstoken"
	"registry-backend/ent/predicate"
	"registry-backend/ent/publisher"
	"registry-backend/ent/publisherpermission"
	"registry-backend/ent/schema"
	"registry-backend/ent/user"
	"registry-backend/gateways/algolia"
	"registry-backend/gateways/discord"
	gateway "registry-backend/gateways/slack"
	"registry-backend/gateways/storage"
	"registry-backend/mapper"
	drip_metric "registry-backend/server/middleware/metric"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/Masterminds/semver/v3"
	"google.golang.org/protobuf/proto"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type RegistryService struct {
	storageService storage.StorageService
	slackService   gateway.SlackService
	algolia        algolia.AlgoliaService
	discordService discord.DiscordService
	config         *config.Config
}

func NewRegistryService(storageSvc storage.StorageService, slackSvc gateway.SlackService, discordSvc discord.DiscordService, algoliaSvc algolia.AlgoliaService, config *config.Config) *RegistryService {
	return &RegistryService{
		storageService: storageSvc,
		slackService:   slackSvc,
		discordService: discordSvc,
		algolia:        algoliaSvc,
		config:         config,
	}
}

type PublisherFilter struct {
	UserID string
}

// NodeFilter holds optional parameters for filtering node results
type NodeFilter struct {
	PublisherID   string
	Search        string
	IncludeBanned bool
}

type NodeVersionFilter struct {
	NodeId   string
	Status   []schema.NodeVersionStatus
	MinAge   time.Duration
	PageSize int
	Page     int
}

type NodeData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PublisherID string `json:"publisherId"`
}

// ListNodesResult is the structure that holds the paginated result of nodes
type ListNodesResult struct {
	Total      int         `json:"total"`
	Nodes      []*ent.Node `json:"nodes"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"totalPages"`
}

type ListNodeVersionsResult struct {
	Total        int                `json:"total"`
	NodeVersions []*ent.NodeVersion `json:"nodes"`
	Page         int                `json:"page"`
	Limit        int                `json:"limit"`
	TotalPages   int                `json:"totalPages"`
}

// ListNodes retrieves a paginated list of nodes with optional filtering.
func (s *RegistryService) ListNodes(ctx context.Context, client *ent.Client, page, limit int, filter *NodeFilter) (*ListNodesResult, error) {
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

	// Count total nodes
	total, err := query.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count nodes: %w", err)
	}

	// Fetch nodes with pagination
	nodes, err := query.
		WithVersions(func(q *ent.NodeVersionQuery) {
			q.Modify(func(s *sql.Selector) {
				s.Where(sql.ExprP(
					`(node_id, create_time) IN (
						SELECT node_id, MAX(create_time)
						FROM node_versions
						GROUP BY node_id
					)`,
				))
			})
		}).
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	// Calculate total pages
	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	// Return the result
	return &ListNodesResult{
		Total:      total,
		Nodes:      nodes,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// ListPublishers queries the Publisher table with an optional user ID filter via PublisherPermission
func (s *RegistryService) ListPublishers(ctx context.Context, client *ent.Client, filter *PublisherFilter) ([]*ent.Publisher, error) {
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
	log.Ctx(ctx).Info().Msgf("updating publisher fields: %v", update.Mutation().Fields())
	publisher, err := update.Save(ctx)
	log.Ctx(ctx).Info().Msgf("success: updated publisher: %v", publisher)
	if err != nil {
		return nil, fmt.Errorf("failed to create publisher: %w", err)
	}

	return publisher, nil
}

func (s *RegistryService) GetPublisher(ctx context.Context, client *ent.Client, publisherID string) (*ent.Publisher, error) {
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
	pats, err := client.PersonalAccessToken.Query().
		Where(personalaccesstoken.PublisherIDEQ(publisherID)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list personal access tokens: %w", err)
	}
	return pats, nil
}

func (s *RegistryService) DeletePersonalAccessToken(ctx context.Context, client *ent.Client, tokenID uuid.UUID) error {
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
	var node *ent.Node
	err := db.WithTx(ctx, client, func(tx *ent.Tx) (err error) {
		update := updateFunc(tx.Client())
		log.Ctx(ctx).Info().Msgf("updating node fields: %v", update.Mutation().Fields())

		node, err = update.Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to update node: %w", err)
		}

		err = s.algolia.IndexNodes(ctx, node)
		if err != nil {
			return fmt.Errorf("failed to index node: %w", err)
		}

		return err
	})
	return node, err
}

func (s *RegistryService) GetNode(ctx context.Context, client *ent.Client, nodeID string) (*ent.Node, error) {
	log.Ctx(ctx).Info().Msgf("getting node: %v", nodeID)
	node, err := client.Node.Get(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}
	return node, nil
}

func (s *RegistryService) CreateNodeVersion(
	ctx context.Context,
	client *ent.Client,
	publisherID, nodeID string,
	nodeVersion *drip.NodeVersion) (*NodeVersionCreation, error) {
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
		objectPath := fmt.Sprintf("%s/%s/%s/%s", publisherID, nodeID, *nodeVersion.Version, "node.tar.gz")
		storageFile := tx.StorageFile.Create().
			SetBucketName(bucketName).
			SetFilePath(objectPath).
			SetFileType("zip").
			// Sample URL: https://storage.googleapis.com/comfy-registry/james-test-publisher/comfyui-inspire-pack/1.0.0/node.tar.gz
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

		message := fmt.Sprintf("Version %s of node %s was published successfully. Publisher: %s. https://registry.comfy.org/nodes/%s", createdNodeVersion.Version, createdNodeVersion.NodeID, publisherID, nodeID)
		slackErr := s.slackService.SendRegistryMessageToSlack(message)
		s.discordService.SendSecurityCouncilMessage(message)
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

func (s *RegistryService) ListNodeVersions(ctx context.Context, client *ent.Client, filter *NodeVersionFilter) (*ListNodeVersionsResult, error) {
	query := client.NodeVersion.Query().
		WithStorageFile().
		Order(ent.Desc(nodeversion.FieldCreateTime))

	if filter.NodeId != "" {
		log.Ctx(ctx).Info().Msgf("listing node versions: %v", filter.NodeId)
		query.Where(nodeversion.NodeIDEQ(filter.NodeId))
	}

	if filter.Status != nil && len(filter.Status) > 0 {
		log.Ctx(ctx).Info().Msgf("listing node versions with status: %v", filter.Status)
		query.Where(nodeversion.StatusIn(filter.Status...))
	}

	if filter.MinAge > 0 {
		query.Where(nodeversion.CreateTimeLT(time.Now().Add(-filter.MinAge)))
	}

	if filter.Page > 0 && filter.PageSize > 0 {
		query.Offset((filter.Page - 1) * filter.PageSize)
		query.Limit(filter.PageSize)
	}
	total, err := query.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count node versions: %w", err)
	}
	versions, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list node versions: %w", err)
	}

	totalPages := 0
	if total > 0 && filter.PageSize > 0 {
		totalPages = total / filter.PageSize

		if total%filter.PageSize != 0 {
			totalPages += 1
		}
	}
	return &ListNodeVersionsResult{
		Total:        total,
		NodeVersions: versions,
		Page:         filter.Page,
		Limit:        filter.PageSize,
		TotalPages:   totalPages,
	}, nil
}

func (s *RegistryService) AddNodeReview(ctx context.Context, client *ent.Client, nodeId, userID string, star int) (nv *ent.Node, err error) {
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

		nv, err = s.GetNode(ctx, tx.Client(), nodeId)
		if err != nil {
			return fmt.Errorf("fail to fetch node s")
		}

		err = s.algolia.IndexNodes(ctx, nv)
		if err != nil {
			return fmt.Errorf("failed to index node: %w", err)
		}

		return nil
	})

	return
}

func (s *RegistryService) GetNodeVersionByVersion(ctx context.Context, client *ent.Client, nodeId, nodeVersion string) (*ent.NodeVersion, error) {
	log.Ctx(ctx).Info().Msgf("getting node version %v@%v", nodeId, nodeVersion)
	return client.NodeVersion.
		Query().
		Where(nodeversion.VersionEQ(nodeVersion)).
		Where(nodeversion.NodeIDEQ(nodeId)).
		WithStorageFile().
		Only(ctx)
}

func (s *RegistryService) UpdateNodeVersion(ctx context.Context, client *ent.Client, update *ent.NodeVersionUpdateOne) (*ent.NodeVersion, error) {
	log.Ctx(ctx).Info().Msgf("updating node version fields: %v", update.Mutation().Fields())
	node, err := update.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update node version: %w", err)
	}
	return node, nil
}

func (s *RegistryService) RecordNodeInstalation(ctx context.Context, client *ent.Client, node *ent.Node) (*ent.Node, error) {
	var n *ent.Node
	err := db.WithTx(ctx, client, func(tx *ent.Tx) (err error) {
		node, err = tx.Node.UpdateOne(node).AddTotalInstall(1).Save(ctx)
		if err != nil {
			return err
		}
		err = s.algolia.IndexNodes(ctx, node)
		if err != nil {
			return fmt.Errorf("failed to index node: %w", err)
		}
		return
	})
	return n, err
}

func (s *RegistryService) GetLatestNodeVersion(ctx context.Context, client *ent.Client, nodeId string) (*ent.NodeVersion, error) {
	log.Ctx(ctx).Info().Msgf("Getting latest version of node: %v", nodeId)
	nodeVersion, err := client.NodeVersion.
		Query().
		Where(nodeversion.NodeIDEQ(nodeId)).
		//Where(nodeversion.StatusEQ(schema.NodeVersionStatusActive)).
		Order(ent.Desc(nodeversion.FieldCreateTime)).
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

func (s *RegistryService) AssertPublisherPermissions(ctx context.Context,
	client *ent.Client,
	publisherID string,
	userID string,
	permissions []schema.PublisherPermissionType,
) (err error) {
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
	log.Ctx(ctx).Info().Msgf("deleting node: %v", nodeID)
	db.WithTx(ctx, client, func(tx *ent.Tx) error {
		err := tx.Client().Node.DeleteOneID(nodeID).Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete node: %w", err)
		}

		if err = s.algolia.DeleteNode(ctx, &ent.Node{ID: nodeID}); err != nil {
			return fmt.Errorf("fail to delete node from algolia: %w", err)
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

		nodes, err := tx.Node.Query().Where(node.PublisherID(id)).All(ctx)
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
	log.Ctx(ctx).Info().Msgf("banning publisher node: %v %v", publisherid, id)

	return db.WithTx(ctx, client, func(tx *ent.Tx) error {
		n, err := tx.Node.Query().Where(node.And(
			node.IDEQ(id),
			node.PublisherIDEQ(publisherid),
		)).Only(ctx)
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
	log.Ctx(ctx).Info().Msgf("reindexing nodes")
	nodes, err := client.Node.Query().All(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch all nodes: %w", err)
	}

	log.Ctx(ctx).Info().Msgf("reindexing %d number of nodes", len(nodes))
	err = s.algolia.IndexNodes(ctx, nodes...)
	if err != nil {
		return fmt.Errorf("failed to reindex all nodes: %w", err)
	}
	return nil
}

func (s *RegistryService) PerformSecurityCheck(ctx context.Context, client *ent.Client, nodeVersion *ent.NodeVersion) error {
	log.Ctx(ctx).Info().Msgf("scanning node %s@%s", nodeVersion.NodeID, nodeVersion.Version)

	if (nodeVersion.Edges.StorageFile == nil) || (nodeVersion.Edges.StorageFile.FileURL == "") {
		return fmt.Errorf("node version %s@%s does not have a storage file", nodeVersion.NodeID, nodeVersion.Version)
	}

	issues, err := sendScanRequest(s.config.SecretScannerURL, nodeVersion.Edges.StorageFile.FileURL)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Ctx(ctx).Info().Msgf("Node zip file doesn’t exist %s@%s. Updating to deleted.", nodeVersion.NodeID, nodeVersion.Version)
			err := nodeVersion.Update().SetStatus(schema.NodeVersionStatusDeleted).SetStatusReason("Node zip file doesn’t exist").Exec(ctx)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msgf("failed to update node version status to active")
			}
		}
		return err
	}

	if issues != "" {
		log.Ctx(ctx).Info().Msgf("No security issues found in node %s@%s. Updating to active.", nodeVersion.NodeID, nodeVersion.Version)
		err := nodeVersion.Update().SetStatus(schema.NodeVersionStatusActive).SetStatusReason("Passed automated checks").Exec(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msgf("failed to update node version status to active")
		}
		err = s.discordService.SendSecurityCouncilMessage(fmt.Sprintf("Node %s@%s has passed automated scans. Changing status to active.", nodeVersion.NodeID, nodeVersion.Version))
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msgf("failed to send message to discord")
		}
	} else {
		log.Ctx(ctx).Info().Msgf("Security issues found in node %s@%s. Updating to flagged.", nodeVersion.NodeID, nodeVersion.Version)
		err := nodeVersion.Update().SetStatus(schema.NodeVersionStatusFlagged).SetStatusReason(issues).Exec(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msgf("failed to update node version status to security issue")
		}
		err = s.discordService.SendSecurityCouncilMessage(fmt.Sprintf("Security issues were found in node %s@%s. Status is flagged. Please check it here: https://registry.comfy.org/admin/nodes/%s/versions/%s", nodeVersion.NodeID, nodeVersion.Version, nodeVersion.NodeID, nodeVersion.Version))
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
