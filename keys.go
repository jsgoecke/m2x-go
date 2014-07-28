// Copyright (c) 2014 Jason Goecke
// keys.go

package m2x

import (
	"encoding/json"
)

// Keys represents a Keys response from the M2X API (https://m2x.att.com/developer/documentation/keys)
type Keys struct {
	Keys        []Key `json:"keys"`
	Total       int   `json:"total"`
	Pages       int   `json:"pages"`
	Limit       int   `json:"limit"`
	CurrentPage int   `json:"current_page"`
}

// Key represents a single Key
type Key struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Key         string   `json:"key"`
	Master      bool     `json:"master"`
	Feed        string   `json:"feed"`
	Stream      string   `json:"stream"`
	ExpiresAt   string   `json:"expires_at"`
	Expired     string   `json:"expired"`
	Permissions []string `json:"permissions"`
}

// CreateKey creates a key
//
// 		keyData := make(map[string]interface{})
// 		name := "Go Created Key"
// 		keyData["name"] = name
// 		keyData["permissions"] = [...]string{"GET", "PUT"}
// 		key, err := client.CreateKey(keyData)
func (c *Client) CreateKey(key map[string]interface{}) (*Key, *ErrorMessage) {
	data, err := json.Marshal(key)
	if err != nil {
		return nil, simpleErrorMessage(err, 0)
	}

	result, statusCode, postErr := post(c.APIBase+"/keys", data)
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

// DeleteKey deletes a key
//
//		err := client.DeleteKey("1234")
func (c *Client) DeleteKey(id string) *ErrorMessage {
	result, statusCode, err := delete(c.APIBase+"/keys", id)
	if err != nil {
		return simpleErrorMessage(err, statusCode)
	}
	if statusCode == 204 {
		return nil
	}
	return generateErrorMessage(result, statusCode)
}

// Keys gets a list of keys from the /keys resource
//
//		keys, err := client.Keys()
func (c *Client) Keys() (*Blueprints, *ErrorMessage) {
	result, statusCode, err := get(c.APIBase + "/keys")
	if err != nil {
		return nil, simpleErrorMessage(err, statusCode)
	}
	data, err := parseBlueprints(result)
	if err != nil {
		return nil, simpleErrorMessage(err, 0)
	}
	return data, nil
}

// Key gets a list of blueprints from the /key resource
//
//		key, err := client.Key()
func (c *Client) Key(id string) (*Key, *ErrorMessage) {
	result, statusCode, err := get(c.APIBase + "/keys/" + id)
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

// UpdateKey updates a key
//
// 		keyData["name"] = "Go key"
// 		err := client.UpdateKey("/feeds/1234", keyData)
func (c *Client) UpdateKey(id string, updateData map[string]interface{}) *ErrorMessage {
	data, err := json.Marshal(updateData)
	if err != nil {
		return simpleErrorMessage(err, 0)
	}
	result, statusCode, postErr := put(c.APIBase+"/keys/"+id, data)
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
