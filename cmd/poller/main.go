package main

import (
	"content-hub/config"
	"content-hub/internal/infrastructure/mysql"
	"content-hub/internal/infrastructure/rabbitmq"
	"content-hub/internal/logger"
	"content-hub/internal/outbox"
	mysqlRepo "content-hub/internal/repository/mysql"
)

func main() {

	logger.Init()

	cfg := config.Load()

	db, err := mysql.NewDB(
		cfg.DBUrl,
	)

	if err != nil {

		logger.Log.Fatal().
			Err(err).
			Msg("failed connect mysql")
	}

	rmqConn, err := rabbitmq.NewRabbitMQ(
		cfg.RabbitMQUrl,
	)

	if err != nil {

		logger.Log.Fatal().
			Err(err).
			Msg("failed connect rabbitmq")
	}

	publisher, err := rabbitmq.NewPublisher(
		rmqConn,
	)

	if err != nil {

		logger.Log.Fatal().
			Err(err).
			Msg("failed create rabbitmq publisher")
	}

	outboxRepo := mysqlRepo.NewOutboxRepository(
		db,
	)

	logger.Log.Info().
		Msg("starting outbox poller")

	poller := outbox.NewPoller(
		outboxRepo,
		publisher,
	)

	poller.Start()
}
