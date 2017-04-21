package app

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
	"github.com/lucapette/deloominator/pkg/db"
)

const Name = "deloominator"

type Options struct {
	Port      int      `default:"3000"`
	Sources   []string `envconfig:"data_sources"`
	LogFormat string   `default:"JSON" split_words:"true"`
	Debug     bool     `default:"false"`
}

type App struct {
	Name        string
	Opts        Options
	dataSources db.DataSources
}

func NewApp() *App {
	var opts Options

	err := envconfig.Process(Name, &opts)

	if err != nil {
		log.Fatal(err.Error())
	}

	switch opts.LogFormat {
	case "JSON":
		log.SetFormatter(&log.JSONFormatter{})
	case "TEXT":
		log.SetFormatter(&log.TextFormatter{})
	default:
		log.Fatalf("unknown log format: %s\n", opts.LogFormat)
	}

	return &App{
		Name: Name,
		Opts: opts,
	}
}
