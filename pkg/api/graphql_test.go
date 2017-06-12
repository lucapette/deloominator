package api_test

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

	"github.com/lucapette/deloominator/pkg/api"
	"github.com/lucapette/deloominator/pkg/db"
	"github.com/lucapette/deloominator/pkg/testutil"
)

var update = flag.Bool("update", false, "update golden files")

func doRequest(t *testing.T, dataSources db.DataSources, code int, query string) string {
	json, err := json.Marshal(api.GraphqlPayload{Query: query})
	if err != nil {
		t.Fatalf(err.Error())
	}

	req := httptest.NewRequest("POST", "http://example.com/graphql", strings.NewReader(string(json)))
	w := httptest.NewRecorder()

	api.GraphQLHandler(dataSources)(w, req)

	resp, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}

	if w.Code != code {
		t.Fatalf("expected request to return %d, but got: %d", code, w.Code)
	}

	return string(resp)
}

func loadOrUpdateFixture(t *testing.T, fixture string, dataSource *db.DataSource, actual string) string {
	out := &bytes.Buffer{}

	testutil.ParseFixture(t, out, fixture, testutil.DBTemplate{Name: dataSource.Name()})

	if *update {
		testutil.WriteFixture(t, fixture, strings.Replace(actual, dataSource.Name(), "{{.Name}}", -1))

	}
	return strings.TrimSuffix(out.String(), "\n")
}

func TestGraphQLNotAValidGraphQLQuery(t *testing.T) {
	actual := doRequest(t, db.DataSources{}, 400, `{ notAValidQuery }`)

	expected := testutil.LoadFixture(t, "not_a_valid_query.json")
	if *update {
		testutil.WriteFixture(t, "not_a_valid_query.json", actual)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Unexpected result, diff: %v", testutil.Diff(expected, actual))
	}
}

func TestGraphQLDataSources(t *testing.T) {
	type response struct {
		Data struct {
			DataSources []api.DataSource `json:"dataSources"`
		} `json:"data"`
	}
	dataSources, cleanup := testutil.SetupDataSources(t)
	defer func() {
		cleanup()
	}()

	jsonResp := doRequest(t, dataSources, 200, `{ dataSources { name tables { name } } }`)
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
		found := &api.DataSource{}
		for _, ds := range actual {

			if strings.Compare(dataSource.Name(), ds.Name) == 0 {
				found = &ds
				break
			}
		}

		if found == nil {
			t.Fatalf("expected to find %s in %v, but did not", dataSource.Name(), actual)
		}

		// this is annoyingly long.
		// Look for improvements.
		queryResult, err := dataSource.Tables()
		if err != nil {
			t.Fatalf("cannot load tables: %v", err)
		}
		rows := make([]string, len(queryResult.Rows))
		for i, r := range queryResult.Rows {
			rows[i] = r[0].Value
		}
		sort.Strings(rows)

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
	defer func() {
		cleanup()
	}()

	for _, dataSource := range dataSources {
		t.Run(dataSource.Driver, func(t *testing.T) {
			actual := doRequest(t, dataSources, 200, fmt.Sprintf(graphQLQuery, dataSource.Name(), `select * from table_that_does_not_exist`))

			expected := loadOrUpdateFixture(t, fmt.Sprintf("query_error_%s.json", dataSource.Driver), dataSource, actual)

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
	defer func() {
		cleanup()
	}()

	for _, dataSource := range dataSources {
		testutil.LoadData(t, dataSource, "actor", rows)

		for _, test := range queryTests {
			t.Run(fmt.Sprintf("%s/%s", dataSource.Driver, test.fixture), func(t *testing.T) {
				actual := doRequest(t, dataSources, 200, fmt.Sprintf(graphQLQuery, dataSource.Name(), test.query))

				expected := strings.TrimSuffix(loadOrUpdateFixture(t, test.fixture, dataSource, actual), "\n")

				if !reflect.DeepEqual(expected, actual) {
					t.Fatalf("Unexpected result, diff: %v", testutil.Diff(expected, actual))
				}
			})
		}
	}
}
