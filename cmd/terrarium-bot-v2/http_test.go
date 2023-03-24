package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestSendRequest(t *testing.T) {
	t.Run("Valid URL", func(t *testing.T) {
		expected := map[string]interface{}{
			"key1": "value1",
			"key2": float64(2),
			"key3": []interface{}{"element1", "element2"},
		}
		mockResponse := `{"key1": "value1", "key2": 2, "key3": ["element1", "element2"]}`

		// Mock http request response
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, mockResponse)
		})
		server := httptest.NewServer(mux)
		defer server.Close()

		// Send request to mocked server URL
		res, err := SendRequest(server.URL, false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Assert response is as expected
		if !reflect.DeepEqual(res, expected) {
			t.Errorf("Unexpected response: got %v, want %v", res, expected)
		}
	})

	t.Run("Invalid URL", func(t *testing.T) {
		_, err := SendRequest("invalid_url", false)
		if err == nil {
			t.Error("Expected error, but got nil")
		}
	})
}
