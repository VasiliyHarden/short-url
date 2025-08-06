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
	store := shortener.NewMemoryStorage()

	sh := shortener.NewService(cfg.BaseURL, gen, store)

	router := handler.NewRouter(sh)

	if err := http.ListenAndServe(cfg.Addr, router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
