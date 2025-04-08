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

// NewWs crea una nueva conexión WebSocket que se mantendrá abierta.
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

	var conn *websocket.Conn
	for {
		log.Printf("Conectando a WebSocket en: %s", parsedURL.String())
		conn, _, err = dialer.Dial(parsedURL.String(), nil)
		if err != nil {
			log.Printf("Error conectando al WebSocket: %v. Reintentando en 5 segundos...", err)
			time.Sleep(5 * time.Second) // Esperar 5 segundos antes de intentar reconectar
			continue // Reintentar la conexión
		}
		break // Salir del bucle si la conexión fue exitosa
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

// Mantener la conexión abierta y leer mensajes de forma continua.
func (ws *WS) KeepAlive() {
	go func() {
		for {
			_, _, err := ws.conn.ReadMessage()
			if err != nil {
				log.Printf("Error al leer mensaje del WebSocket: %v. Intentando reconectar...", err)
				ws.conn.Close() // Cerrar la conexión actual si hubo un error
				ws.conn = nil   // Establecer la conexión como nula para forzar la reconexión
				ws = NewWs()    // Intentar reconectar
			}
		}
	}()
}

// Opcional: Método para cerrar la conexión
func (ws *WS) Close() error {
	if ws.conn != nil {
		return ws.conn.Close()
	}
	return nil
}
