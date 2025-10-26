package mq

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// NewConnection establishes a RabbitMQ connection using the provided URI.
// When uri is empty a local default is used.
func NewConnection(uri string) (*amqp.Connection, error) {
	if uri == "" {
		uri = "amqp://guest:guest@localhost:5672/"
	}

	cfg := amqp.Config{
		Dial: amqp.DefaultDial(10 * time.Second),
	}

	return amqp.DialConfig(uri, cfg)
}
