// Copyright (c) 2014 Jason Goecke
// datasources.go

package m2x

import (
	"encoding/json"
)

// Blueprints represents a collection of blueprints (https://m2x.att.com/developer/documentation/datasource)
type Blueprints struct {
	Blueprints  []Blueprint `json:"blueprints"`
	Total       int         `json:"total"`
	Pages       int         `json:"pages"`
	Limit       int         `json:"limit"`
	CurrentPage int         `json:"current_page"`
}

// Blueprint represents a single blueprint
type Blueprint struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Visibility  string     `json:"visibility"`
	Serial      string     `json:"serial"`
	Status      string     `json:"status"`
	Feed        string     `json:"feed"`
	URL         string     `json:"url"`
	Key         string     `json:"key"`
	Tags        []string   `json:"tags"`
	Created     string     `json:"created"`
	Updated     string     `json:"updated"`
	Datasources Datasource `json:"datasources"`
}

// Batches represents a collection of batches
type Batches struct {
	Batches     []Batch `json:"batches"`
	Total       int     `json:"total"`
	Pages       int     `json:"pages"`
	Limit       int     `json:"limit"`
	CurrentPage int     `json:"current_page"`
}

// Batch represents a single batch
type Batch struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Visibility  string     `json:"visibility"`
	Serial      string     `json:"serial"`
	Status      string     `json:"status"`
	Feed        string     `json:"feed"`
	URL         string     `json:"url"`
	Key         string     `json:"key"`
	Tags        []string   `json:"tags"`
	Created     string     `json:"created"`
	Updated     string     `json:"updated"`
	Datasources Datasource `json:"datasources"`
}

// Datasource represents a single datasource
type Datasource struct {
	Total        int `json:"total"`
	Registered   int `json:"registered"`
	Unregistered int `json:"unregistered"`
}

// CreateBlueprint creates a new blueprint
//
//		blueprintData := make(map[string]string)
// 		blueprintData["name"] = "Go Blueprint"
// 		blueprintData["description"] = "A blueprint for the Go lib for M2X"
// 		blueprintData["visibility"] = "private"
// 		blueprint, err := client.CreateBlueprint(blueprintData)
func (c *Client) CreateBlueprint(blueprint map[string]string) (*Blueprint, *ErrorMessage) {
	data, err := json.Marshal(blueprint)

	result, statusCode, err := post(c.APIBase+"/blueprints", data)
	if err != nil {
		return nil, simpleErrorMessage(err, statusCode)
	}

	if statusCode == 201 {
		newBlueprint := &Blueprint{}
		err = json.Unmarshal(result, &newBlueprint)
		if err != nil {
			return nil, simpleErrorMessage(err, statusCode)
		}
		return newBlueprint, nil
	}
	return nil, generateErrorMessage(result, statusCode)
}

// DeleteBlueprint deletes a blueprint
//
//		err := client.DeleteBlueprint(blueprint.ID)
func (c *Client) DeleteBlueprint(id string) *ErrorMessage {
	result, statusCode, err := delete(c.APIBase+"/blueprints", id)
	if err != nil {
		return simpleErrorMessage(err, statusCode)
	}
	if statusCode == 204 {
		return nil
	}
	return generateErrorMessage(result, statusCode)
}

// Blueprints gets a list of blueprints
//
//		blueprints, err := client.Blueprints()
func (c *Client) Blueprints() (*Blueprints, *ErrorMessage) {
	result, statusCode, err := get(c.APIBase + "/blueprints")
	if err != nil {
		return nil, simpleErrorMessage(err, statusCode)
	}
	if statusCode == 200 {
		data, err := parseBlueprints(result)
		if err != nil {
			return nil, simpleErrorMessage(err, statusCode)
		}
		return data, nil
	}
	return nil, generateErrorMessage(result, statusCode)
}

// Blueprint gets a blueprint
//
//		blueprint, err := client.Blueprint("1234")
func (c *Client) Blueprint(id string) (*Blueprint, *ErrorMessage) {
	result, statusCode, err := get(c.APIBase + "/blueprints/" + id)
	if err != nil {
		return nil, simpleErrorMessage(err, statusCode)
	}
	if statusCode == 200 {
		data, err := parseBlueprint(result)
		if err != nil {
			return nil, simpleErrorMessage(err, statusCode)
		}
		return data, nil
	}
	return nil, generateErrorMessage(result, statusCode)
}

