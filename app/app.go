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

type App struct {
	Name        string
	Opts        Options
	DataSources db.Loaders
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

func (app *App) Sources() db.Loaders {
	if len(app.DataSources) > 0 {
		return app.DataSources
	}

	ds, err := db.NewLoaders(app.Opts.Loaders)
	if err != nil {
		log.Fatal(err.Error())
	}

	app.DataSources = ds

	return app.DataSources
}

func (app *App) Shutdown() {
	for _, loader := range app.DataSources {
		err := loader.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
