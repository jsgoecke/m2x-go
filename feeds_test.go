// Copyright (c) 2014 Jason Goecke
// feeds_test.go

package m2x

import (
	"os"
	"testing"
	"time"
)

func TestParseFeeds(t *testing.T) {
	data := `
	{ "feeds": [
  { "id": "65b89448f954f49e42b746d73b385cbb",
    "name": "Sample Blueprint",
    "description": "Longer description for Sample Blueprint",
    "visibility": "public",
    "status": "active",
    "type": "blueprint",
    "tags": [],
    "url": "/feeds/65b89448f954f49e42b746d73b385cbb",
    "key": "dafbf3c924f027ff22733635c116d06b",
    "location": {},
    "created": "2013-09-01T10:00:00Z",
    "updated": "2013-09-02T10:00:00Z" },
  { "id": "9033bda03e2cad5cb757d024aa4a8462",
    "name": "Another Blueprint",
    "description": "Longer description for Another Blueprint",
    "visibility": "private",
    "status": "active",
    "type": "blueprint",
    "tags": [ "lorem", "ipsum" ],
    "url": "/feeds/9033bda03e2cad5cb757d024aa4a8462",
    "key": "1fb821d61a07a9b99c7cb10db64aead1",
    "location": {
      "name": "Storage Room",
      "latitude": "-37.9788423562422",
      "longitude": "-57.5478776916862",
      "elevation": "5",
      "waypoints": [ { "timestamp": "2013-09-10T19:15:00Z",
                       "latitude": "-37.9788423562422",
                       "longitude": "-57.5478776916862",
                       "elevation": "5" } ] },
    "created": "2013-09-01T10:00:00Z",
    "updated": "2013-09-02T10:00:00Z" } ] }`

	result, _ := parseFeeds([]byte(data))

	if len(result.Feeds) != 2 {
		t.Errorf("Two feeds should have been returned")
	}

	if result.Feeds[0].Id != "65b89448f954f49e42b746d73b385cbb" || result.Feeds[1].Id != "9033bda03e2cad5cb757d024aa4a8462" {
		t.Errorf("Ids did not parse properly")
	}

	if result.Feeds[0].Name != "Sample Blueprint" || result.Feeds[1].Name != "Another Blueprint" {
		t.Errorf("Names did not parse properly")
	}

	if result.Feeds[1].Location.Name != "Storage Room" {
		t.Errorf("Location name did not parse properly")
	}

	if result.Feeds[1].Location.Waypoints[0].Timestamp != "2013-09-10T19:15:00Z" {
		t.Errorf("Waypoints timestamp did not parse properly")
	}
}

func TestParseFeed(t *testing.T) {
	data := `
  { "id": "a4f919d931c265ddd7b76649eac22f7e",
  "name": "Sample Datasource",
  "description": "Sample Datasource",
  "visibility": "public",
  "status": "enabled",
  "type": "datasource",
  "tags": [ "lorem" ],
  "url": "/feeds/a4f919d931c265ddd7b76649eac22f7e",
  "key": "dafbf3c924f027ff22733635c116d06b",
  "location": {
    "name": "Storage Room",
    "latitude": "-37.9788423562422",
    "longitude": "-57.5478776916862",
    "elevation": "5",
    "waypoints": [ { "timestamp": "2013-09-10T19:15:00Z",
                     "latitude": "-37.9788423562422",
                     "longitude": "-57.5478776916862",
                     "elevation": "5" } ] },
  "streams": [
    { "name": "temperature",
      "value": "32",
      "min": "22",
      "max": "34",
      "unit": { "label": "celcius", "symbol": "C" },
      "url": "/feeds/a4f919d931c265ddd7b76649eac22f7e/streams/temperature",
      "created": "2013-09-09T19:15:00Z",
      "updated": "2013-09-10T19:15:00Z" },
    { "name": "humidity",
      "value": "80",
      "min": "40",
      "max": "83",
      "unit": { "label": "percent", "symbol": "%" },
      "url": "/feeds/a4f919d931c265ddd7b76649eac22f7e/streams/humidity",
      "created": "2013-09-09T19:15:00Z",
      "updated": "2013-09-10T19:15:00Z" }
  ],
  "triggers": [
    { "id": "1234",
      "name": "high-temperature",
      "stream": "temperature",
      "condition": ">",
      "value": "30",
      "callback_url": "http://example.com",
      "url": "/feeds/a4f919d931c265ddd7b76649eac22f7e/triggers/high-temperature",
      "status": "enabled",
      "created": "2013-09-09T19:15:00Z",
      "updated": "2013-09-10T19:15:00Z" },
    { "id": "1235",
      "name": "low-temperature",
      "stream": "temperature",
      "condition": "<",
      "value": "5",
      "callback_url": "http://example.com",
      "url": "/feeds/a4f919d931c265ddd7b76649eac22f7e/triggers/low-temperature",
      "status": "enabled",
      "created": "2013-09-09T19:15:00Z",
      "updated": "2013-09-10T19:15:00Z" }
  ],
  "created": "2013-09-10T13:00:00Z",
  "updated": "2013-09-10T14:10:00Z"
}`

	feed, _ := parseFeed([]byte(data))

	if feed.Name != "Sample Datasource" {
		t.Errorf("Name did not parse properly")
	}

	if feed.Streams[0].Name != "temperature" || feed.Streams[1].Name != "humidity" {
		t.Errorf("Name did not parse properly")
	}

	if feed.Triggers[0].Id != "1234" || feed.Triggers[1].Id != "1235" {
		t.Errorf("Triggers ids did not parse properly")
	}
}