// UpdateBlueprint updates a blueprint
//
//		blueprintData["description"] = "A blueprint for the Go lib for AT&T M2X"
//		err := client.UpdateBlueprint(blueprint.ID, blueprintData)
func (c *Client) UpdateBlueprint(id string, updateData map[string]string) *ErrorMessage {
	data, err := json.Marshal(updateData)
	if err != nil {
		return simpleErrorMessage(err, 0)
	}
	result, statusCode, postErr := put(c.APIBase+"/blueprints/"+id, data)
	if postErr != nil {
		return simpleErrorMessage(err, statusCode)
	}

	if statusCode == 204 {
		return nil
	}
	return generateErrorMessage(result, statusCode)
}

// CreateBatch creates a new batch
//
//		batchData := make(map[string]string)
// 		batchData["name"] = "Go Batch"
// 		batchData["description"] = "A batch for the Go lib for M2X"
// 		batchData["visibility"] = "private"
// 		batch, err := client.CreateBatch(batch)
func (c *Client) CreateBatch(batch map[string]string) (*Batch, *ErrorMessage) {
	data, err := json.Marshal(batch)
	if err != nil {
		return nil, simpleErrorMessage(err, 0)
	}

	result, statusCode, err := post(c.APIBase+"/batches", data)
	if err != nil {
		return nil, simpleErrorMessage(err, statusCode)
	}

	if statusCode == 201 {
		newBatch := &Batch{}
		err = json.Unmarshal(result, &newBatch)
		if err != nil {
			return nil, simpleErrorMessage(err, statusCode)
		}
		return newBatch, nil
	}
	return nil, generateErrorMessage(result, statusCode)
}

// DeleteBatch deletes a batch
//
//		err := client.DeleteBatch(batch.ID)
func (c *Client) DeleteBatch(id string) (*Batch, *ErrorMessage) {
	result, statusCode, err := delete(c.APIBase+"/batches", id)
	if err != nil {
		return nil, simpleErrorMessage(err, statusCode)
	}
	if statusCode == 204 {
		return nil, nil
	}
	return nil, generateErrorMessage(result, statusCode)
}

// Batches gets a list of batches
//
//		batches, err := client.Batches()
func (c *Client) Batches() (*Batches, *ErrorMessage) {
	result, statusCode, err := get(c.APIBase + "/batches")
	if err != nil {
		return nil, simpleErrorMessage(err, statusCode)
	}
	if statusCode == 200 {
		data, err := parseBatches(result)
		if err != nil {
			return nil, simpleErrorMessage(err, statusCode)
		}
		return data, nil
	}
	return nil, generateErrorMessage(result, statusCode)
}

// Batch gets a batch
//
//		batch, err := client.Batch("1234")
func (c *Client) Batch(id string) (*Batch, *ErrorMessage) {
	result, statusCode, err := get(c.APIBase + "/batches/" + id)
	if err != nil {
		return nil, simpleErrorMessage(err, statusCode)
	}
	if statusCode == 200 {
		data, err := parseBatch(result)
		if err != nil {
			return data, simpleErrorMessage(err, statusCode)
		}
		return data, nil
	}
	return nil, generateErrorMessage(result, statusCode)
}

// UpdateBatch updates a batch
//
//		batchData["description"] = "A batch for the Go lib for AT&T M2X"
//		err := client.UpdateBatch(batch.ID, batchData)
func (c *Client) UpdateBatch(id string, updateData map[string]string) *ErrorMessage {
	data, err := json.Marshal(updateData)
	if err != nil {
		return simpleErrorMessage(err, 0)
	}
	result, statusCode, postErr := put(c.APIBase+"/batches/"+id, data)
	if postErr != nil {
		return simpleErrorMessage(postErr, statusCode)
	}

	if statusCode == 204 {
		return nil
	}
	return generateErrorMessage(result, statusCode)
}

// Parses the JSON of blueprints request into the appropriate struct
func parseBlueprints(data []byte) (*Blueprints, error) {
	blueprints := &Blueprints{}
	err := json.Unmarshal(data, &blueprints)
	if err != nil {
		return nil, err
	}
	return blueprints, nil
}

// Parses the JSON of a single blueprint request into the appropriate struct
func parseBlueprint(data []byte) (*Blueprint, error) {
	blueprint := &Blueprint{}
	err := json.Unmarshal(data, &blueprint)
	if err != nil {
		return nil, err
	}
	return blueprint, nil
}

// Parses the JSON of batches request into the appropriate struct
func parseBatches(data []byte) (*Batches, error) {
	batches := &Batches{}
	err := json.Unmarshal(data, &batches)
	if err != nil {
		return nil, err
	}
	return batches, nil
}

// Parses the JSON of a single batch request into the appropriate struct
func parseBatch(data []byte) (*Batch, error) {
	batch := &Batch{}
	err := json.Unmarshal(data, &batch)
	if err != nil {
		return nil, err
	}
	return batch, nil
}
