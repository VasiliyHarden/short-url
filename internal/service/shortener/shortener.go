package shortener

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

func (s *Service) Resolve(code string) (string, bool) {
	return s.store.Get(code)
}
