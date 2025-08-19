package shortener

import "sync"

type memoryStorage struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewMemoryStorage() Storage {
	return &memoryStorage{
		data: make(map[string]string),
	}
}

func (m *memoryStorage) Save(code, originalURL string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[code] = originalURL
	return nil
}

func (m *memoryStorage) Get(code string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok := m.data[code]
	return value, ok
}

func (m *memoryStorage) Close() error {
	return nil
}
