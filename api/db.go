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
	UserID     uuid.UUID
	RegionName string
	Country    string
	City       string
	Lat        float32
	Lon        float32
	JoinedAt   string
}

type InteractionRecord struct {
	EventID   uuid.UUID
	UserID    uuid.UUID
	EventType string
	SentAt    string
	Payload   interface{}
}

func ConnectDB() (*sql.DB, error) {
	dbpwd := os.Getenv("DB_PWD")
	dbuser := os.Getenv("DB_USER")
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	if dbname == "" {
		dbname = "postgres"
	}

	if dbpwd == "" || dbuser == "" || dbhost == "" || dbport == "" {
		return nil, errors.New("Error loading all environment variables")
	}
	connStr := fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=require binary_parameters=yes",
		dbuser, dbpwd, dbhost, dbport, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func insertNewUser(db *sql.DB, data *GeoData) (*UserRecord, error) {
	if data == nil {
		return nil, errors.New("Data should not be an empty object")
	}

	query := `
  INSERT INTO users (user_id, region_name, country, city, lat, lon, joined_at) VALUES
    (DEFAULT, $1, $2, $3, $4, $5, DEFAULT) 
    RETURNING *
  `

	var row UserRecord
	err := db.QueryRow(query, data.RegionName, data.Country, data.City, data.Lat, data.Lon).
		Scan(&row.UserID, &row.RegionName, &row.City, &row.Country, &row.Lat, &row.Lon, &row.JoinedAt)

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
		err = db.QueryRow(query, data.UserID, data.EventType, string(data.Payload)).
			Scan(&row.EventID, &row.UserID, &row.EventType, &row.SentAt, &row.Payload)
	}

	if err != nil {
		return nil, err
	}
	return &row, nil
}

func insertBugReport(db *sql.DB, data *BugPayload, userID uuid.UUID) error {
	if data == nil {
		return errors.New("Bug payload is null")
	}
	query := `
  INSERT INTO bugs (user_id, report, email) VALUES ($1, $2, $3)
  `
	return db.QueryRow(query, userID, data.Report, data.Email).Err()
}

func insertSurvey1(db *sql.DB, data *Survey1Payload, userID uuid.UUID) error {
	if data == nil {
		return errors.New("Bug payload is null")
	}

	query := `
		INSERT INTO survey1 (
			user_id, answer1, answer2, answer3, answer4, answer5,
			answer6, answer7, answer8, answer9, answer10, gender_details
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12
		)
	`
	return db.QueryRow(query, userID,
		data.Answers[0],
		data.Answers[1],
		data.Answers[2],
		data.Answers[3],
		data.Answers[4],
		data.Answers[5],
		data.Answers[6],
		data.Answers[7],
		data.Answers[8],
		data.Answers[9],
		data.Gender).Err()
}

func insertSurvey2(db *sql.DB, data *Survey2Payload, userID uuid.UUID) error {
	if data == nil {
		return errors.New("Bug payload is null")
	}

	query := `
    INSERT INTO survey2 (
      user_id, answer1, answer2, answer3, answer4, extra_info
    ) VALUES (
      $1, $2, $3, $4, $5, $6
    )
  `
	return db.QueryRow(query,
		userID,
		data.Answers[0],
		data.Answers[1],
		data.Answers[2],
		data.Answers[3],
		data.Info).Err()
}
