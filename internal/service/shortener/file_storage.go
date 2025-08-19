package shortener

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"
)

type fileStorage struct {
	data map[string]string
	path string
	mu   sync.RWMutex
	f    *os.File
}

type record struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewFileStorage(path string) (Storage, error) {
	fs := &fileStorage{
		data: make(map[string]string),
		path: path,
	}

	readFile, err := os.Open(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	} else {
		defer readFile.Close()

		dec := json.NewDecoder(readFile)
		for {
			var rec record
			if err := dec.Decode(&rec); err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}

			if rec.ShortURL != "" && rec.OriginalURL != "" {
				fs.data[rec.ShortURL] = rec.OriginalURL
			}
		}
	}

	writeFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	fs.f = writeFile

	return fs, nil
}

func (fs *fileStorage) Save(code, originalURL string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fs.data[code] = originalURL

	rec := record{
		ShortURL:    code,
		OriginalURL: originalURL,
	}

	if err := json.NewEncoder(fs.f).Encode(&rec); err != nil {
		return err
	}

	return fs.f.Sync()
}

func (fs *fileStorage) Get(code string) (string, bool) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	value, ok := fs.data[code]
	return value, ok
}

func (fs *fileStorage) Close() error {
	err := fs.f.Close()

	return err
}
