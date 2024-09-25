package main

import (
	"net/http"

	"github.com/LucaIgnatescu/FlatEarthBackend/api"
)

func main() {
	router := api.CreateRouter()
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	server.ListenAndServe()
}
