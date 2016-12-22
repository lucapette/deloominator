package main

import (
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
	"github.com/lucapette/deluminator/api"
)

func main() {
	log.Info("starting deluminator")

	api.Start(3000)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}
