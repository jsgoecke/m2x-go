// Copyright (c) 2014 Jason Goecke
// datasources_test.go

package m2x

import (
	"os"
	"testing"
	"time"
)

func TestParseBlueprints(t *testing.T) {
	// Example blueprints
	data := `
	{ "blueprints": [
  { "id": "6b541a4462878411b945dc88c848f7e7",
    "name": "Sample Blueprint",
    "description": "Longer description for Sample Blueprint",
    "visibility": "public",
    "serial": null,
    "status": "enabled",
    "feed": "/feeds/6b541a4462878411b945dc88c848f7e7",
    "url": "/blueprints/6b541a4462878411b945dc88c848f7e7",
    "key": "c325f5c1aeff96af6492ae622263b97d",
    "created": "2013-09-01T10:00:00Z",
    "updated": "2013-09-02T10:00:00Z" },
  { "id": "06f2e5e74472b0b32c5330e71231cc99",
    "name": "Another Blueprint",
    "description": "Longer description for Another Blueprint",
    "visibility": "private",
    "serial": null,
    "status": "disabled",
    "feed": "/feeds/06f2e5e74472b0b32c5330e71231cc99",
    "url": "/blueprints/06f2e5e74472b0b32c5330e71231cc99",
    "key": "ada3211b1bbeaf5bef45a1821311ece7",
    "created": "2013-09-01T10:00:00Z",
    "updated": "2013-09-02T10:00:00Z" } ] }`

	result, _ := parseBlueprints([]byte(data))

	if result.Blueprints[0].Name != "Sample Blueprint" || result.Blueprints[1].Name != "Another Blueprint" {
		t.Errorf("Names did not parse properly")
	}

	if result.Blueprints[0].Key != "c325f5c1aeff96af6492ae622263b97d" || result.Blueprints[1].Key != "ada3211b1bbeaf5bef45a1821311ece7" {
		t.Errorf("Keys did not parse properly")
	}
}

func TestParseBlueprint(t *testing.T) {
	// Example blueprint
	data := `
    { "id":"3281fe1067ce3d7ae276ad360066e7be",
      "name": "Sample Blueprint",
      "description": "Longer description for Sample Blueprint",
      "visibility": "public",
      "serial": null,
      "status": "enabled",
      "tags": [],
      "feed": "/feeds/3281fe1067ce3d7ae276ad360066e7be",
      "url": "/blueprints/3281fe1067ce3d7ae276ad360066e7be",
      "key": "d66ae7e3d39fad2290ad90577b7a9963",
      "created": "2013-09-01T10:00:00Z",
      "updated": "2013-09-02T10:00:00Z" }`

	blueprint, _ := parseBlueprint([]byte(data))
	if blueprint.ID != "3281fe1067ce3d7ae276ad360066e7be" {
		t.Errorf("Id did not parse properly")
	}
}

func TestParseBatches(t *testing.T) {
	// Example batches
	data := `
    { "batches": [
  { "id": "1b3ba972fcf92a156fc8c0ca1554434c",
    "name": "Sample Batch",
    "description": "Longer description for Sample Batch",
    "visibility": "public",
    "serial": null,
    "status": "enabled",
    "tags": [],
    "feed": "/feeds/1b3ba972fcf92a156fc8c0ca1554434c",
    "url": "/batches/1b3ba972fcf92a156fc8c0ca1554434c",
    "key": "8eb857a0672879dcad16f198aaafbc5b",
    "created": "2013-10-09T17:09:43Z",
    "updated": "2013-10-09T17:09:43Z",
    "datasources": {
      "total": 10,
      "registered": 8,
      "unregistered": 2
    }
  },
  { "id": "4fc642808294b2f8af137bb655fa0e18",
    "name": "Sample Batch 2",
    "description": "Longer description for Sample Batch",
    "visibility": "private",
    "serial": null,
    "status": "enabled",
    "tags": [],
    "feed": "/feeds/4fc642808294b2f8af137bb655fa0e18",
    "url": "/batches/4fc642808294b2f8af137bb655fa0e18",
    "key": "df17a48c5208d8db584b6b994dc6726b",
    "created": "2013-10-10T18:07:54Z",
    "updated": "2013-10-10T18:07:54Z",
    "datasources": {
      "total": 10,
      "registered": 10,
      "unregistered": 0
    }
  } ] }`

	result, _ := parseBatches([]byte(data))
	if result.Batches[0].Name != "Sample Batch" || result.Batches[1].Name != "Sample Batch 2" {
		t.Errorf("Names did not parse properly")
	}

	if result.Batches[0].ID != "1b3ba972fcf92a156fc8c0ca1554434c" || result.Batches[1].ID != "4fc642808294b2f8af137bb655fa0e18" {
		t.Errorf("Names did not parse properly")
	}
}

