package api

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type UserRecord struct {
	UserID   uuid.UUID
	Region   string
	Lat      float32
	Lon      float32
	JoinedAt string
}

type InteractionRecord struct {
	EventID   int8
	UserID    uuid.UUID
	EventType string
	SentAt    string
	Payload   interface{}
}

func connectDB() (*sql.DB, error) {
	dbpwd := os.Getenv("DB_PWD")
	dbuser := os.Getenv("DB_USER")
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")

	if dbpwd == "" || dbuser == "" || dbhost == "" || dbport == "" {
		return nil, errors.New("Error loading all environment variables")
	}
	connStr := fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=postgres",
		dbuser, dbpwd, dbhost, dbport)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func insertNewUser(db *sql.DB, data *GeoData) (*UserRecord, error) {
	if data == nil {
		return nil, errors.New("Data object is null") // TODO: Support this
	}

	query := `
  INSERT INTO users (user_id, region, lat, lon, joined_at) VALUES
    (DEFAULT, $1, $2, $3, DEFAULT) 
    RETURNING *
  `
	var row UserRecord
	err := db.QueryRow(query, data.Region, data.Lat, data.Lon).
		Scan(&row.UserID, &row.Region, &row.Lat, &row.Lon, &row.JoinedAt)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// NOTE: This function does not check if data.Payload is valid json directly
// If it isn't, postgres will error at the insert
func insertEvent(db *sql.DB, data *Interaction) (*InteractionRecord, error) {
	if data == nil {
		return nil, errors.New("Data object is null")
	}

	var row InteractionRecord
	query := `
  INSERT INTO interactions (event_id, user_id, event_type, sent_at, payload) VALUES
  (DEFAULT, $1, $2, DEFAULT, $3) RETURNING *
  `
	var err error
	if data.Payload == nil {
		err = db.QueryRow(query, data.UserID, data.EventType, nil).
			Scan(&row.EventID, &row.UserID, &row.EventType, &row.SentAt, &row.Payload)
	} else {
		err = db.QueryRow(query, data.UserID, data.EventType, data.Payload).
			Scan(&row.EventID, &row.UserID, &row.EventType, &row.SentAt, &row.Payload)
	}

	if err != nil {
		return nil, err
	}
	return &row, nil
}
