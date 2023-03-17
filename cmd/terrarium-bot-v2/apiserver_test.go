package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEndpoint(t *testing.T) {
	apiServer := APIServer{}
	healthPaths := []string{"/health/live", "/health/ready"}

	for _, healthPath := range healthPaths {
		req, err := http.NewRequest("GET", healthPath, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(apiServer.Endpoint)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		if rr.Body.String() != "ok" {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "ok")
		}
	}
}
