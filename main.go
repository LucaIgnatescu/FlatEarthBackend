package main

import (
	"log"
	"net/http"

	"github.com/LucaIgnatescu/FlatEarthBackend/api"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	router := api.CreateRouter()
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Println("Listening on port 8080...")
	server.ListenAndServe()
}
