package handler

import (
	"errors"
	"github.com/VasiliyHarden/short-url/internal/service/shortener"
	"io"
	"mime"
	"net/http"
)

func Shorten(sh *shortener.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || (mediaType != "text/plain" && mediaType != "application/x-gzip") {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		shortURL, err := sh.Generate(string(body))
		if err != nil {
			if errors.Is(err, shortener.ErrDuplicate) {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusConflict)
				_, _ = io.WriteString(w, shortURL)
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, shortURL)

		_ = r.Body.Close()
	}
}
