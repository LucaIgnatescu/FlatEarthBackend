package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/LucaIgnatescu/FlatEarthBackend/api"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/joho/godotenv"
)

func runServer(router http.Handler) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Listening on port %s...", port)
	err := server.ListenAndServe()
	log.Fatal(err)
}

func runLambda(router http.Handler) {
	lambda.Start(httpadapter.New(router).ProxyWithContext)
}

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
	deployMode := os.Getenv("DEPLOY_MODE")

	if deployMode == "serverless" {
		runLambda(router)
	} else {
		runServer(router)
	}
}
