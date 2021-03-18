package socketManager

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/zhirnovv/gochat/api/user"
)

// Client represents a user connected to a certain manager via a websocket
type Client struct {
	UserData        *user.User
	assignedSocket  *websocket.Conn
	assignedManager *Manager
	MessagesToSend  chan []byte // Messages to forward to client.
}

func NewClient(userData *user.User, socket *websocket.Conn, manager *Manager) *Client {
	return &Client{UserData: userData, assignedSocket: socket, assignedManager: manager, MessagesToSend: make(chan []byte)}
}

func (client *Client) ListenToSocket() {
	defer func() {
		client.assignedManager.ClientsToDetach <- client
		client.Kill()
	}()

	for {
		_, data, err := client.assignedSocket.ReadMessage()
		if err != nil {
			// TODO: error handling here
			break
		}

		var dataJSON map[string]interface{}
		parseErr := json.Unmarshal(data, &dataJSON)

		if parseErr != nil {
			// TODO: error handling here
			break
		}

		messageType, messageTypeExists := dataJSON["type"].(string)

		if !messageTypeExists {
			break
		}

		messageData, messageDataExists := dataJSON["data"]

		if !messageDataExists {
			break
		}

		client.assignedManager.MessagesToBroadcast <- NewMessage(messageType, messageData)
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
