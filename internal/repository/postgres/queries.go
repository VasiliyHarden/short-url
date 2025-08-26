package postgres

const insertOnConflictShortURL = `
INSERT INTO short_urls (code, original_url)
VALUES ($1, $2)
ON CONFLICT (original_url)
DO UPDATE SET code = short_urls.code
RETURNING code, (xmax = 0) AS inserted;
`
