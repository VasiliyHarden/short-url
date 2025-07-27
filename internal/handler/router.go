package handler

import (
	"github.com/VasiliyHarden/short-url/internal/service/shortener"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewRouter(sh *shortener.Service) http.Handler {
	r := chi.NewRouter()
	r.Get("/{code}", Resolve(sh))
	r.Post("/", Shorten(sh))

	return r
}
