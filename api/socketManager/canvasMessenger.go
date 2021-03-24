package socketManager

import (
	"fmt"

	"github.com/go-playground/validator"
	"github.com/zhirnovv/canvas/api/bodyParser"
)

// CanvasMessenger stores and validates all canvas-related messages
type CanvasMessenger struct {
	Messages []Message `json:"messages"`
}

func NewCanvasMessenger() *CanvasMessenger {
	return &CanvasMessenger{
		Messages: make([]Message, 0),
	}
}

// OnOpen() passes all previous messages to a new client.
func (messenger *CanvasMessenger) OnOpen() (message Message, isError bool, shouldBroadcast bool) {
	message = NewMessage("/canvas/all", map[string]interface{}{
		"messages": messenger.Messages,
	})

	return message, false, false
}

var validate = validator.New()

// Parse parses, validates and optionally transforms all canvas-related messages.
func (messenger *CanvasMessenger) Parse(message Message) (parsedMessage Message, isError bool, shouldBroadcast bool) {
	switch message.Type {
	case "/canvas/resize":
		{
			message := NewMessage("/canvas/all", map[string]interface{}{
				"messages": messenger.Messages,
			})

			return message, false, false
		}
	case "/canvas/add/line":
		{
			var parsedMessagePayload AddLineMessagePayload
			err := bodyParser.DecodeAndValidate(message.Payload, &parsedMessagePayload)

			if err != nil {
				errorMessage := NewMessage("/error/canvas/add/line/incorrectPayload", map[string]interface{}{
					"error": fmt.Sprintf("Incorrect payload format: %s", err.Error()),
				})
				return errorMessage, true, false
			}

			messenger.Messages = append(messenger.Messages, message)

			return message, false, true
		}
	default:
		{
			errorMessage := NewMessage("/error/parse/type", map[string]interface{}{
				"error": "Incorrect message type.",
			})

			return errorMessage, true, false
		}
	}
}

type Coordinates struct {
	X float32 `mapstructure:"x" validate:"required,numeric,lte=1.0,gte=0.0"`
	Y float32 `mapstructure:"y" validate:"required,numeric,lte=1.0,gte=0.0"`
}

type LineData struct {
	Start Coordinates `mapstructure:"start" validate:"required"`
	End   Coordinates `mapstructure:"start" validate:"required"`
}

type AddLineMessagePayload struct {
	LineData `mapstructure:"lineData" validate:"required"`
	Stroke   int    `mapstructure:"stroke" validate:"required,gte=1,lte=2"`
	Color    string `mapstructure:"color" validate:"required,hexcolor"`
}
