package db_test

import (
	"reflect"
	"testing"

	"github.com/lucapette/deluminator/db"
	"github.com/lucapette/deluminator/testutil"
)

func TestMyTables(t *testing.T) {
	dsn, cleanup := testutil.SetupDB(db.MySQLDriver, t)
	my, err := db.NewMySQL(dsn)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := my.Close()
		if err != nil {
			t.Fatal(err)
		}
		cleanup()
	}()

	expected := []string{"event_types", "user_events", "users"}
	actual, err := my.Tables()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}
