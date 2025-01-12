package server

import (
	"context"
	"registry-backend/config"
	generated "registry-backend/drip"
	"registry-backend/ent"
	"registry-backend/gateways/algolia"
	"registry-backend/gateways/discord"
	"registry-backend/gateways/pubsub"
	"registry-backend/gateways/slack"
	"registry-backend/gateways/storage"
	handler "registry-backend/server/handlers"
	"registry-backend/server/implementation"
	"registry-backend/server/middleware"
	"registry-backend/server/middleware/authentication"
	drip_authorization "registry-backend/server/middleware/authorization"
	"registry-backend/server/middleware/metric"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"github.com/labstack/echo/v4"
	labstack_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

type ServerDependencies struct {
	StorageService   storage.StorageService
	PubSubService    pubsub.PubSubService
	SlackService     slack.SlackService
	AlgoliaService   algolia.AlgoliaService
	DiscordService   discord.DiscordService
	MonitoringClient monitoring.MetricClient
}

type Server struct {
	Client       *ent.Client
	Config       *config.Config
	Dependencies *ServerDependencies
}

func NewServer(client *ent.Client, config *config.Config) (*Server, error) {
	deps, err := initializeDependencies(config)
	if err != nil {
		return nil, err
	}
	return &Server{
		Client:       client,
		Config:       config,
		Dependencies: deps,
	}, nil
}

func initializeDependencies(config *config.Config) (*ServerDependencies, error) {
	storageService, err := storage.NewStorageService(config)
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize storage service")
		return nil, err
	}

	pubsubService, err := pubsub.NewPubSubService(config)
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize pub/sub service")
		return nil, err
	}

	slackService := slack.NewSlackService(config)

	algoliaService, err := algolia.NewAlgoliaService(config)
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize Algolia service")
		return nil, err
	}

	discordService := discord.NewDiscordService(config)

	mon, err := monitoring.NewMetricClient(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize monitoring client")
		return nil, err
	}

	return &ServerDependencies{
		StorageService:   storageService,
		SlackService:     slackService,
		PubSubService:    pubsubService,
		AlgoliaService:   algoliaService,
		DiscordService:   discordService,
		MonitoringClient: *mon,
	}, nil
}

func (s *Server) Start() error {
	e := echo.New()
	e.HideBanner = true

	// Apply middleware
	e.Use(middleware.TracingMiddleware)
	e.Use(labstack_middleware.CORSWithConfig(labstack_middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowOriginFunc: func(origin string) (bool, error) {
			return true, nil
		},
		AllowCredentials: true,
	}))
	e.Use(middleware.RequestLoggerMiddleware())
	e.Use(middleware.ResponseLoggerMiddleware())
	e.Use(metric.MetricsMiddleware(&s.Dependencies.MonitoringClient, s.Config))
	e.Use(authentication.FirebaseAuthMiddleware(s.Client))
	e.Use(authentication.ServiceAccountAuthMiddleware())
	e.Use(authentication.JWTAdminAuthMiddleware(s.Client, s.Config.JWTSecret))

	// Attach implementation of the generated OAPI strict server
	impl := implementation.NewStrictServerImplementation(
		s.Client, s.Config, s.Dependencies.StorageService, s.Dependencies.PubSubService,
		s.Dependencies.SlackService,
		s.Dependencies.DiscordService, s.Dependencies.AlgoliaService)

	// Define middleware for authorization
	authorizationManager := drip_authorization.NewAuthorizationManager(s.Client, impl.RegistryService)
	middlewares := []generated.StrictMiddlewareFunc{
		authorizationManager.AuthorizationMiddleware(),
	}

	// Create the strict handler with middlewares
	wrapped := generated.NewStrictHandler(impl, middlewares)

	// Register routes
	generated.RegisterHandlers(e, wrapped)

	// Define public routes
	e.GET("/openapi", handler.SwaggerHandler)
	e.GET("/health", s.HealthCheckHandler)

	// Start the server
	return e.Start(":8080")
}

// HealthCheckHandler performs health checks on the critical dependencies
func (s *Server) HealthCheckHandler(c echo.Context) error {
	// This could be extended to check storage, slack, and other dependencies
	return c.String(200, "OK")
}
