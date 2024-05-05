package mocks

import (
	"bytes"
	"encoding/json"
	"testing"
)

// Compare deep equality between two json objects.
func AssertJSON(received interface{}, expected interface{}, t *testing.T) {
	expectedJson, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	receivedJson, err := json.Marshal(received)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling receive json data", err)
	}

	if !bytes.Equal(expectedJson, receivedJson) {
		t.Errorf("the expected json: %s is different from received %s", expectedJson, receivedJson)
	}
}
