package main

import (
	"log"

	"github.com/thedevsaddam/gojsonq"
)

func getJsonValue(json string, path string) interface{} {
	return gojsonq.New().FromString(json).Find(path)
}

func Debug(s string, v ...any) {
	if config.Debug {
		log.Printf(s, v...)
	}
}
