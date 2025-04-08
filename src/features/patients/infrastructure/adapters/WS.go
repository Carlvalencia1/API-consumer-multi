package adapters

import (
	"apiconsumer/src/features/patients/domain/entities"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketAdapter struct {
	conn *websocket.Conn
}

func NewWs() *WebSocketAdapter {
	// Configurar timeout
	timeoutSec, _ := strconv.Atoi(os.Getenv("WS_TIMEOUT_SECONDS"))
	if timeoutSec == 0 {
		timeoutSec = 10 // Valor por defecto
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: time.Duration(timeoutSec) * time.Second,
	}

	// Construir URL desde variables de entorno
	wsURL := os.Getenv("WS_SERVER_URL") + os.Getenv("WS_PATIENTS_ENDPOINT")

	// Establecer la conexión WebSocket
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		log.Panicf("Error conectando al WebSocket: %v", err)
	}

	// Configurar ping/pong para mantener viva la conexión
	conn.SetPingHandler(func(appData string) error {
		return conn.WriteControl(
			websocket.PongMessage,
			[]byte(appData),
			time.Now().Add(time.Duration(timeoutSec)*time.Second),
		)
	})

	log.Println("Conexión WebSocket establecida exitosamente")
	return &WebSocketAdapter{conn: conn}
}

func (ws *WebSocketAdapter) SendMessage(patients *entities.Patients) error {
	// Convertir el mensaje de pacientes a JSON
	message, err := json.Marshal(patients)
	if err != nil {
		log.Println("Error convirtiendo mensaje a JSON:", err)
		return err
	}

	// Enviar mensaje por WebSocket
	err = ws.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("Error enviando mensaje:", err)
		return err
	}

	log.Println("Mensaje enviado exitosamente")
	return nil
}

func (ws *WebSocketAdapter) Close() error {
	if ws.conn != nil {
		return ws.conn.Close()
	}
	return nil
}
