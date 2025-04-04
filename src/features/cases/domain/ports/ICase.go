package ports

import "apiconsumer/src/features/cases/domain/entities"

type ICase interface {
	FindID(idExpediente int) (*entities.MedicalCase, error)
}
