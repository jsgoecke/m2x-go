// Copyright (c) 2014 Jason Goecke
// triggers.go

package m2x

import (
	"encoding/json"
)

// Represents a collection of triggers (https://m2x.att.com/developer/documentation/feed#List-Triggers)
type Triggers struct {
	Triggers []Trigger
}

// Represents a trigger
type Trigger struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Stream      string `json:"stream"`
	Condition   string `json:"condition"`
	Value       string `json:"value"`
	CallbackUrl string `json:"callback_url"`
	Url         string `json:"url"`
	Status      string `json:"status"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
}

// Represents a trigger event
type TriggerEvent struct {
	Id          string  `json:"feed_id"`
	Stream      string  `json:"stream"`
	Name        string  `json:"trigger_name"`
	Description string  `json:"trigger_description"`
	Condition   string  `json:"condition"`
	Threshold   string  `json:"threshold"`
	Value       float32 `json:"value"`
	At          string  `json:"at"`
}

// Creates a trigger on a feed stream
//
// 		triggerData := make(map[string]string)
// 		triggerData["name"] = "foobar"
// 		triggerData["stream"] = "temperature"
// 		triggerData["condition"] = ">"
// 		triggerData["value"] = "30"
// 		triggerData["callback_url"] = "http://45bad07a.ngrok.com/streamEvent"
// 		triggerData["status"] = "enabled"
// 		trigger, err := client.CreateTrigger(blueprint.Feed, triggerData)
func (m2xClient *M2xClient) CreateTrigger(resource string, trigger map[string]string) (*Trigger, *ErrorMessage) {
	data, err := json.Marshal(trigger)
	if err != nil {
		return nil, simpleErrorMessage(err, 0)
	}

	result, statusCode, postErr := post(m2xClient.ApiBase+resource+"/triggers", data)
	if postErr != nil {
		return nil, simpleErrorMessage(postErr, statusCode)
	}

	if statusCode == 201 {
		newTrigger := &Trigger{}
		unmarshalErr := json.Unmarshal(result, &newTrigger)
		if unmarshalErr != nil {
			return nil, simpleErrorMessage(err, statusCode)
		}
		return newTrigger, nil
	}
	return nil, generateErrorMessage(result, statusCode)
}

// Deletes a trigger from a feed stream
//
//		err := client.DeleteTrigger("/feeds/1234", "1235")
func (m2xClient *M2xClient) DeleteTrigger(resource string, id string) *ErrorMessage {
	result, statusCode, err := delete(m2xClient.ApiBase+resource+"/triggers/", id)
	if err != nil {
		return simpleErrorMessage(err, statusCode)
	}
	if statusCode == 200 {
		return nil
	}
	return generateErrorMessage(result, statusCode)
}

// Lists a collection of triggers on a feed stream
//
//		triggers, err := client.Triggers()
func (m2xClient *M2xClient) Triggers(resource string) (*Triggers, *ErrorMessage) {
	result, statusCode, err := get(m2xClient.ApiBase + resource + "/triggers")
	if err != nil {
		return nil, simpleErrorMessage(err, statusCode)
	}
	if statusCode == 200 {
		data, err := parseTriggers(result)
		if err != nil {
			return nil, simpleErrorMessage(err, statusCode)
		}
		return data, nil
	}
	return nil, generateErrorMessage(result, statusCode)
}

// Lists a trigger on a feed stream
//
//		trigger, err := client.Trigger("/feeds/1234", "1235")
func (m2xClient *M2xClient) Trigger(resource string, id string) (*Trigger, *ErrorMessage) {
	result, statusCode, err := get(m2xClient.ApiBase + resource + "/triggers/" + id)
	if err != nil {
		return nil, simpleErrorMessage(err, statusCode)
	}
	if statusCode == 200 {
		data, err := parseTrigger(result)
		if err != nil {
			return nil, simpleErrorMessage(err, statusCode)
		}
		return data, nil
	}
	return nil, generateErrorMessage(result, statusCode)
}

// Updates a trigger on a feed stream
//
// 		triggerData["callback_url"] = "http://host.com/streamEvent"
// 		triggerData["status"] = "disabled"
// 		err := client.UpdateTrigger("/feeds/1234", "1235", triggerData)
func (m2xClient *M2xClient) UpdateTrigger(resource string, id string, updateData map[string]string) *ErrorMessage {
	data, err := json.Marshal(updateData)
	if err != nil {
		return simpleErrorMessage(err, 0)
	}
	result, statusCode, postErr := put(m2xClient.ApiBase+resource+"/triggers/"+id, data)
	if postErr != nil {
		return simpleErrorMessage(postErr, statusCode)
	}

	if statusCode == 204 {
		return nil
	}
	return generateErrorMessage(result, statusCode)
}

// Tests a trigger
//
//	err := client.TestTrigger("/feeds/1234", "foobar")
func (m2xClient *M2xClient) TestTrigger(resource string, name string) *ErrorMessage {
	var empty []byte
	result, statusCode, postErr := post(m2xClient.ApiBase+resource+"/triggers/"+name+"/test", empty)
	if postErr != nil {
		return simpleErrorMessage(postErr, statusCode)
	}

	if statusCode == 204 {
		return nil
	}
	return generateErrorMessage(result, statusCode)
}

// Parses the JSON for a collection of triggers
func parseTriggers(data []byte) (*Triggers, error) {
	triggers := &Triggers{}
	err := json.Unmarshal(data, &triggers)
	if err != nil {
		return nil, err
	}
	return triggers, nil
}

// Parses the JSON for a trigger
func parseTrigger(data []byte) (*Trigger, error) {
	trigger := &Trigger{}
	err := json.Unmarshal(data, &trigger)
	if err != nil {
		return nil, err
	}
	return trigger, nil
}

// Parses the JSON for an event returned by a trigger
func ParseTriggerEvent(data []byte) (*TriggerEvent, error) {
	triggerEvent := &TriggerEvent{}
	err := json.Unmarshal(data, &triggerEvent)
	if err != nil {
		return nil, err
	}
	return triggerEvent, nil
}
