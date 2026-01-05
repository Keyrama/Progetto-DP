package main

import (
	"html/template"
	"log"
	"net/http"
	"progetto/restaurant/server/database"
	"progetto/restaurant/server/router_mux"
)

func main() {
	// initialize database
	database.InitDatabase("./server/restaurant.db")
	log.Println("Database initialized")

	defer database.CloseDatabase()

	templates, err := template.ParseGlob("server/templates/*.html")
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}

	router_mux.SetTemplates(templates)
	r := router_mux.InitRouter()

	log.Println("Server in esecuzione su http://localhost:8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
