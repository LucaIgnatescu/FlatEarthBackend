package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type GeoData struct {
	RegionName string  `json:"regionName"`
	Country    string  `json:"country"`
	City       string  `json:"string"`
	Lat        float32 `json:"lat"`
	Lon        float32 `json:"lon"`
	Status     string  `json:"status"`
}

func getData(r *http.Request) (*GeoData, error) {
	ip, ok := r.Context().Value("ip").(string)
	if !ok {
		log.Println("Could not retrieve ip")
		return nil, errors.New("invalid ip address")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprint("http://ip-api.com/json/%s", ip), nil)
	if err != nil {
		log.Println("Error creating request")
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(res.Body)

	var geoData GeoData

	err = json.Unmarshal(data, &geoData)
	if err != nil {
		log.Println("Could not unmarshal response")
		return nil, err
	}

	if geoData.Status == "fail" {
		log.Println("Api request failed")
	}

	return &geoData, nil
}
