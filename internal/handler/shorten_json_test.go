package handler

import (
	"encoding/json"
	"github.com/VasiliyHarden/short-url/internal/service/shortener"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortenJSON(t *testing.T) {
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
			contentType: "application/json; charset=utf-8",
			method:      http.MethodPost,
			want: want{
				code:        http.StatusCreated,
				contentType: "application/json",
			},
		},
		{
			name:        "wrong method",
			contentType: "application/json; charset=utf-8",
			method:      http.MethodGet,
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name:        "wrong content-type",
			contentType: "text/plain",
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
			router.Post("/api/shorten", ShortenJSON(sh))

			r := httptest.NewRequest(tc.method, "/api/shorten", strings.NewReader(`{"url":"https://example.com"}`))
			r.Header.Set("Content-Type", tc.contentType)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)

			assert.Equal(t, tc.want.code, w.Code)

			if tc.want.code == http.StatusCreated {
				var resp struct {
					Result string `json:"result"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.True(t, strings.HasPrefix(resp.Result, baseURL+"/"))
				assert.Equal(t, tc.want.contentType, w.Header().Get("Content-Type"))
			}

			if tc.want.code == http.StatusBadRequest {
				assert.Equal(t, http.StatusText(tc.want.code), strings.TrimSpace(w.Body.String()))
				assert.Equal(t, tc.want.contentType, w.Header().Get("Content-Type"))
			}
		})
	}
}
