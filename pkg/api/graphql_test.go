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
	Rows:    []db.Row{{db.Cell{Value: "42"}, db.Cell{Value: "Anna"}, db.Cell{Value: "Torv"}, db.Cell{Value: "2016-01-01"}}},
	Columns: []db.Column{{Name: "actor_id"}, {Name: "first_name"}, {Name: "last_name"}, {Name: "last_update"}},
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
	dsnPG, cleanupPG := testutil.SetupPG(t)
	cfg := testutil.InitConfig(t, map[string]string{
		"DATA_SOURCES": dsnPG,
	})
	dataSources, err := db.NewDataSources(cfg.Sources)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		dataSources.Shutdown()
		cleanupPG()
	}()

	tests := []test{
		{query: "{ notAQuery }", code: 400, fixture: "wrong_query.json"},
		{query: "{ dataSources {name} }", code: 200, fixture: "data_sources.json"},
		{query: "{ dataSources {name tables {name} } }", code: 200, fixture: "data_sources_with_tables.json"},
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

func TestGraphQLQueryError(t *testing.T) {
	dsnPG, cleanupPG := testutil.SetupPG(t)
	dsnMySQL, cleanupMySQL := testutil.SetupMySQL(t)

	cfg := testutil.InitConfig(t, map[string]string{
		"DATA_SOURCES": fmt.Sprintf("%s,%s", dsnPG, dsnMySQL),
	})
	dataSources, err := db.NewDataSources(cfg.Sources)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		dataSources.Shutdown()
		cleanupPG()
		cleanupMySQL()
	}()

	for _, dataSource := range dataSources {
		t.Run(dataSource.Driver, func(t *testing.T) {
			testutil.LoadData(t, dataSource, "actor", rows)

			query := graphqlPayload(t, fmt.Sprintf(graphQLQuery, dataSource.Name(), `select * from table_that_does_not_exist`))
			req := httptest.NewRequest("POST", "http://example.com/graphql", strings.NewReader(query))
			w := httptest.NewRecorder()

			api.GraphQLHandler(dataSources)(w, req)

			resp, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Fatal(err)
			}
			actual := string(resp)

			if w.Code != 200 {
				t.Fatalf("expected code 200, got: %d. Resp: %s", w.Code, actual)
			}

			fixture := fmt.Sprintf("query_error_%s.json", dataSource.Driver)
			var expected bytes.Buffer
			testutil.ParseFixture(t, &expected, fixture, testutil.DBTemplate{Name: dataSource.Name()})
			if *update {
				testutil.WriteFixture(t, fixture, actual)
			}

			if !reflect.DeepEqual(strings.TrimSuffix(expected.String(), "\n"), actual) {
				t.Fatalf("Unexpected result, diff: %v", testutil.Diff(expected.String(), actual))
			}
		})
	}
}

var queryTests = []test{
	{
		query:   `select actor_id, first_name, last_name from actor`,
		fixture: "query_raw_results.json",
	},
	{
		query:   `select substr(first_name, 1, 1) initial, count(*) total  from actor group by 1`,
		fixture: "query_simple_bar_detected.json",
	},
	{
		query:   `select last_update, count(*) total from actor group by 1`,
		fixture: "query_simple_line_detected.json",
	},
}

func TestGraphQLQuery(t *testing.T) {
	dsnPG, cleanupPG := testutil.SetupPG(t)
	dsnMySQL, cleanupMySQL := testutil.SetupMySQL(t)
	cfg := testutil.InitConfig(t, map[string]string{
		"DATA_SOURCES": fmt.Sprintf("%s,%s", dsnPG, dsnMySQL),
	})
	dataSources, err := db.NewDataSources(cfg.Sources)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		dataSources.Shutdown()
		cleanupPG()
		cleanupMySQL()
	}()

	for _, dataSource := range dataSources {
		testutil.LoadData(t, dataSource, "actor", rows)

		for _, test := range queryTests {
			t.Run(fmt.Sprintf("%s/%s", dataSource.Driver, test.fixture), func(t *testing.T) {
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
					t.Fatalf("expected code 200, got: %d. Resp: %s", w.Code, actual)
				}

				expected := testutil.LoadFixture(t, test.fixture)
				if *update {
					testutil.WriteFixture(t, test.fixture, actual)
				}

				if !reflect.DeepEqual(strings.TrimSuffix(expected, "\n"), actual) {
					t.Fatalf("Unexpected result, diff: %v", testutil.Diff(expected, actual))
				}
			})
		}
	}
}
