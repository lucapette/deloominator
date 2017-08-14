package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/Sirupsen/logrus"
	"github.com/lucapette/deloominator/pkg/api"
	"github.com/lucapette/deloominator/pkg/config"
	"github.com/lucapette/deloominator/pkg/db"
	"github.com/lucapette/deloominator/pkg/db/storage"
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

	var s *storage.Storage
	if cfg.Storage != "" {
		s, err = storage.NewStorage(cfg.Storage)
		if err != nil {
			logrus.Warnf("cannot create storage from %v: %v", cfg.Storage, err)
		}
	}

	options := []api.Option{
		api.Port(cfg.Port),
		api.Debug(cfg.Debug),
		api.DataSources(dataSources),
	}

	if s != nil {
		if err := s.AutoUpgrade(); err != nil {
			logrus.Printf("could not upgrade %s storage: %v", config.BinaryName, err)
		} else {
			options = append(options, api.Storage(s))
		}
	}

	server := api.NewServer(options)

	server.Start()

	sgn := make(chan os.Signal, 1)
	signal.Notify(sgn, os.Interrupt, os.Kill)
	<-sgn

	logrus.WithField("port", cfg.Port).Infof("stopping %s", config.BinaryName)

	dataSources.Close()
}
