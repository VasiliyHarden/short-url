package main

import (
	"net/http"

	"github.com/VasiliyHarden/short-url/internal/handler"
)

func main() {
	router := handler.NewRouter()

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
