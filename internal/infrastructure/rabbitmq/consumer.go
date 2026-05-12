package rabbitmq

import (
	"time"

	"content-hub/internal/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	channel *amqp.Channel
}

func NewConsumer(
	conn *amqp.Connection,
) (*Consumer, error) {

	start := time.Now()

	ch, err := conn.Channel()

	duration := time.Since(start).Milliseconds()

	if err != nil {

		logger.Error(err).
			Str("service", "rabbitmq").
			Str("event", "consumer_channel_create_failed").
			Int64("duration_ms", duration).
			Msg("failed create rabbitmq consumer channel")

		return nil, err
	}

	logger.Info().
		Str("service", "rabbitmq").
		Str("event", "consumer_initialized").
		Int64("duration_ms", duration).
		Msg("rabbitmq consumer initialized")

	return &Consumer{
		channel: ch,
	}, nil
}

func (c *Consumer) Consume(
	queue string,
) (<-chan amqp.Delivery, error) {

	start := time.Now()

	err := DeclareTopology(
		c.channel,
		queue,
	)

	if err != nil {

		duration := time.Since(start).Milliseconds()

		logger.Error(err).
			Str("service", "rabbitmq").
			Str("event", "topology_declare_failed").
			Str("queue", queue).
			Int64("duration_ms", duration).
			Msg("failed declare rabbitmq topology")

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

	duration := time.Since(start).Milliseconds()

	if err != nil {

		logger.Error(err).
			Str("service", "rabbitmq").
			Str("event", "queue_consume_failed").
			Str("queue", queue).
			Int64("duration_ms", duration).
			Msg("failed consume rabbitmq queue")

		return nil, err
	}

	logger.Info().
		Str("service", "rabbitmq").
		Str("event", "queue_subscribed").
		Str("queue", queue).
		Int64("duration_ms", duration).
		Msg("rabbitmq consumer subscribed")

	return msgs, nil
}

func (c *Consumer) Channel() *amqp.Channel {
	return c.channel
}
