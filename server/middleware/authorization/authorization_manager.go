package drip_authorization

import (
	"context"
	"net/http"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/schema"
	drip_authentication "registry-backend/server/middleware/authentication"
	drip_services "registry-backend/services/registry"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	strictecho "github.com/oapi-codegen/runtime/strictmiddleware/echo"
	"github.com/rs/zerolog/log"
)

type Assertor interface {
	AssertPublisherBanned(ctx context.Context, client *ent.Client, publisherID string) error
	AssertPublisherPermissions(ctx context.Context,
		client *ent.Client,
		publisherID string,
		userID string,
		permissions []schema.PublisherPermissionType) (err error)
	IsPersonalAccessTokenValidForPublisher(ctx context.Context,
		client *ent.Client,
		publisherID string,
		accessToken string,
	) (bool, error)
	AssertNodeBelongsToPublisher(ctx context.Context, client *ent.Client, publisherID string, nodeID string) error
	AssertAccessTokenBelongsToPublisher(ctx context.Context, client *ent.Client, publisherID string, tokenId uuid.UUID) error
	AssertNodeBanned(ctx context.Context, client *ent.Client, nodeID string) error
}

// AuthorizationManager manages authorization-related tasks
type AuthorizationManager struct {
	EntClient *ent.Client
	Assertor  Assertor
}

// NewAuthorizationManager creates a new instance of AuthorizationManager
func NewAuthorizationManager(
	entClient *ent.Client, assertor Assertor) *AuthorizationManager {
	return &AuthorizationManager{
		EntClient: entClient,
		Assertor:  assertor,
	}
}

// assertUserBanned checks if the user is banned
func (m *AuthorizationManager) assertUserBanned() drip.StrictMiddlewareFunc {
	return func(f strictecho.StrictEchoHandlerFunc, operationID string) strictecho.StrictEchoHandlerFunc {
		return func(c echo.Context, request interface{}) (response interface{}, err error) {
			ctx := c.Request().Context()
			v := ctx.Value(drip_authentication.UserContextKey)
			userDetails, ok := v.(*drip_authentication.UserDetails)
			if !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "user not found")
			}

			u, err := m.EntClient.User.Get(ctx, userDetails.ID)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "user not found")
			}

			if u.Status == schema.UserStatusTypeBanned {
				return nil, echo.NewHTTPError(http.StatusForbidden, "user/publisher is banned")
			}

			return f(c, request)
		}
	}
}

// assertPublisherPermission checks if the user has the required permissions for the publisher
func (m *AuthorizationManager) assertPublisherPermission(
	permissions []schema.PublisherPermissionType, extractor func(req interface{}) string) drip.StrictMiddlewareFunc {
	return func(f strictecho.StrictEchoHandlerFunc, operationID string) strictecho.StrictEchoHandlerFunc {
		return func(c echo.Context, request interface{}) (response interface{}, err error) {
			ctx := c.Request().Context()
			v := ctx.Value(drip_authentication.UserContextKey)
			userDetails, ok := v.(*drip_authentication.UserDetails)
			if !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "user not found")
			}
			publisherID := extractor(request)

			log.Ctx(ctx).Info().Msgf("Checking if user ID %s has permission "+
				"to update publisher ID %s", userDetails.ID, publisherID)
			err = m.Assertor.AssertPublisherPermissions(ctx, m.EntClient, publisherID, userDetails.ID, permissions)
			switch {
			case ent.IsNotFound(err):
				log.Ctx(ctx).Info().Msgf("Publisher with ID %s not found", publisherID)
				return nil, echo.NewHTTPError(http.StatusNotFound, "Publisher Not Found")

			case drip_services.IsPermissionError(err):
				log.Ctx(ctx).Error().Msgf("Permission denied for user ID %s on "+
					"publisher ID %s w/ err: %v", userDetails.ID, publisherID, err)
				return nil, echo.NewHTTPError(http.StatusForbidden, "Permission denied")

			case err != nil:
				log.Ctx(ctx).Error().Msgf("Failed to assert publisher "+
					"permission %s w/ err: %v", publisherID, err)
				return nil, err
			}

			return f(c, request)
		}
	}
}

