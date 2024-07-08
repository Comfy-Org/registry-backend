package server

import (
	"context"
	"os"
	"registry-backend/config"
	generated "registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/gateways/algolia"
	"registry-backend/gateways/discord"
	gateway "registry-backend/gateways/slack"
	"registry-backend/gateways/storage"
	handler "registry-backend/server/handlers"
	"registry-backend/server/implementation"
	drip_middleware "registry-backend/server/middleware"
	drip_authentication "registry-backend/server/middleware/authentication"
	drip_authorization "registry-backend/server/middleware/authorization"
	drip_metric "registry-backend/server/middleware/metric"
	"strings"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"

	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"
)

type Server struct {
	Client *ent.Client
	Config *config.Config
}

func NewServer(client *ent.Client, config *config.Config) *Server {
	return &Server{
		Client: client,
		Config: config,
	}
}

func (s *Server) Start() error {
	e := echo.New()
	e.HideBanner = true
	e.Use(drip_middleware.TracingMiddleware)
	//e.Use(middleware.Logger()) //  Useful for debugging
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) (bool, error) {
			allowedOrigins := []string{
				".comfyci.org",           // Any subdomain of comfyci.org
				os.Getenv("CORS_ORIGIN"), // Environment-specific allowed origin
				".comfyregistry.org",
				".comfy.org",
			}

			for _, allowed := range allowedOrigins {
				if strings.HasSuffix(origin, allowed) || origin == allowed {
					log.Debug().Msg("[CORSWithConfig] Allowing origin " + origin)
					return true, nil
				}
			}

			log.Debug().Msg("[CORSWithConfig] Rejecting origin " + origin)
			return false, nil
		},
	}))

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			// Ignore when url is path /vm/{sessionId}
			if strings.HasPrefix(c.Request().URL.Path, "/vm/") {
				return nil
			}

			log.Ctx(c.Request().Context()).Debug().
				Str("URI: ", v.URI).
				Int("status", v.Status).Msg("")
			return nil
		},
	}))

	storageService, err := storage.NewGCPStorageService(context.Background())
	if err != nil {
		return err
	}

	slackService := gateway.NewSlackService(s.Config)
	algoliaService, err := algolia.NewFromEnvOrNoop()
	if err != nil {
		return err
	}
	discordService := discord.NewDiscordService(s.Config)

	mon, err := monitoring.NewMetricClient(context.Background())
	if err != nil {
		return err
	}

	// Attach implementation of generated oapi strict server.
	impl := implementation.NewStrictServerImplementation(s.Client, s.Config, storageService, slackService, discordService, algoliaService)

	// Define middlewares in the order of operations
	authorizationManager := drip_authorization.NewAuthorizationManager(s.Client, impl.RegistryService)
	middlewares := []generated.StrictMiddlewareFunc{
		authorizationManager.AuthorizationMiddleware(),
	}
	wrapped := generated.NewStrictHandler(impl, middlewares)

	generated.RegisterHandlers(e, wrapped)

	e.GET("/openapi", handler.SwaggerHandler)
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	// Global Middlewares
	e.Use(drip_metric.MetricsMiddleware(mon, s.Config))
	e.Use(drip_authentication.FirebaseAuthMiddleware(s.Client))
	e.Use(drip_authentication.ServiceAccountAuthMiddleware())
	e.Use(drip_authentication.JWTAdminAuthMiddleware(s.Client, s.Config.JWTSecret))
	e.Use(drip_middleware.ErrorLoggingMiddleware())

	e.Logger.Fatal(e.Start(":8080"))
	return nil
}
