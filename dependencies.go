package main

import (
	"apiconsumer/src/core/middlewares"

	patientsUseCase "apiconsumer/src/features/patients/application"
	patientsInfrastructure "apiconsumer/src/features/patients/infrastructure"
	patientsAdapter "apiconsumer/src/features/patients/infrastructure/adapters"
	patientsController "apiconsumer/src/features/patients/infrastructure/controllers"

	casesUsecase "apiconsumer/src/features/cases/application"
	casesInfrastructure "apiconsumer/src/features/cases/infrastructure"
	casesAdapter "apiconsumer/src/features/cases/infrastructure/adapters"
	casesController "apiconsumer/src/features/cases/infrastructure/controllers"

	"github.com/gin-gonic/gin"
)

func InitDependencies() {
	engine := gin.Default()
	cors := middlewares.NewCorsMiddleware()
	engine.Use(cors)

	// Configuración de dependencias para `patients`
	patientsMySQL := patientsAdapter.NewMysql()
	patientsWS := patientsAdapter.NewWs()

	processPatientsUseCase := patientsUseCase.NewProcessPatientsUseCase(patientsMySQL, patientsWS)
	processPatientsController := patientsController.NewProcessPatientsController(processPatientsUseCase)
	processPatientsRoutes := patientsInfrastructure.NewPatientsRoutes(engine.Group("/patients"), processPatientsController)

	// Configuración de dependencias para `cases`
	casesMySQL := casesAdapter.NewMysql()
	casesWS := casesAdapter.NewWs()

	ProcessCasesUseCase  := casesUsecase.NewProcessCasesUseCase(casesMySQL, casesWS)
	processMessageController := casesController.NewProcessMessageController( *ProcessCasesUseCase )
	processMessageRoutes := casesInfrastructure.NewCasesRoutes(engine.Group("/cases"),*processMessageController)

	// Registrar rutas
	processPatientsRoutes.Run()
	processMessageRoutes.Run()

	// Iniciar el servidor
	engine.Run(":8082")
}
