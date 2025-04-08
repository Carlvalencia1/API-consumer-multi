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

// NewWs crea y mantiene una conexión WebSocket activa.
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

	var conn *websocket.Conn
	for {
		conn, _, err := dialer.Dial(wsURL, nil)
		if err != nil {
			log.Printf("Error connecting to WebSocket server: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second) // Intentar de nuevo después de 5 segundos
			continue
		}
		break // Salir del bucle si la conexión es exitosa
	}

	// Configurar ping-pong para mantener la conexión viva
	conn.SetPingHandler(func(appData string) error {
		return conn.WriteControl(
			websocket.PongMessage,
			[]byte(appData),
			time.Now().Add(time.Duration(timeoutSec)*time.Second),
		)
	})

	// Iniciar goroutine para mantener la conexión viva y leer mensajes
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading WebSocket message: %v. Attempting to reconnect...", err)
				conn.Close() // Cerrar la conexión actual
				conn = nil   // Establecer como nula para forzar la reconexión
				NewWs()      // Intentar reconectar
				break        // Salir del bucle para reconectar
			}
		}
	}()

	return &WebSocketAdapter{conn: conn}
}

// SendMessage envía un mensaje al servidor WebSocket.
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

// Close cierra la conexión WebSocket.
func (ws *WebSocketAdapter) Close() error {
	if ws.conn != nil {
		return ws.conn.Close()
	}
	return nil
}
