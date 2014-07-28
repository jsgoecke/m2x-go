// Copyright (c) 2014 Jason Goecke
// client.go

package m2x

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	USERAGENT = "M2X/1 (Go net/http)"
)

// Represents a client for the M2X API (https://m2x.att.com/developer/documentation/overview)
type M2xClient struct {
	ApiBase string
	Headers map[string]string
}

// Represents a status returned by the /status resource
type Status struct {
	API      string `json:"api"`
	Triggers string `json:"triggers"`
}

var ApiKey string

// Creates a NewClient for the M2X API
//
//		client := NewClient("<API-KEY>")
func NewClient(apiKey string) *M2xClient {
	m2xClient := &M2xClient{}
	m2xClient.ApiBase = "http://api-m2x.att.com/v1"
	m2xClient.Headers = make(map[string]string)
	ApiKey = apiKey
	return m2xClient
}

// Gets the status of the M2X client
//
//		result, err := client.Status()
func (m2xClient *M2xClient) Status() (*Status, error) {
	result, _, err := get(m2xClient.ApiBase + "/status")
	status := &Status{}
	err = json.Unmarshal(result, &status)
	if err != nil {
		return nil, err
	}
	return status, nil
}

// Provides a common facility for doing a DELETE on an M2X API resource
//
//		result, err := delete("http://api-m2x.att.com/v1/feeds", "1234")
func delete(resource string, id string) ([]byte, int, error) {
	httpClient := &http.Client{}
	req, _ := http.NewRequest("DELETE", resource+"/"+id, nil)
	return processRequest(req, httpClient)
}

// Provides a common facility for doing a GET on an M2X API resource
//
//		result, err := get("/status")
func get(resource string) ([]byte, int, error) {
	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", resource, nil)
	return processRequest(req, httpClient)
}

// Provides a common facility for doing a POST on an M2X API resource. Takes
// JSON []byte for the data argument.
//
//		result, err := post("/blueprints", blueprint)
func post(resource string, data []byte) ([]byte, int, error) {
	httpClient := &http.Client{}
	req, _ := http.NewRequest("POST", resource, bytes.NewReader(data))
	return processRequest(req, httpClient)
}

// Provides a common facility for doing a PUT on an M2X API resource. Takes
// JSON []byte for the data argument.
//
//		result, err := put("/blueprints", blueprint)
func put(resource string, data []byte) ([]byte, int, error) {
	httpClient := &http.Client{}
	req, _ := http.NewRequest("PUT", resource, bytes.NewReader(data))
	return processRequest(req, httpClient)
}

func processRequest(req *http.Request, httpClient *http.Client) ([]byte, int, error) {
	setHeaders(req)
	result, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	body, _ := ioutil.ReadAll(result.Body)
	result.Body.Close()
	return body, result.StatusCode, nil
}

// Sets the custom headers required for the M2X API
func setHeaders(req *http.Request) {
	req.Header.Add("User-Agent", USERAGENT)
	req.Header.Add("X-M2X-KEY", ApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
}
