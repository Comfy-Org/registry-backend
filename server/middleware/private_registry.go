package middleware

import (
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Match the operation's position, not a node or Publisher named "feedback".
var versionFeedbackPath = regexp.MustCompile(`^/(?:admin/nodes/[^/]+|publishers/[^/]+/nodes/[^/]+)/versions/[^/]+/feedback(?:/|$)`)

// IsPrivateRegistryPath is also used by body loggers, before authentication runs.
func IsPrivateRegistryPath(path string) bool {
	return strings.HasPrefix(path, "/admin/") || path == "/users/me/node-version-feedback" || versionFeedbackPath.MatchString(path)
}
func PrivateRegistryData() echo.MiddlewareFunc {
	limitBody := echomiddleware.BodyLimit("32K")
	limitWrites := echomiddleware.RateLimiter(echomiddleware.NewRateLimiterMemoryStoreWithConfig(echomiddleware.RateLimiterMemoryStoreConfig{Rate: 1, Burst: 30, ExpiresIn: time.Minute}))
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path
			if IsPrivateRegistryPath(path) {
				c.Response().Header().Set("Cache-Control", "private, no-store")
				c.Response().Header().Add("Vary", "Authorization")
				if versionFeedbackPath.MatchString(path) && c.Request().Method != http.MethodGet {
					return limitBody(limitWrites(next))(c)
				}
			} else if (c.QueryParams().Has("status_reason") || c.QueryParams().Has("feedback_status")) && (path == "/versions" || strings.HasPrefix(path, "/nodes")) {
				return echo.NewHTTPError(400, "Private filters require an admin endpoint")
			}
			return next(c)
		}
	}
}
