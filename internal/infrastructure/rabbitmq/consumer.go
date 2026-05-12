package rabbitmq

import (
	"content-hub/internal/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	channel *amqp.Channel
}

func NewConsumer(
	conn *amqp.Connection,
) (*Consumer, error) {

	ch, err := conn.Channel()

	if err != nil {

		logger.Error(err).
			Msg("failed create rabbitmq consumer channel")

		return nil, err
	}

	logger.Info().
		Msg("rabbitmq consumer initialized")

	return &Consumer{
		channel: ch,
	}, nil
}

func (c *Consumer) Consume(
	queue string,
) (<-chan amqp.Delivery, error) {

	err := DeclareTopology(
		c.channel,
		queue,
	)

	if err != nil {
		return nil, err
	}

	msgs, err := c.channel.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {

		logger.Error(err).
			Str("queue", queue).
			Msg("failed consume rabbitmq queue")

		return nil, err
	}

	logger.Info().
		Str("queue", queue).
		Msg("rabbitmq consumer subscribed")

	return msgs, nil
}

func (c *Consumer) Channel() *amqp.Channel {
	return c.channel
}
