package handler

import (
	"encoding/json"
	"github.com/VasiliyHarden/short-url/internal/service/shortener"
	"mime"
	"net/http"
)

type ShortenRequestPayload struct {
	URL string `json:"url"`
}

type ShortenResponsePayload struct {
	Result string `json:"result"`
}

func ShortenJSON(sh *shortener.Service) http.HandlerFunc {
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

		respBytes, err := json.Marshal(ShortenResponsePayload{Result: sh.Generate(payload.URL)})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(respBytes)
	}
}
