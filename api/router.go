package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type SetupReponse struct {
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

	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(SetupReponse{"success", token})
	if err != nil {
		log.Println("Could not construct response:", err)
		w.WriteHeader(500)
		return
	}
	w.Write(response)
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
	 payload?: number
	}
*/

type Interaction struct {
	EventType string  `json:"event_type"`
	Payload   float32 `json:"payload"`
	UserID    string
}

func logEvent(w http.ResponseWriter, r *http.Request) {
	tokenStr := parseHeader(r)

	if tokenStr == "" {
		w.WriteHeader(401)
		return
	}

	claims, err := parseToken(tokenStr)

	userID := claims.UserID
	log.Println(userID)

	var interaction Interaction
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Could not read body:", err)
		w.WriteHeader(500)
		return
	}
	json.Unmarshal(body, &interaction)
	interaction.UserID = userID
	log.Println(interaction)
}

func CreateRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/setup", registerUser)
	router.HandleFunc("POST /log", logEvent)
	return router
}
