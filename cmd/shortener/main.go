package main

import (
	"github.com/VasiliyHarden/short-url/internal/config"
	"github.com/VasiliyHarden/short-url/internal/handler"
	"github.com/VasiliyHarden/short-url/internal/service/shortener"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()

	gen := shortener.NewHashGenerator()

	var store shortener.Storage
	store, err := shortener.NewFileStorage(cfg.FileStoragePath)
	if err != nil {
		log.Printf("Failed to create file storage: %v, falling back to memory storage", err)
		store = shortener.NewMemoryStorage()
	}
	defer store.Close()

	sh := shortener.NewService(cfg.BaseURL, gen, store)
	logger := config.NewLogger()
	defer logger.Sync()

	router := handler.NewRouter(sh, logger)

	if err := http.ListenAndServe(cfg.Addr, router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
