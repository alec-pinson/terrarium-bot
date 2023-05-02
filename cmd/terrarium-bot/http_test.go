package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestSendRequest(t *testing.T) {
	// Mock http request response
	mockResponse := `{"key1": "value1", "key2": 2, "key3": ["element1", "element2"]}`
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, mockResponse)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("Valid URL", func(t *testing.T) {
		expected := map[string]interface{}{
			"key1": "value1",
			"key2": float64(2),
			"key3": []interface{}{"element1", "element2"},
		}

		// Send request to mocked server URL
		res, respCode, err := SendRequest(server.URL, false, 1, true)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if respCode != 200 {
			t.Fatalf("Expected response code 200, got %v", respCode)
		}

		// Assert response is as expected
		if !reflect.DeepEqual(res, expected) {
			t.Errorf("Unexpected response: got %v, want %v", res, expected)
		}
	})

	t.Run("Invalid URL", func(t *testing.T) {
		_, respCode, err := SendRequest("invalid_url", false, 1, false)
		if err == nil {
			t.Error("Expected error, but got nil")
		}
		if respCode != 0 {
			t.Fatalf("Expected response code 0, got %v", respCode)
		}
	})

	t.Run("Invalid URL with 3 retries", func(t *testing.T) {
		config.Debug = true
		var buf bytes.Buffer
		log.SetOutput(&buf)

		_, respCode, err := SendRequest("invalid_url", false, 3, false)
		if err == nil {
			t.Error("Expected error, but got nil")
		}
		if respCode != 0 {
			t.Fatalf("Expected response code 0, got %v", respCode)
		}

		if got := buf.String(); !strings.Contains(got, "3/3") {
			t.Errorf("Expected 3 retries: %q", got)
		}

		config.Debug = false
		log.SetOutput(os.Stderr)
	})
}
