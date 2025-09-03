package handler

import (
	"github.com/VasiliyHarden/short-url/internal/service/shortener"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShorten(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name        string
		contentType string
		method      string
		want        want
	}{
		{
			name:        "success",
			contentType: "text/plain; charset=utf-8",
			method:      http.MethodPost,
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain",
			},
		},
		{
			name:        "wrong method",
			contentType: "text/plain; charset=utf-8",
			method:      http.MethodGet,
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name:        "wrong content-type",
			contentType: "application/json",
			method:      http.MethodPost,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			const baseURL = "localhost:8080"

			gen := shortener.NewHashGenerator()
			store := shortener.NewMemoryStorage()
			sh := shortener.NewService(baseURL, gen, store)

			router := chi.NewRouter()
			logger := zap.NewNop()
			router.Post("/", Shorten(sh, logger))

			r := httptest.NewRequest(tc.method, "/", strings.NewReader("https://example.com"))
			r.Header.Set("Content-Type", tc.contentType)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)

			assert.Equal(t, tc.want.code, w.Code)

			if tc.want.code == http.StatusCreated {
				assert.True(t, strings.HasPrefix(w.Body.String(), baseURL+"/"))
				assert.Equal(t, tc.want.contentType, w.Header().Get("Content-Type"))
			}

			if tc.want.code == http.StatusBadRequest {
				assert.Equal(t, http.StatusText(tc.want.code), strings.TrimSpace(w.Body.String()))
				assert.Equal(t, tc.want.contentType, w.Header().Get("Content-Type"))
			}
		})
	}
}
