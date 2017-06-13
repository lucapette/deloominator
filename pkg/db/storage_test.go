package db_test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/lucapette/deloominator/pkg/db"
	"github.com/lucapette/deloominator/pkg/testutil"
)

func TestAutoUpgradeBoostrapDB(t *testing.T) {
	dataSources, cleanup := testutil.CreateDataSources(t)
	defer cleanup()

	for _, source := range dataSources {
		storage, err := db.NewStorage(source)
		if err != nil {
			t.Fatalf("could not create Storage from %s: %v", source, err)
		}
		defer storage.Close()

		if err := storage.AutoUpgrade(); err != nil {
			t.Fatalf("could not run AutoUpgrade for %s: %v", source, err)
		}

		dataSource, err := db.NewDataSource(source)
		if err != nil {
			t.Fatalf("could not create DataSource from %s: %v", source, err)
		}
		defer func() {
			if err := dataSource.Close(); err != nil {
				t.Fatalf("could not close DataSource: %v", err)
			}
		}()

		names, err := dataSource.Tables()
		if err != nil {
			t.Fatalf("could not query tables on %s: %v", source, err)
		}

		found := false

		for _, name := range names {
			if name == "migrations" {
				found = true
				break
			}
		}

		if !found {
			t.Fatalf("expected AutoMigrate to create a migrations table, but did not")
		}
	}
}
