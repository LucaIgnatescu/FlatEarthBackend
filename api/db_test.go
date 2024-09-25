package api

import (
	"os"
	"testing"
)

func TestEnv(t *testing.T) {
	loadEnv()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		t.Fatal("Could not retrieve database url from .env")
	}
}

func TestConn(t *testing.T) {
	loadEnv()
	db, err := connectDB()
	if err != nil {
		t.Fatal("Could not connect to db")
	}
	err = db.Ping()
	if err != nil {
		t.Fatal("Could not connect to db")
	}
}

func TestInsertWithData(t *testing.T) {
	loadEnv()
	db, err := connectDB()

	if err != nil {
		t.Fatal(err)
	}
	data := GeoData{"NY", 32.5, 11.1}
	row, err := insertNewUser(db, &data)
	if err != nil {
		t.Fatal(err)
	}

	if row == nil {
		t.Fatal("Nil row returned")
	}

	t.Log(row)
	if row.Region != data.Region || row.Lat != data.Lat || row.Lon != data.Lon {
		t.Fatal("Record not matching provided data")
	}

}

func TestInsertWithoutData(t *testing.T) {
	loadEnv()
	db, err := connectDB()

	if err != nil {
		t.Fatal(err)
	}
	var data *GeoData
	data = nil

	_, err = insertNewUser(db, data)
	if err == nil {
		t.Fatal(err)
	}
}
