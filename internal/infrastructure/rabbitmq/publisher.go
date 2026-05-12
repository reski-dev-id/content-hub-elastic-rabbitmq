package rabbitmq

import (
	"time"

	"content-hub/internal/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch *amqp.Channel
}

func NewPublisher(
	conn *amqp.Connection,
) (*Publisher, error) {

	start := time.Now()

	ch, err := conn.Channel()

	duration := time.Since(start).Milliseconds()

	if err != nil {

		logger.Error(err).
			Str("service", "rabbitmq").
			Str("event", "publisher_channel_create_failed").
			Int64("duration_ms", duration).
			Msg("failed create rabbitmq publisher channel")

		return nil, err
	}

	logger.Info().
		Str("service", "rabbitmq").
		Str("event", "publisher_initialized").
		Int64("duration_ms", duration).
		Msg("rabbitmq publisher initialized")

	return &Publisher{
		ch: ch,
	}, nil
}

func (p *Publisher) Publish(
	queue string,
	body []byte,
) error {

	start := time.Now()

	err := DeclareTopology(
		p.ch,
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

		return err
	}

	err = p.ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	duration := time.Since(start).Milliseconds()

	if err != nil {

		logger.Error(err).
			Str("service", "rabbitmq").
			Str("event", "message_publish_failed").
			Str("queue", queue).
			Int64("duration_ms", duration).
			Msg("failed publish rabbitmq message")

		return err
	}

	logger.Info().
		Str("service", "rabbitmq").
		Str("event", "message_published").
		Str("queue", queue).
		Int64("duration_ms", duration).
		Msg("message published")

	return nil
}
