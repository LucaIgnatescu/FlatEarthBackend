package api

import "testing"

func TestGeoData(t *testing.T) {
	_, err := getData()
	if err != nil {
		t.Fatal(err)
	}
	// TODO: Rewrite this test
}
