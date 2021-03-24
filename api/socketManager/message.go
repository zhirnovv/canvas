package socketManager

import (
	"time"
)

// Message represents a message sent through the websocket.
type Message struct {
	Type      string                 `json:"type" validate:"required" mapstructure:"type"` // Unique message type.
	Payload   map[string]interface{} `json:"payload" validate:"required" mapstructure:"payload"`
	CreatedAt time.Time              `json:"createdAt,omitempty" validate:"omitempty" mapstructure:"createdAt"`
}

func NewMessage(messageType string, data map[string]interface{}) Message {
	return Message{Type: messageType, Payload: data, CreatedAt: time.Now()}
}
