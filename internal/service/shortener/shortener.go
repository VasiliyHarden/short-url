package shortener

import (
	"crypto/sha256"
	"encoding/base64"
)

var BaseURL string

func Init(base string) {
	BaseURL = base
}

var store = make(map[string]string)

func Generate(url string) string {
	hash := sha256.Sum256([]byte(url))
	code := base64.RawURLEncoding.EncodeToString(hash[:6])
	store[code] = url

	return BaseURL + "/" + code
}

func Resolve(code string) (string, bool) {
	value, ok := store[code]

	return value, ok
}
