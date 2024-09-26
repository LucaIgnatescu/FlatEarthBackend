package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"
)

const PORT = ":8080"
const HOST = "http://localhost"

func TestEndpoints(t *testing.T) {
	router := CreateRouter()
	server := http.Server{
		Addr:    PORT,
		Handler: router,
	}

	go func() {
		log.Println("Listening on port", PORT)
		server.ListenAndServe()
	}()

	baseURL := HOST + PORT

	t.Run("TestRegister", func(t *testing.T) {
		r, err := http.Get(baseURL + "/register")
		if err != nil {
			t.Fatal("Could not reach endpoint: ", err)
		}
		body, err := io.ReadAll(r.Body)

		if err != nil {
			t.Fatal("Could not read body: ", err)
		}

		var response RegisterResponse
		err = json.Unmarshal(body, &response)
		if err != nil {

			t.Fatal("Could not decode body: ", err)
		}
	})

	var response RegisterResponse
	t.Run("TestRegister", func(t *testing.T) {
		r, err := http.Get(baseURL + "/register")
		if err != nil {
			t.Fatal("Could not reach endpoint: ", err)
		}

		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)

		if err != nil {
			t.Fatal("Could not read body: ", err)
		}

		err = json.Unmarshal(body, &response)

		if err != nil {
			t.Fatal("Could not decode body: ", err)
		}
	})

	t.Run("TestLog", func(t *testing.T) { // NOTE: Relies on response being properly set by TestRegister
		if response.Token == "" {
			t.Fatal("Token incorrectly returned in previous test")
		}

		req, err := http.NewRequest("POST", baseURL+"/log", nil)
		if err != nil {
			t.Fatal("Could not construct request: ", err)
		}

		req.Header.Set("Authorization", "Bearer "+response.Token)

		res, err := http.DefaultClient.Do(req)

		if err != nil {
			t.Fatal("Could not send post request: ", err)
		}

		if res.StatusCode != http.StatusOK {
			t.Fatal("Received response ", res.StatusCode)
		}
	})

	t.Run("TestUnauthorized", func(t *testing.T) {
		req, err := http.NewRequest("POST", baseURL+"/log", nil)
		if err != nil {
			t.Fatal("Could not construct request: ", err)
		}

		req.Header.Set("Authorization", "Bearer badtoken")

		res, err := http.DefaultClient.Do(req)

		if err != nil {
			t.Fatal("Could not send post request: ", err)
		}

		if res.StatusCode == http.StatusOK {
			t.Fatal("Managed to post with incorrect token")
		}
	})
}
