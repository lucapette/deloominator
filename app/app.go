package app

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
	"github.com/lucapette/deluminator/db"
)

const Name = "deluminator"

type Options struct {
	Port      int      `default:"3000"`
	Loaders   []string `envconfig:"data_sources"`
	LogFormat string   `default:"JSON" split_words:"true"`
	Debug     bool     `default: "false"`
}

var Opts Options
var dataSources db.Loaders

func Init() {
	err := envconfig.Process(Name, &Opts)

	if err != nil {
		log.Fatal(err.Error())
	}

	switch Opts.LogFormat {
	case "JSON":
		log.SetFormatter(&log.JSONFormatter{})
	case "TEXT":
		log.SetFormatter(&log.TextFormatter{})
	default:
		log.Fatalf("unknown log format: %s\n", Opts.LogFormat)
	}
}

func Sources() db.Loaders {
	if len(dataSources) > 0 {
		return dataSources
	}

	ds, err := db.NewLoaders(Opts.Loaders)
	if err != nil {
		log.Fatal(err.Error())
	}

	dataSources = ds

	return dataSources
}

func Shutdown() {
	for _, loader := range dataSources {
		err := loader.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
