package graphql_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/lucapette/deloominator/pkg/db"

	"github.com/lucapette/deloominator/pkg/testutil"
)

const SaveQuestion = `
mutation SaveQuestion($title: String!, $query: String!, $dataSource: String!, $variables: String) {
  saveQuestion(title: $title, query: $query, dataSource: $dataSource, variables: $variables) {
    id
    title
	query
	variables
  } 
}`

func TestGraphQL_SaveQuestion(t *testing.T) {
	tests := []struct {
		query, fixture string
		variables      vars
	}{
		{
			SaveQuestion,
			"save_question_success.json",
			vars{
				"title":      "the answer is 42",
				"query":      "select * from answer",
				"dataSource": "dataSource",
				"variables":  `{"grace": "hopper"}`,
			},
		},
	}

	storages := testutil.NewStorages(t)

	for _, s := range storages {
		testServer := NewTestServer(t, db.DataSources{}, s)
		for _, tc := range tests {
			t.Run(fmt.Sprintf("%s/%s", s, tc.fixture), func(t *testing.T) {
				actual := testServer.do(tc.query, tc.variables)

				tf := testutil.NewGoldenFile(t, tc.fixture)

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
}
