package rabbitmq

import (
	"time"

	"content-hub/internal/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

const MaxRetry = 3

func GetRetryCount(
	msg amqp.Delivery,
) int {

	if val, ok := msg.Headers["x-retry"]; ok {

		if retry, ok := val.(int32); ok {
			return int(retry)
		}
	}

	return 0
}

func RetryMessage(
	ch *amqp.Channel,
	queue string,
	msg amqp.Delivery,
) error {

	start := time.Now()

	retryCount := GetRetryCount(msg)

	if retryCount >= MaxRetry {

		logger.Error(nil).
			Str("service", "rabbitmq").
			Str("event", "message_moved_to_dlq").
			Str("queue", queue).
			Int("retry_count", retryCount).
			Msg("message moved to dead letter queue")

		err := ch.Publish(
			"",
			queue+".dlq",
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        msg.Body,
				Headers: amqp.Table{
					"x-retry": retryCount,
				},
			},
		)

		duration := time.Since(start).Milliseconds()

		if err != nil {

			logger.Error(err).
				Str("service", "rabbitmq").
				Str("event", "dlq_publish_failed").
				Str("queue", queue).
				Int("retry_count", retryCount).
				Int64("duration_ms", duration).
				Msg("failed publish message to dlq")

			return err
		}

		logger.Info().
			Str("service", "rabbitmq").
			Str("event", "dlq_publish_success").
			Str("queue", queue).
			Int("retry_count", retryCount).
			Int64("duration_ms", duration).
			Msg("message published to dlq")

		return nil
	}

	retryCount++

	logger.Warn().
		Str("service", "rabbitmq").
		Str("event", "message_retry").
		Str("queue", queue).
		Int("retry_count", retryCount).
		Msg("retrying message")

	err := ch.Publish(
		"",
		queue+".retry",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg.Body,
			Headers: amqp.Table{
				"x-retry": retryCount,
			},
		},
	)

	duration := time.Since(start).Milliseconds()

	if err != nil {

		logger.Error(err).
			Str("service", "rabbitmq").
			Str("event", "retry_publish_failed").
			Str("queue", queue).
			Int("retry_count", retryCount).
			Int64("duration_ms", duration).
			Msg("failed publish retry message")

		return err
	}

	logger.Info().
		Str("service", "rabbitmq").
		Str("event", "retry_publish_success").
		Str("queue", queue).
		Int("retry_count", retryCount).
		Int64("duration_ms", duration).
		Msg("retry message published")

	return nil
}
