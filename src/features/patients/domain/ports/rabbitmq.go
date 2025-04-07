package ports

import (
	"apiconsumer/src/features/patients/domain/entities"
	"github.com/streadway/amqp" // Añade este import
)

type RabbitMQ interface {
	Publish(patient *entities.Patients) error
	Consume() (<-chan amqp.Delivery, error)
}