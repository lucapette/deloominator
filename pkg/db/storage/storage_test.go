package storage_test

import (
	"testing"

	"github.com/lucapette/deloominator/pkg/db"
	"github.com/lucapette/deloominator/pkg/db/storage"

	"github.com/lucapette/deloominator/pkg/testutil"
)

func TestStorage_AutoUpgrade_DBExists(t *testing.T) {
	sources, cleanup := testutil.CreateDataSources(t)
	defer cleanup()

	for _, source := range sources {
		t.Run(source, func(t *testing.T) {
			ds, err := db.NewDataSource(source)
			if err != nil {
				t.Errorf("could not create data source: %v", err)
			}
			defer ds.Close()

			s, err := storage.NewStorage(source)
			if err != nil {
				t.Errorf("could not create storage: %v", err)
			}
			defer s.Close()

			err = s.AutoUpgrade()
			if err != nil {
				t.Errorf("expected AutoUpgrade to not return error, but did: %v", err)
			}

			qr, err := ds.Query(db.Input{Query: ds.TablesQuery()})
			if err != nil {
				t.Errorf("could not query tables: %v", err)
			}

			if len(qr.Rows) == 0 {
				t.Fatalf("expected %v to migrate, but did not", source)
			}
		})
	}
}

func TestStorage_AutoUpgrade_DBNotExist(t *testing.T) {
	names := testutil.CreateDSN(t)

	for _, dsn := range names {
		t.Run(dsn, func(t *testing.T) {
			ds, err := db.NewDataSource(dsn)
			if err != nil {
				t.Errorf("could not create data source: %v", err)
			}
			defer ds.Close()

			s, err := storage.NewStorage(dsn)
			if err != nil {
				t.Errorf("could not create storage: %v", err)
			}
			defer s.Close()

			err = s.AutoUpgrade()
			if err != nil {
				t.Errorf("expected AutoUpgrade to not return error, but did: %v", err)
			}

			qr, err := ds.Query(db.Input{Query: ds.TablesQuery()})
			if err != nil {
				t.Errorf("could not query tables: %v", err)
			}

			if len(qr.Rows) == 0 {
				t.Fatalf("expected %v to migrate, but did not", dsn)
			}
		})
	}
}
