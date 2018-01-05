package graphql_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/lucapette/deloominator/pkg/api/graphql"
	"github.com/lucapette/deloominator/pkg/db"
	"github.com/lucapette/deloominator/pkg/db/storage"
	"github.com/lucapette/deloominator/pkg/testutil"
)

var update = flag.Bool("update", false, "update golden files")

type payload struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName,omitempty"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
}

type testServer struct {
	s           *storage.Storage
	dataSources db.DataSources
	tt          *testing.T
}
type vars map[string]interface{}

func NewTestServer(t *testing.T, dataSources db.DataSources, s *storage.Storage) *testServer {
	return &testServer{tt: t, s: s, dataSources: dataSources}
}

func (ts *testServer) do(query string, variables vars) string {
	json, err := json.Marshal(payload{Query: query, Variables: variables})
	if err != nil {
		ts.tt.Fatalf("could not marshal payload: %v", err)
	}
	req := httptest.NewRequest("POST", "http://example.com/graphql", bytes.NewReader(json))
	w := httptest.NewRecorder()

	graphql.Handler(ts.dataSources, ts.s)(w, req)

	resp, err := ioutil.ReadAll(w.Body)
	if err != nil {
		ts.tt.Fatalf("could not read response: %v", err)
	}
	return string(resp)
}

func TestGraphQLNotAValidGraphQLQuery(t *testing.T) {
	testServer := NewTestServer(t, db.DataSources{}, nil)
	actual := testServer.do("{ notAValidQuery }", vars{})

	testFile := testutil.NewGoldenFile(t, "not_a_valid_query.json")
	expected := testFile.Load()
	if *update {
		testFile.Write(actual)
		expected = actual
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Unexpected result, diff: %v", testutil.Diff(expected, actual))
	}
}
