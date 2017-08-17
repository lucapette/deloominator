package config

import (
	"github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

// BinaryName is the name of the program
const BinaryName = "deloominator"

// Config represents the available options `deloominator` can read from the
// environment
type Config struct {
	Port      int      `default:"3000"`
	Sources   []string `envconfig:"data_sources" required:"true"`
	Storage   string   `envconfig:"storage"`
	LogFormat string   `default:"JSON" split_words:"true"`
	Debug     bool     `default:"false"`
	SkipOpen  bool     `default:"false" split_words:"true"`
}

// GetConfig reads the configuration from the environment and populates a Config
// structure
func GetConfig() (*Config, error) {
	cfg := Config{}
	err := envconfig.Process(BinaryName, &cfg)
	if err != nil {
		return &cfg, err
	}

	switch cfg.LogFormat {
	case "JSON":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case "TEXT":
		logrus.SetFormatter(&logrus.TextFormatter{})
	default:
		logrus.Fatalf("unknown log format: %s\n", cfg.LogFormat)
	}

	return &cfg, err
}

func Help() error {
	return envconfig.Usage(BinaryName, &Config{})
}
