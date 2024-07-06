package drip_authentication

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
		"/security-scan": {"GET"},
		"/nodes/reindex": {"POST"},
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// Check if the request path and method are in the checklist
			path := ctx.Request().URL.Path
			method := ctx.Request().Method

			methods, ok := checklist[path]
			if !ok {
				return next(ctx)
			}

			for _, m := range methods {
				if method == m {
					ok = true
					break
				}
			}
			if !ok {
				return next(ctx)
			}

			// validate token
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
			if err != nil {
				log.Ctx(ctx.Request().Context()).Error().Err(err).Msg("Invalid token")
				return ctx.JSON(http.StatusUnauthorized, "Invalid token")
			}

			email, _ := payload.Claims["email"].(string)
			if email != "cloud-scheduler@dreamboothy.iam.gserviceaccount.com" {
				log.Ctx(ctx.Request().Context()).Error().Err(err).Msg("Invalid email")
				return ctx.JSON(http.StatusUnauthorized, "Invalid email")
			}

			log.Ctx(ctx.Request().Context()).Info().Msgf("Service Account Email: %s", email)
			return next(ctx)
		}
	}
}
