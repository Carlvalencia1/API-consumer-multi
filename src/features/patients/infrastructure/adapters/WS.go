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

	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		log.Panicf("Error connecting to WebSocket server: %v", err)
	}

	// Configurar heartbeat
	conn.SetPingHandler(func(appData string) error {
		return conn.WriteControl(
			websocket.PongMessage, 
			[]byte(appData), 
			time.Now().Add(time.Duration(timeoutSec)*time.Second),
		)
	})

	return &WebSocketAdapter{conn: conn}
}

func (ws *WebSocketAdapter) SendMessage(patients *entities.Patients) error {
	message, err := json.Marshal(patients)
	if err != nil {
		log.Println("Error converting message to JSON:", err)
		return err
	}

	err = ws.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("Error sending message:", err)
		return err
	}

	log.Println("Message sent successfully")
	return nil
}

func (ws *WebSocketAdapter) Close() error {
	if ws.conn != nil {
		return ws.conn.Close()
	}
	return nil
}