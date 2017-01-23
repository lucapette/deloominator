package db_test

import (
	"fmt"
	"testing"

	"github.com/lucapette/deluminator/pkg/db"
	"github.com/lucapette/deluminator/pkg/testutil"
)

func TestNewDataSources(t *testing.T) {
	pg, pgClean := testutil.SetupDB(db.PostgresDriver, t)
	my, myClean := testutil.SetupDB(db.MySQLDriver, t)
	sources := []string{
		pg.Format(),
		fmt.Sprintf("mysql://%s:%s@%s/%s", my.Username, my.Pass, my.Host, my.DBName), // Format() does not work both ways yet
	}

	dataSources, err := db.NewDataSources(sources)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		for _, ds := range dataSources {
			err = ds.Close()

			if err != nil {
				t.Fatal(err)
			}
		}

		pgClean()
		myClean()
	}()

	actual := len(dataSources)
	expected := len(sources)

	if actual != expected {
		t.Fatalf("Expected %d, got %d", expected, actual)
	}
}
