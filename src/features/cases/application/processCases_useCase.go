package application

import (
	"apiconsumer/src/features/cases/domain/entities"
	"apiconsumer/src/features/cases/domain/ports"
	"encoding/json"
	"log"
)

type ProcessCasesUseCase struct {
	casesRepository ports.ICase
	wsRepository    ports.WS
	rabbitMQ        ports.RabbitMQ
}

func NewProcessCasesUseCase(casesRepository ports.ICase, wsRepository ports.WS, rabbitMQ ports.RabbitMQ) *ProcessCasesUseCase {
	return &ProcessCasesUseCase{
		casesRepository: casesRepository,
		wsRepository:    wsRepository,
		rabbitMQ:        rabbitMQ,
	}
}

func (uc *ProcessCasesUseCase) Run(cases *entities.MedicalCase) error {
	_, errSearch := uc.casesRepository.FindID(cases.IDExpediente)
	if errSearch != nil {
		return errSearch
	}

	// Publicar en RabbitMQ
	if err := uc.rabbitMQ.Publish(cases); err != nil {
		log.Printf("Error publicando en RabbitMQ: %v", err)
		return err
	}

	// Enviar por WebSocket
	if err := uc.wsRepository.SendMessage(cases); err != nil {
		log.Printf("Error enviando mensaje por WS: %v", err)
		return err
	}

	return nil
}

func (uc *ProcessCasesUseCase) StartConsumer() {
	msgs, err := uc.rabbitMQ.Consume()
	if err != nil {
		log.Fatalf("Error al consumir mensajes: %v", err)
	}

	go func() {
		for msg := range msgs {
			var medicalCase entities.MedicalCase
			if err := json.Unmarshal(msg.Body, &medicalCase); err != nil {
				log.Printf("Error decodificando mensaje: %v", err)
				continue
			}

			if err := uc.wsRepository.SendMessage(&medicalCase); err != nil {
				log.Printf("Error enviando mensaje por WS: %v", err)
			}
		}
	}()
}
