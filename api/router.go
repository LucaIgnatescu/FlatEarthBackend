package api

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type RouterDependencies struct {
	db *sql.DB
}

type RegisterResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

func (app *RouterDependencies) registerUser(w http.ResponseWriter, r *http.Request) {
	db := app.db
	data, err := getData(r)

	if err != nil {
		log.Println("Could not retrieve geo data:", err)
		data = &GeoData{
			RegionName: "Unknown",
			Country:    "Unknown",
			City:       "Unknown",
			Lat:        0,
			Lon:        0,
			Status:     "",
		}
	}

	user, err := insertNewUser(db, data)
	if err != nil {
		log.Println("Could not insert record:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	token, err := generateToken(user.UserID.String())
	if err != nil {
		log.Println("Could not generate token:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(RegisterResponse{"success", token})
	if err != nil {
		log.Println("Could not construct response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(response)
	return
}

type Interaction struct {
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
	UserID    string
}

func (app *RouterDependencies) LogEvent(w http.ResponseWriter, r *http.Request) {
	db := app.db

	claims, ok := r.Context().Value("claims").(*UserClaims)
	if !ok || claims == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID := claims.UserID

	var interaction Interaction

	const maxBodySize = 1 << 20

	limitReader := http.MaxBytesReader(w, r.Body, maxBodySize)
	defer limitReader.Close()

	body, err := io.ReadAll(limitReader)
	if err != nil {
		log.Println("Could not read body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &interaction)
	if err != nil {
		log.Println("Could not unmarshal body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	interaction.UserID = userID

	_, err = insertEvent(db, &interaction)

	if err != nil {
		log.Println("Could not insert interaction record:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

/*
Payload structure:

	{
	  report: string,
	  email: string
	}
*/

type BugPayload struct {
	Report string `json:"report"`
	Email  string `json:"email"`
}

func (app *RouterDependencies) LogReport(w http.ResponseWriter, r *http.Request) {
	const reportLength = 2000
	const emailLength = 30
	claims, ok := r.Context().Value("claims").(*UserClaims)
	if !ok || claims == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	const maxBodySize = 1 << 20

	limitReader := http.MaxBytesReader(w, r.Body, maxBodySize)
	defer limitReader.Close()

	body, err := io.ReadAll(limitReader)
	if err != nil {
		log.Println("Could not read body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID := claims.UserID
	userUUID, err := uuid.Parse(userID)

	if err != nil {
		log.Println("Malformed userid:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var payload BugPayload

	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Println("Could not unmarshal body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(payload.Email) > emailLength || len(payload.Report) > reportLength {
		log.Println("Payload too large:", err)
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	err = insertBugReport(app.db, &payload, userUUID)
	if err != nil {
		log.Println("Could not insert bug report:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

/*
Payload structure:

		{
	    answers: int[10]
	    gender_detail: string
		}
*/

type Survey1Payload struct {
	Answers [10]int `json:"answers"`
	Gender  string  `json:"gender_detail"`
}

func (app *RouterDependencies) LogSurvey1(w http.ResponseWriter, r *http.Request) {
	const genderLength = 30
	claims, ok := r.Context().Value("claims").(*UserClaims)
	if !ok || claims == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	const maxBodySize = 1 << 20

	limitReader := http.MaxBytesReader(w, r.Body, maxBodySize)
	defer limitReader.Close()

	body, err := io.ReadAll(limitReader)
	if err != nil {
		log.Println("Could not read body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID := claims.UserID
	userUUID, err := uuid.Parse(userID)

	if err != nil {
		log.Println("Malformed userid:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var payload Survey1Payload

	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Println("Could not unmarshal body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(payload.Gender) > genderLength {
		log.Println("Payload too large:", err)
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	err = insertSurvey1(app.db, &payload, userUUID)
	if err != nil {
		log.Println("Could not insert bug report:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

/*
Payload structure:

		{
	    answers: int[4]
	    text: string
		}
*/

type Survey2Payload struct {
	Answers [4]int `json:"answers"`
	Info    string `json:"text"`
}

func (app *RouterDependencies) LogSurvey2(w http.ResponseWriter, r *http.Request) {
	const infoLength = 1000
	claims, ok := r.Context().Value("claims").(*UserClaims)
	if !ok || claims == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	const maxBodySize = 1 << 20

	limitReader := http.MaxBytesReader(w, r.Body, maxBodySize)
	defer limitReader.Close()

	body, err := io.ReadAll(limitReader)
	if err != nil {
		log.Println("Could not read body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID := claims.UserID
	userUUID, err := uuid.Parse(userID)

	if err != nil {
		log.Println("Malformed userid:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var payload Survey2Payload

	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Println("Could not unmarshal body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(payload.Info) > infoLength {
		log.Println("Payload too large:", err)
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	err = insertSurvey2(app.db, &payload, userUUID)
	if err != nil {
		log.Println("Could not insert bug report:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (*RouterDependencies) HandleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Lambda!"))
	return
}

func CreateRouter(db *sql.DB) http.Handler {
	router := http.NewServeMux()
	logRouter := http.NewServeMux()
	app := RouterDependencies{db}

	logRouter.HandleFunc("POST /report", app.LogReport)
	logRouter.HandleFunc("POST /event", app.LogEvent)
	logRouter.HandleFunc("POST /survey1", app.LogSurvey1)
	logRouter.HandleFunc("POST /survey2", app.LogSurvey2)

	wrappedLogRouter := ApplyMiddleware(logRouter, AuthMiddleware)

	router.Handle("/log/", http.StripPrefix("/log", wrappedLogRouter))
	router.HandleFunc("/", app.HandleIndex)
	router.HandleFunc("/register", app.registerUser)

	wrapperRouter := ApplyMiddleware(router, IPMiddleware, LogMiddleware, CorsMiddleware)

	return wrapperRouter
}
