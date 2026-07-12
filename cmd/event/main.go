package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/Yoshikrit/reservation/config"
	"github.com/Yoshikrit/reservation/internal/controller/event"
	"github.com/Yoshikrit/reservation/internal/pkg/logger"
	"github.com/Yoshikrit/reservation/internal/pkg/telemetry"
)

func main() {
	logger.Init()
	log.Info().Msg("reservation-event: starting")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("reservation-event: failed to load config")
	}

	db, err := config.InitDatabase(cfg.DatabaseConfig.DatabaseUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("reservation-event: failed to connect database")
	}
	if err := config.MigrateDatabase(db); err != nil {
		log.Fatal().Err(err).Msg("reservation-event: failed to migrate database")
	}
	config.SeedDatabase(db)

	producer := config.InitKafkaProducer(cfg.KafkaConfig)
	if producer != nil {
		defer producer.Close()
	}

	shutdown, err := telemetry.Init(context.Background(), cfg.TelemetryConfig.OtelServiceName, cfg.TelemetryConfig.OtelEndpoint)
	if err != nil {
		log.Warn().Err(err).Msg("reservation-event: telemetry init failed, continuing without tracing")
	} else {
		defer shutdown(context.Background())
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	event.StartEvent(ctx, db, producer)

	log.Info().Msg("reservation-event: relay started, waiting for shutdown")

	<-ctx.Done()
	stop()

	log.Info().Msg("reservation-event: shutting down")
}
