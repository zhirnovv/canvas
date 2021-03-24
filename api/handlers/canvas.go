package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/zhirnovv/canvas/api/jsonResponse"
	"github.com/zhirnovv/canvas/api/middleware"
	"github.com/zhirnovv/canvas/api/socketManager"
	"github.com/zhirnovv/canvas/api/user"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func AddClientHandler(userStorage *user.UserStorage, manager *socketManager.Manager, messenger socketManager.Messenger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response jsonResponse.JSONResponse

		userUUID, uuidIsCorrect := r.Context().Value(middleware.UUIDKey).(uuid.UUID)

		if !uuidIsCorrect {
			response = jsonResponse.NewErrorResponse("/canvas/client/noUserUUID", http.StatusUnauthorized, "No user UUID was found in authentication token", "", r.RequestURI)
			response.WriteTo(w)
			return
		}

		user, userDoesNotExist := userStorage.Read(userUUID)

		if userDoesNotExist != nil {
			response = jsonResponse.NewErrorResponse("/canvas/client/noUser", http.StatusForbidden, "User does not exist", fmt.Sprintf("User with id %s does not exist in userStorage", userUUID), r.RequestURI)
			response.WriteTo(w)
			return
		}

		upgrader.CheckOrigin = func(r *http.Request) bool {
			if r.Header.Get("Origin") == "http://dev.domain.com:8080" {
				return true
			}
			return false
		}

		socket, socketErr := upgrader.Upgrade(w, r, nil)

		if socketErr != nil {
			return
		}

		client := socketManager.NewClient(user, socket, manager, messenger)

		manager.ClientsToAttach <- client

		go client.WriteToSocket()
		go client.ListenToSocket()
	}
}
