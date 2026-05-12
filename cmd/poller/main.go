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

		logger.Fatal(err).
			Msg("failed to connect mysql")
	}

	rmqConn, err := rabbitmq.NewRabbitMQ(
		cfg.RabbitMQUrl,
	)

	if err != nil {

		logger.Fatal(err).
			Msg("failed to connect rabbitmq")
	}

	publisher, err := rabbitmq.NewPublisher(
		rmqConn,
	)

	if err != nil {

		logger.Fatal(err).
			Msg("failed to create rabbitmq publisher")
	}

	outboxRepo := mysqlRepo.NewOutboxRepository(
		db,
	)

	logger.Info().
		Msg("starting outbox poller")

	poller := outbox.NewPoller(
		outboxRepo,
		publisher,
	)

	poller.Start()
}
