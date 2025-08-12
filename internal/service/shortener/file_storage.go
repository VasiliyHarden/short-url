package shortener

import (
	"encoding/json"
	"io"
	"os"
	"sync"
)

type fileStorage struct {
	data map[string]string
	path string
	mu   sync.RWMutex
}

type record struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewFileStorage(path string) (Storage, error) {
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fs := &fileStorage{
		data: make(map[string]string),
		path: path,
	}

	var records []record

	if err := json.NewDecoder(f).Decode(&records); err != nil && err != io.EOF {
		return nil, err
	}

	for _, rec := range records {
		if rec.ShortURL != "" && rec.OriginalURL != "" {
			fs.data[rec.ShortURL] = rec.OriginalURL
		}
	}

	return fs, nil
}

func (fs *fileStorage) Save(code, originalURL string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fs.data[code] = originalURL

	records := make([]record, 0, len(fs.data))
	for short, orig := range fs.data {
		records = append(records, record{ShortURL: short, OriginalURL: orig})
	}

	f, err := os.OpenFile(fs.path, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(records); err != nil {
		return err
	}

	return f.Sync()
}

func (fs *fileStorage) Get(code string) (string, bool) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	value, ok := fs.data[code]
	return value, ok
}
