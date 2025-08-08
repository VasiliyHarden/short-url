package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr    string
	BaseURL string
}

const (
	defaultAddr    = ":8080"
	defaultBaseURL = "http://localhost:8080"
)

func Load() *Config {
	cfg := &Config{
		Addr:    defaultAddr,
		BaseURL: defaultBaseURL,
	}

	flag.StringVar(&cfg.Addr, "a", cfg.Addr, "HTTP server address")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base URL for short links")
	flag.Parse()

	parseEnv(&cfg.Addr, "SERVER_ADDRESS")
	parseEnv(&cfg.BaseURL, "BASE_URL")

	return cfg
}

func parseEnv(target *string, key string) {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		*target = v
	}
}
