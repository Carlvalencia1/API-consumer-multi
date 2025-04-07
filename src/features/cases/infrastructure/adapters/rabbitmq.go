package adapters

import (
	"encoding/json"
	"os"

	"apiconsumer/src/features/cases/domain/entities"
	"github.com/streadway/amqp"
)

type RabbitMQAdapter struct {
	channel *amqp.Channel
}

func NewRabbitMQAdapter(ch *amqp.Channel) *RabbitMQAdapter {
	return &RabbitMQAdapter{channel: ch}
}

func (r *RabbitMQAdapter) Publish(c *entities.MedicalCase) error {
	body, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return r.channel.Publish(
		os.Getenv("RABBITMQ_EXCHANGE_CASES"),     // Exchange específico para cases
		os.Getenv("RABBITMQ_ROUTING_KEY_CASES"),  // Routing key específica
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

func (r *RabbitMQAdapter) Consume() (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		os.Getenv("RABBITMQ_QUEUE_CASES"),
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
}
