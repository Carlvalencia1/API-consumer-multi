package core

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ() (*RabbitMQ, error) {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		return nil, fmt.Errorf("RABBITMQ_URL no definida en .env")
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("error conectando a RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("error creando canal: %v", err)
	}

	// Declarar exchange y cola
	err = ch.ExchangeDeclare(
		os.Getenv("RABBITMQ_EXCHANGE"),
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error declarando exchange: %v", err)
	}

	_, err = ch.QueueDeclare(
		os.Getenv("RABBITMQ_QUEUE"),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error declarando cola: %v", err)
	}

	err = ch.QueueBind(
		os.Getenv("RABBITMQ_QUEUE"),
		os.Getenv("RABBITMQ_ROUTING_KEY"),
		os.Getenv("RABBITMQ_EXCHANGE"),
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error enlazando cola: %v", err)
	}

	return &RabbitMQ{Conn: conn, Channel: ch}, nil
}