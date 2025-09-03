package handler

import (
	"encoding/json"
	"github.com/VasiliyHarden/short-url/internal/service/shortener"
	"go.uber.org/zap"
	"mime"
	"net/http"
)

type ShortenRequestPayload struct {
	URL string `json:"url"`
}

type ShortenResponsePayload struct {
	Result string `json:"result"`
}

func writeShortURLJSON(w http.ResponseWriter, status int, shortURL string) {
	response, _ := json.Marshal(ShortenResponsePayload{Result: shortURL})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(response)
}

func ShortenJSON(sh *shortener.Service, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || (mediaType != "application/json" && mediaType != "application/x-gzip") {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		var payload ShortenRequestPayload
		if err := dec.Decode(&payload); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		shortURL, err := sh.Generate(payload.URL)
		respondShortURL(w, shortURL, err, writeShortURLJSON, logger)
	}
}
