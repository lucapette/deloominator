package api_test

import (
	"io/ioutil"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/lucapette/deluminator/api"
	"github.com/lucapette/deluminator/app"
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

	app := app.NewApp()
	api.Start(app)

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

func init() {
	log.SetOutput(ioutil.Discard)
}
