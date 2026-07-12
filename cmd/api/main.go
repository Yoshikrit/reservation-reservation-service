package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"

	"github.com/Yoshikrit/reservation/config"
	"github.com/Yoshikrit/reservation/internal/controller/rest"
	"github.com/Yoshikrit/reservation/internal/pkg/logger"
	"github.com/Yoshikrit/reservation/internal/pkg/telemetry"
)

func main() {
	logger.Init()
	log.Info().Msg("reservation-api: starting")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("reservation-api: failed to load config")
	}

	db, err := config.InitDatabase(cfg.DatabaseConfig.DatabaseUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("reservation-api: failed to connect database")
	}
	if err := config.MigrateDatabase(db); err != nil {
		log.Fatal().Err(err).Msg("reservation-api: failed to migrate database")
	}
	config.SeedDatabase(db)

	redis := config.InitRedis(cfg.RedisConfig)

	config.InitKafkaProducer(cfg.KafkaConfig)

	grpcConns, err := config.NewGrpcConns(cfg.GrpcConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("reservation-api: failed to connect gRPC services")
	}
	defer grpcConns.Close()

	telemShutdown, err := telemetry.Init(context.Background(), cfg.TelemetryConfig.OtelServiceName, cfg.TelemetryConfig.OtelEndpoint)
	if err != nil {
		log.Warn().Err(err).Msg("reservation-api: telemetry init failed, continuing without tracing")
	} else {
		defer telemShutdown(context.Background())
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app := fiber.New(config.NewRestConfig(rest.ErrorHandler()))
	rest.NewRestRouter(app, db, redis, grpcConns)

	go func() {
		if err := app.Listen(":" + cfg.RestConfig.RestPort); err != nil {
			log.Error().Err(err).Msg("reservation-api: server stopped")
		}
	}()

	log.Info().Str("port", cfg.RestConfig.RestPort).Msg("reservation-api: listening")

	<-ctx.Done()
	stop()

	log.Info().Msg("reservation-api: shutting down")
	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = app.ShutdownWithContext(shutCtx)
}
