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
	Port      int      `default:"3000"`
	Loaders   []string `envconfig:"data_sources" required:"true"`
	LogFormat string   `default:"JSON" split_words:"true"`
	Debug     bool     `default: "false"`
}

var c Config

func main() {
	log.WithField("port", c.Port).Info("starting deluminator")

	sources, err := db.NewLoaders(c.Loaders)

	if err != nil {
		log.Fatal(err.Error())
	}

	api.Start(&api.Config{
		Port:    c.Port,
		Loaders: sources,
		Debug:   c.Debug,
	})

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)
	<-s
}

func init() {
	err := envconfig.Process("deluminator", &c)

	if err != nil {
		log.Fatal(err.Error())
	}

	switch c.LogFormat {
	case "JSON":
		log.SetFormatter(&log.JSONFormatter{})
	case "TEXT":
		log.SetFormatter(&log.TextFormatter{})
	default:
		log.Fatalf("unknown log format: %s\n", c.LogFormat)
	}
}
