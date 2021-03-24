package handlers

import (
	"net/http"
	"time"

	// "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/zhirnovv/canvas/api/auth"
	"github.com/zhirnovv/canvas/api/bodyParser"
	"github.com/zhirnovv/canvas/api/jsonResponse"
	"github.com/zhirnovv/canvas/api/middleware"
	"github.com/zhirnovv/canvas/api/user"
)

type signupForm struct {
	Name string `mapstructure:"name" validate:"required,gte=3,alphanum"`
}

// SignupHandler creates a new user in provided UserStorage.
// TODO: Add proper authentication at this layer. Creating a websocket client without a secure JWT would be impossible if authentication was to be implemented. In this case only userId is required to have a JWT issued.
func SignupHandler(userStorage *user.UserStorage, authenticator auth.Authenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response jsonResponse.JSONResponse
		var form signupForm

		parseErr := bodyParser.ParseJSON(r, &form)

		if parseErr != nil {
			parseErr.WriteTo(w)
			return
		}

		name := form.Name

		newUser, userCreationError := userStorage.Create(name)

		if userCreationError != nil {
			response = jsonResponse.NewErrorResponse("/signup/userCreation", http.StatusInternalServerError, "Error while creating user", userCreationError.Error(), r.RequestURI)
			response.WriteTo(w)
			return
		}

		userAuthToken, tokenIssueErr := authenticator.Issue(newUser.ID)

		if tokenIssueErr != nil {
			response = jsonResponse.NewErrorResponse("/signup/tokenGeneration", http.StatusInternalServerError, "Failed to create user authentication token", tokenIssueErr.Error(), r.RequestURI)
			response.WriteTo(w)
			return
		}

		authTokenCookie := &http.Cookie{
			Name:     "user_auth_token",
			Value:    userAuthToken,
			Expires:  time.Now().Add(time.Hour * 24),
			Domain:   "dev.domain.com",
			HttpOnly: true,
		}

		http.SetCookie(w, authTokenCookie)

		response = jsonResponse.NewSuccessResponse(
			map[string]interface{}{
				"data": map[string]string{
					"message": "Successfully created user.",
				},
			},
			http.StatusOK,
			r.RequestURI,
		)

		response.WriteTo(w)
	}
}

func AuthorizationHandler(storage *user.UserStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := storage.Read(r.Context().Value(middleware.UUIDKey).(uuid.UUID))

		if err == nil {
			successResponse := jsonResponse.NewSuccessResponse(
				map[string]interface{}{
					"user": map[string]interface{}{
						"name": user.Name,
					},
				},
				http.StatusOK,
				r.RequestURI,
			)

			successResponse.WriteTo(w)
		}
	}
}
