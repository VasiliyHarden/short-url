package handler

import (
	"github.com/VasiliyHarden/short-url/internal/service/shortener"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Resolve(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	url, ok := shortener.Resolve(code)
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
