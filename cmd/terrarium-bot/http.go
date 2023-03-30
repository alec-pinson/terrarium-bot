package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

func SendRequest(url string, insecure bool, retries int) (map[string]interface{}, int, error) {
	result := map[string]interface{}{}
	var resp *http.Response
	var err error

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}
	client := &http.Client{Transport: tr}

	// retry x times if an error occurs, sleep 1 second each time
	for i := 0; i < retries; i++ {
		Debug("Request attempt %v/%v", i+1, retries)
		resp, err = client.Get(url)
		if err == nil && resp.StatusCode == 200 {
			continue
		}
		time.Sleep(1 * time.Second)
	}
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
