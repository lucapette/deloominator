package api_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/lucapette/deloominator/pkg/api"
	"github.com/lucapette/deloominator/pkg/db"
	"github.com/lucapette/deloominator/pkg/testutil"
)

var update = flag.Bool("update", false, "update golden files")

var rows = db.QueryResult{
	Rows:    []db.Row{{db.Cell{Value: "42"}, db.Cell{Value: "Anna"}, db.Cell{Value: "Torv"}}},
	Columns: []db.Column{{Name: "actor_id"}, {Name: "first_name"}, {Name: "last_name"}},
}

type test struct {
	query   string
	code    int
	fixture string
}

func graphqlPayload(t *testing.T, query string) string {
	payload := struct {
		Query string `json:"query"`
	}{Query: query}

	json, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf(err.Error())
	}

	return string(json)
}

func TestGraphQLDataSources(t *testing.T) {
	dsn, cleanup := testutil.SetupPG(t)
	cfg := testutil.InitConfig(t, map[string]string{
		"DATA_SOURCES": dsn,
	})
	dataSources, err := db.NewDataSources(cfg.Sources)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		dataSources.Shutdown()
		cleanup()
	}()

	tests := []test{
		{query: "{ notAQuery }", code: 400, fixture: "wrong_query.json"},
		{query: "{ dataSources {name} }", code: 200, fixture: "data_sources.json"},
		{query: "{ dataSources {name tables {name}}}", code: 200, fixture: "data_sources_with_tables.json"},
	}

	for _, dataSource := range dataSources {
		testutil.LoadData(t, dataSource, "actor", rows)

		for _, test := range tests {
			t.Run(test.fixture, func(t *testing.T) {
				req := httptest.NewRequest("POST", "http://example.com/graphql", strings.NewReader(graphqlPayload(t, test.query)))
				w := httptest.NewRecorder()

				api.GraphQLHandler(dataSources)(w, req)

				resp, err := ioutil.ReadAll(w.Body)
				if err != nil {
					t.Fatal(err)
				}
				actual := string(resp)

				if w.Code != test.code {
					t.Fatalf("expected code %d, got: %d. Resp: %s", test.code, w.Code, actual)
				}

				var expected bytes.Buffer
				testutil.ParseFixture(t, &expected, test.fixture, testutil.DBTemplate{Name: dataSource.Name()})
				if *update {
					testutil.WriteFixture(t, test.fixture, actual)
				}

				if !reflect.DeepEqual(strings.TrimSuffix(expected.String(), "\n"), actual) {
					t.Fatalf("Unexpected result, diff: %v", testutil.Diff(expected.String(), actual))
				}
			})
		}
	}
}

var graphQLQuery = `
{ query(source: "%s", input: "%s") {
	... on results {
		chartName
	    columns { name type }
	    rows { cells { value } }
    }
    ... on queryError { message }
  }
}
`

var queryTests = []test{
	{
		query:   `select * from table_that_does_not_exist`,
		fixture: "query_error.json",
	},
	{
		query:   `select actor_id, first_name, last_name from actor`,
		fixture: "query_raw_results.json",
	},
	{
		query:   `select substr(first_name, 1, 1) initial, count(*)  from actor group by 1`,
		fixture: "query_simple_bar_detected.json",
	},
}

func TestGraphQLQuery(t *testing.T) {
	dsn, cleanup := testutil.SetupPG(t)
	cfg := testutil.InitConfig(t, map[string]string{
		"DATA_SOURCES": dsn,
	})
	dataSources, err := db.NewDataSources(cfg.Sources)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		dataSources.Shutdown()
		cleanup()
	}()

	for _, dataSource := range dataSources {
		testutil.LoadData(t, dataSource, "actor", rows)

		for _, test := range queryTests {
			t.Run(test.fixture, func(t *testing.T) {
				query := graphqlPayload(t, fmt.Sprintf(graphQLQuery, dataSource.Name(), test.query))
				req := httptest.NewRequest("POST", "http://example.com/graphql", strings.NewReader(query))
				w := httptest.NewRecorder()

				api.GraphQLHandler(dataSources)(w, req)

				resp, err := ioutil.ReadAll(w.Body)
				if err != nil {
					t.Fatal(err)
				}
				actual := string(resp)

				if w.Code != 200 {
					t.Fatalf("expected code %d, got: %d. Resp: %s", 200, w.Code, actual)
				}

				var expected bytes.Buffer
				testutil.ParseFixture(t, &expected, test.fixture, testutil.DBTemplate{Name: dataSource.Name()})
				if *update {
					testutil.WriteFixture(t, test.fixture, actual)
				}

				if !reflect.DeepEqual(strings.TrimSuffix(expected.String(), "\n"), actual) {
					t.Fatalf("Unexpected result, diff: %v", testutil.Diff(expected.String(), actual))
				}
			})
		}
	}
}
