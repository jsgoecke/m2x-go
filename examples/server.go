package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/jsgoecke/m2x-go"
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
