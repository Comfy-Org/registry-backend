package drip_jobs

import (
	"context"
	"fmt"
	"registry-backend/ent"
	drip_metric "registry-backend/server/middleware/metric"
	drip_services "registry-backend/services/registry"

	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog/log"
)

func ReindexAllNodes(ctx context.Context, client *ent.Client, registry drip_services.RegistryService, crontab string) error {
	s, _ := gocron.NewScheduler()
	defer func() { _ = s.Shutdown() }()

	if crontab == "" {
		crontab = "0 0 * * *"
	}

	_, err := s.NewJob(
		gocron.CronJob(crontab, false),
		gocron.NewTask(func() {
			err := registry.ReindexAllNodes(ctx, client)
			if err != nil {
				log.Error().Msgf("Failed to reindex all nodes w/ err: %v", err)
				drip_metric.IncrementCustomCounterMetric(ctx, drip_metric.CustomCounterIncrement{
					Type:   "reindex-error",
					Val:    1,
					Labels: map[string]string{},
				})
			}
		}),
	)
	if err != nil {
		return fmt.Errorf("Failed to create job w/ err: %v", err)
	}

	log.Info().Msgf("Starting the scheduled jobs for reindexing nodes")
	s.Start()

	<-ctx.Done()
	log.Info().Msgf("Stoping the scheduled jobs for reindexing nodes")

	return err

}
