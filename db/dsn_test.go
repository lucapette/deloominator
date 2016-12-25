package db_test

import (
	"reflect"
	"testing"

	"github.com/lucapette/deluminator/db"
)

func TestParseDSN(t *testing.T) {
	sourceTests := []struct {
		source   string
		expected *db.DSN
	}{
		{ // full
			"postgresql://grace:hopper@hal9000:3000/test?foo=bar",
			&db.DSN{
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
			&db.DSN{
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
			&db.DSN{
				Driver:  "postgresql",
				Host:    "hal9000",
				Port:    3000,
				DBName:  "test",
				Options: "foo=bar",
			},
		},
	}

	for _, test := range sourceTests {
		actual, err := db.ParseDSN(test.source)

		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(actual, test.expected) {
			t.Fatalf("Expected %v to equal %v", actual, test.expected)
		}
	}
}
