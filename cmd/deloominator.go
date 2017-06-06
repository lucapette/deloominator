package main

import (
	"fmt"
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
	"github.com/lucapette/deloominator/pkg/api"
	"github.com/lucapette/deloominator/pkg/config"
	"github.com/lucapette/deloominator/pkg/db"
	flag "github.com/spf13/pflag"
)

func help() {
	err := config.Help()
	if err != nil {
		fmt.Printf("cannot read configuration %v", err)
	}

	os.Exit(1)
}

func main() {
	var helpFlag = flag.Bool("help", false, "show help")

	flag.Parse()

	if *helpFlag {
		help()
	}

	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("cannot read configuration: %v\n", err)

		help()
	}

	log.WithField("port", cfg.Port).
		Infof("starting %s", config.BinaryName)

	dataSources, err := db.NewDataSources(cfg.Sources)
	if err != nil {
		log.Fatalf("cannot create DataSources: %v", err)
	}

	api.Start(cfg, dataSources)

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)
	<-s

	log.WithField("port", cfg.Port).
		Infof("stopping %s", config.BinaryName)

	dataSources.Shutdown()
}
