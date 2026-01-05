package main

import (
	"log"
	"net/http"
	"progetto/notification/handler"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/notification", handler.NotificationHandler).Methods("POST")

	log.Println("Notification microservice listening on :8081...")

	log.Fatal(http.ListenAndServe(":8081", r))
}
