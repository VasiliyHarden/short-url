package main

import (
	"github.com/VasiliyHarden/short-url/internal/config"
	"github.com/VasiliyHarden/short-url/internal/handler"
	"github.com/VasiliyHarden/short-url/internal/repository/postgres"
	"github.com/VasiliyHarden/short-url/internal/service/shortener"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()

	gen := shortener.NewHashGenerator()

	db, err := postgres.NewRepository(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	var store shortener.Storage
	store, err = shortener.NewFileStorage(cfg.FileStoragePath)
	if err != nil {
		log.Printf("Failed to create file storage: %v, falling back to memory storage", err)
		store = shortener.NewMemoryStorage()
	}
	defer store.Close()

	sh := shortener.NewService(cfg.BaseURL, gen, store)
	logger := config.NewLogger()
	defer logger.Sync()

	router := handler.NewRouter(sh, logger, db)

	if err := http.ListenAndServe(cfg.Addr, router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
