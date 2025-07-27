package shortener

type Storage interface {
	Save(code, originalURL string)
	Get(code string) (string, bool)
}

type MemoryStorage struct {
	data map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]string),
	}
}

func (m *MemoryStorage) Save(code, originalURL string) {
	m.data[code] = originalURL
}

func (m *MemoryStorage) Get(code string) (string, bool) {
	value, ok := m.data[code]

	return value, ok
}
