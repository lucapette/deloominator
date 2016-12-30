package db_test

import (
	"testing"

	"github.com/lucapette/deluminator/db"
)

func TestNewLoaders(t *testing.T) {
	pg, cleanup := setupDB("postgres", t)
	dataSources := []string{
		pg.String(),
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

		cleanup()
	}()

	actual := len(loaders)
	expected := len(dataSources)

	if actual != expected {
		t.Fatalf("Expected %d, got %d", expected, actual)
	}
}
