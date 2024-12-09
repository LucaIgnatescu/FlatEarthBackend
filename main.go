package main

import (
	"log"
	"net/http"
	"os"

	"github.com/LucaIgnatescu/FlatEarthBackend/api"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	router := api.CreateRouter()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	log.Printf("Listening on port %s...", port)
	err := server.ListenAndServe()
	log.Fatal(err)
}
