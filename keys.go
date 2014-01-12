// Copyright (c) 2014 Jason Goecke
// keys.go

package m2x

import (
	"encoding/json"
)

// Represents a Keys response from the M2X API (https://m2x.att.com/developer/documentation/keys)
type Keys struct {
	Keys        []Key `json:"keys"`
	Total       int   `json:"total"`
	Pages       int   `json:"pages"`
	Limit       int   `json:"limit"`
	CurrentPage int   `json:"current_page"`
}

// Represents a single Key
type Key struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Key         string   `json:"key"`
	Master      bool     `json:"master"`
	Feed        string   `json:"feed"`
	Stream      string   `json:"stream"`
	ExpiresAt   string   `json:"expires_at"`
	Expired     string   `json:"expired"`
	Permissions []string `json:"permissions"`
}

// Creates a key
//
// 		keyData := make(map[string]interface{})
// 		name := "Go Created Key"
// 		keyData["name"] = name
// 		keyData["permissions"] = [...]string{"GET", "PUT"}
// 		key, err := client.CreateKey(keyData)
func (m2xClient *M2xClient) CreateKey(key map[string]interface{}) (*Key, *ErrorMessage) {
	data, err := json.Marshal(key)
	if err != nil {
		return nil, simpleErrorMessage(err, 0)
	}

	result, statusCode, postErr := post(m2xClient.ApiBase+"/keys", data)
	if postErr != nil {
		return nil, simpleErrorMessage(postErr, statusCode)
	}
	if statusCode == 201 {
		newKey := &Key{}
		unmarshalErr := json.Unmarshal(result, &newKey)
		if unmarshalErr != nil {
			return nil, simpleErrorMessage(err, statusCode)
		}
		return newKey, nil
	}
	return nil, generateErrorMessage(result, statusCode)
}

// Deletes a key
//
//		err := client.DeleteKey("1234")
func (m2xClient *M2xClient) DeleteKey(id string) *ErrorMessage {
	result, statusCode, err := delete(m2xClient.ApiBase+"/keys", id)
	if err != nil {
		return simpleErrorMessage(err, statusCode)
	}
	if statusCode == 204 {
		return nil
	}
	return generateErrorMessage(result, statusCode)
}

// Gets a list of keys from the /keys resource
//
//		keys, err := client.Keys()
func (m2xClient *M2xClient) Keys() (*Blueprints, *ErrorMessage) {
	result, statusCode, err := get(m2xClient.ApiBase + "/keys")
	if err != nil {
		return nil, simpleErrorMessage(err, statusCode)
	}
	data, err := parseBlueprints(result)
	if err != nil {
		return nil, simpleErrorMessage(err, 0)
	}
	return data, nil
}

// Gets a list of blueprints from the /key resource
//
//		key, err := client.Key()
func (m2xClient *M2xClient) Key(id string) (*Key, *ErrorMessage) {
	result, statusCode, err := get(m2xClient.ApiBase + "/keys/" + id)
	if err != nil {
		return nil, simpleErrorMessage(err, statusCode)
	}
	if statusCode == 200 {
		data, err := parseKey(result)
		if err != nil {
			return nil, simpleErrorMessage(err, statusCode)
		}
		return data, nil
	}
	return nil, generateErrorMessage(result, statusCode)
}

// Updates a key
//
// 		keyData["name"] = "Go key"
// 		err := client.UpdateKey("/feeds/1234", keyData)
func (m2xClient *M2xClient) UpdateKey(id string, updateData map[string]interface{}) *ErrorMessage {
	data, err := json.Marshal(updateData)
	if err != nil {
		return simpleErrorMessage(err, 0)
	}
	result, statusCode, postErr := put(m2xClient.ApiBase+"/keys/"+id, data)
	if postErr != nil {
		return simpleErrorMessage(postErr, statusCode)
	}

	if statusCode == 204 {
		return nil
	}
	return generateErrorMessage(result, statusCode)
}

// Parses the JSON of keys request into the appropriate struct
func parseKeys(data []byte) (*Keys, error) {
	keys := &Keys{}
	err := json.Unmarshal(data, &keys)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// Parses the JSON of a single key request into the appropriate struct
func parseKey(data []byte) (*Key, error) {
	key := &Key{}
	err := json.Unmarshal(data, &key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
