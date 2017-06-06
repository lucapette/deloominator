package testutil

import (
	"github.com/kelseyhightower/envconfig"
)

type testConfig struct {
	MysqlUser string `default:"deloominator"`
	MysqlPass string `default:"secret"`
	PGUser    string `default:"deloominator"`
	PGPass    string `default:"secret"`
}

func getTestConfig() (*testConfig, error) {
	cfg := &testConfig{}

	err := envconfig.Process("test", cfg)
	return cfg, err
}
