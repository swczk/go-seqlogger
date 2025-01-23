package seqlogger

import (
	"log/slog"
	"time"
)

type Config struct {
	Endpoint      string
	APIKey        string
	LogLevel      slog.Level
	AddSource     bool
	RequestIDKey  string
	ClientTimeout time.Duration
}

func DefaultConfig(endpoint string) Config {
	return Config{
		Endpoint:      endpoint,
		LogLevel:      slog.LevelInfo,
		AddSource:     false,
		ClientTimeout: 5 * time.Second,
	}
}

func (c Config) WithAPIKey(apiKey string) Config {
	c.APIKey = apiKey
	return c
}

func (c Config) WithLogLevel(level slog.Level) Config {
	c.LogLevel = level
	return c
}

func (c Config) WithSourceTracking() Config {
	c.AddSource = true
	return c
}

func (c Config) WithRequestIDKey(key string) Config {
	c.RequestIDKey = key
	return c
}
