package db

import (
	"os"
	"testing"
)

func TestEnv(t *testing.T) {
	loadEnv()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		t.Errorf("Could not retrieve database url from .env")
	}
}

func TestConn(t *testing.T) {
	loadEnv()
	db, err := connectDb()
	if err != nil {
		t.Errorf("Could not connect to db")
	}
	err = db.Ping()
	if err != nil {
		t.Errorf("Could not connect to db")
	}
}
