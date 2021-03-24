package socketManager

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

func getClientDoesNotExistError(id uuid.UUID) error {
	return errors.New(fmt.Sprintf("Client with id %s is not attached to the canvas.", id.String()))
}

// Manager represents a socket manager. A socket manager can accept, validate and rebroadcast messages from its clients.
type Manager struct {
	clients             map[uuid.UUID]*Client // Client map by according user ID.
	ClientsToAttach     chan *Client
	ClientsToDetach     chan *Client
	MessagesToBroadcast chan Message
}

// NewManager() is a constructor function for a Manager.
func NewManager() *Manager {
	return &Manager{
		clients:             make(map[uuid.UUID]*Client),
		ClientsToAttach:     make(chan *Client),
		ClientsToDetach:     make(chan *Client),
		MessagesToBroadcast: make(chan Message),
	}
}

// AttachClient() attaches a provided client to the manager.
func (manager *Manager) attachClient(client *Client) {
	manager.clients[client.ID] = client

	message, isError, shouldBroadcast := client.assignedMessenger.OnOpen()

	if isError {
		manager.detachClient(client)
		return
	}

	if shouldBroadcast {
		manager.broadcastMessage(message)
	} else {
		manager.sendMessage(client.ID, message)
	}
}

// DetachClient() detaches a client from the manager.
func (manager *Manager) detachClient(client *Client) error {
	id := client.ID

	_, exists := manager.clients[id]
	if !exists {
		return getClientDoesNotExistError(id)
	}

	delete(manager.clients, id)

	return nil
}

func (manager *Manager) broadcastMessage(message Message) {
	messageJSON, marshalErr := json.Marshal(message)

	if marshalErr != nil {
		return
	}

	for _, client := range manager.clients {
		client.MessagesToSend <- messageJSON
	}
}

func (manager *Manager) sendMessage(clientId uuid.UUID, message Message) {
	client, exists := manager.clients[clientId]

	if !exists {
		return
	}

	messageJSON, marshalErr := json.Marshal(message)

	if marshalErr != nil {
		errorMessage := NewMessage("/error/message/invalid", map[string]interface{}{
			"error": marshalErr,
		})

		errorMessageJSON, _ := json.Marshal(errorMessage)

		client.MessagesToSend <- errorMessageJSON
	}

	client.MessagesToSend <- messageJSON

	return
}

// Run() initializes the manager.
func (manager *Manager) Run() {
	for {
		select {
		case client := <-manager.ClientsToAttach:
			manager.attachClient(client)
		case client := <-manager.ClientsToDetach:
			manager.detachClient(client)
		case messageToBroadcast := <-manager.MessagesToBroadcast:
			manager.broadcastMessage(messageToBroadcast)
		}
	}
}
