-- name: CreateBook :one
INSERT INTO books (id, title, author, published_date, isbn, collection_id, created_at, updated_at, library_id)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5,
    NOW(),
    NOW(),
    $6
)
RETURNING *;

-- name: GetBookByID :one
SELECT * FROM books
WHERE id = $1;

-- name: GetAllBooksFromCollection :many
SELECT * FROM books
WHERE collection_id = $1
ORDER BY created_at DESC;

-- name: UpdateBook :one
UPDATE books
SET title = $2,
    author = $3,
    published_date = $4,
    isbn = $5,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteBook :exec
DELETE FROM books
WHERE id = $1;