package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zhirnovv/gochat/api/handling"
	"log"
	"net/http"
)

var address = "127.0.0.1:8000"

func main() {

	coreRouter := mux.NewRouter()

	authSubrouter := coreRouter.PathPrefix("/auth").Subrouter()
	authSubrouter.HandleFunc("/", handling.SignupHandler).Methods("POST").Headers("Content-Type", "application/json")

	http.Handle("/", coreRouter)

	server := http.Server{
		Addr:    address,
		Handler: coreRouter,
	}

	log.Fatal(server.ListenAndServe())
	defer fmt.Println("Server started on" + address)
}
