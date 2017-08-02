package graphql_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/lucapette/deloominator/pkg/api/graphql"
	"github.com/lucapette/deloominator/pkg/db"
	"github.com/lucapette/deloominator/pkg/db/storage"
	"github.com/lucapette/deloominator/pkg/testutil"
)

var update = flag.Bool("update", false, "update golden files")

func doRequest(t *testing.T, dataSources db.DataSources, storage *storage.Storage, code int, query string) string {
	json, err := json.Marshal(graphql.Payload{Query: query})
	if err != nil {
		t.Fatalf(err.Error())
	}

	req := httptest.NewRequest("POST", "http://example.com/graphql", strings.NewReader(string(json)))
	w := httptest.NewRecorder()

	graphql.Handler(dataSources, storage)(w, req)

	resp, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}

	if w.Code != code {
		t.Fatalf("expected request to return %d, but got: %d", code, w.Code)
	}

	return string(resp)
}

func loadOrUpdateGoldenFile(t *testing.T, fixture string, dataSource *db.DataSource, actual string) string {
	out := &bytes.Buffer{}

	testFile := testutil.NewGoldenFile(t, fixture)

	testFile.Parse(out, dataSource.Name())

	if *update {
		testFile.Write(strings.Replace(actual, dataSource.Name(), "{{.}}", -1))
	}

	return strings.TrimSuffix(out.String(), "\n")
}

func TestGraphQLNotAValidGraphQLQuery(t *testing.T) {
	actual := doRequest(t, db.DataSources{}, nil, 400, `{ notAValidQuery }`)

	testFile := testutil.NewGoldenFile(t, "not_a_valid_query.json")
	expected := testFile.Load()
	if *update {
		testFile.Write(actual)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Unexpected result, diff: %v", testutil.Diff(expected, actual))
	}
}

func TestGraphQLDataSources(t *testing.T) {
	type response struct {
		Data struct {
			DataSources []graphql.DataSource `json:"dataSources"`
		} `json:"data"`
	}
	dataSources, cleanup := testutil.SetupDataSources(t)
	defer cleanup()

	jsonResp := doRequest(t, dataSources, nil, 200, `{ dataSources { name tables { name } } }`)
	resp := response{}
	err := json.Unmarshal([]byte(jsonResp), &resp)
	if err != nil {
		t.Fatalf("cannot unmarshal response: %v", err)
	}
	actual := resp.Data.DataSources

	if len(dataSources) != len(actual) {
		t.Fatalf("expected response to have %d dataSources, but got %d", len(dataSources), len(actual))
	}

	for _, dataSource := range dataSources {
		found := &graphql.DataSource{}
		for _, ds := range actual {

			if dataSource.Name() == ds.Name {
				found = &ds
				break
			}
		}

		if found == nil {
			t.Fatalf("expected to find %s in %v, but did not", dataSource.Name(), actual)
		}

		// this is annoyingly long.
		// Look for improvements.
		rows, err := dataSource.Tables()
		if err != nil {
			t.Fatalf("cannot load tables: %v", err)
		}

		tables := make([]string, len(found.Tables))
		for i, r := range found.Tables {
			tables[i] = r.Name
		}
		sort.Strings(tables)

		if !reflect.DeepEqual(rows, tables) {
			t.Fatalf("Unexpected result, diff: %v", testutil.Diff(rows, tables))
		}
	}
}

const graphQLQuery = `
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
	dataSources, cleanup := testutil.SetupDataSources(t)
	defer cleanup()

	for _, dataSource := range dataSources {
		t.Run(dataSource.Driver, func(t *testing.T) {
			actual := doRequest(t, dataSources, nil, 200, fmt.Sprintf(graphQLQuery, dataSource.Name(), `select * from table_that_does_not_exist`))

			expected := loadOrUpdateGoldenFile(t, fmt.Sprintf("query_error_%s.json", dataSource.Driver), dataSource, actual)

			if !reflect.DeepEqual(expected, actual) {
				t.Fatalf("Unexpected result, diff: %v", testutil.Diff(expected, actual))
			}
		})
	}
}

var queryTests = []struct{ query, fixture string }{
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
	{
		query:   `select first_name, last_name, count(*) total from actor group by 1, 2`,
		fixture: "query_grouped_bar_detected.json",
	},
	{
		query:   `select last_update, first_name, count(*) total from actor group by 1, 2`,
		fixture: "query_multi_line_detected.json",
	},
}

func TestGraphQLQuery(t *testing.T) {
	rows := db.QueryResult{
		Rows:    []db.Row{{{Value: "42"}, {Value: "Anna"}, {Value: "Torv"}, {Value: "2016-01-01"}}},
		Columns: []db.Column{{Name: "actor_id"}, {Name: "first_name"}, {Name: "last_name"}, {Name: "last_update"}},
	}

	dataSources, cleanup := testutil.SetupDataSources(t)
	defer cleanup()

	for _, dataSource := range dataSources {
		testutil.LoadData(t, dataSource, "actor", rows)

		for _, test := range queryTests {
			t.Run(fmt.Sprintf("%s/%s", dataSource.Driver, test.fixture), func(t *testing.T) {
				actual := doRequest(t, dataSources, nil, 200, fmt.Sprintf(graphQLQuery, dataSource.Name(), test.query))

				expected := strings.TrimSuffix(loadOrUpdateGoldenFile(t, test.fixture, dataSource, actual), "\n")

				if !reflect.DeepEqual(expected, actual) {
					t.Fatalf("Unexpected result, diff: %v", testutil.Diff(expected, actual))
				}
			})
		}
	}
}
