package shortener

import (
	"context"
	"errors"
)

type Service struct {
	gen     Generator
	store   Storage
	baseURL string
}

var ErrDuplicate = errors.New("duplicate")

func NewService(baseURL string, gen Generator, store Storage) *Service {
	return &Service{
		gen:     gen,
		store:   store,
		baseURL: baseURL,
	}
}

func (s *Service) Generate(url string) (string, error) {
	code := s.gen.Generate(url)
	existingCode, inserted, err := s.store.Save(code, url)
	if err != nil {
		return "", err
	}
	if !inserted {
		return s.baseURL + "/" + existingCode, ErrDuplicate
	}

	return s.baseURL + "/" + code, nil
}

func (s *Service) GenerateBatch(ctx context.Context, urls []string) ([]string, error) {
	batch := make([]BatchItem, len(urls))
	for i, url := range urls {
		code := s.gen.Generate(url)
		batch[i] = BatchItem{Code: code, OriginalURL: url}
	}

	result := make([]string, len(batch))
	for i, batchItem := range batch {
		result[i] = s.baseURL + "/" + batchItem.Code
	}

	err := s.store.SaveBatch(ctx, batch)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) Resolve(code string) (string, bool) {
	return s.store.Get(code)
}
