-- name: CreateLibrary :one
INSERT INTO library (name, description)
VALUES ($1, $2)
RETURNING *;

-- name: GetLibraryByID :one
SELECT * FROM library
WHERE id = $1;