-- name: CreateCollection :one
INSERT INTO collections (id, name, description, created_at, updated_at, library_id)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    NOW(),
    NOW(),
    $3
)
RETURNING *;
UPDATE library
SET collection_count = collection_count + 1
WHERE id = $3;

-- name: GetCollectionByID :one
SELECT * FROM collections
WHERE id = $1;

-- name: GetAllCollectionsFromLibrary :many
SELECT * FROM collections
WHERE library_id = $1
ORDER BY created_at DESC;

-- name: UpdateCollection :one
UPDATE collections
SET name = $2,
    description = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCollection :exec
UPDATE library
SET collection_count = collection_count - 1
WHERE id = (SELECT b.library_id FROM collections b WHERE b.id = $1);
DELETE FROM collections
WHERE id = $1;