func TestListFeeds(t *testing.T) {
	client := NewClient(os.Getenv("M2X_API_KEY"))
	result, _ := client.Feeds()
	if result.CurrentPage != 1 {
		t.Errorf("Listing the feeds did not work properly")
	}
}

func TestListBadFeedId(t *testing.T) {
	client := NewClient(os.Getenv("M2X_API_KEY"))
	_, errorMessage := client.Feed("/feeds/1234")
	if errorMessage.StatusCode != 404 || errorMessage.Message != "The specified feed does not exist" {
		t.Errorf("We did not get the proper error message or code back")
	}
}

func TestFeedLocation(t *testing.T) {
	client := NewClient(os.Getenv("M2X_API_KEY"))

	// Create a new blueprint
	blueprint := make(map[string]string)
	theTime := time.Now()
	name := "Go Created Blueprint - " + theTime.Format("20060102150405")
	blueprint["name"] = name
	blueprint["description"] = "Unit testing Go lib for M2X"
	blueprint["visibility"] = "private"
	result, _ := client.CreateBlueprint(blueprint)
	if result.Description != "Unit testing Go lib for M2X" || result.Visibility != "private" {
		t.Errorf("Did create a new blueprint properly")
	}

	// Create a new location
	loc := make(map[string]interface{})
	loc["name"] = "Storage Room"
	loc["latitude"] = "-37.9788423562422"
	loc["longitude"] = "-57.5478776916862"
	loc["elevation"] = "5"
	errorMessage := client.UpdateFeedLocation(result.Feed, loc)
	if errorMessage != nil {
		t.Errorf("Did not update the location properly")
	}

	// Delete the blueprint
	errorMessage = client.DeleteBlueprint(result.Id)
	if errorMessage != nil {
		t.Errorf("Did not delete feed properly")
	}
}

func TestFeedStream(t *testing.T) {
	client := NewClient(os.Getenv("M2X_API_KEY"))

	// Create a new blueprint
	blueprint := make(map[string]string)
	theTime := time.Now()
	name := "Go Created Blueprint - " + theTime.Format("20060102150405")
	blueprint["name"] = name
	blueprint["description"] = "Unit testing Go lib for M2X"
	blueprint["visibility"] = "private"
	result, _ := client.CreateBlueprint(blueprint)
	if result.Description != "Unit testing Go lib for M2X" || result.Visibility != "private" {
		t.Errorf("Did create a new blueprint properly")
	}

	// Create a new stream
	stream := make(map[string]interface{})
	unit := make(map[string]string)
	unit["label"] = "celcius"
	unit["symbol"] = "C"
	stream["unit"] = unit
	errorMessage := client.UpdateFeedStream(result.Feed, "temperature", stream)
	if errorMessage != nil {
		t.Errorf("Did not update the stream properly")
	}

	response, _ := client.FeedStream(result.Feed, "temperature")
	if response.Name != "temperature" {
		t.Errorf("Listing the stream did not work properly")
	}

	// Update values on a stream
	values := make(map[string]interface{})
	values["values"] = []*Value{
		{"2013-09-09T19:15:00Z", "32"},
		{"2013-09-09T19:16:00Z", "28 "},
		{"2013-09-09T19:17:00Z", "25"},
	}
	errorMessage = client.UpdateFeedStreamValues(result.Feed, "temperature", values)
	if errorMessage != nil {
		t.Errorf("Updating the values did not work")
	}

	// List the values of a stream
	streamValues, _ := client.FeedStreamValues(result.Feed, "temperature")
	if streamValues.Limit != 100 || streamValues.Values[0].Value != "25" {
		t.Errorf("Listing the values did not work")
	}

	// Request the log
	requestResult, errorMessage := client.RequestLog(result.Feed)
	if requestResult.Requests[0].Status != 200 {
		t.Errorf("RequestLog did not work")
	}

	// Delete the blueprint
	errorMessage = client.DeleteBlueprint(result.Id)
	if errorMessage != nil {
		t.Errorf("Did not delete feed properly")
	}
}
