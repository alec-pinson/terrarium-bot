package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"strings"
	"testing"
)

func TestGetJsonValue(t *testing.T) {
	jsonData := `{
		"person": {
			"name": "Alice",
			"age": 25,
			"address": {
				"city": "New York",
				"state": "NY",
				"country": "USA"
			},
			"email": "alice@example.com"
		}
	}`

	tests := []struct {
		path    string
		want    interface{}
		wantErr bool
	}{
		{
			path:    "person.name",
			want:    "Alice",
			wantErr: false,
		},
		{
			path:    "person.address.city",
			want:    "New York",
			wantErr: false,
		},
		{
			path:    "person.age",
			want:    25.0,
			wantErr: false,
		},
		{
			path:    "person.email",
			want:    "alice@example.com",
			wantErr: false,
		},
		{
			path:    "person.address.invalid_field",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := getJsonValue(jsonData, tt.path)

			if !jsonEqual(got, tt.want) {
				t.Errorf("getJsonValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func jsonEqual(a, b interface{}) bool {
	aBytes, err := json.Marshal(a)
	if err != nil {
		return false
	}

	bBytes, err := json.Marshal(b)
	if err != nil {
		return false
	}

	return string(aBytes) == string(bBytes)
}

func TestDebug(t *testing.T) {
	// Test case 1: debug mode disabled, nothing gets logged
	config.Debug = false
	var buf bytes.Buffer
	log.SetOutput(&buf)
	Debug("test %s %d %v", "one", 2, true)
	if buf.String() != "" {
		t.Errorf("unexpected log output: %q", buf.String())
	}

	// Test case 2: debug mode enabled, log message gets output
	config.Debug = true
	buf.Reset()
	Debug("test %s %d %v", "one", 2, true)
	if got := buf.String(); !strings.Contains(got, "test one 2 true") {
		t.Errorf("unexpected log output: %q", got)
	}

	// reset
	config.Debug = false
	log.SetOutput(os.Stderr)
}
