package db_test

import (
	"reflect"
	"testing"

	"github.com/lucapette/deluminator/db"
)

func TestNewSource(t *testing.T) {
	sourceTests := []struct {
		source   string
		expected *db.DataSource
	}{
		{ // full
			"postgresql://grace:hopper@hal9000:3000/test?foo=bar",
			&db.DataSource{
				Driver:   "postgresql",
				Username: "grace",
				Pass:     "hopper",
				Host:     "hal9000",
				Port:     3000,
				DBName:   "test",
				Options:  "foo=bar",
			},
		},
		{ // no port
			"postgresql://grace:hopper@hal9000/test?foo=bar",
			&db.DataSource{
				Driver:   "postgresql",
				Username: "grace",
				Pass:     "hopper",
				Host:     "hal9000",
				DBName:   "test",
				Options:  "foo=bar",
			},
		},
		{ // no credentials
			"postgresql://hal9000:3000/test?foo=bar",
			&db.DataSource{
				Driver:  "postgresql",
				Host:    "hal9000",
				Port:    3000,
				DBName:  "test",
				Options: "foo=bar",
			},
		},
	}

	for _, test := range sourceTests {
		actual, err := db.NewDataSource(test.source)

		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(actual, test.expected) {
			t.Fatalf("Expected %v to equal %v", actual, test.expected)
		}
	}
}

func TestNewSources(t *testing.T) {
	dataSources := []string{"postgresql://localhost/test"}
	sources, err := db.NewDataSources(dataSources)

	if err != nil {
		t.Fatal(err)
	}

	actual := len(sources)
	expected := len(dataSources)

	if actual != expected {
		t.Fatalf("Expected %d, got %d", expected, actual)
	}
}
