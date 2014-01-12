// Copyright (c) 2014 Jason Goecke
// client_test.go

package m2x

import (
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient(os.Getenv("M2X_API_KEY"))
	if client.ApiBase != "http://api-m2x.att.com/v1" {
		t.Errorf("ApiBase was not set properly")
	}
}

func TestStatus(t *testing.T) {
	client := NewClient(os.Getenv("M2X_API_KEY"))
	result, err := client.Status()
	if err != nil || result.Status != "OK" {
		t.Errorf("Status did not return OK")
	}
}
