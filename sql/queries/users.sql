-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, apikey)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByApiKey :one
SELECT * FROM users
WHERE apikey = $1 LIMIT 1;