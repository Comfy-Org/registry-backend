package drip_services

import (
	"context"
	"errors"
	"fmt"
	"registry-backend/db"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/node"
	"registry-backend/ent/nodeversion"
	"registry-backend/ent/personalaccesstoken"
	"registry-backend/ent/publisher"
	"registry-backend/ent/publisherpermission"
	"registry-backend/ent/schema"
	gateway "registry-backend/gateways/slack"
	"registry-backend/gateways/storage"
	"registry-backend/mapper"

	"github.com/Masterminds/semver/v3"
	"google.golang.org/protobuf/proto"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type RegistryService struct {
	storageService storage.StorageService
	slackService   gateway.SlackService
}

func NewRegistryService(storageSvc storage.StorageService, slackSvc gateway.SlackService) *RegistryService {
	return &RegistryService{
		storageService: storageSvc,
		slackService:   slackSvc,
	}
}

type PublisherFilter struct {
	UserID string
}

// NodeFilter holds optional parameters for filtering node results
type NodeFilter struct {
	PublisherID string
	// Add more filter fields here
}

type NodeData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PublisherID string `json:"publisherId"`
	// Add other fields as necessary
}

// ListNodesResult is the structure that holds the paginated result of nodes
type ListNodesResult struct {
	Total      int         `json:"total"`
	Nodes      []*ent.Node `json:"nodes"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"totalPages"`
}

// ListNodes retrieves a paginated list of nodes with optional filtering.
func (s *RegistryService) ListNodes(ctx context.Context, client *ent.Client, page, limit int, filter *NodeFilter) (*ListNodesResult, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	query := client.Node.Query().WithPublisher().WithVersions(
		func(q *ent.NodeVersionQuery) {
			q.Order(ent.Desc(nodeversion.FieldCreateTime))
		},
	)
	if filter != nil {
		if filter.PublisherID != "" {
			query.Where(node.PublisherID(filter.PublisherID))
		}
	}
	offset := (page - 1) * limit
	total, err := query.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count nodes: %w", err)
	}

	nodes, err := query.
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages += 1
	}

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

	createNode, err := mapper.ApiCreateNodeToDb(publisherId, node, client)
	log.Ctx(ctx).Info().Msgf("creating node with fields: %v", createNode.Mutation().Fields())
	if err != nil {
		return nil, fmt.Errorf("failed to map node: %w", err)
	}

	createdNode, err := createNode.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create node: %w", err)
	}

	return createdNode, nil
}

func (s *RegistryService) UpdateNode(ctx context.Context, client *ent.Client, update *ent.NodeUpdateOne) (*ent.Node, error) {
	log.Ctx(ctx).Info().Msgf("updating node fields: %v", update.Mutation().Fields())
	node, err := update.
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to update node: %w", err)
	}
	return node, nil
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

		slackErr := s.slackService.SendRegistryMessageToSlack(fmt.Sprintf("Version %s of node %s was published successfully. Publisher: %s. https://comfyregistry.org/nodes/%s", createdNodeVersion.Version, createdNodeVersion.NodeID, publisherID, nodeID))
		if slackErr != nil {
			log.Ctx(ctx).Error().Msgf("Failed to send message to Slack w/ err: %v", slackErr)
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

func (s *RegistryService) ListNodeVersions(ctx context.Context, client *ent.Client, nodeID string) ([]*ent.NodeVersion, error) {
	log.Ctx(ctx).Info().Msgf("listing node versions: %v", nodeID)
	versions, err := client.NodeVersion.Query().
		Where(nodeversion.NodeIDEQ(nodeID)).
		WithStorageFile().
		Order(ent.Desc(nodeversion.FieldCreateTime)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list node versions: %w", err)
	}
	return versions, nil
}

func (s *RegistryService) GetNodeVersion(ctx context.Context, client *ent.Client, nodeId, nodeVersion string) (*ent.NodeVersion, error) {
	log.Ctx(ctx).Info().Msgf("getting node version: %v", nodeVersion)
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

func (s *RegistryService) GetLatestNodeVersion(ctx context.Context, client *ent.Client, nodeId string) (*ent.NodeVersion, error) {
	log.Ctx(ctx).Info().Msgf("getting latest version of node: %v", nodeId)
	nodeVersion, err := client.NodeVersion.
		Query().
		Where(nodeversion.NodeIDEQ(nodeId)).
		Order(ent.Desc(nodeversion.FieldCreateTime)).
		WithStorageFile().
		First(ctx)

	if err != nil {
		if ent.IsNotFound(err) {

			return nil, nil
		}
		return nil, err
	}
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
	err := client.Node.DeleteOneID(nodeID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}
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
	var e *errorPermission
	return errors.As(err, &e)
}