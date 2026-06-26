package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	defaultDataDir = ".scout"
	defaultBaseURL = "https://huggingface.co"
)

type Config struct {
	DataDir   string
	BaseURL   string
	RunDate   string
	Timeout   time.Duration
	Retries   int
	RetryWait time.Duration
	UserAgent string
}

func FromEnv() (Config, error) {
	cfg := Config{
		DataDir:   envOrDefault("SCOUT_DATA_DIR", defaultDataDir),
		BaseURL:   envOrDefault("SCOUT_BASE_URL", defaultBaseURL),
		RunDate:   os.Getenv("SCOUT_RUN_DATE"),
		Timeout:   durationFromEnv("SCOUT_HTTP_TIMEOUT", 30*time.Second),
		Retries:   intFromEnv("SCOUT_HTTP_RETRIES", 3),
		RetryWait: durationFromEnv("SCOUT_HTTP_RETRY_WAIT", 2*time.Second),
		UserAgent: envOrDefault("SCOUT_HTTP_USER_AGENT", "scout/1"),
	}
	if cfg.Retries < 1 {
		cfg.Retries = 1
	}
	if cfg.DataDir == "" {
		cfg.DataDir = defaultDataDir
	}
	abs, err := filepath.Abs(cfg.DataDir)
	if err != nil {
		return Config{}, fmt.Errorf("resolve data dir: %w", err)
	}
	cfg.DataDir = abs
	return cfg, nil
}

func (c Config) ResolveRunDate(now time.Time) (time.Time, error) {
	if c.RunDate == "" {
		return now.UTC(), nil
	}
	date, err := time.Parse("2006-01-02", c.RunDate)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse run date %q: %w", c.RunDate, err)
	}
	return date.UTC(), nil
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func durationFromEnv(name string, fallback time.Duration) time.Duration {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return d
}

func intFromEnv(name string, fallback int) int {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return i
}
