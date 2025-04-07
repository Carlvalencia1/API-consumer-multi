package application

import (
	"apiconsumer/src/features/patients/domain/entities"
	"apiconsumer/src/features/patients/domain/ports"
	"log"
	"encoding/json" 
)

type ProcessPatientsUseCase struct {
	patientsRepository ports.IPatients
	WS                 ports.WS
	RabbitMQ           ports.RabbitMQ // Nuevo
}

func NewProcessPatientsUseCase(patientsRepository ports.IPatients, ws ports.WS, rmq ports.RabbitMQ) *ProcessPatientsUseCase {
	return &ProcessPatientsUseCase{
		patientsRepository: patientsRepository,
		WS:                 ws,
		RabbitMQ:           rmq,
	}
}

func (uc *ProcessPatientsUseCase) Run(patients *entities.Patients) error {
	errSearch := uc.patientsRepository.FindID(patients.IDUsuario)
	if errSearch != nil {
		log.Println("Error finding patient ID")
		return errSearch
	}

	// Publicar en RabbitMQ
	if err := uc.RabbitMQ.Publish(patients); err != nil {
		log.Printf("Error publishing to RabbitMQ: %v", err)
		return err
	}

	// Enviar por WebSocket
	if err := uc.WS.SendMessage(patients); err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	return nil
}

// Nuevo método para consumir mensajes
func (uc *ProcessPatientsUseCase) StartConsumer() {
	msgs, err := uc.RabbitMQ.Consume()
	if err != nil {
		log.Fatalf("Error al consumir mensajes: %v", err)
	}

	go func() {
		for msg := range msgs {
			var patient entities.Patients
			if err := json.Unmarshal(msg.Body, &patient); err != nil {
				log.Printf("Error decodificando mensaje: %v", err)
				continue
			}
			// Procesar mensaje recibido
			if err := uc.WS.SendMessage(&patient); err != nil {
				log.Printf("Error enviando mensaje por WS: %v", err)
			}
		}
	}()
}