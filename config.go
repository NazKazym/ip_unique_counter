package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/viper"
)

// Global defaults & shared constants
const (
	defaultConfigFileName = "config.yaml"
	defaultSourceTypeFile = "file"

	defaultLogLevel     = "info"
	defaultMaxRetries   = 3
	defaultRetryDelayMs = 500

	// ===== новые дефолты для счётчика уникальных IP =====
	defaultCounterStorage         = "auto" // "auto" | "map" | "bitmap"
	defaultCounterBitmapThreshold = 5_000_000
)

// CounterConfig описывает поведение счётчика уникальных IP.
type CounterConfig struct {
	// Storage:
	//   - "auto"   — выбор по порогу;
	//   - "map"    — всегда map;
	//   - "bitmap" — всегда bitmap.
	Storage string `mapstructure:"storage"`

	// BitmapThreshold — если Storage="auto":
	//   если EstimatedItems <= BitmapThreshold -> map
	//   если > BitmapThreshold                -> bitmap
	BitmapThreshold uint64 `mapstructure:"bitmap_threshold"`
}

type Config struct {
	SourceType string `mapstructure:"source_type"` // "file" for now
	SourceURI  string `mapstructure:"source_uri"`  // path to file

	NumWorkers int `mapstructure:"num_workers"`

	LogLevel     string `mapstructure:"log_level"`
	MaxRetries   int    `mapstructure:"max_retries"`
	RetryDelayMs int    `mapstructure:"retry_delay_ms"`

	// ===== новый блок конфигурации для счётчика =====
	Counter CounterConfig `mapstructure:"counter"`
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

	if cfg.Counter.Storage == "" {
		cfg.Counter.Storage = defaultCounterStorage
	}
	if cfg.Counter.BitmapThreshold == 0 {
		cfg.Counter.BitmapThreshold = defaultCounterBitmapThreshold
	}

	return cfg, nil
}
