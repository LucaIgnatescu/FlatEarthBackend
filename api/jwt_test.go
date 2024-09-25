package api

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestJwt(t *testing.T) {
	godotenv.Load()
	userID := "cheg"
	token, err := generateToken(userID)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := parseToken(token)
	if err != nil {
		t.Fatal(err)
	}

	if claims.UserID != userID {
		t.Fatal("Incorrect userid")
	}

}
