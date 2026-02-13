package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port       int
	CacheTTL   time.Duration
	JuheAPIKey string
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

	juheAPIKey := os.Getenv("JUHE_API_KEY")

	return &Config{
		Port:       port,
		CacheTTL:   cacheTTL,
		JuheAPIKey: juheAPIKey,
	}
}
