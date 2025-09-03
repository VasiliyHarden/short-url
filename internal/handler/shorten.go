package handler

import (
	"github.com/VasiliyHarden/short-url/internal/service/shortener"
	"go.uber.org/zap"
	"io"
	"mime"
	"net/http"
)

func writeShortURLText(w http.ResponseWriter, status int, url string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(url))
}

func Shorten(sh *shortener.Service, logger *zap.Logger) http.HandlerFunc {
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
		defer r.Body.Close()

		shortURL, err := sh.Generate(string(body))
		respondShortURL(w, shortURL, err, writeShortURLText, logger)
	}
}
