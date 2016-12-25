package api_test

import (
	"testing"

	"github.com/lucapette/deluminator/api"
	"github.com/sclevine/agouti"
)

func TestServerStart(t *testing.T) {
	driver := agouti.PhantomJS()
	if err := driver.Start(); err != nil {
		t.Fatal("Failed to start PhantomJS:", err)
	}

	page, err := driver.NewPage(agouti.Browser("firefox"))
	if err != nil {
		t.Fatal("Failed to open page:", err)
	}

	api.Start(&api.Config{Port: 3000})

	if err := page.Navigate("http://localhost:3000"); err != nil {
		t.Fatal("Failed to navigate:", err)
	}

	expected := "It works!"

	actual, err := page.Find("#welcome").Text()
	if err != nil {
		t.Fatal("Failed to load home page", err)
	}

	if actual != expected {
		t.Fatalf("%s expected but got %s", expected, actual)
	}

	if err := driver.Stop(); err != nil {
		t.Fatal("Failed to close pages and stop WebDriver:", err)
	}
}
