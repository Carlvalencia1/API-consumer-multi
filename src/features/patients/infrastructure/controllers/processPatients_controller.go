package controllers

import (
	"apiconsumer/src/features/patients/application"
	"apiconsumer/src/features/patients/domain/entities"
	"log"

	"github.com/gin-gonic/gin"
)

type ProcessPatientsController struct {
	UseCase *application.ProcessCasesUseCase
}

func NewProcessPatientsController(uc *application.ProcessCasesUseCase) *ProcessPatientsController {
	return &ProcessPatientsController{UseCase: uc}
}

func (ctr *ProcessPatientsController) FindByID(ctx *gin.Context) {
	var patients entities.Patients

	if err := ctx.ShouldBindJSON(&patients); err != nil {
		log.Println("Error binding JSON:", err)
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	errSend := ctr.UseCase.Run(&patients)

	log.Println("Error sending message:", errSend)

	if errSend != nil {
		ctx.JSON(500, gin.H{"error": "Failed to process patient"})
		return
	}

}

func (ctr *ProcessPatientsController) Run(ctx *gin.Context) {
	var patients entities.Patients

	if err := ctx.ShouldBindJSON(&patients); err != nil {
		log.Println("Error binding JSON:", err)
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	errSend := ctr.UseCase.Run(&patients)

	log.Println("Error sending message:", errSend)

	if errSend != nil {
		ctx.JSON(500, gin.H{"error": "Failed to process patient"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Patient processed successfully"})
}
