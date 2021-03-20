package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/zhirnovv/canvas/api/error"
	"github.com/zhirnovv/canvas/api/socketManager"
	"github.com/zhirnovv/canvas/api/user"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func AddClientHandler(userStorage *user.UserStorage, manager *socketManager.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("Authorization")

		if userId == "" {
			noIdError := error.NewAPIError("/canvas/client/noId", http.StatusBadRequest, "User ID was not provided", "Please provide a user id.", r.RequestURI)
			noIdError.WriteTo(w)
			return
		}

		uuid, uuidErr := uuid.Parse(userId)

		if uuidErr != nil {
			badUUIDErr := error.NewAPIError("/canvas/client/badUuid", http.StatusBadRequest, "UUID is incorrect", "Please provide a valid uuid", r.RequestURI)
			badUUIDErr.WriteTo(w)
			return
		}

		user, userDoesNotExist := userStorage.Read(uuid)

		if userDoesNotExist != nil {
			noUserError := error.NewAPIError("/canvas/client/noUser", http.StatusForbidden, "User does not exist", fmt.Sprintf("User with id %s does not exist in userStorage", userId), r.RequestURI)
			noUserError.WriteTo(w)
			return
		}

		socket, socketErr := upgrader.Upgrade(w, r, nil)

		if socketErr != nil {
			socketCreationErr := error.NewAPIError("/canvas/client/socketErr", http.StatusInternalServerError, "Failed to create socket", socketErr.Error(), r.RequestURI)
			socketCreationErr.WriteTo(w)
			return
		}

		client := socketManager.NewClient(user, socket, manager)

		fmt.Println(client)

		manager.ClientsToAttach <- client

		go client.WriteToSocket()
		go client.ListenToSocket()
	}
}
