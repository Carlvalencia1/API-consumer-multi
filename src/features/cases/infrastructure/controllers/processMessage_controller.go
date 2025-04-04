package controllers

import (
	"github.com/gin-gonic/gin"
	"apiconsumer/src/features/cases/application"
	"apiconsumer/src/features/cases/domain/entities"
)

type ProcessMessageController struct {
	useCase application.ProcessCasesUseCase
}

func NewProcessMessageController(useCase application.ProcessCasesUseCase) *ProcessMessageController {
	return &ProcessMessageController{
		useCase: useCase,
	}
}

func (c *ProcessMessageController) CreateMessage(ctx *gin.Context) {
	var medicalCase entities.MedicalCase

	// Validación del cuerpo del request
	if err := ctx.ShouldBindJSON(&medicalCase); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Ejecuta el caso de uso y obtiene el resultado
	result := c.useCase.Run(&medicalCase)

	// Verifica si el resultado es nulo
	if result == nil {
		ctx.JSON(500, gin.H{"error": "Processing failed"})
		return
	}

	// Responde con éxito
	ctx.JSON(200, gin.H{"message": "Message sent", "result": result})
}
