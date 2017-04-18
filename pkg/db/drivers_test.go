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
				Rows:    []db.Row{db.Row{db.Cell{Value: "42"}, db.Cell{Value: "Anna"}, db.Cell{Value: "Torv"}}},
				Columns: []db.Column{db.Column{Name: "actor_id"}, db.Column{Name: "first_name"}, db.Column{Name: "last_name"}},
			}
			testutil.LoadData(t, driver, "actor", expected)

			actual, err := driver.Query("select actor_id, first_name, last_name from actor")
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("Expected %v, got %v", expected, actual)
			}
		})
	}
}
