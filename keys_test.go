// Copyright (c) 2014 Jason Goecke
// keys_test.go

package m2x

import (
	"os"
	"testing"
	"time"
)

func TestParseKeys(t *testing.T) {
	data := `
	{
  "keys": [
    {
      "name": "Master Key",
      "key": "a58c276d47d5a3eb34343cd2e6aebaf2",
      "master": true,
      "feed": null,
      "stream": null,
      "expires_at": null,
      "expired": null,
      "permissions": [
        "DELETE",
        "GET",
        "POST",
        "PUT"
      ]
    },
    {
      "name": "Primary Key",
      "key": "e829bd941c0b8b0381b93e98c3d971b5",
      "master": true,
      "feed": null,
      "stream": null,
      "expires_at": null,
      "expired": null,
      "permissions": [
        "DELETE",
        "GET",
        "POST",
        "PUT"
      ]
    }
  ]
}`

	result, _ := parseKeys([]byte(data))

	if result.Keys[0].Name != "Master Key" || result.Keys[1].Name != "Primary Key" {
		t.Errorf("Names did not parse properly")
	}

	if result.Keys[0].Key != "a58c276d47d5a3eb34343cd2e6aebaf2" || result.Keys[1].Key != "e829bd941c0b8b0381b93e98c3d971b5" {
		t.Errorf("Keys did not parse properly")
	}
}

func TestParseKey(t *testing.T) {
	data := `
    {
  "name": "New Key",
  "key": "ab9702fe90c644d953cbe8817dfaa2a0",
  "master": true,
  "feed": null,
  "stream": null,
  "expires_at": null,
  "expired": null,
  "permissions": [
    "GET"
  ]
}`

	key, _ := parseKey([]byte(data))
	if key.Name != "New Key" {
		t.Errorf("Name did not parse properly")
	}
}

func TestListKeys(t *testing.T) {
	client := NewClient(os.Getenv("M2X_API_KEY"))
	result, _ := client.Keys()
	if result.CurrentPage != 0 {
		t.Errorf("Listing the feeds did not work properly")
	}
}

func TestListBadKeyId(t *testing.T) {
	client := NewClient(os.Getenv("M2X_API_KEY"))
	_, errorMessage := client.Key("1234")
	if errorMessage.StatusCode != 404 || errorMessage.Message != "The specified key does not exist" {
		t.Errorf("We did not get the proper error message or code back")
	}
}

func TestCreateAndListAndUpdateAndDeleteKey(t *testing.T) {
	client := NewClient(os.Getenv("M2X_API_KEY"))

	// Create a new key
	key := make(map[string]interface{})
	theTime := time.Now()
	name := "Go Created Key - " + theTime.Format("20060102150405")
	key["name"] = name
	key["permissions"] = [...]string{"GET", "PUT"}
	result, errorMessage := client.CreateKey(key)
	if result.Name != name {
		t.Errorf("Did create a new key properly")
	}
	id := result.Key

	// List the key
	result, errorMessage = client.Key(id)
	if result.Name != name {
		t.Errorf("Did not fetch key properly")
	}

	// Update a key
	name = "Go key - " + theTime.Format("20060102150405")
	key["name"] = name
	errorMessage = client.UpdateKey(id, key)
	if errorMessage != nil {
		t.Errorf("Did not update key properly")
	}

	// Delete the key
	errorMessage = client.DeleteKey(id)
	if errorMessage != nil {
		t.Errorf("Did not delete key properly")
	}
}
