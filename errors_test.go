// Copyright (c) 2014 Jason Goecke
// errors_test.go

package m2x

import (
	"testing"
)

func TestParseErrorMessage(t *testing.T) {
	data := `{ "message": "Account must be active in order to perform this action" }`
	result, _ := parseErrorMessage([]byte(data))
	if result.Message != "Account must be active in order to perform this action" {
		t.Errorf("Error message did not parse properly")
	}
}
