-- name: CreateUser :one
INSERT INTO users (id, username, email, password_hash, created_at, updated_at, library_id)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    NOW(),
    NOW(),
    $4
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $1, updated_at = NOW()
WHERE id = $2;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;