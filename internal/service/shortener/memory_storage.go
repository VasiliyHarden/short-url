package shortener

type memoryStorage struct {
	data map[string]string
}

func NewMemoryStorage() Storage {
	return &memoryStorage{
		data: make(map[string]string),
	}
}

func (m *memoryStorage) Save(code, originalURL string) error {
	m.data[code] = originalURL

	return nil
}

func (m *memoryStorage) Get(code string) (string, bool) {
	value, ok := m.data[code]

	return value, ok
}
