package handler

import (
	"errors"
	"github.com/VasiliyHarden/short-url/internal/service/shortener"
	"go.uber.org/zap"
	"net/http"
)

type shortURLWriter func(w http.ResponseWriter, status int, shortURL string)

func respondShortURL(w http.ResponseWriter, shortURL string, err error, write shortURLWriter, logger *zap.Logger) {
	if err != nil {
		if errors.Is(err, shortener.ErrDuplicate) {
			logger.Info("duplicate URL encountered", zap.String("shortURL", shortURL))
			write(w, http.StatusConflict, shortURL)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	write(w, http.StatusCreated, shortURL)
}
