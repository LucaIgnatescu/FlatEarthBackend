package api

import (
	"database/sql"
	"errors"
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
	EventID    int8
	UserID     uuid.UUID
	event_type string
	sent_at    string
	payload    int
}

func connectDB() (*sql.DB, error) {
	connStr := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func insertNewUser(db *sql.DB, data *GeoData) (*UserRecord, error) {
	if data == nil {
		return nil, errors.New("data object is null") // TODO: Support this
	}

	query := `
  INSERT INTO users (user_id, region, lat, lon, joined_at) VALUES
    (DEFAULT, $1, $2, $3, DEFAULT) 
    RETURNING *
  `
	var row UserRecord
	err := db.QueryRow(query, data.Region, data.Lat, data.Lon).Scan(&row.UserID, &row.Region, &row.Lat, &row.Lon, &row.JoinedAt)
	if err != nil {
		return nil, err
	}
	return &row, nil
}
