package api_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/lucapette/deluminator/api"
)

func TestGraphQLGETQuery(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/graphql?query={hello}", nil)
	w := httptest.NewRecorder()
	api.GraphQLHandler(w, req)
	if w.Code != 200 {
		t.Fatalf("expected code 200, got: %v", w.Code)
	}

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}

	var resp struct {
		Data map[string]string
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Data) == 0 {
		t.Fatal("Expected data in response but got none")
	}

	world := resp.Data["hello"]

	if world != "world" {
		t.Fatalf("Expected hello to return world but got %v", world)
	}
}
