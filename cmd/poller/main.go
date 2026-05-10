package main

import (
	"log"

	"content-hub/config"
	"content-hub/internal/infrastructure/mysql"
	"content-hub/internal/infrastructure/rabbitmq"
	"content-hub/internal/outbox"
	mysqlRepo "content-hub/internal/repository/mysql"
)

func main() {
	cfg := config.Load()

	db, err := mysql.NewDB(cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}

	rmqConn, err := rabbitmq.NewConnection(
		cfg.RabbitMQUrl,
	)

	if err != nil {
		log.Fatal(err)
	}

	publisher, err := rabbitmq.NewPublisher(rmqConn)
	if err != nil {
		log.Fatal(err)
	}

	outboxRepo := mysqlRepo.NewOutboxRepository(db)

	poller := outbox.NewPoller(
		outboxRepo,
		publisher,
	)

	poller.Start()
}
