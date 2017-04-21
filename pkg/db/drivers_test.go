package db_test

import (
	"reflect"
	"testing"

	"github.com/lucapette/deloominator/pkg/db"
	"github.com/lucapette/deloominator/pkg/testutil"
)

var drivers = []struct {
	name   string
	driver db.DriverType
}{
	{"postgres", db.PostgresDriver},
	{"mysql", db.MySQLDriver},
}

func TestDriversTables(t *testing.T) {
	for _, test := range drivers {
		t.Run(test.name, func(t *testing.T) {
			dsn, cleanup := testutil.SetupDB(t, test.driver)
			driver, err := db.NewDataSource(dsn)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				err = driver.Close()
				if err != nil {
					t.Fatal(err)
				}
				cleanup()
			}()

			expected := []string{"actor", "address", "category", "city", "country", "customer", "film", "film_actor", "film_category", "inventory", "language", "payment", "rental", "staff", "store"}
			qr, err := driver.Tables()
			if err != nil {
				t.Fatal(err)
			}

			actual := make([]string, len(qr.Rows))

			for i, row := range qr.Rows {
				actual[i] = row[0].Value
			}

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("Expected %v, got %v", expected, actual)
			}
		})
	}
}

func TestDriversQuery(t *testing.T) {
	for _, test := range drivers {
		t.Run(test.name, func(t *testing.T) {
			dsn, cleanup := testutil.SetupDB(t, test.driver)
			driver, err := db.NewDataSource(dsn)
			if err != nil {
				t.Fatal(err)
			}

			defer func() {
				err = driver.Close()
				if err != nil {
					t.Fatal(err)
				}
				cleanup()
			}()

			expected := db.QueryResult{
				Columns: []db.Column{
					{Name: "film_id", Type: db.Number},
					{Name: "title", Type: db.Text},
					{Name: "rental_rate", Type: db.Number},
					{Name: "last_update", Type: db.Time},
				},
				Rows: []db.Row{
					{
						db.Cell{Value: "42"},
						db.Cell{Value: "Back to the future"},
						db.Cell{Value: "4.99"},
						db.Cell{Value: "1985-10-21 07:28:00"},
					},
				},
			}

			testutil.LoadData(t, driver, "film", expected)

			actual, err := driver.Query("select film_id, title, rental_rate, last_update from film")
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("Expected %v, got %v", expected, actual)
			}
		})
	}
}
