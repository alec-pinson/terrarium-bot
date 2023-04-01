package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/thedevsaddam/gojsonq"
)

func getJsonValue(json string, path string) interface{} {
	return gojsonq.New().FromString(json).Find(path)
}

func writeResponse(w http.ResponseWriter, response any, Error bool) {
	if Error {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(response)
}

func Debug(s string, v ...any) {
	if config.Debug {
		log.Printf(s, v...)
	}
}
