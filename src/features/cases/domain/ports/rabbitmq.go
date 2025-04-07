package ports

import (
	"apiconsumer/src/features/cases/domain/entities"
	"github.com/streadway/amqp"
)

type RabbitMQ interface {
	Publish(c *entities.MedicalCase) error
	Consume() (<-chan amqp.Delivery, error)
}
