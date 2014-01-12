// Copyright (c) 2014 Jason Goecke
// feeds.go

package m2x

import (
	"encoding/json"
)

// Represents a collection of feeds resource (https://m2x.att.com/developer/documentation/feed)
type Feeds struct {
	Feeds       []Feed `json:"feeds"`
	Total       int    `json:"total"`
	Pages       int    `json:"pages"`
	Limit       int    `json:"limit"`
	CurrentPage int    `json:"current_page"`
}

// Represents a location
type Location struct {
	Name      string     `json:"name"`
	Latitude  string     `json:"latitude"`
	Longitude string     `json:"longitude"`
	Elevation string     `json:"elevation"`
	Waypoints []Waypoint `json:"waypoints"`
}

// Represents a waypoint
type Waypoint struct {
	Timestamp string `json:"timestamp"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Elevation string `json:"elevation"`
}

// Represents an individual feed
type Feed struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Visibility  string    `json:"visibility"`
	Status      string    `json:"status"`
	Type        string    `json:"type"`
	Tags        []string  `json:"tags"`
	Url         string    `json:"url"`
	Key         string    `json:"key"`
	Created     string    `json:"created"`
	Updated     string    `json:"updated"`
	Location    Location  `json:"location"`
	Streams     []Stream  `json:"streams"`
	Triggers    []Trigger `json:"triggers"`
}

// Represents a tream
type Stream struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Min     string `json:"min"`
	Max     string `json:"max"`
	Unit    Unit   `json:"unit"`
	Url     string `json:"url"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}

// Represents a collection of values
type Values struct {
	Start  string `json:"start"`
	End    string `json:"end"`
	Limit  int    `json:"limit"`
	Values []Value
}

// Represents a value
type Value struct {
	At    string `json:"at"`
	Value string `json:"value"`
}

// Represents a request
type Unit struct {
	Label  string `json:"label"`
	Symbol string `json:"symbol"`
}

// Represents a collection of request
type Requests struct {
	Requests []Request
}

// Represents a request
type Request struct {
	At     string `json:"at"`
	Status int    `json:"status"`
	Method string `json:"method"`
	Path   string `json:"path"`
}

// Gets a list of feeds
//
//		feeds, err := client.Feeds()
func (m2xClient *M2xClient) Feeds() (*Feeds, *ErrorMessage) {
	result, statusCode, err := get(m2xClient.ApiBase + "/feeds")
	if err != nil {
		return nil, simpleErrorMessage(err, statusCode)
	}
	if statusCode == 200 {
		data, err := parseFeeds(result)
		if err != nil {
			return nil, simpleErrorMessage(err, statusCode)
		}
		return data, nil
	}
	return nil, generateErrorMessage(result, statusCode)
}

// Gets a feed
//
//		feed, err := client.Feed("/feeds/1234")
func (m2xClient *M2xClient) Feed(resource string) (*Feed, *ErrorMessage) {
	result, statusCode, err := get(m2xClient.ApiBase + resource)
	if err != nil {
		return nil, simpleErrorMessage(err, statusCode)
	}
	if statusCode == 200 {
		data, _ := parseFeed(result)
		return data, nil
	}
	return nil, generateErrorMessage(result, statusCode)
}

// Gets a feed location
//
//		feed, err := client.FeedLocation("/feeds/1234")
func (m2xClient *M2xClient) FeedLocation(resource string) (*Location, *ErrorMessage) {
	result, statusCode, err := get(m2xClient.ApiBase + resource + "/location")
	if err != nil {
		return nil, simpleErrorMessage(err, statusCode)
	}
	if statusCode == 200 {
		location := &Location{}
		err := json.Unmarshal(result, &location)
		if err != nil {
			return nil, simpleErrorMessage(err, statusCode)
		} else {
			return location, nil
		}
	}
	return nil, generateErrorMessage(result, statusCode)
}

// Create/Update a feed location
//
// 		loc := make(map[string]interface{})
// 		loc["name"] = "Storage Room in Sevilla, Spain"
// 		loc["latitude"] = "37.383055"
// 		loc["longitude"] = "-5.996392"
// 		loc["elevation"] = "5"
// 		err := client.UpdateFeedLocation("/feeds/1234", loc)
func (m2xClient *M2xClient) UpdateFeedLocation(resource string, updateData map[string]interface{}) *ErrorMessage {
	data, err := json.Marshal(updateData)
	if err != nil {
		return simpleErrorMessage(err, 0)
	}
	result, statusCode, putErr := put(m2xClient.ApiBase+resource+"/location", data)
	if putErr != nil {
		return simpleErrorMessage(putErr, statusCode)
	}

	if statusCode == 204 {
		return nil
	}
	return generateErrorMessage(result, statusCode)
}

