package main

import (
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
	server.ListenAndServe()
}
