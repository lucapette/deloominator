package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/Sirupsen/logrus"
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
	helpFlag := flag.Bool("help", false, "show help")

	flag.Parse()

	if *helpFlag {
		help()
	}

	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("cannot read configuration: %v\n", err)

		help()
	}

	logrus.WithField("port", cfg.Port).
		Infof("starting %s", config.BinaryName)

	dataSources, err := db.NewDataSources(cfg.Sources)
	if err != nil {
		logrus.Fatalf("cannot create dataSources from %v: %v", cfg.Sources, err)
	}

	storage, err := db.NewStorage(cfg.Storage)
	if err != nil {
		logrus.Printf("cannot create storage from %v: %v\n", cfg.Storage, err)
	}

	options := []api.Option{
		api.Port(cfg.Port),
		api.Debug(cfg.Debug),
		api.DataSources(dataSources),
	}

	if storage != nil {
		if err := storage.AutoUpgrade(); err != nil {
			logrus.Printf("could not upgrade %s storage: %v", config.BinaryName, err)
		} else {
			options = append(options, api.Storage(storage))
		}
	}

	server := api.NewServer(options)

	server.Start()

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)
	<-s

	logrus.WithField("port", cfg.Port).
		Infof("stopping %s", config.BinaryName)

	dataSources.Close()
}
