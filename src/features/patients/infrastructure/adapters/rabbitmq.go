package adapters

import (
	"encoding/json"
	"os"

	"apiconsumer/src/features/patients/domain/entities"
	"github.com/streadway/amqp" // Import esencial
)

type RabbitMQAdapter struct {
	channel *amqp.Channel
}

func NewRabbitMQAdapter(ch *amqp.Channel) *RabbitMQAdapter {
	return &RabbitMQAdapter{channel: ch}
}

func (r *RabbitMQAdapter) Publish(patient *entities.Patients) error {
	body, err := json.Marshal(patient)
	if err != nil {
		return err
	}

	return r.channel.Publish(
		os.Getenv("RABBITMQ_EXCHANGE"),
		os.Getenv("RABBITMQ_ROUTING_KEY"),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

func (r *RabbitMQAdapter) Consume() (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		os.Getenv("RABBITMQ_QUEUE"),
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}