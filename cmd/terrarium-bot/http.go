package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

var httpClientPool []*HttpClientPool

type HttpClientPool struct {
	Hostname string
	Client   http.Client
}

func getClient(address string, insecure bool) *http.Client {
	url, err := url.Parse(address)
	if err != nil {
		log.Fatal(err)
	}
	hostname := strings.TrimPrefix(url.Hostname(), "www.")

	// if there's an existing client, return it
	for _, pool := range httpClientPool {
		if pool.Hostname == hostname {
			return &pool.Client
		}
	}
	// otherwise create a new client and return that
	Debug("Creating http pool for %s", hostname)
	client := http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}}
	httpClientPool = append(httpClientPool, &HttpClientPool{
		Hostname: hostname,
		Client:   client,
	})
	return &client
}

func SendRequest(url string, insecure bool, retries int, decodeJson bool) (map[string]interface{}, int, error) {
	result := map[string]interface{}{}
	var req *http.Request
	var resp *http.Response
	var err error

	if isTesting {
		// client trace to log whether the request's underlying tcp connection was re-used
		traceCtx := httptrace.WithClientTrace(context.Background(), &httptrace.ClientTrace{
			GotConn: func(info httptrace.GotConnInfo) { log.Printf("conn was reused: %t", info.Reused) },
		})

		req, err = http.NewRequestWithContext(traceCtx, http.MethodGet, url, nil)
	} else {
		req, err = http.NewRequest(http.MethodGet, url, nil)
	}

	// set http insecure mode true/false
	client := getClient(url, insecure)

	// retry x times if an error occurs, sleep 1 second each time
	for i := 0; i < retries; i++ {
		Debug("Request attempt %v/%v", i+1, retries)
		resp, err = client.Do(req)
		if err == nil && resp.StatusCode == 200 {
			break
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

	if decodeJson {
		// attempt decode and return response
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&result); err != nil {
			log.Printf("Request to %v: could not decode response: %v", url, err)
			return nil, resp.StatusCode, err
		}
	} else {
		// need to read the resp body in order for the connection to be reused https://golang.cafe/blog/how-to-reuse-http-connections-in-go.html
		io.Copy(ioutil.Discard, resp.Body)
	}
	return result, resp.StatusCode, nil
}
