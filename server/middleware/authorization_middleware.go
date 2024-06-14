package drip_middleware

import (
	"net/http"
	"registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/ent/schema"

	"github.com/labstack/echo/v4"
	strictecho "github.com/oapi-codegen/runtime/strictmiddleware/echo"
)

func AuthorizationMiddleware(entClient *ent.Client) drip.StrictMiddlewareFunc {
	restrictedOperationsForBannedUsers := map[string]struct{}{
		"CreatePublisher":           {},
		"UpdatePublisher":           {},
		"CreateNode":                {},
		"DeleteNode":                {},
		"UpdateNode":                {},
		"PublishNodeVersion":        {},
		"UpdateNodeVersion":         {},
		"DeleteNodeVersion":         {},
		"CreatePersonalAccessToken": {},
		"DeletePersonalAccessToken": {},
	}
	return func(f strictecho.StrictEchoHandlerFunc, operationID string) strictecho.StrictEchoHandlerFunc {
		return func(c echo.Context, request interface{}) (response interface{}, err error) {
			// Bypass authorization for non-write operations
			if _, ok := restrictedOperationsForBannedUsers[operationID]; !ok {
				return f(c, request)
			}

			// Get user details from the context
			ctx := c.Request().Context()
			v := ctx.Value(UserContextKey)
			userDetails, ok := v.(*UserDetails)
			if !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "user not found")
			}

			u, err := entClient.User.Get(ctx, userDetails.ID)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "user not found")
			}

			if _, ok := restrictedOperationsForBannedUsers[operationID]; ok && u.Status == schema.UserStatusTypeBanned {
				return nil, echo.NewHTTPError(http.StatusForbidden, "user/publisher is banned")
			}

			return f(c, request)
		}

	}
}
