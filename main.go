package main

import (
	"fmt"
	"os"
	"registry-backend/config"
	"registry-backend/ent"
	drip_logging "registry-backend/logging"
	"registry-backend/server"

	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq"
)

func main() {
	drip_logging.SetGlobalLogLevel(os.Getenv("LOG_LEVEL"))

	connection_string := os.Getenv("DB_CONNECTION_STRING")

	config := config.Config{
		ProjectID: os.Getenv("PROJECT_ID"),
		DripEnv:   os.Getenv("DRIP_ENV"),
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
	//if os.Getenv("DRIP_ENV") == "localdev" || os.Getenv("DRIP_ENV") == "staging" {
	//	log.Info().Msg("Running migrations")
	//	if err := client.Schema.Create(context.Background(), migrate.WithDropIndex(true),
	//		migrate.WithDropColumn(true)); err != nil {
	//		log.Fatal().Err(err).Msg("failed creating schema resources.")
	//	}
	//}

	server := server.NewServer(client, &config)
	server.Start()
}
