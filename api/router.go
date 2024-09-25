package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type SetupReponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

// TODO: More robust failures
func Setup(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		log.Print("Could not connect to db:", err)
		w.WriteHeader(500)
		return
	}

	data, err := getData()
	if err != nil {
		log.Print("Could not retrieve geo data:", err)
		data = &GeoData{}
	}

	user, err := insertNewUser(db, data)

	if err != nil {
		log.Print("Could not insert record:", err)
		w.WriteHeader(500)
		return
	}
	token, err := generateToken(user.UserID.String())
	if err != nil {
		log.Fatal("Could not generate token:", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(SetupReponse{"success", token})
	if err != nil {
		log.Print("Could not construct response:", err)
		w.WriteHeader(500)
		return
	}
	w.Write(response)
}

func CreateRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/setup", Setup)
	return router
}
