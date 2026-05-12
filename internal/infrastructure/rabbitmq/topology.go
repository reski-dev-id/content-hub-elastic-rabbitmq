package rabbitmq

import (
	"time"

	"content-hub/internal/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareTopology(
	ch *amqp.Channel,
	queue string,
) error {

	start := time.Now()

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

		duration := time.Since(start).Milliseconds()

		logger.Error(err).
			Str("service", "rabbitmq").
			Str("event", "main_queue_declare_failed").
			Str("queue", queue).
			Int64("duration_ms", duration).
			Msg("failed declare main queue")

		return err
	}

	logger.Info().
		Str("service", "rabbitmq").
		Str("event", "main_queue_declared").
		Str("queue", queue).
		Int64("duration_ms", time.Since(start).Milliseconds()).
		Msg("main queue declared")

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

		duration := time.Since(start).Milliseconds()

		logger.Error(err).
			Str("service", "rabbitmq").
			Str("event", "retry_queue_declare_failed").
			Str("queue", queue+".retry").
			Int64("duration_ms", duration).
			Msg("failed declare retry queue")

		return err
	}

	logger.Info().
		Str("service", "rabbitmq").
		Str("event", "retry_queue_declared").
		Str("queue", queue+".retry").
		Int64("duration_ms", time.Since(start).Milliseconds()).
		Msg("retry queue declared")

	// DLQ

	_, err = ch.QueueDeclare(
		queue+".dlq",
		true,
		false,
		false,
		false,
		nil,
	)

	duration := time.Since(start).Milliseconds()

	if err != nil {

		logger.Error(err).
			Str("service", "rabbitmq").
			Str("event", "dlq_declare_failed").
			Str("queue", queue+".dlq").
			Int64("duration_ms", duration).
			Msg("failed declare dead letter queue")

		return err
	}

	logger.Info().
		Str("service", "rabbitmq").
		Str("event", "dlq_declared").
		Str("queue", queue+".dlq").
		Int64("duration_ms", duration).
		Msg("dead letter queue declared")

	logger.Info().
		Str("service", "rabbitmq").
		Str("event", "topology_declared").
		Str("queue", queue).
		Int64("duration_ms", duration).
		Msg("rabbitmq topology declared")

	return nil
}
