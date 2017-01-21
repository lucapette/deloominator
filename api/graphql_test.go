package api_test

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/lucapette/deluminator/api"
	"github.com/lucapette/deluminator/db"
	"github.com/lucapette/deluminator/testutil"
)

var update = flag.Bool("update", false, "update golden files")

type test struct {
	query   string
	code    int
	fixture string
}

func TestGraphQLQueries(t *testing.T) {
	dsn, cleanup := testutil.SetupDB(db.PostgresDriver, t)
	app := testutil.InitApp(t, map[string]string{
		"DATA_SOURCES": dsn.Format(),
	})
	defer func() {
		app.Shutdown()
		cleanup()
	}()

	tests := []test{
		{query: "{notAQuery}", code: 400, fixture: "wrong_query.json"},
		{query: "{DataSources {name}}", code: 200, fixture: "data_sources.json"},
		{query: "{DataSources {name tables {name}}}", code: 200, fixture: "data_sources_with_tables.json"},
		{
			query: fmt.Sprintf(`{
				                  Query(source: "%s", input: "select * from users") {
			                        ... on RawResults {
									  total
								    }
		                          }
	                            }`, dsn.DBName),
			code:    200,
			fixture: "query_raw_results.json",
		},
	}

	for _, test := range tests {
		t.Run(test.fixture, func(t *testing.T) {
			req := httptest.NewRequest("POST", "http://example.com/graphql", strings.NewReader(test.query))

			w := httptest.NewRecorder()

			api.GraphQLHandler(app)(w, req)

			resp, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Fatal(err)
			}
			actual := string(resp)

			if w.Code != test.code {
				t.Fatalf("expected code %d, got: %d. Resp: %s", test.code, w.Code, actual)
			}

			var expected bytes.Buffer
			testutil.ParseFixture(t, &expected, test.fixture, testutil.DBTemplate{Name: dsn.DBName})
			if *update {
				testutil.WriteFixture(t, test.fixture, actual)
			}

			if !reflect.DeepEqual(expected.String(), actual) {
				t.Fatalf("Unexpected result, diff: %v", testutil.Diff(expected.String(), actual))
			}
		})
	}
}
