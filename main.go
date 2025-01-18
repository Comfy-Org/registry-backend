package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"registry-backend/config"
	"registry-backend/ent"
	"registry-backend/ent/migrate"
	drip_logging "registry-backend/logging"
	"registry-backend/server"
	"syscall"

	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq"
)

// validateEnvVars ensures all required environment variables are set based on the environment.
func validateEnvVars(env string) {
	// Variables mandatory for all environments
	mandatoryVars := []string{
		"DB_CONNECTION_STRING",
		"PROJECT_ID",
		"DRIP_ENV",
		"JWT_SECRET",
	}

	// TODO: Add staging specific variables

	// Additional variables mandatory for production and staging
	prodVars := []string{
		"CLOUD_STORAGE_BUCKET_NAME",
		"SLACK_REGISTRY_CHANNEL_WEBHOOK",
		"SECRET_SCANNER_URL",
		"SECURITY_COUNCIL_DISCORD_WEBHOOK",
		"ALGOLIA_APP_ID",
		"ALGOLIA_API_KEY",
		"ID_TOKEN_AUDIENCE",
		"PUBSUB_TOPIC",
	}

	// Add production specific variables
	if env == "prod" {
		mandatoryVars = append(mandatoryVars, prodVars...)
	}

	// Validate that all mandatory environment variables are set
	var missingVars []string
	for _, key := range mandatoryVars {
		if os.Getenv(key) == "" {
			missingVars = append(missingVars, key)
		}
	}

	// Log and terminate if mandatory variables are missing
	if len(missingVars) > 0 {
		log.Fatal().Msgf("Missing mandatory environment variables for '%s': %v", env, missingVars)
	}
}

func main() {
	// Retrieve the current environment
	dripEnv := os.Getenv("DRIP_ENV")
	if dripEnv == "" {
		log.Fatal().Msg("Environment variable DRIP_ENV is not set.")
	}

	// Validate environment variables based on the current environment
	validateEnvVars(dripEnv)

	// Set global log level based on the LOG_LEVEL environment variable
	drip_logging.SetGlobalLogLevel(os.Getenv("LOG_LEVEL"))

	// Retrieve the database connection string
	connectionString := os.Getenv("DB_CONNECTION_STRING")

	// Build the application configuration
	appConfig := config.Config{
		ProjectID:                            os.Getenv("PROJECT_ID"),
		DripEnv:                              dripEnv,
		SlackRegistryChannelWebhook:          os.Getenv("SLACK_REGISTRY_CHANNEL_WEBHOOK"),
		JWTSecret:                            os.Getenv("JWT_SECRET"),
		SecretScannerURL:                     os.Getenv("SECRET_SCANNER_URL"),
		DiscordSecurityChannelWebhook:        os.Getenv("SECURITY_COUNCIL_DISCORD_WEBHOOK"),
		DiscordSecurityPrivateChannelWebhook: os.Getenv("SECURITY_COUNCIL_DISCORD_PRIVATE_WEBHOOK"),
		AlgoliaAppID:                         os.Getenv("ALGOLIA_APP_ID"),
		AlgoliaAPIKey:                        os.Getenv("ALGOLIA_API_KEY"),
		IDTokenAudience:                      os.Getenv("ID_TOKEN_AUDIENCE"),
		CloudStorageBucketName:               os.Getenv("CLOUD_STORAGE_BUCKET_NAME"),
		PubSubTopic:                          os.Getenv("PUBSUB_TOPIC"),
		NewRelicLicenseKey:                   os.Getenv("NEW_RELIC_LICENSE_KEY"),
	}

	// Construct the database connection string
	var dsn string
	if dripEnv == "localdev" {
		// For local development, disable SSL for easier setup
		dsn = fmt.Sprintf("%s sslmode=disable", connectionString)
	} else {
		// Use the connection string as-is for non-development environments
		dsn = connectionString
	}

	// Initialize the database client
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to establish a connection to the PostgreSQL database.")
	}
	defer client.Close() // Ensure the database client is closed when the application exits

	// Run database migrations in local development to keep the schema up to date
	if dripEnv == "localdev" {
		log.Info().Msg("Running migrations for local development.")
		if err := client.Schema.Create(context.Background(),
			migrate.WithDropIndex(true), migrate.WithDropColumn(true)); err != nil {
			log.Fatal().Err(err).Msg("Failed to create schema resources during migration.")
		}
	}

	// Handle graceful shutdown
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		log.Info().Msg("Shutting down server gracefully.")
		client.Close()
		os.Exit(0)
	}()

	// Initialize and start the server
	registryServer, err := server.NewServer(client, &appConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize the server.")
	}

	log.Info().Msg("Starting the server.")
	log.Fatal().Err(registryServer.Start()).Msg("Server has stopped unexpectedly.")
}
