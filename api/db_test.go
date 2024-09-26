package api

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestEnv(t *testing.T) {
	godotenv.Load()
	dbpwd := os.Getenv("DB_PWD")
	if dbpwd == "" {
		t.Fatalf("Error loading db password")
	}
	dbuser := os.Getenv("DB_USER")
	if dbuser == "" {
		t.Fatalf("Error loading db user")
	}
	dbhost := os.Getenv("DB_HOST")
	if dbhost == "" {
		t.Fatalf("Error loading db host")
	}
	dbport := os.Getenv("DB_PORT")
	if dbport == "" {
		t.Fatalf("Error loading db port")
	}
}

func TestConn(t *testing.T) {
	godotenv.Load()
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
	godotenv.Load()
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
	godotenv.Load()
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

func TestInteraction(t *testing.T) {
	godotenv.Load()
	db, err := connectDB()

	if err != nil {
		t.Fatal(err)
	}

	interaction := Interaction{"test", 0, "0da01d08-7f22-48d6-b8c9-1ca0a683713e"}

	row, err := insertEvent(db, &interaction)
	if err != nil {
		log.Fatal("Insertion faied: ", err)
	}
	t.Log(row)
}
