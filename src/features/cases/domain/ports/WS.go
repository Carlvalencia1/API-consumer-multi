package ports

import "apiconsumer/src/features/cases/domain/entities"

type WS interface {
    SendMessage(message *entities.MedicalCase) error
}