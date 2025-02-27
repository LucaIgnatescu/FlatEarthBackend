package api

import (
	"encoding/json"
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
	db, err := ConnectDB()
	if err != nil {
		t.Fatal("Could not Connect to db")
	}
	err = db.Ping()
	if err != nil {
		t.Fatal("Could not Connect to db")
	}
}

func TestInserts(t *testing.T) {
	godotenv.Load()
	db, err := ConnectDB()

	if err != nil {
		t.Fatal(err)
	}
	data := GeoData{"A", "A", "A", 0, 0, "A"}
	row, err := insertNewUser(db, &data)
	if err != nil {
		t.Fatal(err)
	}

	if row == nil {
		t.Fatal("Nil row returned")
	}

	if row.Region != data.RegionName || row.Lat != data.Lat || row.Lon != data.Lon {
		t.Fatal("Record not matching provided data")
	}

	interaction := Interaction{"test", nil, row.UserID.String()}

	_, err = insertEvent(db, &interaction)

	if err != nil {
		log.Fatal("Insertion with empty payload failed: ", err)
	}
	interaction = Interaction{"test", json.RawMessage(`{"test": "test"}`), row.UserID.String()}

	_, err = insertEvent(db, &interaction)

	if err != nil {
		log.Fatal("Insertion with nonempty payload failed: ", err)
	}
}

func TestInsertWithoutData(t *testing.T) {
	godotenv.Load()
	db, err := ConnectDB()

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

func TestMalformedInsert(t *testing.T) {
	godotenv.Load()
	db, err := ConnectDB()

	if err != nil {
		t.Fatal(err)
	}
	interaction := Interaction{"test", json.RawMessage("abc"), "malformed-uuid"}

	_, err = insertEvent(db, &interaction)
	if err == nil {
		t.Fatal("Malformed input inserted into db")
	}
}
