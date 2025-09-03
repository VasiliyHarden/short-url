package shortener

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"
)

type fileStorage struct {
	codeToURL map[string]string
	urlToCode map[string]string
	path      string
	mu        sync.RWMutex
	f         *os.File
}

type record struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewFileStorage(path string) (Storage, error) {
	fs := &fileStorage{
		codeToURL: make(map[string]string),
		urlToCode: make(map[string]string),
		path:      path,
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
				fs.codeToURL[rec.ShortURL] = rec.OriginalURL
				fs.urlToCode[rec.OriginalURL] = rec.ShortURL
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

func (fs *fileStorage) Save(code, originalURL string) (retCode string, inserted bool, err error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if existingCode, ok := fs.urlToCode[originalURL]; ok {
		return existingCode, false, nil
	}

	fs.codeToURL[code] = originalURL
	fs.urlToCode[originalURL] = code

	rec := record{
		ShortURL:    code,
		OriginalURL: originalURL,
	}

	if err := json.NewEncoder(fs.f).Encode(&rec); err != nil {
		return "", false, err
	}

	return code, true, fs.f.Sync()
}

func (fs *fileStorage) SaveBatch(_ context.Context, batch []BatchItem) error {
	for _, batchItem := range batch {
		_, _, err := fs.Save(batchItem.Code, batchItem.OriginalURL)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fs *fileStorage) Get(code string) (string, bool) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	value, ok := fs.codeToURL[code]
	return value, ok
}

func (fs *fileStorage) Close() error {
	err := fs.f.Close()

	return err
}
