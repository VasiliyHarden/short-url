package handler

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/{code}", Resolve)
	r.Post("/", Shorten)

	return r
}
