-- +goose up
CREATE TABLE IF NOT EXISTS collections (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    book_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    library_id UUID NOT NULL,
    FOREIGN KEY (library_id) REFERENCES library(id) ON DELETE CASCADE
);

ALTER TABLE books
ADD COLUMN collection_id UUID,
ADD FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE SET NULL;

ALTER TABLE library
ADD COLUMN collection_count INTEGER NOT NULL DEFAULT 0;

-- +goose down
ALTER TABLE library
DROP COLUMN IF EXISTS collection_count;

ALTER TABLE books
DROP COLUMN IF EXISTS collection_id;

DROP TABLE IF EXISTS collections;