package api

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
)

// TODO: Redo this
func extractIP(r *http.Request) (string, error) {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return strings.Split(forwarded, ",")[0], nil
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", errors.New("Error parsing IP")
	}
	return host, nil
}

type GeoData struct {
	Region string  `json:"region"`
	Lat    float32 `json:"lat"`
	Lon    float32 `json:"lon"`
}

func getData() (*GeoData, error) {
	res, err := http.Get("http://ip-api.com/json/") // TODO: Check timeout timer
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(res.Body)

	var geoData GeoData

	err = json.Unmarshal(data, &geoData)
	if err != nil {
		return nil, err
	}

	return &geoData, nil
}
