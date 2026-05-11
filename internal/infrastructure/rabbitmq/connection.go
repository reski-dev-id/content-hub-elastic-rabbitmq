package rabbitmq

import (
	"content-hub/internal/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQ(url string) (*amqp.Connection, error) {

	conn, err := amqp.Dial(url)

	if err != nil {

		logger.Log.Error().
			Err(err).
			Msg("failed connect rabbitmq")

		return nil, err
	}

	logger.Log.Info().
		Msg("rabbitmq connected")

	return conn, nil
}
