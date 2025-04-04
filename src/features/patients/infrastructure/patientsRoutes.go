package infrastructure 

import (
	"apiconsumer/src/features/patients/infrastructure/controllers"

	"github.com/gin-gonic/gin"
)

type PatientsRoutes struct {
	router *gin.RouterGroup
	processController *controllers.ProcessPatientsController
}

func NewPatientsRoutes(router *gin.RouterGroup, processController *controllers.ProcessPatientsController) *PatientsRoutes {
	return &PatientsRoutes{
		router: router,
		processController: processController,
	}
}

func (r *PatientsRoutes) Run() {
	PatientsRoutes := r.router.Group("/patients")
	{
		PatientsRoutes.POST("/", r.processController.Run)
	}
	}