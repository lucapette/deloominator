package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/lucapette/deloominator/pkg/api"
	"github.com/lucapette/deloominator/pkg/db"
	"github.com/lucapette/deloominator/pkg/testutil"
)

var update = flag.Bool("update", false, "update golden files")

var data = db.QueryResult{
	Columns: []db.Column{
		{Name: "title", Type: db.Text},
		{Name: "rating", Type: db.Text},
	},
	Rows: []db.Row{
		{
			{Value: "Back to the future"},
			{Value: "PG"},
		},
	},
}

func TestExportCSV(t *testing.T) {
	dataSources, cleanup := testutil.SetupDataSources(t)
	defer cleanup()

	options := []api.Option{
		api.Port(3000),
		api.DataSources(dataSources),
	}
	server := api.NewServer(options)
	server.Start()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer server.Stop(ctx)

	tests := []struct {
		name   string
		query  string
		golden string
	}{
		{"export-csv", "select title, rating from film", "film-rating.csv"},
	}

	for dataSourceName, dataSource := range dataSources {
		testutil.LoadData(t, dataSource, "film", data)
		for _, tc := range tests {
			t.Run(dataSource.Dialect.DriverName()+"-"+tc.name, func(t *testing.T) {
				testFile := testutil.NewGoldenFile(t, tc.golden)
				expected := testFile.Load()
				json, err := json.Marshal(api.ExportPayload{
					Source: dataSourceName,
					Query:  tc.query,
				})
				if err != nil {
					t.Fatalf("could not marshal payload: %v", err)
				}
				resp, err := http.Post("http://localhost:3000/export/csv", "application/json", bytes.NewReader(json))
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
}
