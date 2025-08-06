package handler

import (
	"github.com/VasiliyHarden/short-url/internal/service/shortener"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResolve(t *testing.T) {
	type want struct {
		code     int
		location string
	}
	tests := []struct {
		name   string
		method string
		code   string
		want   want
	}{
		{
			name:   "success",
			method: http.MethodGet,
			want: want{
				code:     http.StatusTemporaryRedirect,
				location: "https://example.com",
			},
		},
		{
			name:   "non-existing code",
			method: http.MethodGet,
			code:   "/dummy",
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			name:   "wrong method",
			method: http.MethodPost,
			code:   "/dummy",
			want: want{
				code: http.StatusNotFound,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			const baseURL = "localhost:8080"

			gen := shortener.NewHashGenerator()
			store := shortener.NewMemoryStorage()
			sh := shortener.NewService(baseURL, gen, store)

			var code string
			if tc.code != "" {
				code = tc.code
			} else {
				full := sh.Generate(tc.want.location)
				code = strings.TrimPrefix(full, baseURL+"/")
			}

			router := chi.NewRouter()
			router.Get("/{code}", Resolve(sh))

			r := httptest.NewRequest(tc.method, "/"+code, nil)
			r.Header.Set("Content-Type", "text/plain")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)

			assert.Equal(t, tc.want.code, w.Code)

			if tc.want.code == http.StatusTemporaryRedirect {
				assert.Equal(t, w.Header().Get("Location"), tc.want.location)
			}
		})
	}
}
