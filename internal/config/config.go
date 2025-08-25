package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr            string
	BaseURL         string
	FileStoragePath string
	DatabaseDSN     string
}

const (
	defaultAddr            = ":8080"
	defaultBaseURL         = "http://localhost:8080"
	defaultFileStoragePath = "./storage.json"
)

func Load() *Config {
	cfg := &Config{
		Addr:            defaultAddr,
		BaseURL:         defaultBaseURL,
		FileStoragePath: defaultFileStoragePath,
	}

	flag.StringVar(&cfg.Addr, "a", cfg.Addr, "HTTP server address")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base URL for short links")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "File storage path")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "Database DSN")
	flag.Parse()

	parseEnv(&cfg.Addr, "SERVER_ADDRESS")
	parseEnv(&cfg.BaseURL, "BASE_URL")
	parseEnv(&cfg.FileStoragePath, "FILE_STORAGE_PATH")
	parseEnv(&cfg.DatabaseDSN, "DATABASE_DSN")

	return cfg
}

func parseEnv(target *string, key string) {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		*target = v
	}
}
