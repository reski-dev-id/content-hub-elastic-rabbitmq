package rabbitmq

import (
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

	retryCount := GetRetryCount(msg)

	if retryCount >= MaxRetry {

		logger.Log.Error().
			Str("queue", queue).
			Int("retry_count", retryCount).
			Msg("message moved to dead letter queue")

		return ch.Publish(
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
	}

	retryCount++

	logger.Log.Warn().
		Str("queue", queue).
		Int("retry_count", retryCount).
		Msg("retrying message")

	return ch.Publish(
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
}
