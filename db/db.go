package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func loadEnv() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Couldn't get working directory")
	}

	path := filepath.Join(filepath.Dir(wd), ".env")
	godotenv.Load(path)
}

func connectDb() (*sql.DB, error) {
	connStr := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}
