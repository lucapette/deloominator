package handlers_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/lucapette/deloominator/pkg/api/handlers"
	"github.com/lucapette/deloominator/pkg/testutil"
)

func TestQueryEvaluatorNoVars(t *testing.T) {
	dataSources, cleanup := testutil.SetupDataSources(t)
	defer cleanup()

	cleanupSrv := testutil.SetupServer(dataSources)
	defer cleanupSrv()

	for _, dataSource := range dataSources {
		t.Run(dataSource.Dialect.DriverName()+"-no variables", func(t *testing.T) {
			testFile := testutil.NewGoldenFile(t, "no-variables.json")
			expected := testFile.Load()
			json, err := json.Marshal(handlers.QueryPayload{Query: "select * from film"})
			if err != nil {
				t.Fatalf("could not marshal payload: %v", err)
			}
			resp, err := http.Post("http://localhost:3000/query/evaluate", "application/json", bytes.NewReader(json))
			if err != nil {
				t.Errorf("cannot create request: %v", err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}
			actual := string(body)
			if *update {
				testFile.Write(actual)
				expected = actual
			}

			if resp.StatusCode != 200 {
				t.Errorf("expected 200, but got %v", resp.StatusCode)
			}

			if !reflect.DeepEqual(expected, actual) {
				t.Fatalf("Unexpected result, diff: %v", testutil.Diff(expected, actual))
			}
		})
	}
}

func TestQueryEvaluatorVariables(t *testing.T) {
	dataSources, cleanup := testutil.SetupDataSources(t)
	defer cleanup()

	cleanupSrv := testutil.SetupServer(dataSources)
	defer cleanupSrv()

	for _, dataSource := range dataSources {
		t.Run(dataSource.Dialect.DriverName()+"-variables", func(t *testing.T) {
			jsonPayload, err := json.Marshal(handlers.QueryPayload{
				Query: "select * from film where last_udpate > {timestamp}",
			})
			if err != nil {
				t.Fatalf("could not marshal payload: %v", err)
			}
			resp, err := http.Post("http://localhost:3000/query/evaluate", "application/json", bytes.NewReader(jsonPayload))
			if err != nil {
				t.Errorf("cannot create request: %v", err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}
			if resp.StatusCode != 200 {
				t.Errorf("expected 200, but got %v", resp.StatusCode)
			}
			type payload struct {
				Variables struct {
					Timestamp time.Time `json:"{timestamp}"`
				} `json:"variables"`
			}

			p := payload{}
			err = json.Unmarshal(body, &p)
			if err != nil {
				t.Errorf("could not parse json response: %v", err)
			}

			if p.Variables.Timestamp.IsZero() {
				t.Errorf("expected timestamp to have a value, but did not")
			}
		})
	}
}
