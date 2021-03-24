package socketManager

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/zhirnovv/canvas/api/bodyParser"
	"github.com/zhirnovv/canvas/api/user"
)

// Client represents a user connected to a certain manager via a websocket
type Client struct {
	ID                uuid.UUID
	UserData          *user.User
	assignedSocket    *websocket.Conn
	assignedManager   *Manager
	assignedMessenger Messenger
	MessagesToSend    chan []byte
}

func NewClient(userData *user.User, socket *websocket.Conn, manager *Manager, messenger Messenger) *Client {
	return &Client{
		ID:                uuid.New(),
		UserData:          userData,
		assignedSocket:    socket,
		assignedManager:   manager,
		assignedMessenger: messenger,
		MessagesToSend:    make(chan []byte),
	}
}

func (client *Client) ListenToSocket() {
	defer func() {
		client.assignedManager.ClientsToDetach <- client
		client.Kill()
	}()

	for {
		_, data, err := client.assignedSocket.ReadMessage()

		if err != nil {
			errorMessage := NewMessage("/error/message/failedToRead", map[string]interface{}{
				"error": "Could not read message.",
			})

			client.assignedManager.sendMessage(client.ID, errorMessage)
			break
		}

		var unmarshalledMessage map[string]interface{}
		parseErr := json.Unmarshal(data, &unmarshalledMessage)

		if parseErr != nil {
			errorMessage := NewMessage("/error/message/invalidJson", map[string]interface{}{
				"error": parseErr.Error(),
			})
			client.assignedManager.sendMessage(client.ID, errorMessage)

			break
		}

		var decodedMessage Message
		decodeErr := bodyParser.DecodeAndValidate(unmarshalledMessage, &decodedMessage)

		if decodeErr != nil {
			errorMessage := NewMessage("/error/message/invalidJson", map[string]interface{}{
				"error": parseErr.Error(),
			})
			client.assignedManager.sendMessage(client.ID, errorMessage)
		}

		parsedMessage, isError, shouldBroadcast := client.assignedMessenger.Parse(decodedMessage)

		if isError {
			client.assignedManager.sendMessage(client.ID, parsedMessage)
		} else if shouldBroadcast {
			client.assignedManager.MessagesToBroadcast <- parsedMessage
		} else {
			client.assignedManager.sendMessage(client.ID, parsedMessage)
		}

	}
}

func (client *Client) WriteToSocket() {
	for {
		select {
		case message, ok := <-client.MessagesToSend:
			if !ok {
				client.assignedSocket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := client.assignedSocket.WriteMessage(websocket.TextMessage, message)

			if err != nil {
				return
			}
		}
	}
}

func (client *Client) Kill() error {
	err := client.assignedSocket.Close()
	close(client.MessagesToSend)

	if err != nil {
		return err
	}

	return nil
}
