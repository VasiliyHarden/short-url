package handler

import (
	"encoding/json"
	"github.com/VasiliyHarden/short-url/internal/service/shortener"
	"mime"
	"net/http"
)

type ShortenBatchRequestPayload struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ShortenBatchResponsePayload struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func ShortenBatchJSON(sh *shortener.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || (mediaType != "application/json" && mediaType != "application/x-gzip") {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		var payload []ShortenBatchRequestPayload
		if err := dec.Decode(&payload); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		originalURLs := make([]string, len(payload))
		for i, item := range payload {
			originalURLs[i] = item.OriginalURL
		}

		result, err := sh.GenerateBatch(r.Context(), originalURLs)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		response := make([]ShortenBatchResponsePayload, len(payload))
		for i, item := range payload {
			response[i] = ShortenBatchResponsePayload{
				CorrelationID: item.CorrelationID,
				ShortURL:      result[i],
			}
		}

		respBytes, err := json.Marshal(response)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(respBytes)
	}
}
