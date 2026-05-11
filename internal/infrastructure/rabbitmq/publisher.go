package rabbitmq

import (
	"content-hub/internal/logger"

	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch *amqp091.Channel
}

func NewPublisher(
	conn *amqp091.Connection,
) (*Publisher, error) {

	ch, err := conn.Channel()

	if err != nil {

		logger.Log.Error().
			Err(err).
			Msg("failed create rabbitmq publisher channel")

		return nil, err
	}

	logger.Log.Info().
		Msg("rabbitmq publisher initialized")

	return &Publisher{ch}, nil
}

func (p *Publisher) Publish(
	queue string,
	body []byte,
) error {

	_, err := p.ch.QueueDeclare(
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

		return err
	}

	err = p.ch.Publish(
		"",
		queue,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {

		logger.Log.Error().
			Err(err).
			Str("queue", queue).
			Msg("failed publish rabbitmq message")

		return err
	}

	logger.Log.Info().
		Str("queue", queue).
		Msg("message published to rabbitmq")

	return nil
}
