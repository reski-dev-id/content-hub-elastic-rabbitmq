package main

import (
	"time"

	"content-hub/config"
	"content-hub/internal/infrastructure/mysql"
	"content-hub/internal/infrastructure/rabbitmq"
	"content-hub/internal/logger"
	"content-hub/internal/outbox"
	mysqlRepo "content-hub/internal/repository/mysql"
)

func main() {

	start := time.Now()

	logger.Init()

	cfg := config.Load()

	db, err := mysql.NewDB(
		cfg.DBUrl,
	)

	if err != nil {

		logger.Fatal(err).
			Str("service", "poller").
			Str("event", "mysql_connection_failed").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to connect mysql")
	}

	rmqConn, err := rabbitmq.NewRabbitMQ(
		cfg.RabbitMQUrl,
	)

	if err != nil {

		logger.Fatal(err).
			Str("service", "poller").
			Str("event", "rabbitmq_connection_failed").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to connect rabbitmq")
	}

	publisher, err := rabbitmq.NewPublisher(
		rmqConn,
	)

	if err != nil {

		logger.Fatal(err).
			Str("service", "poller").
			Str("event", "rabbitmq_publisher_create_failed").
			Int64(
				"duration_ms",
				time.Since(start).Milliseconds(),
			).
			Msg("failed to create rabbitmq publisher")
	}

	outboxRepo := mysqlRepo.NewOutboxRepository(
		db,
	)

	logger.Info().
		Str("service", "poller").
		Str("event", "poller_started").
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("starting outbox poller")

	poller := outbox.NewPoller(
		outboxRepo,
		publisher,
	)

	logger.Info().
		Str("service", "poller").
		Str("event", "poller_running").
		Dur(
			"duration",
			time.Since(start),
		).
		Int64(
			"duration_ms",
			time.Since(start).Milliseconds(),
		).
		Msg("outbox poller running")

	poller.Start()
}
