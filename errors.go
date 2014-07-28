// Copyright (c) 2014 Jason Goecke
// errors.go

package m2x

import (
	"encoding/json"
	"errors"
)

// ErrorMessage represents an API error message
type ErrorMessage struct {
	Message    string `json:"message"`
	StatusCode int
	Errors     Error `json:"errors"`
	Error      error
}

// Error is the simple error message
type Error struct {
	Name []string `json:"name"`
}

// Generates an error message based on a JSON return from the API
func generateErrorMessage(data []byte, statusCode int) *ErrorMessage {
	errorMessage := &ErrorMessage{}
	json.Unmarshal(data, &errorMessage)
	errorMessage.StatusCode = statusCode
	errorMessage.Error = errors.New(errorMessage.Message)
	return errorMessage
}

// Generates an error message without a JSON return from the API
func simpleErrorMessage(err error, statusCode int) *ErrorMessage {
	errorMessage := &ErrorMessage{}
	errorMessage.Message = err.Error()
	errorMessage.StatusCode = statusCode
	errorMessage.Error = err
	return errorMessage
}

// Parses the JSON from an error message
func parseErrorMessage(data []byte) (*ErrorMessage, error) {
	errorMessage := &ErrorMessage{}
	err := json.Unmarshal(data, &errorMessage)
	if err != nil {
		return nil, err
	}
	return errorMessage, nil
}
