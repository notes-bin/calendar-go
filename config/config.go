package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port     int
	CacheTTL time.Duration
}

func Load() *Config {
	port := 8080
	if portStr := os.Getenv("PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	cacheTTL := 6 * time.Hour
	if ttlStr := os.Getenv("CACHE_TTL"); ttlStr != "" {
		if d, err := time.ParseDuration(ttlStr); err == nil {
			cacheTTL = d
		}
	}

	return &Config{
		Port:     port,
		CacheTTL: cacheTTL,
	}
}
