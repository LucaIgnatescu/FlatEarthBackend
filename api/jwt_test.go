package api

import (
	"testing"
)

func TestJwt(t *testing.T) {
	loadEnv()
	userID := "cheg"
	token, err := createToken(userID)
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
