package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

// A little experiment in documentation automation The plan is to make it part
// of the build and automate screenshot creation. Right now, it is hard to run
// and I would not consider it as an "official" part of deloominator as it is
// not ready yet.

func saveScreenshot(wd selenium.WebDriver, path string, filename string) {
	if err := wd.Get("http://localhost:3000" + path); err != nil {
		log.Fatalf("could not connect to application: %v", err)
	}

	if err := wd.MaximizeWindow(""); err != nil {
		log.Fatalf("could not maximize window: %v", err)
	}

	time.Sleep(1 * time.Second)

	output, err := wd.Screenshot()
	if err != nil {
		log.Fatalf("could not take screenshot of the playground: %v", err)
	}

	if err := ioutil.WriteFile("./img/"+filename, output, 0644); err != nil {
		log.Fatalf("could not write image: %v", err)
	}
}

func main() {
	chromeOptions := chrome.Capabilities{
		MobileEmulation: &chrome.MobileEmulation{
			DeviceMetrics: chrome.DeviceMetrics{Width: 1200, Height: 1400},
		},
	}
	capabilities := selenium.Capabilities{
		"browserName":   "chrome",
		"chromeOptions": chromeOptions,
	}

	// Connect to the WebDriver instance running locally.
	wd, err := selenium.NewRemote(capabilities, "http://localhost:9515")
	if err != nil {
		log.Fatalf("could not connect to webdriver: %v", err)
	}
	defer wd.Quit()

	saveScreenshot(wd, "/playground", "playground.png")
	saveScreenshot(wd, "/questions", "questions.png")
	saveScreenshot(wd, "/questions/1", "question.png")
}
