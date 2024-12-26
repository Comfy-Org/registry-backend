package drip_authentication

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/idtoken"
)

func ServiceAccountAuthMiddleware() echo.MiddlewareFunc {
	// Handlers in here should be checked by this middleware.
	var checklistRegex = map[*regexp.Regexp][]string{
		regexp.MustCompile(`^/security-scan$`):                          {"GET"},
		regexp.MustCompile(`^/nodes/reindex$`):                          {"POST"},
		regexp.MustCompile(`^/nodes/[^/]+/versions/[^/]+/comfy-nodes$`): {"POST"},
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// Check if the request reqPath and method are in the checklist
			reqPath := ctx.Request().URL.Path
			reqMethod := ctx.Request().Method
			match := false
			for pathRe, methods := range checklistRegex {
				if !pathRe.MatchString(reqPath) {
					continue
				}

				for _, method := range methods {
					if method != "ANY" && reqMethod != method {
						continue
					}

					match = true
					break
				}
			}
			if !match {
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

			log.Ctx(ctx.Request().Context()).Info().Msgf("Validating google id token %s for path %s and method %s", token, reqPath, reqMethod)

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
