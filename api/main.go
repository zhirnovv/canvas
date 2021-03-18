package main

import (
	// "github.com/gorilla/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/zhirnovv/gochat/api/handlers"
	"github.com/zhirnovv/gochat/api/socketManager"
	"github.com/zhirnovv/gochat/api/user"
)

var address = "127.0.0.1:8000"
var secretToken = "supersecret"

func main() {
	userStorage := user.NewUserStorage()
	manager := socketManager.NewManager()

	rootRouter := mux.NewRouter()

	authSubrouter := rootRouter.PathPrefix("/auth").Subrouter()
	authSubrouter.HandleFunc("", handlers.SignupHandler(userStorage)).Methods("POST", "OPTIONS")

	canvasSubrouter := rootRouter.PathPrefix("/canvas").Subrouter()
	canvasSubrouter.HandleFunc("/client", handlers.AddClientHandler(userStorage, manager)).Methods("GET", "OPTIONS")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://dev.domain.com:8080"},
		AllowedMethods:   []string{"POST", "PUT", "GET", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	})

	server := http.Server{
		Addr:    address,
		Handler: c.Handler(rootRouter),
	}

	log.Fatal(server.ListenAndServe())
}
