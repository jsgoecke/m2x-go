// Copyright (c) 2014 Jason Goecke
// triggers_test.go

package m2x

import (
	// "log"
	"os"
	"testing"
	"time"
)

func TestParseTriggers(t *testing.T) {
	data := `
	{ "triggers": [
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
    "updated": "2013-09-10T19:15:00Z" } ] }`

	result, _ := parseTriggers([]byte(data))

	if result.Triggers[0].Name != "high-temperature" || result.Triggers[1].Name != "low-temperature" {
		t.Errorf("Names did not parse properly")
	}

	if result.Triggers[0].ID != "1234" || result.Triggers[1].ID != "1235" {
		t.Errorf("Keys did not parse properly")
	}
}

func TestParseTrigger(t *testing.T) {
	data := `
  { "id": "1234",
  "name": "high-temperature",
  "stream": "temperature",
  "condition": ">",
  "value": "30",
  "callback_url": "http://example.com",
  "url": "/feeds/a4f919d931c265ddd7b76649eac22f7e/triggers/1234",
  "status": "enabled",
  "created": "2013-09-09T19:15:00Z",
  "updated": "2013-09-10T19:15:00Z" }`

	feed, _ := parseTrigger([]byte(data))

	if feed.ID != "1234" {
		t.Errorf("Id did not parse properly")
	}

	if feed.Name != "high-temperature" {
		t.Errorf("Name did not parse properly")
	}

	if feed.Stream != "temperature" {
		t.Errorf("Stream did not parse properly")
	}
}

func TestParseTriggerEvent(t *testing.T) {
	data := `
	{
	    "feed_id": "a65689ce7a9a69291c6ed2deda1affad",
	    "stream": "temperature",
	    "trigger_name": "foobar",
	    "trigger_description": "temperature \u003e 30",
	    "condition": "\u003e",
	    "threshold": "30",
	    "value": 31.5,
	    "at": "2014-01-11T16:14:14Z"
	}`

	triggerEvent, err := ParseTriggerEvent([]byte(data))
	if err != nil || triggerEvent["feed_id"] != "a65689ce7a9a69291c6ed2deda1affad" {
		t.Errorf("Failed to parse trigger event")
	}
}

func TestCreateAndListAndUpdateAndDeleteTrigger(t *testing.T) {
	client := NewClient(os.Getenv("M2X_API_KEY"))

	// Create a new batch
	blueprintData := make(map[string]string)
	theTime := time.Now()
	name := "Go Created Blueprint - " + theTime.Format("20060102150405")
	blueprintData["name"] = name
	blueprintData["description"] = "Unit testing Go lib for M2X"
	blueprintData["visibility"] = "private"
	blueprint, err := client.CreateBlueprint(blueprintData)
	if err != nil || blueprint.Description != "Unit testing Go lib for M2X" || blueprint.Visibility != "private" {
		t.Errorf("Did create a new batch properly")
	}

	streamData := make(map[string]interface{})
	unit := make(map[string]string)
	unit["label"] = "celcius"
	unit["symbol"] = "C"
	streamData["unit"] = unit
	err = client.UpdateFeedStream(blueprint.Feed, "temperature", streamData)
	if err != nil {
		t.Errorf("Did not update the stream properly")
	}

	stream, err := client.FeedStream(blueprint.Feed, "temperature")
	if err != nil || stream.Name != "temperature" {
		t.Errorf("Listing the stream did not work properly")
	}

	triggerData := make(map[string]string)
	triggerData["name"] = "foobar"
	triggerData["stream"] = "temperature"
	triggerData["condition"] = ">"
	triggerData["value"] = "30"
	triggerData["callback_url"] = "http://foobar.com"
	triggerData["status"] = "enabled"
	trigger, _ := client.CreateTrigger(blueprint.Feed, triggerData)
	if trigger.Name != "foobar" || trigger.Stream != "temperature" {
		t.Errorf("Did not create a trigger properly")
	}

	// Update the trigger
	triggerData["name"] = "barfoo"
	triggerData["condition"] = ">"
	triggerData["value"] = "25"
	triggerData["callback_url"] = "http://barfoo.com"
	triggerData["status"] = "disabled"
	err = client.UpdateTrigger(blueprint.Feed, trigger.ID, triggerData)
	trigger, err = client.Trigger(blueprint.Feed, trigger.ID)
	if trigger.Name != "barfoo" || trigger.CallbackURL != "http://barfoo.com" {
		t.Errorf("Did not update a trigger properly")
	}

	err = client.TestTrigger(blueprint.Feed, trigger.ID)
	if err != nil {
		t.Errorf("Did not test trigger properly")
	}

	err = client.DeleteTrigger(blueprint.Feed, trigger.ID)
	if err != nil {
		t.Errorf("Did not delete trigger properly")
	}

	// Delete the blueprint
	err = client.DeleteBlueprint(blueprint.ID)
	if err != nil {
		t.Errorf("Did not delete blueprint properly")
	}
}
