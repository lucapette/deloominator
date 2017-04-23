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

const appName = "deloominator"

var helpFlag *bool = flag.Bool("help", false, "show help")

func help() {
	err := config.Help()
	if err != nil {
		fmt.Println(err.Error())
	}

	os.Exit(1)
}

func main() {
	if *helpFlag {
		help()
	}

	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println(err.Error())

		help()
	}

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

func init() {
	flag.Parse()
}
