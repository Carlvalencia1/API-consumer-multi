package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"

	"apiconsumer/src/core"
	"apiconsumer/src/core/middlewares"

	// Patients
	patientsUseCase "apiconsumer/src/features/patients/application"
	patientsInfrastructure "apiconsumer/src/features/patients/infrastructure"
	patientsAdapter "apiconsumer/src/features/patients/infrastructure/adapters"
	patientsController "apiconsumer/src/features/patients/infrastructure/controllers"

	// Cases
	casesUseCase "apiconsumer/src/features/cases/application"
	casesInfrastructure "apiconsumer/src/features/cases/infrastructure"
	casesAdapter "apiconsumer/src/features/cases/infrastructure/adapters"
	casesController "apiconsumer/src/features/cases/infrastructure/controllers"
)

func InitDependencies() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error cargando archivo .env: %v", err)
	}

	engine := gin.Default()
	engine.Use(middlewares.NewCorsMiddleware())

	// 🟠 RabbitMQ (compartido para patients y cases)
	rabbitMQ, err := core.NewRabbitMQ()
	if err != nil {
		log.Fatalf("Error al conectar con RabbitMQ: %v", err)
	}
	defer func() {
		if err := rabbitMQ.Conn.Close(); err != nil {
			log.Printf("Error al cerrar conexión RabbitMQ: %v", err)
		}
	}()

	// 🟢 Patients
	patientsMySQL := patientsAdapter.NewMysql()
	patientsWS := patientsAdapter.NewWs()
	patientsRabbit := patientsAdapter.NewRabbitMQAdapter(rabbitMQ.Channel)

	processPatientsUseCase := patientsUseCase.NewProcessPatientsUseCase(
		patientsMySQL,
		patientsWS,
		patientsRabbit,
	)
	go processPatientsUseCase.StartConsumer()

	patientsController := patientsController.NewProcessPatientsController(processPatientsUseCase)
	patientsRoutes := patientsInfrastructure.NewPatientsRoutes(engine.Group(""), patientsController) // 👈 sin doble /patients

	// 🔵 Cases
	casesMySQL := casesAdapter.NewMysql()
	casesWS := casesAdapter.NewWs()
	casesRabbit := casesAdapter.NewRabbitMQAdapter(rabbitMQ.Channel)

	processCasesUseCase := casesUseCase.NewProcessCasesUseCase(
		casesMySQL,
		casesWS,
		casesRabbit,
	)
	go processCasesUseCase.StartConsumer()

	casesController := casesController.NewProcessMessageController(*processCasesUseCase)
	casesRoutes := casesInfrastructure.NewCasesRoutes(engine.Group(""), *casesController) // 👈 sin doble /cases

	// 🚀 Iniciar rutas
	patientsRoutes.Run()
	casesRoutes.Run()

	// Servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Servidor iniciado en :%s", port)
	if err := engine.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}
