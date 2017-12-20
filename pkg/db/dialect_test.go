package db_test

import (
	"reflect"
	"testing"

	"github.com/lucapette/deloominator/pkg/db"
	"github.com/lucapette/deloominator/pkg/testutil"
)

func TestDialectsTables(t *testing.T) {
	expected := []string{"actor", "address", "category", "city", "country", "customer", "film", "film_actor", "film_category", "inventory", "language", "payment", "rental", "staff", "store"}

	dataSources, cleanup := testutil.SetupDataSources(t)
	defer cleanup()

	for _, dataSource := range dataSources {
		t.Run(dataSource.DriverName(), func(t *testing.T) {
			actual, err := dataSource.Tables()
			if err != nil {
				t.Fatalf("could not retrive tables for %s: %v", dataSource.DBName(), err)
			}

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("Unexpected result, diff: %v", testutil.Diff(expected, actual))
			}
		})
	}
}

func TestDialectsQuery(t *testing.T) {
	expected := db.QueryResult{
		Columns: []db.Column{
			{Name: "film_id", Type: db.Number},
			{Name: "title", Type: db.Text},
			{Name: "rental_rate", Type: db.Number},
			{Name: "last_update", Type: db.Time},
		},
		Rows: []db.Row{
			{
				{Value: "42"},
				{Value: "Back to the future"},
				{Value: "4.99"},
				{Value: "1985-10-21 07:28:00"},
			},
		},
	}

	dataSources, cleanup := testutil.SetupDataSources(t)
	defer cleanup()

	for _, dataSource := range dataSources {
		t.Run(dataSource.DriverName(), func(t *testing.T) {
			testutil.LoadData(t, dataSource, "film", expected)

			actual, err := dataSource.Query(db.Input{Query: "select film_id, title, rental_rate, last_update from film"})
			if err != nil {
				t.Fatalf("could not query %s: %v", dataSource.DBName(), err)
			}

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("Unexpected result, diff: %v", testutil.Diff(expected, actual))
			}
		})
	}
}