// List a feed stream
//
//		stream, err := client.FeedStream("/feeds/1234", "temperature")
func (m2xClient *M2xClient) FeedStream(resource string, name string) (*Stream, *ErrorMessage) {
	result, statusCode, err := get(m2xClient.ApiBase + resource + "/streams/" + name)
	if err != nil {
		return nil, simpleErrorMessage(err, statusCode)
	}
	if statusCode == 200 {
		data, err := parseStream(result)
		if err != nil {
			return nil, simpleErrorMessage(err, statusCode)
		}
		return data, nil
	}
	return nil, generateErrorMessage(result, statusCode)
}

// Update a feed stream
//
// 		streamData := make(map[string]interface{})
// 		unit := make(map[string]string)
// 		unit["label"] = "celcius"
// 		unit["symbol"] = "C"
// 		streamData["unit"] = unit
// 		err := client.UpdateFeedStream("/feeds/1234", "temperature", streamData)
func (m2xClient *M2xClient) UpdateFeedStream(resource string, name string, updateData map[string]interface{}) *ErrorMessage {
	data, err := json.Marshal(updateData)
	if err != nil {
		return simpleErrorMessage(err, 0)
	}
	result, statusCode, putErr := put(m2xClient.ApiBase+resource+"/streams/"+name, data)
	if putErr != nil {
		return simpleErrorMessage(putErr, 0)
	}

	if statusCode == 201 {
		return nil
	}
	return generateErrorMessage(result, statusCode)
}

// List the feeds stream values
//
//		values, err := client.FeedStreamValues("/feeds/1234", "temperature")
func (m2xClient *M2xClient) FeedStreamValues(resource string, name string) (*Values, *ErrorMessage) {
	result, statusCode, err := get(m2xClient.ApiBase + resource + "/streams/" + name + "/values")
	if err != nil {
		return nil, simpleErrorMessage(err, statusCode)
	}
	if statusCode == 200 {
		data, err := parseValues(result)
		if err != nil {
			return nil, simpleErrorMessage(err, statusCode)
		}
		return data, nil
	}
	return nil, generateErrorMessage(result, statusCode)
}

// Update feeds stream values
//
// 		values := make(map[string]interface{})
// 		values["values"] = []*m2x.Value{
// 			{"2013-09-09T19:15:00Z", "32"},
// 			{"2013-09-09T19:16:00Z", "28 "},
// 			{"2013-09-09T19:17:00Z", "25"},
// 			{"2013-09-09T19:17:00Z", "40"},
// 		}
// 		err := client.UpdateFeedStreamValues("/feeds/1234", "temperature", values)
func (m2xClient *M2xClient) UpdateFeedStreamValues(resource string, name string, updateData map[string]interface{}) *ErrorMessage {
	data, err := json.Marshal(updateData)
	if err != nil {
		return simpleErrorMessage(err, 0)
	}
	result, statusCode, putErr := post(m2xClient.ApiBase+resource+"/streams/"+name+"/values", data)
	if putErr != nil {
		return simpleErrorMessage(putErr, statusCode)
	}

	if statusCode == 204 {
		return nil
	}
	return generateErrorMessage(result, statusCode)
}

// Request a log
//
//		requests, err := RequestLog("/feeds/1234")
func (m2xClient *M2xClient) RequestLog(resource string) (*Requests, *ErrorMessage) {
	result, statusCode, err := get(m2xClient.ApiBase + resource + "/log")
	if err != nil {
		return nil, simpleErrorMessage(err, statusCode)
	}
	if statusCode == 200 {
		data, err := parseRequests(result)
		if err != nil {
			return nil, simpleErrorMessage(err, statusCode)
		}
		return data, nil
	}
	return nil, generateErrorMessage(result, statusCode)
}

// Parses the Feeds JSON
func parseFeeds(data []byte) (*Feeds, error) {
	feeds := &Feeds{}
	err := json.Unmarshal(data, &feeds)
	if err != nil {
		return nil, err
	}
	return feeds, nil
}

// Parses the Feed JSON
func parseFeed(data []byte) (*Feed, error) {
	feed := &Feed{}
	err := json.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}
	return feed, nil
}

// Parses the Stream JSON
func parseStream(data []byte) (*Stream, error) {
	stream := &Stream{}
	err := json.Unmarshal(data, &stream)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

// Parses the Values JSON
func parseValues(data []byte) (*Values, error) {
	values := &Values{}
	err := json.Unmarshal(data, &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}

// Parses the Requests JSON
func parseRequests(data []byte) (*Requests, error) {
	requests := &Requests{}
	err := json.Unmarshal(data, &requests)
	if err != nil {
		return nil, err
	}
	return requests, nil
}
