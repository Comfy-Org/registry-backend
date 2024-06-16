package drip_authentication

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"registry-backend/ent"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTAdminAuthMiddleware checks for a JWT token in the Authorization header,
// verifies it using the provided secret, and adds user details to the context if valid.
//
// This check is only performed for specific admin protected endpoints.
func JWTAdminAuthMiddleware(entClient *ent.Client, secret string) echo.MiddlewareFunc {
	// Key function to validate the JWT token
	keyfunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	}

	// Define the regex patterns for the protected endpoints
	protectedEndpoints := []*regexp.Regexp{
		regexp.MustCompile(`^/publishers/[^/]+/ban$`),
		regexp.MustCompile(`^/publishers/[^/]+/nodes/[^/]+/ban$`),
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqPath := c.Request().URL.Path

			// Check if the request path matches any of the protected endpoints
			isProtected := false
			for _, pattern := range protectedEndpoints {
				if pattern.MatchString(reqPath) {
					isProtected = true
					break
				}
			}

			if !isProtected {
				// If the request is not for a protected endpoint, skip this middleware
				return next(c)
			}

			// Get the Authorization header
			header := c.Request().Header.Get("Authorization")
			if header == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid jwt token")
			}

			// Extract the JWT token from the header
			splitToken := strings.Split(header, "Bearer ")
			if len(splitToken) != 2 {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid jwt token")
			}
			token := splitToken[1]

			// Parse and validate the JWT token
			tokenData, err := jwt.Parse(token, keyfunc)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid jwt token")
			}

			// Extract claims from the token
			claims, ok := tokenData.Claims.(jwt.MapClaims)
			if !ok || !tokenData.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid jwt token")
			}

			// Get the subject (user ID) from the claims
			sub, err := claims.GetSubject()
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing sub claim")
			}

			// Retrieve the user from the database
			user, err := entClient.User.Get(c.Request().Context(), sub)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid user")
			}

			// Add user details to the request context
			authContext := context.WithValue(c.Request().Context(), UserContextKey, &UserDetails{
				ID:    user.ID,
				Email: user.Email,
				Name:  user.Name,
			})
			c.SetRequest(c.Request().WithContext(authContext))

			// Call the next handler
			return next(c)
		}
	}
}
