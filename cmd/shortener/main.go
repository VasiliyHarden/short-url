package main

import (
	"github.com/VasiliyHarden/short-url/internal/config"
	"github.com/VasiliyHarden/short-url/internal/handler"
	"github.com/VasiliyHarden/short-url/internal/service/shortener"
	"net/http"
)

func main() {
	cfg := config.Load()
	shortener.Init(cfg.BaseURL)
	router := handler.NewRouter()

	err := http.ListenAndServe(cfg.Addr, router)
	if err != nil {
		panic(err)
	}
}
