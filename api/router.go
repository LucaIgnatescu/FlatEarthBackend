package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
)

type RegisterResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		log.Println("Could not connect to db:", err)
		w.WriteHeader(500)
		return
	}

	data, err := getData()
	if err != nil {
		log.Println("Could not retrieve geo data:", err)
		data = &GeoData{}
	}

	user, err := insertNewUser(db, data)

	if err != nil {
		log.Println("Could not insert record:", err)
		w.WriteHeader(500)
		return
	}
	token, err := generateToken(user.UserID.String())
	if err != nil {
		log.Println("Could not generate token:", err)
		w.WriteHeader(500)
		return
	}

	response, err := json.Marshal(RegisterResponse{"success", token})
	if err != nil {
		log.Println("Could not construct response:", err)
		w.WriteHeader(500)
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

func LogEvent(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*UserClaims)
	if !ok || claims == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID := claims.UserID

	var interaction Interaction
	body, err := io.ReadAll(r.Body)
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
	db, err := connectDB()
	if err != nil {
		log.Println("Could not connect to db:", err)
		w.WriteHeader(500)
		return
	}
	_, err = insertEvent(db, &interaction)

	if err != nil {
		log.Println("Could not insert interaction record:", err)
		w.WriteHeader(500)
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

func LogReport(w http.ResponseWriter, r *http.Request) {
	const reportLength = 2000
	const emailLength = 30
	claims, ok := r.Context().Value("claims").(*UserClaims)
	if !ok || claims == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
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

	db, err := connectDB()
	if err != nil {
		log.Println("Could not connect to db:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = insertBugReport(db, &payload, userUUID)
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

func LogSurvey1(w http.ResponseWriter, r *http.Request) {
	const genderLength = 30
	claims, ok := r.Context().Value("claims").(*UserClaims)
	if !ok || claims == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
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

	db, err := connectDB()
	if err != nil {
		log.Println("Could not connect to db:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = insertSurvey1(db, &payload, userUUID)
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

func LogSurvey2(w http.ResponseWriter, r *http.Request) {
	const infoLength = 1000
	claims, ok := r.Context().Value("claims").(*UserClaims)
	if !ok || claims == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
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

	db, err := connectDB()
	if err != nil {
		log.Println("Could not connect to db:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = insertSurvey2(db, &payload, userUUID)
	if err != nil {
		log.Println("Could not insert bug report:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func CreateRouter() http.Handler {
	router := http.NewServeMux()
	logRouter := http.NewServeMux()

	logRouter.HandleFunc("POST /report", LogReport)
	logRouter.HandleFunc("POST /event", LogEvent)
	logRouter.HandleFunc("POST /survey1", LogSurvey1)
	logRouter.HandleFunc("POST /survey2", LogSurvey2)

	wrappedLogRouter := ApplyMiddleware(logRouter, AuthMiddleware, RateLimitMiddleware)

	router.Handle("/log/", http.StripPrefix("/log", wrappedLogRouter))
	router.HandleFunc("/register", registerUser)

	wrapperRouter := ApplyMiddleware(router, LogMiddleware, CorsMiddleware)

	return wrapperRouter
}
