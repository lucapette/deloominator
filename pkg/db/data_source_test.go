package db_test

import (
	"testing"

	"github.com/lucapette/deloominator/pkg/db"
	"github.com/lucapette/deloominator/pkg/testutil"
)

func TestNewDataSources(t *testing.T) {
	pg, pgClean := testutil.SetupDB(t, db.PostgresDriver)
	my, myClean := testutil.SetupDB(t, db.MySQLDriver)
	sources := []string{pg, my}

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
