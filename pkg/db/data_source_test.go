package db_test

import (
	"testing"

	"github.com/lucapette/deloominator/pkg/db"
	"github.com/lucapette/deloominator/pkg/testutil"
)

func TestNewDataSources(t *testing.T) {
	pg, pgClean := testutil.SetupPG(t)
	my, myClean := testutil.SetupMySQL(t)
	sources := []string{pg, my}

	dataSources, err := db.NewDataSources(sources)
	if err != nil {
		t.Fatalf("could not create dataSources from %v: %v", sources, err)
	}

	defer func() {
		for _, ds := range dataSources {
			if err = ds.Close(); err != nil {
				t.Fatalf("could not close %s: %v", ds.Name(), err)
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
