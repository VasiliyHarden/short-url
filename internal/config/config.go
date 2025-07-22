package config

import (
	"flag"
)

type Config struct {
	Addr    string
	BaseURL string
}

func Load() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Addr, "a", ":8080", "HTTP server address")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Base URL for short links")

	flag.Parse()

	return cfg
}
