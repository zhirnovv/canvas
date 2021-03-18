package socketManager

import (
	"time"
)

// Message represents a message sent through the websocket.
type Message struct {
	Type      string      `json:"type"` // Unique message type.
	Payload   interface{} `json:"payload"`
	CreatedAt time.Time   `json:"createdAt"`
}

func NewMessage(messageType string, data interface{}) Message {
	return Message{Type: messageType, Payload: data, CreatedAt: time.Now()}
}
