package db_test

import (
	"fmt"
	"testing"

	"github.com/lucapette/deluminator/db"
	"github.com/lucapette/deluminator/testutil"
)

func TestNewLoaders(t *testing.T) {
	pg, pgClean := testutil.SetupDB(db.Postgres, t)
	my, myClean := testutil.SetupDB(db.MySQL, t)
	dataSources := []string{
		pg.Format(),
		fmt.Sprintf("mysql://%s:%s@%s/%s", my.Username, my.Pass, my.Host, my.DBName), // Format() does not work both ways yet
	}

	loaders, err := db.NewLoaders(dataSources)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		for _, loader := range loaders {
			err = loader.Close()

			if err != nil {
				t.Fatal(err)
			}
		}

		pgClean()
		myClean()
	}()

	actual := len(loaders)
	expected := len(dataSources)

	if actual != expected {
		t.Fatalf("Expected %d, got %d", expected, actual)
	}
}
