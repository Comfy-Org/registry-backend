package main

import (
	"context"
	"fmt"
	"os"
	"registry-backend/config"
	"registry-backend/ent"
	"registry-backend/ent/migrate"
	drip_logging "registry-backend/logging"
	"registry-backend/server"

	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq"
)

func main() {
	drip_logging.SetGlobalLogLevel(os.Getenv("LOG_LEVEL"))

	connection_string := os.Getenv("DB_CONNECTION_STRING")

	config := config.Config{
		ProjectID:                     os.Getenv("PROJECT_ID"),
		DripEnv:                       os.Getenv("DRIP_ENV"),
		SlackRegistryChannelWebhook:   os.Getenv("SLACK_REGISTRY_CHANNEL_WEBHOOK"),
		JWTSecret:                     os.Getenv("JWT_SECRET"),
		SecretScannerURL:              os.Getenv("SECRET_SCANNER_URL"),
		DiscordSecurityChannelWebhook: os.Getenv("SECURITY_COUNCIL_DISCORD_WEBHOOK"),
	}

	var dsn string
	if os.Getenv("DRIP_ENV") == "localdev" {
		dsn = fmt.Sprintf("%s sslmode=disable", connection_string)
	} else {
		dsn = connection_string
	}

	client, err := ent.Open("postgres", dsn)

	if err != nil {
		log.Fatal().Err(err).Msg("failed opening connection to postgres.")
	}
	defer client.Close()
	// Run the auto migration tool for localdev.
	if os.Getenv("DRIP_ENV") == "localdev" {
		log.Info().Msg("Running migrations")
		if err := client.Schema.Create(context.Background(), migrate.WithDropIndex(true),
			migrate.WithDropColumn(true)); err != nil {
			log.Fatal().Err(err).Msg("failed creating schema resources.")
		}
	}

	server := server.NewServer(client, &config)
	log.Fatal().Err(server.Start()).Msg("Server stopped")
}
