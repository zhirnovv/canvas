package handlers

import (
	"encoding/json"
	"github.com/zhirnovv/canvas/api/bodyParser"
	"github.com/zhirnovv/canvas/api/error"
	"github.com/zhirnovv/canvas/api/user"
	"net/http"
)

// SignupHandler creates a new user in provided UserStorage.
// TODO: Add proper authentication at this layer. Creating a websocket client without a secure JWT would be impossible if authentication was to be implemented. In this case only userId is required to have a JWT issued.
func SignupHandler(userStorage *user.UserStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, parseErr := bodyParser.ParseJSON(r)

		if parseErr != nil {
			parseErr.WriteTo(w)
			return
		}

		name, nameExists := body["name"].(string)

		if !nameExists {
			noNameError := error.NewAPIError("/signup/noName", http.StatusBadRequest, "Username was not provided", "Please provide a username", r.RequestURI)
			noNameError.WriteTo(w)
			return
		}

		newUser, userCreationError := userStorage.Create(name)

		if userCreationError != nil {
			creationError := error.NewAPIError("/signup/userCreation", http.StatusInternalServerError, "Error while creating user", userCreationError.Error(), r.RequestURI)
			creationError.WriteTo(w)
			return
		}

		newUserJson, marshalErr := json.Marshal(newUser)

		if marshalErr != nil {
			marshalError := error.NewAPIError("/signup/marshalError", http.StatusInternalServerError, "Could not marshal new user to JSON format", marshalErr.Error(), r.RequestURI)
			marshalError.WriteTo(w)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write(newUserJson)
	}
}
