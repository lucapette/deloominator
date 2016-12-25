package main

import (
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
	"github.com/lucapette/deluminator/api"
	"github.com/lucapette/deluminator/db"
)

type Config struct {
	Port        int `default:"3000"`
	DataSources []string
}

func main() {
	var c Config
	err := envconfig.Process("deluminator", &c)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.WithField("port", c.Port).Info("starting deluminator")

	sources, err := db.NewSources(c.DataSources)

	if err != nil {
		log.Fatal(err.Error())
	}

	api.Start(&api.Config{
		Port:    c.Port,
		Sources: sources,
	})

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)
	<-s
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}
