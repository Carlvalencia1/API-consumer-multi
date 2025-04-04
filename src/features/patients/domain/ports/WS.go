package ports

import "apiconsumer/src/features/patients/domain/entities"

type WS interface {
	SendMessage(patients *entities.Patients) error
}