package adapters

import (
    "apiconsumer/src/features/cases/domain/entities"
    "encoding/json"
    "log"
    "net/url"

    "github.com/gorilla/websocket"
)

type WS struct {
    conn *websocket.Conn
}

func NewWs() *WS {
    url := url.URL{
        Scheme: "ws",
        Host:   "localhost:8081",
        Path:   "cases/",
    }
    conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)

    if err != nil {
        log.Panicf("Error establishing WebSocket connection: %v", err)
    }

    return &WS{
        conn: conn,
    }
}

func (ws *WS) SendMessage(medicalCase *entities.MedicalCase) error {
    // Convertir la estructura MedicalCase a JSON
    message, err := json.Marshal(medicalCase)
    if err != nil {
        log.Printf("Error converting MedicalCase to JSON: %v", err)
        return err
    }

    // Enviar el mensaje JSON a través del WebSocket
    errMessage := ws.conn.WriteMessage(websocket.TextMessage, message)
    if errMessage != nil {
        log.Printf("Error sending message: %v", errMessage)
        return errMessage
    }

    log.Printf("Message sent successfully!")
    return nil
}