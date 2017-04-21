package config

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

const BinaryName = "deloominator"

type Config struct {
	Port      int      `default:"3000"`
	Sources   []string `envconfig:"data_sources"`
	LogFormat string   `default:"JSON" split_words:"true"`
	Debug     bool     `default:"false"`
}

func GetConfig() *Config {
	var cfg Config

	err := envconfig.Process(BinaryName, &cfg)

	if err != nil {
		log.Fatal(err.Error())
	}

	switch cfg.LogFormat {
	case "JSON":
		log.SetFormatter(&log.JSONFormatter{})
	case "TEXT":
		log.SetFormatter(&log.TextFormatter{})
	default:
		log.Fatalf("unknown log format: %s\n", cfg.LogFormat)
	}

	return &cfg
}
