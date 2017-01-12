package db_test

import (
	"reflect"
	"testing"

	"github.com/lucapette/deluminator/db"
	"github.com/lucapette/deluminator/testutil"
)

func TestPGTables(t *testing.T) {
	dsn, cleanup := testutil.SetupDB(db.Postgres, t)
	pg, err := db.NewPGLoader(dsn)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = pg.Close()
		if err != nil {
			t.Fatal(err)
		}
		cleanup()
	}()

	expected := []string{"event_types", "user_events", "users"}
	actual, err := pg.Tables()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}
