package api_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/kr/pretty"
	"github.com/lucapette/deluminator/api"
)

func TestGraphQLPOSTQuery(t *testing.T) {
	req := httptest.NewRequest("POST", "http://example.com/graphql", strings.NewReader("{hello}"))

	w := httptest.NewRecorder()

	api.GraphQLHandler(w, req)

	if w.Code != 200 {
		t.Fatalf("expected code 200, got: %v", w.Code)
	}

	resp, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}

	var actual *graphql.Result

	err = json.Unmarshal(resp, &actual)
	if err != nil {
		t.Fatal(err)
	}

	expected := &graphql.Result{
		Data: map[string]interface{}{
			"hello": "world",
		},
	}

	if actual.HasErrors() {
		t.Fatalf("wrong result, unexpected errors: %v", actual.Errors)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Unexpected result, Diff: %v", pretty.Diff(expected, actual))
	}
}

func TestGraphQLPOSTWrongQuery(t *testing.T) {
	req := httptest.NewRequest("POST", "http://example.com/graphql", strings.NewReader("{notAQuery}"))
	w := httptest.NewRecorder()
	api.GraphQLHandler(w, req)

	if w.Code != 400 {
		t.Fatalf("expected code 400, got: %v", w.Code)
	}

	resp, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}

	var actual *graphql.Result

	err = json.Unmarshal(resp, &actual)
	if err != nil {
		t.Fatal(err)
	}

	if !actual.HasErrors() {
		t.Fatal("Expected error but got none")
	}
}
