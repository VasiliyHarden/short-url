package shortener

import (
	"crypto/sha256"
	"encoding/base64"
)

type Generator interface {
	Generate(originalURL string) string
}

type HashGenerator struct{}

func NewHashGenerator() *HashGenerator {
	return &HashGenerator{}
}

func (g *HashGenerator) Generate(originalURL string) string {
	hash := sha256.Sum256([]byte(originalURL))
	code := base64.RawURLEncoding.EncodeToString(hash[:6])

	return code
}
