package rabbitmq

import (
	"content-hub/internal/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	channel *amqp.Channel
}

func NewConsumer(conn *amqp.Connection) (*Consumer, error) {

	ch, err := conn.Channel()

	if err != nil {

		logger.Log.Error().
			Err(err).
			Msg("failed create rabbitmq channel")

		return nil, err
	}

	logger.Log.Info().
		Msg("rabbitmq consumer initialized")

	return &Consumer{
		channel: ch,
	}, nil
}

func (c *Consumer) Consume(
	queue string,
) (<-chan amqp.Delivery, error) {

	_, err := c.channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {

		logger.Log.Error().
			Err(err).
			Str("queue", queue).
			Msg("failed declare rabbitmq queue")

		return nil, err
	}

	messages, err := c.channel.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {

		logger.Log.Error().
			Err(err).
			Str("queue", queue).
			Msg("failed consume rabbitmq queue")

		return nil, err
	}

	logger.Log.Info().
		Str("queue", queue).
		Msg("rabbitmq consumer started")

	return messages, nil
}
