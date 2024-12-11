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
	log.Println("Server starting... Connecting to DB")
	db, err := api.ConnectDB()

	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	router := api.CreateRouter(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	log.Printf("Listening on port %s...", port)
	err = server.ListenAndServe()
	log.Fatal(err)
}
