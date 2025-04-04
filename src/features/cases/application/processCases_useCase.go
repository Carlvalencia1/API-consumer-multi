package application

import (
	"apiconsumer/src/features/cases/domain/entities"
	"apiconsumer/src/features/cases/domain/ports"
)

type ProcessCasesUseCase struct {
	casesRepository ports.ICase
	wsRepository    ports.WS
}

func NewProcessCasesUseCase(casesRepository ports.ICase, wsRepository ports.WS) *ProcessCasesUseCase {
	return &ProcessCasesUseCase{
		casesRepository: casesRepository,
		wsRepository:    wsRepository,
	}
}

func (uc *ProcessCasesUseCase) Run(cases *entities.MedicalCase) error {
	// Buscar el caso médico por ID
	foundCase, errSearch := uc.casesRepository.FindID(cases.IDExpediente)
	if errSearch != nil {
		return errSearch
	}

	// Opcional: puedes usar `foundCase` para realizar validaciones adicionales
	_ = foundCase

	// Enviar el mensaje a través del WebSocket
	errSendMessage := uc.wsRepository.SendMessage(cases)
	if errSendMessage != nil {
		return errSendMessage
	}

	return nil
}
