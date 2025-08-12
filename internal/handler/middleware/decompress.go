package middleware

import (
	"github.com/VasiliyHarden/short-url/internal/handler/middleware/encoding"
	"net/http"
	"strings"
)

func Decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")

		if sendsGzip {
			gr, err := encoding.NewGzipReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = gr
			r.Header.Del("Content-Encoding")
			r.Header.Del("Content-Length")

			defer gr.Close()
		}

		next.ServeHTTP(w, r)
	})
}
