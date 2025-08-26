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

func TestShortenBatchJSON(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}

	tests := []struct {
		name        string
		contentType string
		method      string
		body        string
		want        want
	}{
		{
			name:        "success",
			contentType: "application/json; charset=utf-8",
			method:      http.MethodPost,
			body: `[
                {"correlation_id":"1","original_url":"https://example.com/a"},
                {"correlation_id":"2","original_url":"https://example.com/b"}
            ]`,
			want: want{
				code:        http.StatusCreated,
				contentType: "application/json",
			},
		},
		{
			name:        "wrong method",
			contentType: "application/json; charset=utf-8",
			method:      http.MethodGet,
			body:        `[]`,
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name:        "wrong content-type",
			contentType: "text/plain",
			method:      http.MethodPost,
			body:        `[]`,
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
			router.Post("/api/shorten/batch", ShortenBatchJSON(sh))

			r := httptest.NewRequest(tc.method, "/api/shorten/batch", strings.NewReader(tc.body))
			r.Header.Set("Content-Type", tc.contentType)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)

			assert.Equal(t, tc.want.code, w.Code)

			if tc.want.code == http.StatusCreated {
				var resp []struct {
					CorrelationID string `json:"correlation_id"`
					ShortURL      string `json:"short_url"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)

				require.Len(t, resp, 2)
				assert.Equal(t, "1", resp[0].CorrelationID)
				assert.Equal(t, "2", resp[1].CorrelationID)
				assert.True(t, strings.HasPrefix(resp[0].ShortURL, baseURL+"/"))
				assert.True(t, strings.HasPrefix(resp[1].ShortURL, baseURL+"/"))
				assert.Equal(t, tc.want.contentType, w.Header().Get("Content-Type"))
			}

			if tc.want.code == http.StatusBadRequest {
				assert.Equal(t, http.StatusText(tc.want.code), strings.TrimSpace(w.Body.String()))
				assert.Equal(t, tc.want.contentType, w.Header().Get("Content-Type"))
			}
		})
	}
}
