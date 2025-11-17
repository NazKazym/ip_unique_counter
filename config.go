package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	SourceURI string `mapstructure:"source_uri"`
}

const defaultConfigFileName = "config.yaml"

// LoadConfig reads YAML (via viper) and applies minimal validation.
func LoadConfig(path string) (Config, error) {
	var cfg Config

	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	if cfg.SourceURI == "" {
		return cfg, fmt.Errorf("source_uri must be provided in config")
	}

	return cfg, nil
}
