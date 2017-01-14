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
	"github.com/lucapette/deluminator/db"
	"github.com/lucapette/deluminator/testutil"
)

func TestDataSources(t *testing.T) {
	dsn, cleanup := testutil.SetupDB(db.Postgres, t)
	app := testutil.InitApp(t, map[string]string{
		"DATA_SOURCES": dsn.Format(),
	})

	defer func() {
		app.Shutdown()
		cleanup()
	}()

	req := httptest.NewRequest("POST", "http://example.com/graphql", strings.NewReader("{DataSources {name}}"))

	w := httptest.NewRecorder()

	api.GraphQLHandler(app)(w, req)

	resp, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}

	if w.Code != 200 {
		t.Fatalf("expected code 200, got: %v. Resp: %v", w.Code, string(resp))
	}

	var actual *graphql.Result

	err = json.Unmarshal(resp, &actual)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("result %v", actual)

	expected := &graphql.Result{
		Data: map[string]interface{}{
			"DataSources": []interface{}{
				map[string]interface{}{
					"name": dsn.DBName,
				},
			},
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
	app := testutil.InitApp(t, map[string]string{}) // ugly signature. Needs fixing.

	defer func() {
		app.Shutdown()
	}()

	req := httptest.NewRequest("POST", "http://example.com/graphql", strings.NewReader("{notAQuery}"))
	w := httptest.NewRecorder()
	api.GraphQLHandler(app)(w, req)

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
