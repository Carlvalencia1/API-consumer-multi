package adapters

import (
	"apiconsumer/src/features/patients/domain/entities"
	"encoding/json"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type WebSocketAdapter struct {
	conn *websocket.Conn
}

func NewWs() *WebSocketAdapter {
	url := url.URL{
		Scheme: "ws",
		Host:   "localhost:8081",
		Path:   "/patients",
	}
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.Panicf("Error connecting to WebSocket server: %v", err)
	}
	return &WebSocketAdapter{conn: conn}
}

func (ws *WebSocketAdapter) SendMessage(patients *entities.Patients) error {
	// Convert the Patients struct to JSON
	message, err := json.Marshal(patients)
	if err != nil {
		log.Println("Error converting message to JSON:", err)
		return err
	}

	// Send the JSON message over the WebSocket connection
	err = ws.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("Error sending message:", err)
		return err
	}

	log.Println("Message sent successfully")
	return nil
}
