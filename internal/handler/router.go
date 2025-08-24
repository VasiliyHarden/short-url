package handler

import (
	"github.com/VasiliyHarden/short-url/internal/handler/middleware"
	"github.com/VasiliyHarden/short-url/internal/service/shortener"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

func NewRouter(sh *shortener.Service, logger *zap.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Decompress)
	r.Use(middleware.Logging(logger))
	r.Use(middleware.Compress)

	r.Get("/{code}", Resolve(sh))
	r.Post("/", Shorten(sh))
	r.Post("/api/shorten", ShortenJSON(sh))

	return r
}
