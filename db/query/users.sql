-- name: CreateUser :one
INSERT INTO users (
  id, first_name, last_name, email, hash_password, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, first_name, last_name, email, hash_password
FROM users
WHERE email = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
set first_name = $2,
last_name = $3,
email = $4,
hash_password = $5
WHERE id = $1
RETURNING *;