func TestParseBatch(t *testing.T) {
	// Example batch
	data := `
    { "id": "1b3ba972fcf92a156fc8c0ca1554434c",
  "name": "Sample Batch",
  "description": "Longer description for Sample Batch",
  "visibility": "public",
  "serial": null,
  "status": "enabled",
  "tags": [],
  "feed": "/feeds/1b3ba972fcf92a156fc8c0ca1554434c",
  "url": "/batches/1b3ba972fcf92a156fc8c0ca1554434c",
  "key": "8eb857a0672879dcad16f198aaafbc5b",
  "created": "2013-10-09T17:09:43Z",
  "updated": "2013-10-09T17:09:43Z",
  "datasources": {
    "total": 0,
    "registered": 0,
    "unregistered": 0
  } }`

	batch, _ := parseBatch([]byte(data))
	if batch.ID != "1b3ba972fcf92a156fc8c0ca1554434c" {
		t.Errorf("Id did not parse properly")
	}

	if batch.Datasources.Total != 0 {
		t.Errorf("Datasources total did not parse properly")
	}
}

func TestListBlueprints(t *testing.T) {
	client := NewClient(os.Getenv("M2X_API_KEY"))
	blueprints, _ := client.Blueprints()
	if blueprints.CurrentPage != 1 {
		t.Errorf("Did not fetch Blueprints properly")
	}
}

func TestCreateAndListAndUpdateAndDeleteBlueprint(t *testing.T) {
	client := NewClient(os.Getenv("M2X_API_KEY"))

	// Create a new blueprint
	blueprintData := make(map[string]string)
	theTime := time.Now()
	name := "Go Created Blueprint - " + theTime.Format("20060102150405")
	blueprintData["name"] = name
	blueprintData["description"] = "Unit testing Go lib for M2X"
	blueprintData["visibility"] = "private"
	blueprint, err := client.CreateBlueprint(blueprintData)
	if blueprint.Description != "Unit testing Go lib for M2X" || blueprint.Visibility != "private" {
		t.Errorf("Did create a new blueprint properly")
	}

	// List the blueprint
	blueprint, err = client.Blueprint(blueprint.ID)
	if err != nil || blueprint.Description != "Unit testing Go lib for M2X" {
		t.Errorf("Did not fetch blueprint properly")
	}

	// Update the blueprint
	blueprintData["name"] = name
	blueprintData["description"] = "Updated description!"
	blueprintData["visibility"] = "private"
	err = client.UpdateBlueprint(blueprint.ID, blueprintData)
	if err != nil {
		t.Errorf("Updating blueprint did not work")
	}
	blueprint, err = client.Blueprint(blueprint.ID)
	if err != nil || blueprint.Description != "Updated description!" {
		t.Errorf("Did not fetch blueprint properly")
	}

	// Delete the blueprint
	err = client.DeleteBlueprint(blueprint.ID)
	if err != nil {
		t.Errorf("Did not delete blueprint properly")
	}
}

func TestListBatches(t *testing.T) {
	client := NewClient(os.Getenv("M2X_API_KEY"))
	blueprints, _ := client.Batches()
	if blueprints.CurrentPage != 1 {
		t.Errorf("Did not fetch Batches properly")
	}
}

func TestCreateAndListAndUpdateAndDeleteBatch(t *testing.T) {
	client := NewClient(os.Getenv("M2X_API_KEY"))

	// Create a new batch
	batchData := make(map[string]string)
	theTime := time.Now()
	name := "Go Created Batch - " + theTime.Format("20060102150405")
	batchData["name"] = name
	batchData["description"] = "Unit testing Go lib for M2X"
	batchData["visibility"] = "private"
	batch, err := client.CreateBatch(batchData)
	if err != nil || batch.Description != "Unit testing Go lib for M2X" || batch.Visibility != "private" {
		t.Errorf("Did create a new batch properly")
	}

	// List the batch
	batch, err = client.Batch(batch.ID)
	if batch.Description != "Unit testing Go lib for M2X" {
		t.Errorf("Did not fetch batch properly")
	}

	// Update the batch
	batchData["name"] = name
	batchData["description"] = "Updated description!"
	batchData["visibility"] = "private"
	err = client.UpdateBatch(batch.ID, batchData)
	if err != nil {
		t.Errorf("Did not update batch properly")
	}
	batch, err = client.Batch(batch.ID)
	if err != nil || batch.Description != "Updated description!" {
		t.Errorf("Did not fetch batch properly")
	}

	// Delete the batch
	batch, err = client.DeleteBatch(batch.ID)
	if batch != nil || err != nil {
		t.Errorf("Did not delete batch properly")
	}
}

func TestListBadBlueprintId(t *testing.T) {
	client := NewClient(os.Getenv("M2X_API_KEY"))
	_, errorMessage := client.Blueprint("1234")
	if errorMessage.StatusCode != 404 || errorMessage.Message != "The specified blueprint does not exist" {
		t.Errorf("We did not get the proper error message or code back")
	}
}

func TestListBadBatchId(t *testing.T) {
	client := NewClient(os.Getenv("M2X_API_KEY"))
	_, errorMessage := client.Batch("1234")
	if errorMessage.StatusCode != 404 || errorMessage.Message != "The specified batch does not exist" {
		t.Errorf("We did not get the proper error message or code back")
	}
}
