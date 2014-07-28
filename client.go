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
	// UserAgent represents the value of the HTTP user agent
	UserAgent = "M2X/1 (Go net/http)"
)

// Client represents a client for the M2X API (https://m2x.att.com/developer/documentation/overview)
type Client struct {
	APIBase string
	Headers map[string]string
}

// Status represents a status returned by the /status resource
type Status struct {
	API      string `json:"api"`
	Triggers string `json:"triggers"`
}

// APIKey is the key for the API
var APIKey string

// NewClient creates a NewClient for the M2X API
//
//		client := NewClient("<API-KEY>")
func NewClient(apiKey string) *Client {
	m2xClient := &Client{
		APIBase: "http://api-m2x.att.com/v1",
		Headers: make(map[string]string),
	}
	APIKey = apiKey
	return m2xClient
}

// Status gets the status of the M2X client
//
//		result, err := client.Status()
func (c *Client) Status() (*Status, error) {
	result, _, err := get(c.APIBase + "/status")
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
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("X-M2X-KEY", APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
}
