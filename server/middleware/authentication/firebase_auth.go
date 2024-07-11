package drip_authentication

import (
	"context"
	"net/http"
	"regexp"
	"registry-backend/db"
	"registry-backend/ent"
	"strings"

	"github.com/rs/zerolog/log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"
)

// FirebaseAuthMiddleware validates and extracts user details from the Firebase token.
// Certain endpoints are allow-listed and bypass this middleware.
func FirebaseAuthMiddleware(entClient *ent.Client) echo.MiddlewareFunc {
	// Handlers in here should bypass this middleware.
	var allowlist = map[*regexp.Regexp][]string{
		regexp.MustCompile(`^/openapi$`):                               {"GET"},
		regexp.MustCompile(`^/security-scan$`):                         {"GET"},
		regexp.MustCompile(`^/users/sessions$`):                        {"DELETE"},
		regexp.MustCompile(`^/vm$`):                                    {"ANY"},
		regexp.MustCompile(`^/health$`):                                {"GET"},
		regexp.MustCompile(`^/upload-artifact$`):                       {"POST"},
		regexp.MustCompile(`^/gitcommit$`):                             {"POST", "GET"},
		regexp.MustCompile(`^/workflowresult/[^/]+$`):                  {"GET"},
		regexp.MustCompile(`^/branch$`):                                {"GET"},
		regexp.MustCompile(`^/publishers/[^/]+/nodes/[^/]+/versions$`): {"POST"},
		regexp.MustCompile(`^/publishers/[^/]+/nodes$`):                {"GET"},
		regexp.MustCompile(`^/publishers/[^/]+$`):                      {"GET"},
		regexp.MustCompile(`^/nodes$`):                                 {"GET"},
		regexp.MustCompile(`^/versions$`):                              {"GET"},
		regexp.MustCompile(`^/nodes/[^/]+$`):                           {"GET"},
		regexp.MustCompile(`^/nodes/[^/]+/versions$`):                  {"GET"},
		regexp.MustCompile(`^/nodes/[^/]+/install$`):                   {"GET"},
		regexp.MustCompile(`^/nodes/reindex$`):                         {"POST"},
		regexp.MustCompile(`^/publishers/[^/]+/ban$`):                  {"POST"},
		regexp.MustCompile(`^/publishers/[^/]+/nodes/[^/]+/ban$`):      {"POST"},
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// Check if the request is in the allow list.
			reqPath := ctx.Request().URL.Path
			reqMethod := ctx.Request().Method
			for basePathRegex, methods := range allowlist {
				if basePathRegex.MatchString(reqPath) {
					for _, method := range methods {
						if method == "ANY" || reqMethod == method {
							log.Ctx(ctx.Request().Context()).Debug().
								Msgf("Letting through %s request to %s", reqMethod, reqPath)
							return next(ctx)
						}
					}
				}
			}

			// If header is present, extract the token and verify it.
			header := ctx.Request().Header.Get("Authorization")
			if header != "" {
				// Extract the JWT token from the header
				splitToken := strings.Split(header, "Bearer ")
				if len(splitToken) != 2 {
					return echo.NewHTTPError(http.StatusUnauthorized, "token is not in Bearer format")
				}
				idToken := splitToken[1]

				// Initialize Firebase
				app, err := firebase.NewApp(context.Background(), nil)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "firebase initialization failed")
				}

				client, err := app.Auth(context.Background())
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "firebase auth client failed")
				}

				// Verify ID token
				token, err := client.VerifyIDToken(context.Background(), idToken)
				if err != nil {
					// print the error
					log.Ctx(ctx.Request().Context()).Error().Err(err).Msg("error verifying ID token")
					return echo.NewHTTPError(http.StatusUnauthorized, "invalid auth token")
				}

				userDetails := extractUserDetails(token)
				log.Ctx(ctx.Request().Context()).Debug().Msg("Authenticated user " + userDetails.Email)

				authContext := context.WithValue(ctx.Request().Context(), UserContextKey, userDetails)
				ctx.SetRequest(ctx.Request().WithContext(authContext))

				newUserError := db.UpsertUser(
					ctx.Request().Context(), entClient, token.UID, userDetails.Email, userDetails.Name)
				if newUserError != nil {
					log.Ctx(ctx.Request().Context()).Error().Err(newUserError).Msg("error User upserted successfully.")
				}

				return next(ctx)
			}

			return echo.NewHTTPError(http.StatusUnauthorized, "missing auth token")
		}
	}
}

type ContextKey string

const UserContextKey ContextKey = "user"

type UserDetails struct {
	ID    string
	Email string
	Name  string
}

func extractUserDetails(token *auth.Token) *UserDetails {
	claims := token.Claims
	email, _ := claims["email"].(string)
	name, _ := claims["name"].(string)

	return &UserDetails{
		ID:    token.UID,
		Email: email,
		Name:  name,
	}
}
