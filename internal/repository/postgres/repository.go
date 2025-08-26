package postgres

func (repo *DB) Save(code, originalURL string) error {
	_, err := repo.db.Exec("INSERT INTO short_urls (code, original_url) VALUES ($1, $2)", code, originalURL)
	return err
}

func (repo *DB) Get(code string) (string, bool) {
	row := repo.db.QueryRow("SELECT original_url FROM short_urls WHERE code = $1", code)
	var originalURL string
	err := row.Scan(&originalURL)

	if err != nil {
		return "", false
	}

	return originalURL, true
}

func (repo *DB) Close() error {
	if repo.db == nil {
		return nil
	}

	err := repo.db.Close()
	if err != nil {
		return err
	}

	return nil
}
