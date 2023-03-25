package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func SendRequest(url string, insecure bool) (map[string]interface{}, int, error) {
	result := map[string]interface{}{}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}
	defer resp.Body.Close()

	// if non-200 status code or debug enabled then dump out request/response info
	if config.Debug || resp.StatusCode != 200 {
		fmt.Printf("Request\n: %v", url)

		responseDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Response\n: %v\n\n", string(responseDump))
	}

	// attempt decode and return response
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&result)
	return result, resp.StatusCode, nil
}
