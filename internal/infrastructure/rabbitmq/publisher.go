package rabbitmq

import (
	"content-hub/internal/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch *amqp.Channel
}

func NewPublisher(
	conn *amqp.Connection,
) (*Publisher, error) {

	ch, err := conn.Channel()

	if err != nil {

		logger.Error(err).
			Msg("failed create rabbitmq publisher channel")

		return nil, err
	}

	logger.Info().
		Msg("rabbitmq publisher initialized")

	return &Publisher{
		ch: ch,
	}, nil
}

func (p *Publisher) Publish(
	queue string,
	body []byte,
) error {

	err := DeclareTopology(
		p.ch,
		queue,
	)

	if err != nil {
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

	if err != nil {

		logger.Error(err).
			Str("queue", queue).
			Msg("failed publish rabbitmq message")

		return err
	}

	logger.Info().
		Str("queue", queue).
		Msg("message published")

	return nil
}
