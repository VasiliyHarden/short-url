package shortener

import "context"

type BatchItem struct {
	Code        string
	OriginalURL string
}

type Storage interface {
	Save(code, originalURL string) (retCode string, inserted bool, err error)
	SaveBatch(ctx context.Context, batch []BatchItem) error
	Get(code string) (string, bool)
	Close() error
}
