package rabbitmq

import "github.com/rabbitmq/amqp091-go"

func NewConnection(url string) (*amqp091.Connection, error) {
	return amqp091.Dial(url)
}
