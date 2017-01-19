package db_test

import (
	"reflect"
	"testing"

	"github.com/lucapette/deluminator/db"
	"github.com/lucapette/deluminator/testutil"
)

var drivers = []struct {
	name    string
	driver  db.DriverType
	factory func(*db.DSN) (db.DataSource, error)
}{
	{"postgres", db.PostgresDriver, func(dsn *db.DSN) (db.DataSource, error) { return db.NewPostgres(dsn) }},
	{"mysql", db.MySQLDriver, func(dsn *db.DSN) (db.DataSource, error) { return db.NewMySQL(dsn) }},
}

func TestDriversTables(t *testing.T) {
	for _, test := range drivers {
		t.Run(test.name, func(t *testing.T) {
			dsn, cleanup := testutil.SetupDB(test.driver, t)
			driver, err := test.factory(dsn)
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

			expected := []string{"event_types", "user_events", "users"}
			tables, err := driver.Tables()
			if err != nil {
				t.Fatal(err)
			}

			actual := make([]string, len(tables))

			for i, row := range tables {
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
			dsn, cleanup := testutil.SetupDB(test.driver, t)
			driver, err := test.factory(dsn)
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

			expected := db.Rows{db.Row{db.Column{Name: "id", Value: "42"}, db.Column{Name: "name", Value: "Grace Hopper"}}}
			actual, err := driver.Query("select id, name from users")
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("Expected %v, got %v", expected, actual)
			}
		})
	}
}
