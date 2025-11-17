package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/viper"
)

// Global defaults & shared constants
const (
	defaultNumBuckets        = 128
	defaultBatchLines        = 10_000
	defaultBucketDir         = "buckets"
	defaultJobChanMultiplier = 2

	defaultBucketChanBufSize = 1024
	bucketFileNamePattern    = "bucket_%04d.tmp"

	defaultConfigFileName = "config.yaml"
	defaultSourceTypeFile = "file"

	defaultLogLevel     = "info"
	defaultMaxRetries   = 3
	defaultRetryDelayMs = 500
)

type Config struct {
	SourceType string `mapstructure:"source_type"` // "file" for now
	SourceURI  string `mapstructure:"source_uri"`  // path to file

	BucketDir   string `mapstructure:"bucket_dir"`
	NumWorkers  int    `mapstructure:"num_workers"`
	NumBuckets  int    `mapstructure:"num_buckets"`
	BatchLines  int    `mapstructure:"batch_lines"`
	RemoveAfter bool   `mapstructure:"remove_after"`

	LogLevel     string `mapstructure:"log_level"`
	MaxRetries   int    `mapstructure:"max_retries"`
	RetryDelayMs int    `mapstructure:"retry_delay_ms"`
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

	// Defaults and validation
	if cfg.SourceType == "" {
		cfg.SourceType = defaultSourceTypeFile
	}
	if cfg.SourceURI == "" {
		return cfg, fmt.Errorf("source_uri must be provided in config")
	}
	if cfg.BucketDir == "" {
		cfg.BucketDir = defaultBucketDir
	}
	if cfg.NumBuckets <= 0 {
		cfg.NumBuckets = defaultNumBuckets
	}
	if cfg.BatchLines <= 0 {
		cfg.BatchLines = defaultBatchLines
	}
	if cfg.NumWorkers <= 0 {
		cfg.NumWorkers = runtime.NumCPU()
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = defaultLogLevel
	}
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = defaultMaxRetries
	}
	if cfg.RetryDelayMs <= 0 {
		cfg.RetryDelayMs = defaultRetryDelayMs
	}

	return cfg, nil
}
