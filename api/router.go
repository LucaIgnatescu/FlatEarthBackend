package api

import (
	"net/http"
)

func CreateRouter() *http.ServeMux {
	router := http.NewServeMux()
	return router
}
