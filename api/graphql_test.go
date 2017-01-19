package api_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/lucapette/deluminator/api"
	"github.com/lucapette/deluminator/db"
	"github.com/lucapette/deluminator/testutil"
)

var update = flag.Bool("update", false, "update golden files")

func TestDataSources(t *testing.T) {
	dsn, cleanup := testutil.SetupDB(db.PostgresDriver, t)
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
	actual := string(resp)

	if w.Code != 200 {
		t.Fatalf("expected code 200, got: %v. Resp: %v", w.Code, actual)
	}

	var expected bytes.Buffer
	testutil.ParseFixture(t, &expected, "data_sources.json", testutil.DBTemplate{Name: dsn.DBName})
	if *update {
		testutil.WriteFixture(t, "data_sources.json", actual)
	}

	if !reflect.DeepEqual(expected.String(), actual) {
		t.Fatalf("Unexpected result, diff: %v", testutil.Diff(expected.String(), actual))
	}
}

func TestDataSourcesWithTables(t *testing.T) {
	dsn, cleanup := testutil.SetupDB(db.PostgresDriver, t)
	app := testutil.InitApp(t, map[string]string{
		"DATA_SOURCES": dsn.Format(),
	})

	defer func() {
		app.Shutdown()
		cleanup()
	}()

	req := httptest.NewRequest("POST", "http://example.com/graphql", strings.NewReader("{DataSources {name tables {name}}}"))

	w := httptest.NewRecorder()

	api.GraphQLHandler(app)(w, req)

	resp, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	actual := string(resp)

	if w.Code != 200 {
		t.Fatalf("expected code 200, got: %v. Resp: %v", w.Code, actual)
	}

	var expected bytes.Buffer
	testutil.ParseFixture(t, &expected, "data_sources_with_tables.json", testutil.DBTemplate{Name: dsn.DBName})
	if *update {
		testutil.WriteFixture(t, "data_sources_with_tables.json", actual)
	}

	if !reflect.DeepEqual(expected.String(), actual) {
		t.Fatalf("Unexpected result, diff: %v", testutil.Diff(expected.String(), actual))
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
