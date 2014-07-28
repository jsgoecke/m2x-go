# m2x-go 
[![wercker status](https://app.wercker.com/status/aa442caf92d85b9debbc71618d7b269a "wercker status")](https://app.wercker.com/project/bykey/aa442caf92d85b9debbc71618d7b269a)

Go library for [AT&T's M2X API](https://m2x.att.com). AT&T's M2X is a cloud-based data storage service and management toolset customized for the internet of things.

## Version

0.2

## Installation

	go get github.com/jsgoecke/m2x-go

## Documentation

[M2X @ Godoc.org](http://godoc.org/github.com/jsgoecke/m2x-go)

## Usage

### M2X Client

```go
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
```
### M2X Event Receiver

```go
package main

import (
	"github.com/jsgoecke/m2x-go"
	"encoding/json"
	"github.com/codegangsta/martini"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	m := martini.Classic()
	m.Post("/streamEvent", streamRequestHandler)
	http.ListenAndServe(":3000", m)
}

func streamRequestHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	triggerEvent, err := m2x.ParseTriggerEvent(body)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Received trigger event!")
		jsonData, _ := json.MarshalIndent(triggerEvent, "", "    ")
		log.Println(string(jsonData[:]))
	}
}
```

## Testing

Right now the tests are a combination of unit tests and functional tests. For the functional
tests to run, you will need to set an environment variable 'M2X_API_KEY' with a valid key. Keep in mind that the tests will add and remove elements from your account, and if a tests fail may orphan
the elements.

	go test

## Todo

	1. Mock out net/http requests in order to move functional tests to unit tests.

## License

MIT