package api

//
// import (
// 	"fmt"
// 	"net/http"
// 	"testing"
// )
//
// const PORT = ":8080"
//
// func TestEndpoints(t *testing.T) {
// 	router := CreateRouter()
// 	server := http.Server{
// 		Addr:    PORT,
// 		Handler: router,
// 	}
//
// 	go func() {
// 		fmt.Println("Listening on port", PORT)
// 		server.ListenAndServe()
// 	}()
//
// 	endpoints := []string{"/"}
// 	failed := false
// 	for _, endpoint := range endpoints {
// 		_, err := http.Get("http://localhost" + PORT + endpoint)
// 		if err != nil {
// 			t.Errorf("Endpoint %v failed with error: %v\n", endpoint, err)
// 			failed = true
// 		}
// 	}
// 	if failed {
// 		t.Fatalf("Api test failed")
// 	}
// }
