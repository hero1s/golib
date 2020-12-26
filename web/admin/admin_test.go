package admin

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	t.Log("Testing the adding of JSON to the response")

	w := httptest.NewRecorder()
	originalBody := []int{1, 2, 3}

	res, _ := json.Marshal(originalBody)

	writeJSON(w, res)

	decodedBody := []int{}
	err := json.NewDecoder(w.Body).Decode(&decodedBody)

	if err != nil {
		t.Fatal("Could not decode response body into slice.")
	}

	for i := range decodedBody {
		if decodedBody[i] != originalBody[i] {
			t.Fatalf("Expected %d but got %d in decoded body slice", originalBody[i], decodedBody[i])
		}
	}
}
