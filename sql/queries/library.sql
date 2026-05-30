-- name: CreateLibrary :one
INSERT INTO library (id, name, description, created_at, updated_at, user_id)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    NOW(),
    NOW(),
    $3
)
RETURNING *;

-- name: GetLibraryByID :one
SELECT * FROM library
WHERE id = $1;

-- name: GetAllLibraries :many
SELECT * FROM library
ORDER BY created_at DESC;

-- name: UpdateLibrary :one
UPDATE library
SET name = $2,
    description = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteLibrary :exec
DELETE FROM library
WHERE id = $1;

-- name: DeleteAllLibraries :exec
DELETE FROM library;