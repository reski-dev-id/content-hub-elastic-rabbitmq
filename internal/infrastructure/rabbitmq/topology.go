package rabbitmq

import (
	"content-hub/internal/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareTopology(
	ch *amqp.Channel,
	queue string,
) error {

	// MAIN QUEUE

	_, err := ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-routing-key": queue + ".retry",
		},
	)

	if err != nil {

		logger.Error(err).
			Str("queue", queue).
			Msg("failed declare main queue")

		return err
	}

	// RETRY QUEUE

	_, err = ch.QueueDeclare(
		queue+".retry",
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-message-ttl":             10000,
			"x-dead-letter-routing-key": queue,
		},
	)

	if err != nil {

		logger.Error(err).
			Str("queue", queue+".retry").
			Msg("failed declare retry queue")

		return err
	}

	// DLQ

	_, err = ch.QueueDeclare(
		queue+".dlq",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {

		logger.Error(err).
			Str("queue", queue+".dlq").
			Msg("failed declare dead letter queue")

		return err
	}

	logger.Info().
		Str("queue", queue).
		Msg("rabbitmq topology declared")

	return nil
}
