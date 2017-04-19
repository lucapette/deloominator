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

func TestGraphQLQueries(t *testing.T) {
	dsn, cleanup := testutil.SetupDB(t, db.PostgresDriver)
	app := testutil.InitApp(t, map[string]string{
		"DATA_SOURCES": dsn,
	})
	driver, err := db.NewDataSource(dsn)
	if err != nil {
		t.Fatal(err)
	}

	rows := db.QueryResult{
		Rows:    []db.Row{db.Row{db.Cell{Value: "42"}, db.Cell{Value: "Anna"}, db.Cell{Value: "Torv"}}},
		Columns: []db.Column{db.Column{Name: "actor_id"}, db.Column{Name: "first_name"}, db.Column{Name: "last_name"}},
	}
	testutil.LoadData(t, driver, "actor", rows)

	defer func() {
		err = driver.Close()
		if err != nil {
			t.Fatal(err)
		}

		app.Shutdown()
		cleanup()
	}()

	tests := []test{
		{query: graphqlPayload(t, "{ notAQuery }"), code: 400, fixture: "wrong_query.json"},
		{query: graphqlPayload(t, "{ dataSources {name} }"), code: 200, fixture: "data_sources.json"},
		{query: graphqlPayload(t, "{ dataSources {name tables {name}}}"), code: 200, fixture: "data_sources_with_tables.json"},
		{
			query: graphqlPayload(t, fmt.Sprintf(`{ query(source: "%s", input: "select actor_id, first_name, last_name from actor") {
			                                          ... on results {
														chartName
									                    total
									                    columns { name type }
									                    rows { cells { value } }
								                      }
		                                            }
	                                              }`, driver.Name())),
			code:    200,
			fixture: "query_raw_results.json",
		},
		{
			query: graphqlPayload(t, fmt.Sprintf(`{ query(source: "%s", input: "select count(*), substr(first_name, 1, 1) initial  from actor group by 2") {
			                                          ... on results {
														chartName
									                    total
									                    columns { name type }
									                    rows { cells { value } }
								                      }
		                                            }
	                                              }`, driver.Name())),
			code:    200,
			fixture: "query_chart_detected.json",
		},
		{
			query: graphqlPayload(t, fmt.Sprintf(`{ query(source: "%s", input: "select * from table_that_does_not_exist") {
			                                          ... on queryError {
								                        message
								                      }
		                                            }
												  }`, driver.Name())),
			code:    200,
			fixture: "query_error.json",
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
			testutil.ParseFixture(t, &expected, test.fixture, testutil.DBTemplate{Name: driver.Name()})
			if *update {
				testutil.WriteFixture(t, test.fixture, actual)
			}

			if !reflect.DeepEqual(strings.TrimSuffix(expected.String(), "\n"), actual) {
				t.Fatalf("Unexpected result, diff: %v", testutil.Diff(expected.String(), actual))
			}
		})
	}
}