// assertNodeBanned checks if the node is banned
func (m *AuthorizationManager) assertNodeBanned(extractor func(req interface{}) string) drip.StrictMiddlewareFunc {
	return func(f strictecho.StrictEchoHandlerFunc, operationID string) strictecho.StrictEchoHandlerFunc {
		return func(c echo.Context, request interface{}) (response interface{}, err error) {
			ctx := c.Request().Context()
			nodeID := extractor(request)
			err = m.Assertor.AssertNodeBanned(ctx, m.EntClient, nodeID)
			switch {
			case drip_services.IsPermissionError(err):
				log.Ctx(ctx).Error().Msgf("Node %s banned", nodeID)
				return nil, echo.NewHTTPError(http.StatusForbidden, "Node Banned")

			case err != nil:
				log.Ctx(ctx).Error().Msgf("Failed to assert node ban status %s w/ err: %v", nodeID, err)
				return nil, err
			}

			return f(c, request)
		}
	}
}

// assertPublisherBanned checks if the publisher is banned
func (m *AuthorizationManager) assertPublisherBanned(extractor func(req interface{}) string) drip.StrictMiddlewareFunc {
	return func(f strictecho.StrictEchoHandlerFunc, operationID string) strictecho.StrictEchoHandlerFunc {
		return func(c echo.Context, request interface{}) (response interface{}, err error) {
			ctx := c.Request().Context()
			publisherID := extractor(request)

			switch err = m.Assertor.AssertPublisherBanned(ctx, m.EntClient, publisherID); {
			case drip_services.IsPermissionError(err):
				log.Ctx(ctx).Error().Msgf("Publisher %s banned", publisherID)
				return nil, echo.NewHTTPError(http.StatusForbidden, "Node Banned")

			case err != nil:
				log.Ctx(ctx).Error().Msgf("Failed to assert publisher ban status %s w/ err: %v", publisherID, err)
				return nil, err
			}

			return f(c, request)
		}
	}
}

// assertPersonalAccessTokenValid check if personal access token is valid for a publisher
func (m *AuthorizationManager) assertPersonalAccessTokenValid(
	extractorPublsherID func(req interface{}) (nodeid string),
	extractorPAT func(req interface{}) (pat string),
) drip.StrictMiddlewareFunc {
	return func(f strictecho.StrictEchoHandlerFunc, operationID string) strictecho.StrictEchoHandlerFunc {
		return func(c echo.Context, request interface{}) (response interface{}, err error) {
			ctx := c.Request().Context()
			pubID := extractorPublsherID(request)
			pat := extractorPAT(request)
			tokenValid, err := m.Assertor.IsPersonalAccessTokenValidForPublisher(
				ctx, m.EntClient, pubID, pat)
			if err != nil {
				log.Ctx(ctx).Error().Msgf("Token validation failed w/ err: %v", err)
				return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to validate token")
			}
			if !tokenValid {
				log.Ctx(ctx).Error().Msg("Invalid personal access token")
				return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid personal access token")
			}

			return f(c, request)
		}
	}
}

// assertNodeBelongsToPublisher check if a node belongs to a publisher
func (m *AuthorizationManager) assertNodeBelongsToPublisher(
	extractorPublsherID func(req interface{}) (nodeid string),
	extractorNodeID func(req interface{}) (nodeid string),
) drip.StrictMiddlewareFunc {
	return func(f strictecho.StrictEchoHandlerFunc, operationID string) strictecho.StrictEchoHandlerFunc {
		return func(c echo.Context, request interface{}) (response interface{}, err error) {
			ctx := c.Request().Context()
			pubID := extractorPublsherID(request)
			nodeID := extractorNodeID(request)

			err = m.Assertor.AssertNodeBelongsToPublisher(ctx, m.EntClient, pubID, nodeID)
			switch {
			case ent.IsNotFound(err):
				return f(c, request)

			case drip_services.IsPermissionError(err):
				log.Ctx(ctx).Error().Msgf(
					"Permission denied for publisher ID %s on node ID %s w/ err: %v", pubID, nodeID, err)
				return nil, echo.NewHTTPError(http.StatusForbidden, "Permission denied")

			case err != nil:
				return nil, err
			}

			return f(c, request)
		}
	}
}
