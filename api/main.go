package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/zhirnovv/canvas/api/auth"
	"github.com/zhirnovv/canvas/api/handlers"
	"github.com/zhirnovv/canvas/api/middleware"
	"github.com/zhirnovv/canvas/api/socketManager"
	"github.com/zhirnovv/canvas/api/user"
)

var address = "127.0.0.1:8000"
var secretToken = "supersecret"

func main() {
	userStorage := user.NewUserStorage()

	var userAuthenticator auth.Authenticator = auth.NewUserAuthenticator(userStorage, secretToken)

	messenger := socketManager.NewCanvasMessenger()
	manager := socketManager.NewManager()
	go manager.Run()

	rootRouter := mux.NewRouter()

	authSubrouter := rootRouter.PathPrefix("/auth").Subrouter()
	authSubrouter.HandleFunc("", handlers.SignupHandler(userStorage, userAuthenticator)).Methods("POST", "OPTIONS")
	authSubrouter.HandleFunc("", middleware.RequireUserAuthentication(userAuthenticator, handlers.AuthorizationHandler(userStorage))).Methods("GET", "OPTIONS")

	canvasSubrouter := rootRouter.PathPrefix("/canvas").Subrouter()
	canvasSubrouter.HandleFunc("/client", middleware.RequireUserAuthentication(userAuthenticator, handlers.AddClientHandler(userStorage, manager, messenger))).Methods("GET", "OPTIONS")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://dev.domain.com:8080"},
		AllowedMethods:   []string{"POST", "PUT", "GET", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Set-Cookie"},
		AllowCredentials: true,
	})

	server := http.Server{
		Addr:    address,
		Handler: c.Handler(rootRouter),
	}

	log.Fatal(server.ListenAndServe())
}
