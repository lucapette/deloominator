package main

import (
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
	"github.com/lucapette/deloominator/pkg/api"
	"github.com/lucapette/deloominator/pkg/app"
	"github.com/lucapette/deloominator/pkg/db"
)

func main() {
	app := app.NewApp()
	log.WithField("port", app.Opts.Port).
		Infof("starting %s", app.Name)

	dataSources, err := db.NewDataSources(app.Opts.Sources)
	if err != nil {
		log.Fatal(err.Error())
	}

	api.Start(app, dataSources)

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)
	<-s

	log.WithField("port", app.Opts.Port).
		Infof("stopping %s", app.Name)

	dataSources.Shutdown()
}
