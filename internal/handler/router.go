package handler

import (
	"net/http"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			Shorten(w, r)
			return
		}

		Resolve(w, r)
	})

	return mux
}
