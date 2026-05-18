-- +goose up
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    library_id UUID NOT NULL,
    FOREIGN KEY (library_id) REFERENCES library(id) ON DELETE CASCADE
);

-- +goose down
DROP TABLE IF EXISTS users;