package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func SendRequest(url string) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}

	// if non-200 status code or debug enabled then dump out request/response info
	if config.Debug || resp.StatusCode != 200 {
		fmt.Printf("Request\n: %v", url)

		responseDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Response\n: %v\n\n", string(responseDump))
	}

	// decode and return response
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}
