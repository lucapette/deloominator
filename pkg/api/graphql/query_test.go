package graphql_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/lucapette/deloominator/pkg/db"
	"github.com/lucapette/deloominator/pkg/db/storage"
	"github.com/lucapette/deloominator/pkg/testutil"
)

func TestGraphQL_DataSources(t *testing.T) {
	type dataSourceType struct {
		Name   string `json:"name"`
		Tables []*struct {
			Name string `json:"name"`
		} `json:"tables"`
	}
	type response struct {
		Data struct {
			DataSources []dataSourceType
		} `json:"data"`
	}
	dataSources, cleanup := testutil.SetupDataSources(t)
	defer cleanup()

	testServer := NewTestServer(t, dataSources, nil)
	jsonResp := testServer.do(`{ dataSources { name tables { name } } }`, vars{})
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
		found := &dataSourceType{}
		for _, ds := range actual {

			if dataSource.DBName() == ds.Name {
				found = &ds
				break
			}
		}

		if found == nil {
			t.Fatalf("expected to find %s in %v, but did not", dataSource.DBName(), actual)
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

const Query = `
  query Query($dataSource: String!, $query: String!) {
    query(dataSource: $dataSource, query: $query) {
      ... on results {
        chartName
        columns {
          name
          type
        }
        rows {
          cells {
            value
          }
        }
      }
      ... on queryError {
        message
      }
    }
  }
`

func TestGraphQL_QueryError(t *testing.T) {
	dataSources, cleanup := testutil.SetupDataSources(t)
	defer cleanup()

	for _, dataSource := range dataSources {
		t.Run(dataSource.DriverName(), func(t *testing.T) {
			testServer := NewTestServer(t, dataSources, nil)
			actual := testServer.do(Query, vars{
				"dataSource": dataSource.DBName(),
				"query":      `select * from table_that_does_not_exist`,
			})

			tf := testutil.NewGoldenFile(t, fmt.Sprintf("query_error_%s.json", dataSource.DriverName()))
			expected := tf.ParseOrUpdate(*update, dataSource.DBName(), actual)

			if !reflect.DeepEqual(expected, actual) {
				t.Fatalf("Unexpected result, diff: %v", testutil.Diff(expected, actual))
			}
		})
	}
}

func TestGraphQL_Query(t *testing.T) {
	tests := []struct{ query, fixture string }{
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
	rows := db.QueryResult{
		Rows:    []db.Row{{{Value: "42"}, {Value: "Anna"}, {Value: "Torv"}, {Value: "2016-01-01"}}},
		Columns: []db.Column{{Name: "actor_id"}, {Name: "first_name"}, {Name: "last_name"}, {Name: "last_update"}},
	}

	dataSources, cleanup := testutil.SetupDataSources(t)
	defer cleanup()

	for _, dataSource := range dataSources {
		testutil.LoadData(t, dataSource, "actor", rows)

		for _, tc := range tests {
			t.Run(fmt.Sprintf("%s/%s", dataSource.DriverName(), tc.fixture), func(t *testing.T) {
				testServer := NewTestServer(t, dataSources, nil)
				actual := testServer.do(Query, vars{
					"dataSource": dataSource.DBName(),
					"query":      tc.query,
				})

				tf := testutil.NewGoldenFile(t, tc.fixture)
				expected := tf.ParseOrUpdate(*update, dataSource.DBName(), actual)

				if !reflect.DeepEqual(expected, actual) {
					t.Fatalf("Unexpected result, diff: %v", testutil.Diff(expected, actual))
				}
			})
		}
	}
}

func TestGraphQL_Question(t *testing.T) {
	question := `
  query Question($id: ID!) {
    question(id: $id) {
      id
      title
			query
			dataSource
			variables {
				name
				value
				isControllable
			}
    }
  }
`

	storages := testutil.NewStorages(t)
	for _, s := range storages {
		t.Run(s.String(), func(t *testing.T) {
			q, err := s.InsertQuestion(&storage.Question{
				Title:      "the answer is 42",
				Query:      "select * from answer",
				DataSource: "source",
				Variables:  `[{"name": "date", "value": "2017-10-15"}]`,
			})
			if err != nil {
				t.Fatalf("could not create question: %v", err)
			}
			testServer := NewTestServer(t, db.DataSources{}, s)
			actual := testServer.do(question, vars{"id": q.ID})

			tf := testutil.NewGoldenFile(t, "question_with_results.json")

			expected := tf.Load()

			if *update {
				expected = actual
				tf.Write(actual)
			}

			if !reflect.DeepEqual(expected, actual) {
				t.Fatalf("Unexpected result, diff: %v", testutil.Diff(expected, actual))
			}
		})
	}
}

func TestGraphQL_Questions(t *testing.T) {
	query := `
  query Questions {
    questions {
      id
      title
			query
			dataSource
			variables {
				name
				value
			}
    }
  }
`

	storages := testutil.NewStorages(t)
	for _, s := range storages {
		t.Run(s.String(), func(t *testing.T) {
			_, err := s.InsertQuestion(&storage.Question{
				Title:      "the answer is 42",
				Query:      "select * from answer",
				DataSource: "source",
				Variables:  `[{"name": "date", "value": "2017-10-15"}]`,
			})
			if err != nil {
				t.Fatalf("could not create question: %v", err)
			}
			testServer := NewTestServer(t, db.DataSources{}, s)
			actual := testServer.do(query, vars{})

			tf := testutil.NewGoldenFile(t, "questions.json")

			expected := tf.Load()

			if *update {
				expected = actual
				tf.Write(actual)
			}

			if !reflect.DeepEqual(expected, actual) {
				t.Fatalf("Unexpected result, diff: %v", testutil.Diff(expected, actual))
			}
		})
	}
}
