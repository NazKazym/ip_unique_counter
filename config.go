package main

import (
	"github.com/spf13/viper"
)

const (
	defaultConfigFileName = "config.yaml"
	defaultBufferSizeMB   = 64
)

type CounterConfig struct {
	BufferSizeMB int `mapstructure:"buffer_size_MB"`
}

type Config struct {
	SourceURI string `mapstructure:"source_uri"`
	Counter   CounterConfig
}

// LoadConfig reads YAML (via viper) and applies defaults & validation.
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

	if cfg.Counter.BufferSizeMB <= 0 {
		cfg.Counter.BufferSizeMB = defaultBufferSizeMB
	}

	return cfg, nil
}
