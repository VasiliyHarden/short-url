package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr    string
	BaseURL string
}

func Load() *Config {
	cfg := &Config{
		Addr:    getEnv("ADDR", ":8080"),
		BaseURL: getEnv("BASE_URL", "http://localhost:8080"),
	}

	flag.StringVar(&cfg.Addr, "a", cfg.Addr, "HTTP server address")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base URL for short links")

	flag.Parse()
	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
