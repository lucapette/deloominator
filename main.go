package main

import (
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
	"github.com/lucapette/deluminator/api"
)

type Config struct {
	Port int `default:"3000"`
}

func main() {
	var c Config
	err := envconfig.Process("deluminator", &c)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.WithField("port", c.Port).Info("starting deluminator")

	api.Start(c.Port)

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)
	<-s
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}
