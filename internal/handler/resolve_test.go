package handler

import (
	"github.com/VasiliyHarden/short-url/internal/service/shortener"
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
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "wrong method",
			method: http.MethodPost,
			code:   "/dummy",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var code string
			if tc.code != "" {
				code = tc.code
			} else {
				full := shortener.Generate(tc.want.location)
				code = strings.TrimPrefix(full, shortener.BaseURL+"/")
			}

			r := httptest.NewRequest(tc.method, "/"+code, nil)
			r.Header.Set("Content-Type", "text/plain")
			w := httptest.NewRecorder()

			Resolve(w, r)

			assert.Equal(t, tc.want.code, w.Code)

			if tc.want.code == http.StatusTemporaryRedirect {
				assert.Equal(t, w.Header().Get("Location"), tc.want.location)
			}
		})
	}
}
