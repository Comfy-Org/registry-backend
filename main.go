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

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	config := config.Config{
		ProjectID:  os.Getenv("PROJECT_ID"),
		DripEnv:    os.Getenv("DRIP_ENV"),
		SelfOrigin: os.Getenv("SELF_ORIGIN"),
	}

	var dsn string
	if os.Getenv("DRIP_ENV") == "localdev" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password)
	} else {
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", host, port, user, dbname, password)
	}

	client, err := ent.Open("postgres", dsn)

	if err != nil {
		log.Fatal().Err(err).Msg("failed opening connection to postgres.")
	}
	defer client.Close()
	// Run the auto migration tool for localdev.
	if os.Getenv("DRIP_ENV") == "localdev" || os.Getenv("DRIP_ENV") == "staging" {
		log.Info().Msg("Running migrations")
		if err := client.Schema.Create(context.Background(), migrate.WithDropIndex(true),
			migrate.WithDropColumn(true)); err != nil {
			log.Fatal().Err(err).Msg("failed creating schema resources.")
		}
	}

	server := server.NewServer(client, &config)
	server.Start()
}
