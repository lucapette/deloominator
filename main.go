package main

import (
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
	"github.com/lucapette/deluminator/api"
	"github.com/lucapette/deluminator/app"
)

func main() {
	app := app.NewApp()
	log.WithField("port", app.Opts.Port).
		Infof("starting %s", app.Name)

	api.Start(app)

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)
	<-s

	log.WithField("port", app.Opts.Port).
		Infof("stopping %s", app.Name)

	app.Shutdown()
}
