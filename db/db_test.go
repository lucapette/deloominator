package db_test

import (
	"testing"

	"github.com/lucapette/deluminator/db"
)

func TestNewSource(t *testing.T) {
	dataSources := []string{"postgresql://localhost/test"}
	sources, err := db.NewSources(dataSources)

	if err != nil {
		t.Fatal(err)
	}

	actual := len(sources)
	expected := len(dataSources)

	if actual != expected {
		t.Fatalf("Expected %d, got %d", expected, actual)
	}
}
