package shortener

type Storage interface {
	Save(code, originalURL string) error
	Get(code string) (string, bool)
}
