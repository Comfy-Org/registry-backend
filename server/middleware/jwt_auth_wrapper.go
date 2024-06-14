package drip_middleware

import (
	"context"
	"fmt"
	"net/http"
	"registry-backend/ent"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTWrapperMiddleware(entClient *ent.Client, secret string, auth echo.MiddlewareFunc) echo.MiddlewareFunc {
	keyfunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if header == "" {
				return auth(next)(c) // pass to next auth middleware
			}

			// Extract the JWT token from the header
			splitToken := strings.Split(header, "Bearer ")
			if len(splitToken) != 2 {
				return auth(next)(c) // pass to next auth middleware
			}
			token := splitToken[1]

			tokendata, err := jwt.Parse(token, keyfunc)
			if err != nil {
				return auth(next)(c) // pass to next auth middleware
			}

			claims, ok := tokendata.Claims.(jwt.MapClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid jwt token")
			}
			sub, err := claims.GetSubject()
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing sub claim")
			}

			user, err := entClient.User.Get(c.Request().Context(), sub)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid user")
			}

			authdCtx := context.WithValue(c.Request().Context(), UserContextKey, &UserDetails{
				ID:    user.ID,
				Email: user.Email,
				Name:  user.Name,
			})
			c.SetRequest(c.Request().WithContext(authdCtx))
			return next(c)
		}
	}
}
