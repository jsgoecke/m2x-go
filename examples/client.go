package main

import (
	"github.com/jsgoecke/m2x-go"
	"log"
	"os"
)

func main() {
	// Create a client
	client := m2x.NewClient(os.Getenv("M2X_API_KEY"))

	// Create a blueprint
	blueprintData := make(map[string]string)
	blueprintData["name"] = "Go Blueprint"
	blueprintData["description"] = "A blueprint for the Go lib for M2X"
	blueprintData["visibility"] = "private"

	blueprint, errorMessage := client.CreateBlueprint(blueprintData)
	if errorMessage != nil {
		log.Println(errorMessage)
	}

	// Update a bluprint
	blueprintData["description"] = "A blueprint for the Go lib for AT&T M2X"
	errorMessage = client.UpdateBlueprint(blueprint.Id, blueprintData)

	// Create a stream
	streamData := make(map[string]interface{})
	unit := make(map[string]string)
	unit["label"] = "celcius"
	unit["symbol"] = "C"
	streamData["unit"] = unit
	errorMessage = client.UpdateFeedStream(blueprint.Feed, "temperature", streamData)
	if errorMessage != nil {
		log.Println("Error creating stream")
	}

	//Update location of the feed stream
	loc := make(map[string]interface{})
	loc["name"] = "Storage Room in Sevilla, Spain"
	loc["latitude"] = "37.383055"
	loc["longitude"] = "-5.996392"
	loc["elevation"] = "5"
	errorMessage = client.UpdateFeedLocation(blueprint.Feed, loc)
	if errorMessage != nil {
		log.Println("Error updating location")
	}

	// Create a trigger for the feed
	triggerData := make(map[string]string)
	triggerData["name"] = "foobar"
	triggerData["stream"] = "temperature"
	triggerData["condition"] = ">"
	triggerData["value"] = "30"
	triggerData["callback_url"] = "http://45bad07a.ngrok.com/streamEvent"
	triggerData["status"] = "enabled"
	_, errorMessage = client.CreateTrigger(blueprint.Feed, triggerData)
	if errorMessage != nil {
		log.Println("Error creating trigger")
	}

	// Update stream with data
	values := make(map[string]interface{})
	values["values"] = []*m2x.Value{
		{"2013-09-09T19:15:00Z", "32"},
		{"2013-09-09T19:16:00Z", "28 "},
		{"2013-09-09T19:17:00Z", "25"},
		{"2013-09-09T19:17:00Z", "40"},
	}
	errorMessage = client.UpdateFeedStreamValues(blueprint.Feed, "temperature", values)

	// Delete the blueprint
	client.DeleteBlueprint(blueprint.Id)
}
