package main

import (
	"github.com/spf13/viper"
)

const (
	defaultConfigFileName = "config.yaml"
	defaultBatchSize      = 500
	defaultBufferSize     = 1024 * 1024
)

type CounterConfig struct {
	BufferSize int `mapstructure:"buffer_size"` // MB
	BatchSize  int `mapstructure:"batch_size"`  // number of lines
}

type Config struct {
	SourceURI         string `mapstructure:"source_uri"` // path to file
	BitmapThresholdMB int    `mapstructure:"bitmapThresholdMB"`
	Counter           CounterConfig
}

// LoadConfig reads YAML (via viper) and applies defaults & validation.
func LoadConfig(path string) (Config, error) {
	var cfg Config

	viper.SetConfigFile(path)
	viper.AutomaticEnv() // allow env overrides if needed

	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	if cfg.Counter.BufferSize <= 0 {
		cfg.Counter.BufferSize = defaultBufferSize
	}
	if cfg.Counter.BatchSize <= 0 {
		cfg.Counter.BatchSize = defaultBatchSize
	}

	return cfg, nil
}
