package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
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

func parseHeader(r *http.Request) string {
	split := strings.Split(r.Header.Get("Authorization"), " ")

	if len(split) != 2 || split[0] != "Bearer" {
		return ""
	}

	return split[1]
}

/*
	NOTE: Request should have the following structure

Authorization: Bearer <token>

	Body: {
	 event_type: string
	 payload?: Object
	}
*/

type Interaction struct {
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
	UserID    string
}

func logEvent(w http.ResponseWriter, r *http.Request) {
	tokenStr := parseHeader(r)

	if tokenStr == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, err := parseToken(tokenStr)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
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

	json.Unmarshal(body, &interaction)
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

func CreateRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/register", registerUser)
	router.HandleFunc("POST /log", logEvent)

	wrapperRouter := ApplyMiddleware(router, LogMiddleware, CorsMiddleware)

	return wrapperRouter
}
