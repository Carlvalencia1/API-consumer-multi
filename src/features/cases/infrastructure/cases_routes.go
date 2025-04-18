package infrastructure

import (
	"apiconsumer/src/features/cases/infrastructure/controllers"
	"github.com/gin-gonic/gin"
)

type CasesRoutes struct {
	engine         *gin.RouterGroup
	processMessage controllers.ProcessMessageController
}

func NewCasesRoutes(engine *gin.RouterGroup, processMessage controllers.ProcessMessageController) *CasesRoutes {
	return &CasesRoutes{
		engine:         engine,
		processMessage: processMessage,
	}
}

func (r *CasesRoutes) Run() {
	casesRoutes := r.engine.Group("/cases")
	{
		// Cambié ProcessMessage por CreateMessage para que coincida con tu controlador
		casesRoutes.POST("/", r.processMessage.CreateMessage)
	}
}
