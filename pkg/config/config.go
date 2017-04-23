package config

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

const BinaryName = "deloominator"

type Config struct {
	Port      int      `default:"3000"`
	Sources   []string `envconfig:"data_sources" required:"true"`
	LogFormat string   `default:"JSON" split_words:"true"`
	Debug     bool     `default:"false"`
}

func GetConfig() (*Config, error) {
	var cfg Config

	err := envconfig.Process(BinaryName, &cfg)
	if err != nil {
		return &cfg, err
	}

	switch cfg.LogFormat {
	case "JSON":
		log.SetFormatter(&log.JSONFormatter{})
	case "TEXT":
		log.SetFormatter(&log.TextFormatter{})
	default:
		log.Fatalf("unknown log format: %s\n", cfg.LogFormat)
	}

	return &cfg, err
}

func Help() error {
	var cfg Config
	return envconfig.Usage(BinaryName, &cfg)
}
