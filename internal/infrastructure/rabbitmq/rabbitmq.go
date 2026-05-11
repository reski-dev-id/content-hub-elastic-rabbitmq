package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQ(url string) (*amqp.Connection, error) {
	return amqp.Dial(url)
}
