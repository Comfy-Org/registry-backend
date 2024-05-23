package drip_middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/idtoken"
)

func ServiceAccountAuthMiddleware() echo.MiddlewareFunc {
	// Handlers in here should be checked by this middleware.
	var checklist = map[string][]string{
		"/users/sessions": {"DELETE"},
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			path := ctx.Request().URL.Path
			method := ctx.Request().Method

			// Check if the request path and method are in the checklist
			if methods, ok := checklist[path]; ok {
				for _, allowMethod := range methods {
					if method == allowMethod {
						authHeader := ctx.Request().Header.Get("Authorization")
						token := ""
						if strings.HasPrefix(authHeader, "Bearer ") {
							token = authHeader[7:] // Skip the "Bearer " part
						}

						if token == "" {
							return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
						}

						log.Ctx(ctx.Request().Context()).Info().Msgf("Validating google id token %s for path %s and method %s", token, path, method)

						payload, err := idtoken.Validate(ctx.Request().Context(), token, "https://api.comfy.org")

						if err == nil {
							if email, ok := payload.Claims["email"].(string); ok {
								log.Ctx(ctx.Request().Context()).Info().Msgf("Service Account Email: %s", email)
								// TODO(robinhuang): Make service account an environment variable.
								if email == "stop-vm-sa@dreamboothy.iam.gserviceaccount.com" {
									return next(ctx)
								}
							}
						}

						log.Ctx(ctx.Request().Context()).Error().Err(err).Msg("Invalid token")
						return ctx.JSON(http.StatusUnauthorized, "Invalid token")
					}
				}
			}

			// Proceed with the next middleware or handler
			return next(ctx)
		}
	}
}
