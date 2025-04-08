package adapters

import (
	"apiconsumer/src/features/cases/domain/entities"
	"encoding/json"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type WS struct {
	conn *websocket.Conn
}

func NewWs() *WS {
	// 1. Obtener URL del WebSocket desde variables de entorno
	wsURL := os.Getenv("WS_SERVER_URL") // Ej: "ws://3.91.184.130:8081"
	if wsURL == "" {
		log.Panicf("Falta la variable WS_SERVER_URL en el .env")
	}

	// 2. Parsear y validar la URL
	parsedURL, err := url.Parse(wsURL)
	if err != nil {
		log.Panicf("URL de WebSocket inválida: %v", err)
	}

	// 3. Configurar timeout (desde variables de entorno)
	timeoutStr := os.Getenv("WS_TIMEOUT_SECONDS")
	timeout := 10 * time.Second // Valor por defecto
	if timeoutStr != "" {
		if customTimeout, err := time.ParseDuration(timeoutStr + "s"); err == nil {
			timeout = customTimeout
		}
	}

	// 4. Configurar dialer con timeout
	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = timeout

	// 5. Establecer conexión WebSocket
	log.Printf("Conectando a WebSocket en: %s", parsedURL.String())
	conn, _, err := dialer.Dial(parsedURL.String(), nil)
	if err != nil {
		log.Panicf("Error conectando al WebSocket: %v", err)
	}

	log.Println("Conexión WebSocket establecida exitosamente")
	return &WS{
		conn: conn,
	}
}

func (ws *WS) SendMessage(medicalCase *entities.MedicalCase) error {
	// 1. Convertir estructura a JSON
	message, err := json.Marshal(medicalCase)
	if err != nil {
		log.Printf("Error serializando caso médico: %v", err)
		return err
	}

	// 2. Enviar mensaje por WebSocket
	err = ws.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Printf("Error enviando mensaje: %v", err)
		return err
	}

	log.Printf("Mensaje enviado exitosamente: %s", string(message))
	return nil
}

// Opcional: Método para cerrar la conexión
func (ws *WS) Close() error {
	if ws.conn != nil {
		return ws.conn.Close()
	}
	return nil
}