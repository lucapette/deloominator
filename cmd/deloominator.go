package main

import (
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
	"github.com/lucapette/deloominator/pkg/api"
	"github.com/lucapette/deloominator/pkg/config"
	"github.com/lucapette/deloominator/pkg/db"
)

const appName = "deloominator"

func main() {
	cfg := config.GetConfig()
	log.WithField("port", cfg.Port).
		Infof("starting %s", config.BinaryName)

	dataSources, err := db.NewDataSources(cfg.Sources)
	if err != nil {
		log.Fatal(err.Error())
	}

	api.Start(cfg, dataSources)

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)
	<-s

	log.WithField("port", cfg.Port).
		Infof("stopping %s", config.BinaryName)

	dataSources.Shutdown()
}
