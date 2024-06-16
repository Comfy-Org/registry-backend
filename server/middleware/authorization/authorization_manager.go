package drip_authorization

import (
	"github.com/labstack/echo/v4"
	strictecho "github.com/oapi-codegen/runtime/strictmiddleware/echo"
	"github.com/rs/zerolog/log"
	"net/http"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/schema"
	drip_authentication "registry-backend/server/middleware/authentication"
	drip_services "registry-backend/services/registry"
)

// AuthorizationManager manages authorization-related tasks
type AuthorizationManager struct {
	EntClient       *ent.Client
	RegistryService *drip_services.RegistryService
}

// NewAuthorizationManager creates a new instance of AuthorizationManager
func NewAuthorizationManager(
	entClient *ent.Client, registryService *drip_services.RegistryService) *AuthorizationManager {
	return &AuthorizationManager{
		EntClient:       entClient,
		RegistryService: registryService,
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
			err = m.RegistryService.AssertPublisherPermissions(ctx, m.EntClient, publisherID, userDetails.ID, permissions)
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
			err = m.RegistryService.AssertNodeBanned(ctx, m.EntClient, nodeID)
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

			pub, _ := m.RegistryService.GetPublisher(ctx, m.EntClient, publisherID)
			if pub != nil && pub.Status == schema.PublisherStatusTypeBanned {
				log.Ctx(ctx).Error().Msgf("Publisher %s banned", publisherID)
				return nil, echo.NewHTTPError(http.StatusForbidden, "Publisher Banned")
			}

			return f(c, request)
		}
	}
}
