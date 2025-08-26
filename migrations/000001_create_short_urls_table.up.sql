CREATE TABLE short_urls (
    id SERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    original_url TEXT NOT NULL
);
