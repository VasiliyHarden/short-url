package middleware

import (
	"github.com/VasiliyHarden/short-url/internal/handler/middleware/encoding"
	"net/http"
	"strings"
)

func Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")

		if supportsGzip {
			gw := encoding.NewGzipWriter(ow)

			ow = gw

			defer gw.Close()
		}

		next.ServeHTTP(ow, r)
	})
}
