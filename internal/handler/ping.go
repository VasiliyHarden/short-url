package handler

import (
	"github.com/VasiliyHarden/short-url/internal/repository/postgres"
	"net/http"
)

func Ping(db *postgres.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := db.Ping(r.Context())

		if err != nil {
			http.Error(w, "db unreachable", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}
}
