-- +goose up
CREATE TABLE IF NOT EXISTS books (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    published_date DATE,
    isbn TEXT UNIQUE,
    library_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (library_id) REFERENCES library(id) ON DELETE CASCADE
);

-- +goose down
DROP TABLE IF EXISTS books;