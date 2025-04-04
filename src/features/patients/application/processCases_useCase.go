package application

import (
	"apiconsumer/src/features/patients/domain/entities"
	"apiconsumer/src/features/patients/domain/ports"
	"log"
)

type ProcessPatientsUseCase struct {
    patientsRepository ports.IPatients
	WS                 ports.WS
}

// Cambiar 'type' por 'func' y devolver un puntero a ProcessCasesUseCase
func NewProcessPatientsUseCase(patientsRepository ports.IPatients, ws ports.WS) *ProcessPatientsUseCase {
	return &ProcessPatientsUseCase{
		patientsRepository: patientsRepository,
		WS:                 ws,
	}
}

func (uc *ProcessPatientsUseCase) Run(patients *entities.Patients) error {
	errSearch := uc.patientsRepository.FindID(patients.IDUsuario)
	if errSearch != nil {
		log.Println("Error finding patient ID")
		return errSearch
	}

	// Cambiar log.Println a log.Printf o eliminar el formato
	log.Printf("message: %+v", patients)

	errSendMessage := uc.WS.SendMessage(patients)
	if errSendMessage != nil {
		log.Printf("Error sending message: %v", errSendMessage)
		return errSendMessage
	}

	return nil

}