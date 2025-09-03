package shortener

import (
	"context"
	"sync"
)

type memoryStorage struct {
	codeToURL map[string]string
	urlToCode map[string]string
	mu        sync.RWMutex
}

func NewMemoryStorage() Storage {
	return &memoryStorage{
		codeToURL: make(map[string]string),
		urlToCode: make(map[string]string),
	}
}

func (m *memoryStorage) Save(code, originalURL string) (retCode string, inserted bool, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if existingCode, ok := m.urlToCode[originalURL]; ok {
		return existingCode, false, nil
	}

	m.codeToURL[code] = originalURL
	m.urlToCode[originalURL] = code
	return code, true, nil
}

func (m *memoryStorage) SaveBatch(_ context.Context, batch []BatchItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, batchItem := range batch {
		m.codeToURL[batchItem.Code] = batchItem.OriginalURL
		m.urlToCode[batchItem.OriginalURL] = batchItem.Code
	}
	return nil
}

func (m *memoryStorage) Get(code string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok := m.codeToURL[code]
	return value, ok
}

func (m *memoryStorage) Close() error {
	return nil
}
