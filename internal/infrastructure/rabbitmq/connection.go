package rabbitmq

import (
	"time"

	"content-hub/internal/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQ(
	url string,
) (*amqp.Connection, error) {

	start := time.Now()

	conn, err := amqp.Dial(
		url,
	)

	duration := time.Since(start).Milliseconds()

	if err != nil {

		logger.Error(err).
			Str("service", "rabbitmq").
			Str("event", "connection_failed").
			Int64("duration_ms", duration).
			Msg("failed connect rabbitmq")

		return nil, err
	}

	logger.Info().
		Str("service", "rabbitmq").
		Str("event", "connected").
		Int64("duration_ms", duration).
		Msg("rabbitmq connected")

	return conn, nil
}
