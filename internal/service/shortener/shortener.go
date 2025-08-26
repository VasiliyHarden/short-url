package shortener

import "context"

type Service struct {
	gen     Generator
	store   Storage
	baseURL string
}

func NewService(baseURL string, gen Generator, store Storage) *Service {
	return &Service{
		gen:     gen,
		store:   store,
		baseURL: baseURL,
	}
}

func (s *Service) Generate(url string) string {
	code := s.gen.Generate(url)
	s.store.Save(code, url)

	return s.baseURL + "/" + code
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
