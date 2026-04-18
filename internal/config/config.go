package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL string
	Port        int
	LogLevel    string
	WorkerID    string
}

// Load reads config from environment variables. DATABASE_URL is required.
func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL: getenv("DATABASE_URL", ""),
		LogLevel:    getenv("LOG_LEVEL", "info"),
		WorkerID:    getenv("WORKER_ID", ""),
	}

	portStr := getenv("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PORT %q: %w", portStr, err)
	}
	cfg.Port = port

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	return cfg, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
