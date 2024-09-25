package api

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

func generateToken(userID string) (string, error) {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		return "", errors.New("Could not retrieve JWT key")
	}
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "flat-earth-challenge",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)), // NOTE: Challenge should not take more than 3 hours to complete
		},
		UserID: userID,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return s, nil
}

func parseToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		key := os.Getenv("JWT_KEY")
		if key == "" {
			return nil, errors.New("Could not retrieve decode secret")
		}

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected encryption algorithm")
		}

		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); !ok {
		return nil, errors.New("Unexpected claims set")
	} else {
		return claims, nil
	}
}
